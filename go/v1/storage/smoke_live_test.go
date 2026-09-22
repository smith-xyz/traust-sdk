package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"
)

// TestLiveCorpusReadsThroughEveryNewQuery is the smoke test for step 9: the
// three Query* methods added here, run against a store built from real
// artifacts, not fixtures. Skipped unless TRAUST_SMOKE_DB names one.
func TestLiveCorpusReadsThroughEveryNewQuery(t *testing.T) {
	path := os.Getenv("TRAUST_SMOKE_DB")
	if path == "" {
		t.Skip("set TRAUST_SMOKE_DB to a real store to run the live smoke")
	}
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	client, err := NewClient(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	scopes := []string{"local"}

	current, err := client.QueryValidationCurrent(ctx, scopes)
	if err != nil {
		t.Fatalf("QueryValidationCurrent: %v", err)
	}
	t.Logf("validation_current rows: %d", len(current))
	verdicts := map[string]int{}
	for _, row := range current {
		if row.Verdict != nil {
			verdicts[*row.Verdict]++
		}
	}
	t.Logf("verdicts: %v", verdicts)

	exposure, err := client.QueryValidationExposure(ctx, scopes)
	if err != nil {
		t.Fatalf("QueryValidationExposure: %v", err)
	}
	t.Logf("validation_exposure rows: %d", len(exposure))

	advisory, err := client.QueryAdvisoryExposure(ctx, scopes)
	if err != nil {
		t.Fatalf("QueryAdvisoryExposure: %v", err)
	}
	t.Logf("advisory_exposure rows: %d", len(advisory))
	if len(advisory) > 0 {
		row := advisory[0]
		t.Logf("first advisory row: advisory=%v repo=%v classification=%v evidence_level=%v",
			show(row.Advisory), show(row.Repo), show(row.Classification), show(row.EvidenceLevel))
	}
}

func show(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
