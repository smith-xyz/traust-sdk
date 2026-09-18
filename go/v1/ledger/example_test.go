package ledger_test

import (
	"context"
	"fmt"

	"github.com/traust-security/traust-sdk/go/v1/ledger"
	"github.com/traust-security/traust-sdk/go/v1/ledger/ingesttest"
)

func ExampleNewClient() {
	provider := ingesttest.NewStaticProvider().
		WithSubmitResponse(ledger.SubmitResponse{ID: "batch-001", Status: "accepted"})
	client := ledger.NewClient(provider)

	resp, err := client.BatchSubmit(context.Background(), "repo-a", ledger.BatchSubmitInput{
		SourceRef:  "findings/repo-a/triage.json",
		RecordedAt: "2026-01-16T00:00:00+00:00",
		Events: []map[string]interface{}{
			{"finding_ref": "FIND-001"},
		},
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("id: %s, status: %s\n", resp.ID, resp.Status)
	// Output: id: batch-001, status: accepted
}

func ExampleNewHTTPClient() {
	client := ledger.NewHTTPClient("https://ledger.example.com",
		ledger.WithBearerToken("my-token"),
	)

	_, _ = client.SubmitCountersign(context.Background(), ledger.CountersignInput{
		EventMeta: ledger.EventMeta{
			LayerID:    "repo-a",
			RecordedAt: "2026-01-16T00:00:00+00:00",
		},
		FindingRef:    "FIND-001",
		Verdict:       "true_positive",
		Justification: "Reviewed source and confirmed exploit path.",
	})
}
