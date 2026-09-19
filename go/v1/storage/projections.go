package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/traust-security/traust-sdk/go/v1/types"
)

type projectionName string

type projectionField string

const (
	projectionFieldOriginalID projectionField = "orig_id"
	projectionFieldRow        projectionField = "row"
)

func indexedProjectionField(field projectionField, index int) projectionField {
	return projectionField(fmt.Sprintf("%s[%d]", field, index))
}

func projectionText(value any) (string, error) {
	text, ok := value.(string)
	if !ok {
		return "", errors.New("expected string")
	}
	return text, nil
}

func optionalProjectionText(value any) (*string, error) {
	if value == nil {
		return nil, nil
	}
	if text, ok := value.(*string); ok {
		return text, nil
	}
	text, err := projectionText(value)
	if err != nil {
		return nil, err
	}
	return &text, nil
}

func projectionJSON(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("encode JSON: %w", err)
	}
	return string(payload), nil
}

func optionalProjectionJSON(value any) (*string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}
	if string(payload) == "null" {
		return nil, nil
	}
	text := string(payload)
	return &text, nil
}

func projectionInteger(value any) (*int64, error) {
	if value == nil {
		return nil, nil
	}
	var text string
	switch number := value.(type) {
	case json.Number:
		text = number.String()
	case float64:
		if math.IsInf(number, 0) || math.IsNaN(number) || math.Trunc(number) != number {
			return nil, errors.New("expected integer")
		}
		text = strconv.FormatFloat(number, 'f', -1, 64)
	case int:
		integer := int64(number)
		return &integer, nil
	case int64:
		return &number, nil
	default:
		return nil, errors.New("expected integer")
	}
	integer, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		floating, floatErr := strconv.ParseFloat(text, 64)
		if floatErr != nil || math.Trunc(floating) != floating || floating > math.MaxInt64 || floating < math.MinInt64 {
			return nil, errors.New("expected integer")
		}
		integer = int64(floating)
	}
	return &integer, nil
}

func optionalInt(value *int) *int64 {
	if value == nil {
		return nil
	}
	converted := int64(*value)
	return &converted
}

func (s *sqlStore) projectLayer(
	ctx context.Context,
	conn *sql.Conn,
	state writeState,
	layer types.Layer,
) error {
	if err := s.queries.layerMetadataUpsert(ctx, conn, layerMetadataUpsertParams{
		bindingId:      state.bindingID,
		artifactDigest: state.digest,
		repo:           &layer.Metadata.Repository,
		createdAt:      &layer.Metadata.Created,
		merkleRoot:     layer.Metadata.MerkleRoot,
		merkleEpoch:    optionalInt(layer.Metadata.MerkleEpoch),
	}); err != nil {
		return projectionError(projectionLayerMetadata, projectionFieldRow, err)
	}
	// The events are the TIME DIMENSION. A Go caller that saves a layer
	// without them writes the merkle root and silently loses every date the
	// trend, MTTR and SLA views are computed from.
	for _, event := range layer.Events {
		source := string(event.Source.Type)
		actorKind := string(event.Source.Actor.Kind)
		var validity, resolution *string
		if event.Disposition.Validity != nil {
			value := string(*event.Disposition.Validity)
			validity = &value
		}
		if event.Disposition.Resolution != nil {
			value := string(*event.Disposition.Resolution)
			resolution = &value
		}
		if err := s.queries.layerEventUpsert(ctx, conn, layerEventUpsertParams{
			bindingId:       state.bindingID,
			artifactDigest:  state.digest,
			eventId:         event.EventId,
			findingRef:      event.FindingRef,
			fingerprint:     event.Fingerprint,
			fingerprintAlgo: event.FingerprintAlgo,
			recordedAt:      event.RecordedAt,
			occurredAt:      event.OccurredAt,
			sourceType:      &source,
			sourceRef:       &event.Source.Ref,
			actorKind:       &actorKind,
			validity:        validity,
			resolution:      resolution,
			evidenceGrade:   event.EvidenceGrade,
			autoAcceptTier:  optionalBoolAsInt(event.AutoAcceptTier),
		}); err != nil {
			return projectionError(projectionLayerEvent, projectionFieldRow, err)
		}
	}
	return nil
}

func (s *sqlStore) projectVulnFindings(
	ctx context.Context,
	conn *sql.Conn,
	state writeState,
	document types.VulnFindings,
) error {
	for index, finding := range document.Findings {
		line, err := projectionInteger(finding.Line)
		if err != nil {
			return projectionError(projectionFinding, indexedProjectionField(projectionFieldLine, index), err)
		}
		category := finding.Category
		if err := s.queries.findingUpsert(ctx, conn, findingUpsertParams{
			bindingId:      state.bindingID,
			artifactDigest: state.digest,
			findingId:      finding.Id,
			target:         document.Target,
			scannedAt:      document.ScannedAt,
			title:          finding.Title,
			severity:       string(finding.Severity),
			description:    finding.Description,
			category:       &category,
			file:           finding.File,
			line:           line,
			cwe:            finding.Cwe,
			recommendation: finding.Recommendation,
			confidence:     finding.Confidence,
		}); err != nil {
			return projectionError(projectionFinding, indexedProjectionField(projectionFieldRow, index), err)
		}
	}
	return nil
}

func (s *sqlStore) projectTriage(
	ctx context.Context,
	conn *sql.Conn,
	state writeState,
	document types.Triage,
) error {
	for index, finding := range document.Findings {
		sourceFindingID, err := optionalProjectionText(finding.OrigId)
		if err != nil {
			return projectionError(projectionTriageVerdict, indexedProjectionField(projectionFieldOriginalID, index), err)
		}
		severity, err := optionalProjectionText(finding.Severity)
		if err != nil {
			return projectionError(projectionTriageVerdict, indexedProjectionField(projectionFieldSeverity, index), err)
		}
		voteBreakdown, err := optionalProjectionJSON(finding.VoteBreakdown)
		if err != nil {
			return projectionError(projectionTriageVerdict, indexedProjectionField(projectionFieldVoteBreakdown, index), err)
		}
		rationale, err := optionalProjectionText(finding.Rationale)
		if err != nil {
			return projectionError(projectionTriageVerdict, indexedProjectionField(projectionFieldRationale, index), err)
		}
		if err := s.queries.triageVerdictUpsert(ctx, conn, triageVerdictUpsertParams{
			bindingId:       state.bindingID,
			artifactDigest:  state.digest,
			findingId:       finding.Id,
			sourceFindingId: sourceFindingID,
			triageCompleted: document.TriageCompleted,
			verdict:         string(finding.Verdict),
			severity:        severity,
			voteBreakdown:   voteBreakdown,
			rationale:       rationale,
		}); err != nil {
			return projectionError(projectionTriageVerdict, indexedProjectionField(projectionFieldRow, index), err)
		}
	}
	return nil
}

// projectCorpusRegistry fans a registry out to one subject_ownership row per
// subject. Hand-written because the generator's one-row projector maps ROOT
// schema properties to columns, and these live inside subjects[].
//
// Ownership is the denominator every dashboard cut divides by, so a Go
// caller that saves a registry without this writes the evidence and leaves
// every owned-scoped query empty.
func (s *sqlStore) projectCorpusRegistry(
	ctx context.Context,
	conn *sql.Conn,
	state writeState,
	registry types.CorpusRegistry,
) error {
	for _, subject := range registry.Subjects {
		var refKind *string
		if subject.RefKind != nil {
			value := string(*subject.RefKind)
			refKind = &value
		}
		if err := s.queries.subjectOwnershipUpsert(ctx, conn, subjectOwnershipUpsertParams{
			bindingId:      state.bindingID,
			artifactDigest: state.digest,
			subjectId:      subject.SubjectId,
			tree:           subject.Tree,
			ownership:      string(subject.Ownership),
			businessUnit:   subject.BusinessUnit,
			label:          subject.Label,
			product:        subject.Product,
			repoUrl:        subject.RepoUrl,
			ref:            subject.Ref,
			refKind:        refKind,
			isBranchAudit:  optionalBoolAsInt(subject.IsBranchAudit),
		}); err != nil {
			return projectionError(projectionSubjectOwnership, projectionFieldRow, err)
		}
	}
	return nil
}

// optionalBoolAsInt keeps absent ABSENT. storage/v1 stores flags as INTEGER
// on both dialects, and "not a branch audit" must not be confused with "the
// registry does not say".
func optionalBoolAsInt(value *bool) *int64 {
	if value == nil {
		return nil
	}
	stored := int64(0)
	if *value {
		stored = 1
	}
	return &stored
}

// projectThreatRegister fans a register out to one `threat` row per threat.
// Hand-written for the same reason projectCorpusRegistry is: the generator's
// one-row projector maps ROOT schema properties to columns, and these live
// inside threats[].
//
// Keyed on Key, never Id: every threat model numbers its threats from T1, so
// Id collides across the whole model set and an Id-keyed write would keep one
// threat per number out of tens of thousands.
func (s *sqlStore) projectThreatRegister(
	ctx context.Context,
	conn *sql.Conn,
	state writeState,
	register types.ThreatRegister,
) error {
	for _, threat := range register.Threats {
		actors, err := optionalProjectionJSON(threat.Actors)
		if err != nil {
			return projectionError(projectionThreat, projectionFieldActors, err)
		}
		// Kept even when empty: an empty list means MODELLED BUT NOT
		// EVIDENCED, which is a different claim from unmitigated, and
		// threat_exposure reports the two apart.
		evidence, err := optionalProjectionJSON(threat.Evidence)
		if err != nil {
			return projectionError(projectionThreat, projectionFieldEvidence, err)
		}
		dimensions, err := optionalProjectionJSON(threat.IsolationDimensions)
		if err != nil {
			return projectionError(projectionThreat, projectionFieldIsolationDimensions, err)
		}
		boundaries, err := optionalProjectionJSON(threat.IsolationBoundaries)
		if err != nil {
			return projectionError(projectionThreat, projectionFieldIsolationBoundaries, err)
		}
		impact := string(threat.Impact)
		likelihood := string(threat.Likelihood)
		status := string(threat.Status)
		statement := threat.Threat
		if err := s.queries.threatUpsert(ctx, conn, threatUpsertParams{
			bindingId:           state.bindingID,
			artifactDigest:      state.digest,
			threatKey:           threat.Key,
			threatId:            threat.Id,
			model:               threat.Model,
			subjectId:           threat.SubjectId,
			product:             threat.Product,
			statement:           &statement,
			surface:             threat.Surface,
			asset:               threat.Asset,
			impact:              &impact,
			likelihood:          &likelihood,
			status:              &status,
			controls:            threat.Controls,
			actors:              actors,
			evidence:            evidence,
			linddun:             optionalBoolAsInt(threat.Linddun),
			score:               optionalIntAsInt64(threat.Score),
			isolationDimensions: dimensions,
			isolationBoundaries: boundaries,
		}); err != nil {
			return projectionError(projectionThreat, projectionFieldRow, err)
		}
	}
	return nil
}

// optionalIntAsInt64 keeps absent ABSENT, matching optionalBoolAsInt: a
// missing score is not a score of zero.
func optionalIntAsInt64(value *int) *int64 {
	if value == nil {
		return nil
	}
	stored := int64(*value)
	return &stored
}
