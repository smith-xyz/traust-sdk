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
