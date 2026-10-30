// Package ops runs traust work for a consumer: it holds an order's lease, calls the consumer's
// Executor, and records and completes the order through core. Core selects and authorizes work.
package ops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/traust-security/traust-sdk/go/v1/api"
)

const (
	defaultRequestTimeout = time.Minute
	maxLeaseSeconds       = 24 * 60 * 60
)

// Executor performs the work. On cancellation it stops the work, including child processes,
// waits for it to stop, and returns.
type Executor interface {
	Execute(context.Context, Execution) (ProducedReport, error)
}

// Execution is what an Executor runs: server context, never credentials.
type Execution struct {
	Order       api.RunOrder     `json:"order"`
	TriageInput *api.TriageInput `json:"triage_input,omitempty"`
	Input       json.RawMessage  `json:"input,omitempty"`
	Argv        []string         `json:"argv,omitempty"`
}

// ProducedReport is the executor's outcome and its original output, byte for byte. RecordingRun
// sends Payload to its record door; Runner sends only Result.
type ProducedReport struct {
	Payload []byte
	Result  api.JobResult
}

// LeaseLostError means the order is no longer this process's: nothing may be recorded or completed.
type LeaseLostError struct{ Err error }

func (e *LeaseLostError) Error() string { return "ops: order lease lost: " + e.Err.Error() }
func (e *LeaseLostError) Unwrap() error { return e.Err }

func lost(err error) bool {
	var lease *LeaseLostError
	return errors.As(err, &lease)
}

func checkOutcome(result api.JobResult) error {
	switch result.Outcome {
	case api.OutcomeSucceeded, api.OutcomeDegraded, api.OutcomeRefused, api.OutcomeFailed:
		return nil
	}
	return fmt.Errorf("ops: executor returned outcome %q", result.Outcome)
}

func leaseDuration(seconds int) time.Duration { return time.Duration(seconds) * time.Second }
