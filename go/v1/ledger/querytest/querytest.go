// Package querytest provides test doubles for the query package.
package querytest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/traust-security/traust-sdk/go/v1/enums"
	"github.com/traust-security/traust-sdk/go/v1/ledger"
	"github.com/traust-security/traust-sdk/go/v1/types"
)

// StaticProvider is a test double that returns pre-configured responses for
// ledger read paths. It implements ledger.Provider.
type StaticProvider struct {
	responses map[string][]byte
	err       error
	calls     []Call
}

// Call records one invocation of Query.
type Call struct {
	Method string
	Path   string
}

// NewStaticProvider creates an empty StaticProvider.
func NewStaticProvider() *StaticProvider {
	return &StaticProvider{responses: make(map[string][]byte)}
}

func (p *StaticProvider) withResponse(path string, resp interface{}) *StaticProvider {
	data, _ := json.Marshal(resp)
	p.responses[path] = data
	return p
}

// WithHealthResponse sets the response for GET /healthz.
func (p *StaticProvider) WithHealthResponse(resp ledger.HealthResponse) *StaticProvider {
	return p.withResponse("/healthz", resp)
}

// WithLayerResponse sets the response for GET /v1/ledger/layer?layer_id={layerID}.
func (p *StaticProvider) WithLayerResponse(layerID string, resp types.Layer) *StaticProvider {
	return p.withResponse(ledger.LayerPath(ledger.LayerOpDocument, layerID, nil), resp)
}

// WithFindingsResponse sets the response for GET /v1/ledger/layer/findings?layer_id={layerID}.
func (p *StaticProvider) WithFindingsResponse(layerID string, resp ledger.FindingsResponse) *StaticProvider {
	return p.withResponse(ledger.LayerPath(ledger.LayerOpFindings, layerID, nil), resp)
}

// WithBulkFindingsResponse sets the response for GET /v1/ledger/findings.
func (p *StaticProvider) WithBulkFindingsResponse(resp ledger.BulkFindingsResponse) *StaticProvider {
	return p.withResponse("/v1/ledger/findings", resp)
}

// WithVerifyResponse sets the response for GET /v1/ledger/layer/verify?layer_id={layerID}.
func (p *StaticProvider) WithVerifyResponse(layerID string, resp ledger.VerifyResponse) *StaticProvider {
	return p.withResponse(ledger.LayerPath(ledger.LayerOpVerify, layerID, nil), resp)
}

// WithEventsResponse sets the response for GET /v1/ledger/layer/events?layer_id={layerID}.
func (p *StaticProvider) WithEventsResponse(layerID string, resp ledger.EventsResponse) *StaticProvider {
	return p.withResponse(ledger.LayerPath(ledger.LayerOpEvents, layerID, nil), resp)
}

// WithLayerListResponse sets the response for GET /v1/ledger/layers.
func (p *StaticProvider) WithLayerListResponse(resp ledger.LayerListResponse) *StaticProvider {
	return p.withResponse("/v1/ledger/layers", resp)
}

// WithFindLayerResponse sets the response for GET /v1/ledger/layers?product_repo_id=.
func (p *StaticProvider) WithFindLayerResponse(productRepoID string, resp ledger.LayerListResponse) *StaticProvider {
	return p.withResponse("/v1/ledger/layers?product_repo_id="+productRepoID, resp)
}

// WithError makes all Query calls return the given error.
func (p *StaticProvider) WithError(err error) *StaticProvider {
	p.err = err
	return p
}

// Query implements ledger.Provider.
func (p *StaticProvider) Query(_ context.Context, method, path string) ([]byte, error) {
	p.calls = append(p.calls, Call{Method: method, Path: path})
	if p.err != nil {
		return nil, p.err
	}
	if resp, ok := p.responses[path]; ok {
		return resp, nil
	}
	if i := strings.Index(path, "?"); i >= 0 {
		// Layer routes are keyed by path plus layer_id; other query options
		// (limit, offset, filters) do not change the canned response.
		key := path[:i]
		if values, err := url.ParseQuery(path[i+1:]); err == nil && values.Has("layer_id") {
			key += "?" + url.Values{"layer_id": {values.Get("layer_id")}}.Encode()
		}
		if resp, ok := p.responses[key]; ok {
			return resp, nil
		}
	}
	return json.Marshal(ledger.HealthResponse{Status: "healthy"})
}

// Calls returns all recorded invocations.
func (p *StaticProvider) Calls() []Call { return p.calls }

// FixtureHealthResponse returns a typical health response.
func FixtureHealthResponse() ledger.HealthResponse {
	return ledger.HealthResponse{Status: "healthy"}
}

// FixtureLayer returns a minimal layer document.
func FixtureLayer() types.Layer {
	return types.Layer{
		Events:      []types.Event{},
		Metadata:    types.LayerMetadata{Repository: "repo-a"},
		NeedsReview: []types.ReviewItem{},
	}
}

// FixtureFindingsResponse returns a typical layer findings response.
func FixtureFindingsResponse() ledger.FindingsResponse {
	return ledger.FindingsResponse{
		Findings: []ledger.FindingDisposition{
			{
				FindingRef: "FIND-001",
				Disposition: types.ReportDisposition{
					Events:      []string{"evt-1"},
					LastUpdated: "2026-01-16T00:00:00+00:00",
					Resolution:  enums.DispositionResolutionOpen,
					Validity:    enums.ValidityConfirmed,
				},
				EventCount: 1,
			},
		},
		Summary: ledger.FindingsSummary{
			ByValidity:   types.DispositionSummaryByValidity{Confirmed: 1},
			ByResolution: types.DispositionSummaryByResolution{Open: 1},
		},
		LedgerOnly: true,
	}
}

// FixtureBulkFindingsResponse returns a typical bulk findings page.
func FixtureBulkFindingsResponse() ledger.BulkFindingsResponse {
	return ledger.BulkFindingsResponse{
		Layers: []ledger.LayerFindings{
			{
				LayerID:    "repo-a",
				Findings:   FixtureFindingsResponse().Findings,
				Summary:    FixtureFindingsResponse().Summary,
				LedgerOnly: true,
			},
		},
		TotalFindings: 1,
		HasMore:       false,
	}
}

// FixtureVerifyResponse returns a typical verify response.
func FixtureVerifyResponse() ledger.VerifyResponse {
	return ledger.VerifyResponse{
		Passed:    true,
		Findings:  []ledger.VerifyFinding{},
		CheckedAt: "2026-01-16T00:00:00+00:00",
	}
}

// FixtureEventsResponse returns a typical layer events response.
func FixtureEventsResponse(layerID string) ledger.EventsResponse {
	return ledger.EventsResponse{
		Events: []types.Event{
			{
				EventId:    "evt-1",
				FindingRef: "FIND-001",
				Rationale:  "initial report",
				RecordedAt: "2026-01-16T00:00:00+00:00",
				Disposition: types.Disposition{
					Validity:   ptr(enums.EventValidityConfirmed),
					Resolution: ptr(enums.DispositionResolutionOpen),
				},
				Source: types.EventSource{
					Ref:  "report-1",
					Type: enums.SourceTypeTriageReport,
					Actor: types.Actor{
						Kind: enums.ActorKindHuman,
					},
				},
			},
		},
		Total:   1,
		LayerID: layerID,
	}
}

func ptr[T any](v T) *T { return &v }

// Ingest implements ledger.Provider for write calls not exercised by this double.
func (p *StaticProvider) Ingest(context.Context, ledger.IngestMeta, []byte) ([]byte, error) {
	return nil, fmt.Errorf("ingest response not configured")
}
func (p *StaticProvider) Submit(context.Context, string, []byte) ([]byte, error) {
	return nil, fmt.Errorf("submit response not configured")
}
func (p *StaticProvider) Post(context.Context, string, []byte) ([]byte, error) {
	return nil, fmt.Errorf("post response not configured")
}
