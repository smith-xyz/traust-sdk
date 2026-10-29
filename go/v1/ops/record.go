package ops

import (
	"context"
	"errors"

	"github.com/traust-security/traust-sdk/go/v1/api"
)

func (run RecordingRun) record(ctx context.Context, original api.RunOrder, current *api.RunOrder, raw []byte) (RecordingReceipt, error) {
	var receipt RecordingReceipt
	if err := usableOrder(current); err != nil {
		return receipt, &LeaseLostError{Err: err}
	}
	if err := ctx.Err(); err != nil {
		return receipt, &LeaseLostError{Err: err}
	}
	requestCtx, cancel := run.ownedContext(ctx, current)
	defer cancel()
	// Conversion occurs only after utf8.Valid in Run; encoding/json therefore
	// carries the exact original bytes instead of replacing malformed sequences.
	payload := string(raw)
	switch {
	case run.Request.Assessment != nil:
		body := *run.Request.Assessment
		provenance, err := bindProvenance(original, body.RunOrderId, body.Provenance)
		if err != nil {
			return receipt, err
		}
		body.RunOrderId, body.Provenance, body.Payload = *original.Id, provenance, payload
		response, err := run.Client.RecordAssessmentWithResponse(requestCtx, body)
		if err != nil {
			return receipt, err
		}
		receipt.Assessment, err = nativeResponse(response.StatusCode(), response.Body, response.JSON200)
		return receipt, err
	case run.Request.Triage != nil:
		body := *run.Request.Triage
		provenance, err := bindProvenance(original, body.RunOrderId, body.Provenance)
		if err != nil {
			return receipt, err
		}
		body.RunOrderId, body.Provenance, body.Payload = *original.Id, provenance, payload
		response, err := run.Client.RecordTriageWithResponse(requestCtx, run.Request.AssessmentID, body)
		if err != nil {
			return receipt, err
		}
		receipt.Triage, err = nativeResponse(response.StatusCode(), response.Body, response.JSON200)
		return receipt, err
	case run.Request.Validation != nil:
		body := *run.Request.Validation
		provenance, err := bindProvenance(original, body.RunOrderId, body.Provenance)
		if err != nil {
			return receipt, err
		}
		body.RunOrderId, body.Provenance, body.Payload = *original.Id, provenance, payload
		response, err := run.Client.RecordValidationWithResponse(requestCtx, body)
		if err != nil {
			return receipt, err
		}
		receipt.Validation, err = nativeResponse(response.StatusCode(), response.Body, response.JSON200)
		return receipt, err
	case run.Request.Verification != nil:
		body := *run.Request.Verification
		provenance, err := bindProvenance(original, body.RunOrderId, body.Provenance)
		if err != nil {
			return receipt, err
		}
		body.RunOrderId, body.Provenance, body.Payload = *original.Id, provenance, payload
		response, err := run.Client.RecordVerificationWithResponse(requestCtx, body)
		if err != nil {
			return receipt, err
		}
		receipt.Verification, err = nativeResponse(response.StatusCode(), response.Body, response.JSON200)
		return receipt, err
	}
	return receipt, errors.New("recording run has no native request template")
}

func (request RecordingRequest) templateInput() (string, api.Provenance) {
	switch {
	case request.Assessment != nil:
		return request.Assessment.Payload, request.Assessment.Provenance
	case request.Triage != nil:
		return request.Triage.Payload, request.Triage.Provenance
	case request.Validation != nil:
		return request.Validation.Payload, request.Validation.Provenance
	case request.Verification != nil:
		return request.Verification.Payload, request.Verification.Provenance
	}
	return "", api.Provenance{}
}
