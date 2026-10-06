package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/traust-security/traust-sdk/go/v1/types"
)

func TestLayerEventProjectsEveryContractColumn(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		payload := replaceJSONField(t, sampleArtifacts(t)["layer"], func(document map[string]any) {
			event := document["events"].([]any)[0].(map[string]any)
			event["rationale"] = "keeps a literal \x00 byte"
			event["harness_version"] = "0.300.0"
			event["evidence_refs"] = []any{"https://tracker.example.test/1"}
			event["source"].(map[string]any)["reported_by"] = "reporter@example.test"
			event["disposition"] = map[string]any{"validity": "confirmed", "severity": "high", "embargo": "not_required"}
			event["risk_weight"] = map[string]any{
				"lambda": 0.25, "weights_version": "v2", "tenancy_profile": "multi_tenant", "profile_source": "profile.json",
			}
			event["alias"] = map[string]any{"new_finding_ref": "FIND-101", "matched_by": "fingerprint"}
			event["finding"] = map[string]any{
				"id": "FIND-001", "title": "Synthetic", "severity": "high", "cwes": []any{"CWE-79"},
				"locations": []any{map[string]any{"path": "a.go"}}, "description": "Synthetic.", "remediation": "Fix.",
			}
		})
		artifact := mustParseArtifact(t, payload, types.ParseLayerArtifact)
		binding := runBinding("layer-test", stringPointer("layer-a"))
		if _, err := client.SaveLayer(ctx, SaveLayerInput{Binding: binding, Artifact: artifact}); err != nil {
			t.Fatal(err)
		}

		var got struct {
			rationale, harness, refs, reportedBy, severity, embargo string
			lambda                                                  float64
			weights, tenancy, profile, alias, finding               string
			restatement                                             sql.NullString
		}
		if err := sqlDB(client).QueryRow(`SELECT rationale, harness_version, evidence_refs, source_reported_by,
			severity, embargo, risk_lambda, risk_weights_version, risk_tenancy_profile, risk_profile_source,
			alias, finding, restatement
			FROM layer_event WHERE finding_ref = 'FIND-001' ORDER BY recorded_at LIMIT 1`).Scan(
			&got.rationale, &got.harness, &got.refs, &got.reportedBy, &got.severity, &got.embargo,
			&got.lambda, &got.weights, &got.tenancy, &got.profile, &got.alias, &got.finding, &got.restatement,
		); err != nil {
			t.Fatal(err)
		}
		if got.rationale != `keeps a literal \u0000 byte` || got.harness != "0.300.0" ||
			got.reportedBy != "reporter@example.test" || got.severity != "high" || got.embargo != "not_required" ||
			got.lambda != 0.25 || got.weights != "v2" || got.tenancy != "multi_tenant" ||
			got.profile != "profile.json" || got.restatement.Valid {
			t.Fatalf("layer_event = %+v", got)
		}
		for column, want := range map[string]struct {
			stored string
			value  any
		}{
			"evidence_refs": {got.refs, []any{"https://tracker.example.test/1"}},
			"alias":         {got.alias, map[string]any{"new_finding_ref": "FIND-101", "matched_by": "fingerprint"}},
		} {
			var decoded any
			if err := json.Unmarshal([]byte(want.stored), &decoded); err != nil || !reflect.DeepEqual(decoded, want.value) {
				t.Errorf("%s = %s, %v", column, want.stored, err)
			}
		}
		var finding map[string]any
		if err := json.Unmarshal([]byte(got.finding), &finding); err != nil || finding["id"] != "FIND-001" {
			t.Errorf("finding = %s, %v", got.finding, err)
		}
	})
}
