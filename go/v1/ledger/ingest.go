package ledger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/traust-security/traust-sdk/go/v1/ledger/internal/check"
	"github.com/traust-security/traust-sdk/go/v1/types"
	"github.com/traust-security/traust-sdk/go/v1/validate"
)

type submissionOp[In any] struct {
	kind IngestKind
}

func (o submissionOp[In]) Submit(ctx context.Context, p Provider, in In) (SubmitResponse, error) {
	var zero SubmitResponse

	raw, err := json.Marshal(in)
	if err != nil {
		return zero, ingestErr(o.kind, PhaseMarshal, err)
	}

	resp, err := p.Ingest(ctx, IngestMeta{
		Kind:             o.kind,
		ContractsVersion: types.ContractsVersion,
	}, raw)
	return decodeSubmitResponse(resp, err, o.kind)
}

type reportSubmissionOp[In any] struct {
	kind         IngestKind
	reportSchema string
	convert      func(In) (ConvertResult, error)
	layerID      func(In) string
	sourceRef    func(In) string
	recordedAt   func(In) string
}

func (o reportSubmissionOp[In]) Submit(ctx context.Context, p Provider, in In) (SubmitResponse, error) {
	var zero SubmitResponse

	raw, err := json.Marshal(in)
	if err != nil {
		return zero, ingestErr(o.kind, PhaseMarshal, err)
	}

	if err := validateReportField(raw, o.reportSchema); err != nil {
		return zero, ingestErr(o.kind, PhaseValidate, err)
	}

	converted, err := o.convert(in)
	if err != nil {
		return zero, ingestErr(o.kind, PhaseConvert, err)
	}

	// Fingerprints are validated on the converted events, not the raw report:
	// report-lane payloads (triage/validation/verification) do not carry
	// fingerprints on their findings — the caller supplies them via the input's
	// FingerprintIndex and conversion stamps each event. Checking here confirms
	// every derived event got its fingerprint.
	if err := check.EventFingerprints(converted.Events); err != nil {
		return zero, ingestErr(o.kind, PhaseFingerprint, err)
	}

	batch, err := json.Marshal(BatchSubmitInput{
		SourceRef:   o.sourceRef(in),
		RecordedAt:  o.recordedAt(in),
		Events:      converted.Events,
		NeedsReview: converted.NeedsReview,
	})
	if err != nil {
		return zero, ingestErr(o.kind, PhaseMarshal, err)
	}

	resp, err := p.Submit(ctx, o.layerID(in), batch)
	return decodeSubmitResponse(resp, err, o.kind)
}

func decodeSubmitResponse(resp []byte, err error, kind IngestKind) (SubmitResponse, error) {
	var zero SubmitResponse
	if err != nil {
		return zero, ingestErr(kind, PhaseIngest, err)
	}

	var out SubmitResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return zero, ingestErr(kind, PhaseDecode, err)
	}
	return out, nil
}

func decodeResolveResponse(resp []byte, err error) (ResolveResponse, error) {
	var zero ResolveResponse
	if err != nil {
		return zero, fmt.Errorf("ingest resolve: %w", err)
	}

	var out ResolveResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return zero, fmt.Errorf("ingest resolve decode: %w", err)
	}
	return out, nil
}

func decodeFingerprintResponse(resp []byte, err error) (FingerprintResponse, error) {
	var zero FingerprintResponse
	if err != nil {
		return zero, fmt.Errorf("ingest fingerprint: %w", err)
	}

	var out FingerprintResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return zero, fmt.Errorf("ingest fingerprint decode: %w", err)
	}
	return out, nil
}

func decodeStampResponse(resp []byte, err error) (StampResponse, error) {
	var zero StampResponse
	if err != nil {
		return zero, fmt.Errorf("ingest stamp: %w", err)
	}
	var out StampResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return zero, fmt.Errorf("ingest stamp decode: %w", err)
	}
	return out, nil
}

func decodeSignResponse(resp []byte, err error) (SignResponse, error) {
	var zero SignResponse
	if err != nil {
		return zero, fmt.Errorf("ingest sign: %w", err)
	}
	var out SignResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return zero, fmt.Errorf("ingest sign decode: %w", err)
	}
	return out, nil
}

func validateReportField(raw []byte, schema string) error {
	var wrapper struct {
		Report json.RawMessage `json:"report"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return fmt.Errorf("decode report field: %w", err)
	}
	if err := validate.ValidateBytes(schema, wrapper.Report); err != nil {
		return fmt.Errorf("validate report against %q: %w", schema, err)
	}
	return nil
}

// BatchSubmit sends pre-formed events to POST /v1/ledger/layers/{layer_id}/submit.
func BatchSubmit(ctx context.Context, p Provider, layerID string, input BatchSubmitInput) (SubmitResponse, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return SubmitResponse{}, ingestErr(KindTriage, PhaseMarshal, err)
	}
	resp, err := p.Submit(ctx, layerID, raw)
	return decodeSubmitResponse(resp, err, KindTriage)
}

// ResolveReviewItem resolves a needs_review queue item.
func ResolveReviewItem(ctx context.Context, p Provider, layerID string, input ResolveInput) (ResolveResponse, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return ResolveResponse{}, fmt.Errorf("ingest resolve marshal: %w", err)
	}
	path := fmt.Sprintf("/v1/ledger/layers/%s/resolve", layerID)
	resp, err := p.Post(ctx, path, raw)
	return decodeResolveResponse(resp, err)
}

// ComputeFingerprints stamps fingerprints on findings via POST /v1/ledger/fingerprint.
func ComputeFingerprints(ctx context.Context, p Provider, input FingerprintInput) (FingerprintResponse, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return FingerprintResponse{}, fmt.Errorf("ingest fingerprint marshal: %w", err)
	}
	resp, err := p.Post(ctx, "/v1/ledger/fingerprint", raw)
	return decodeFingerprintResponse(resp, err)
}

// StampEventIdentities backfills event fingerprints on a layer via
// POST /v1/ledger/layers/{layer_id}/stamp. The ledger stamps each event whose
// finding_ref appears in input.Fingerprints (never overwriting an existing
// fingerprint) and re-signs the layer, returning the new root and count stamped.
func StampEventIdentities(ctx context.Context, p Provider, layerID string, input StampInput) (StampResponse, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return StampResponse{}, fmt.Errorf("ingest stamp marshal: %w", err)
	}
	path := fmt.Sprintf("/v1/ledger/layers/%s/stamp", layerID)
	resp, err := p.Post(ctx, path, raw)
	return decodeStampResponse(resp, err)
}

// SignLayer requests the ledger service to sign a layer's merkle tree.
func SignLayer(ctx context.Context, p Provider, layerID string, opts SignOpts) (SignResponse, error) {
	path := fmt.Sprintf("/v1/ledger/layers/%s/sign", layerID)
	if opts.Rekor {
		path += "?rekor=true"
	}
	resp, err := p.Post(ctx, path, nil)
	return decodeSignResponse(resp, err)
}
