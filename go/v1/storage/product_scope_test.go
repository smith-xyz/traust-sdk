package storage

import (
	"context"
	"testing"
)

func TestGetProductScopes(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		anchored := registerProductRepo(t, client, "acm", "https://github.com/stolostron/console", "release-2.10")
		loose := registerProductRepo(t, client, "insights", "https://github.com/RedHatInsights/lightspeed-stack", "")

		insert := func(bindingID string, productRepoID *string) {
			t.Helper()
			digest := "digest-" + bindingID
			if _, err := sqlDB(client).Exec(
				"INSERT INTO artifact_evidence (digest, byte_size, first_ingested_at) VALUES ($1, 1, '2026-10-13T00:00:00Z')",
				digest,
			); err != nil {
				t.Fatal(err)
			}
			if _, err := sqlDB(client).Exec(
				"INSERT INTO artifact_binding (binding_id, artifact_digest, artifact_name, artifact_role, scope_id, subject_id, run_id, layer_id, bound_at, product_repo_id) VALUES ($1, $2, 'report', 'baseline', 'scope', 'subject', 'run', NULL, '2026-10-13T00:00:00Z', $3)",
				bindingID, digest, productRepoID,
			); err != nil {
				t.Fatal(err)
			}
		}
		insert("binding-anchored", &anchored)
		insert("binding-loose", &loose)
		insert("binding-unanchored", nil)

		scopes, err := client.GetProductScopes(ctx, []string{"binding-anchored", "binding-loose", "binding-unanchored", "binding-missing"})
		if err != nil {
			t.Fatal(err)
		}
		if len(scopes) != 2 {
			t.Fatalf("expected two anchored scopes, got %d: %+v", len(scopes), scopes)
		}
		anchoredScope, ok := scopes["binding-anchored"]
		if !ok || anchoredScope.ProductSlug != "acm" || anchoredScope.Ref != "release-2.10" || anchoredScope.ProductRepoID != anchored {
			t.Fatalf("anchored scope wrong: %+v", anchoredScope)
		}
		looseScope, ok := scopes["binding-loose"]
		if !ok || looseScope.ProductSlug != "insights" || looseScope.Ref != "" || looseScope.ProductRepoID != loose {
			t.Fatalf("loose scope wrong: %+v", looseScope)
		}
		if scopes["binding-unanchored"] != (ProductScope{}) || scopes["binding-missing"] != (ProductScope{}) {
			t.Fatalf("unanchored and missing bindings must be absent: %+v", scopes)
		}
	})
}

func TestGetProductScopesRejectsInvalidInput(t *testing.T) {
	client := openTestStorage(t)
	if _, err := client.GetProductScopes(context.Background(), nil); err == nil {
		t.Fatal("empty binding list must be rejected")
	}
	if _, err := client.GetProductScopes(context.Background(), []string{"ok", ""}); err == nil {
		t.Fatal("empty binding id must be rejected")
	}
}
