package ledger_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/traust-security/traust-sdk/go/v1/enums"
	"github.com/traust-security/traust-sdk/go/v1/ledger"
	"github.com/traust-security/traust-sdk/go/v1/ledger/querytest"
	"github.com/traust-security/traust-sdk/go/v1/types"
)

const productRepoID = "3f2b6c1e-8d4a-4e2b-9c1f-0a1b2c3d4e5f"

func stringPtr(value string) *string { return &value }

func TestInitializeLayerSendsProductRepoAndLayer(t *testing.T) {
	tests := []struct {
		name  string
		owner *string
		want  string
	}{
		{"database ledger", stringPtr(productRepoID), `"` + productRepoID + `"`},
		{"file ledger", nil, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod string
			var body map[string]json.RawMessage
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotMethod = r.URL.Path, r.Method
				raw, _ := io.ReadAll(r.Body)
				mustNoErr(t, json.Unmarshal(raw, &body))
				_, _ = w.Write([]byte(`{"layer_id":"repo-a"}`))
			}))
			defer srv.Close()

			resp, err := ledger.NewHTTPClient(srv.URL).InitializeLayer(context.Background(), "repo-a", ledger.InitializeInput{
				ProductRepoID: tt.owner,
				Layer:         querytest.FixtureLayer(),
			})
			mustNoErr(t, err)
			if gotMethod != http.MethodPost || gotPath != "/v1/ledger/layers/repo-a/initialize" {
				t.Fatalf("request = %s %s", gotMethod, gotPath)
			}
			if string(body["product_repo_id"]) != tt.want {
				t.Fatalf("product_repo_id = %s, want %s", body["product_repo_id"], tt.want)
			}
			var layer map[string]any
			mustNoErr(t, json.Unmarshal(body["layer"], &layer))
			if layer["metadata"] == nil || layer["events"] == nil {
				t.Fatalf("layer = %s", body["layer"])
			}
			if resp.LayerID != "repo-a" {
				t.Fatalf("response = %+v", resp)
			}
		})
	}
}

func TestFindLayerByProductRepo(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		wantFound bool
		wantErr   bool
	}{
		{"found", http.StatusOK, `{"layers":[{"layer_id":"repo-a","product_repo_id":"` + productRepoID + `"}]}`, true, false},
		{"no layer", http.StatusNotFound, `{"detail":"no layer for product_repo"}`, false, false},
		{"invalid id", http.StatusUnprocessableEntity, `{"detail":"invalid"}`, false, true},
		{"ambiguous", http.StatusOK, `{"layers":[{"layer_id":"a"},{"layer_id":"b"}]}`, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotQuery string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotQuery = r.URL.Query().Get("product_repo_id")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			ref, found, err := ledger.NewHTTPClient(srv.URL).FindLayer(context.Background(), productRepoID)
			if gotQuery != productRepoID {
				t.Fatalf("query product_repo_id = %q", gotQuery)
			}
			if (err != nil) != tt.wantErr || found != tt.wantFound {
				t.Fatalf("find = %+v, %v, %v", ref, found, err)
			}
			if tt.wantFound && (ref.LayerID != "repo-a" || ref.ProductRepoID == nil || *ref.ProductRepoID != productRepoID) {
				t.Fatalf("ref = %+v", ref)
			}
			var status *ledger.StatusError
			if tt.name == "invalid id" && !errors.As(err, &status) {
				t.Fatalf("err = %v, want StatusError", err)
			}
		})
	}
}

func TestListLayersDecodesProductRepoPairs(t *testing.T) {
	provider := querytest.NewStaticProvider().WithLayerListResponse(ledger.LayerListResponse{
		Layers: []ledger.LayerRef{
			{LayerID: "db-layer", ProductRepoID: stringPtr(productRepoID)},
			{LayerID: "file-layer"},
		},
	})
	resp, err := ledger.NewClient(provider).ListLayers(context.Background())
	mustNoErr(t, err)
	if len(resp.Layers) != 2 || *resp.Layers[0].ProductRepoID != productRepoID || resp.Layers[1].ProductRepoID != nil {
		t.Fatalf("layers = %+v", resp.Layers)
	}
	if calls := provider.Calls(); len(calls) != 1 || calls[0].Path != "/v1/ledger/layers" {
		t.Fatalf("calls = %+v", calls)
	}
}

func TestListLayersRendersPluggableFilters(t *testing.T) {
	provider := querytest.NewStaticProvider().WithLayerListResponse(ledger.LayerListResponse{})
	_, err := ledger.NewClient(provider).ListLayers(
		context.Background(),
		ledger.ReviewStatusIn("pending", "rejected"),
		ledger.LayerIDs("layer-a", "layer-b"),
		ledger.UpdatedSince(time.Date(2026, 10, 13, 16, 0, 0, 0, time.UTC)),
		ledger.WithReviews(),
	)
	mustNoErr(t, err)
	calls := provider.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls = %+v", calls)
	}
	query := calls[0].Path
	for _, want := range []string{
		"filter=review_status%3Ain%3Apending%2Crejected",
		"filter=layer_id%3Ain%3Alayer-a%2Clayer-b",
		"filter=updated_at%3Agte%3A2026-10-13T16%3A00%3A00Z",
		"include_reviews=true",
	} {
		if !strings.Contains(query, want) {
			t.Fatalf("query %q missing %q", query, want)
		}
	}
}

func TestListLayersRawFilterPassesThrough(t *testing.T) {
	provider := querytest.NewStaticProvider().WithLayerListResponse(ledger.LayerListResponse{})
	_, err := ledger.NewClient(provider).ListLayers(
		context.Background(),
		ledger.RawFilter("merkle_epoch", "gte", "42"),
	)
	mustNoErr(t, err)
	if calls := provider.Calls(); len(calls) != 1 || !strings.Contains(calls[0].Path, "filter=merkle_epoch%3Agte%3A42") {
		t.Fatalf("calls = %+v", calls)
	}
}

func TestListLayersRejectsInvalidFilters(t *testing.T) {
	provider := querytest.NewStaticProvider().WithLayerListResponse(ledger.LayerListResponse{})
	if _, err := ledger.NewClient(provider).ListLayers(context.Background(), ledger.ReviewStatus("")); err == nil {
		t.Fatal("empty eq value should error")
	}
	if _, err := ledger.NewClient(provider).ListLayers(context.Background(), ledger.LayerIDs()); err == nil {
		t.Fatal("empty in-list should error")
	}
	if _, err := ledger.NewClient(provider).ListLayers(context.Background(), ledger.LayerIDs("a", "b,c")); err == nil {
		t.Fatal("separator in in-list member should error")
	}
	if _, err := ledger.NewClient(provider).ListLayers(context.Background(), ledger.UpdatedSince(time.Time{})); err == nil {
		t.Fatal("zero timestamp should error")
	}
	if calls := provider.Calls(); len(calls) != 0 {
		t.Fatalf("invalid filters must not reach the wire: %+v", calls)
	}
}

func TestListLayersDecodesEmbeddedReviewItems(t *testing.T) {
	status := enums.ReviewItemStatusPending
	findingRef := "FIND-001"
	owner := productRepoID
	provider := querytest.NewStaticProvider().WithLayerListResponse(ledger.LayerListResponse{
		Layers: []ledger.LayerRef{
			{
				LayerID:       "db-layer",
				ProductRepoID: &owner,
				NeedsReview: []types.ReviewItem{
					{
						Author:              "reporter@example.test",
						QueuedAt:            "2026-10-13T16:00:00Z",
						Quote:               "suspect",
						SourceRef:           "scan-001",
						Status:              status,
						SuggestedFindingRef: &findingRef,
					},
				},
			},
		},
	})
	resp, err := ledger.NewClient(provider).ListLayers(
		context.Background(),
		ledger.ReviewStatus("pending"),
	)
	mustNoErr(t, err)
	if len(resp.Layers) != 1 || len(resp.Layers[0].NeedsReview) != 1 {
		t.Fatalf("layers = %+v", resp.Layers)
	}
	item := resp.Layers[0].NeedsReview[0]
	if item.Status != status || item.SuggestedFindingRef == nil || *item.SuggestedFindingRef != findingRef {
		t.Fatalf("item = %+v", item)
	}
}
