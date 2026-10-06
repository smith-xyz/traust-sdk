package storage

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"testing"

	"github.com/traust-security/traust-sdk/go/v1/types"
)

func eachDialect(t *testing.T, test func(t *testing.T, client *Client)) {
	t.Run("sqlite", func(t *testing.T) { test(t, openTestStorage(t)) })
	t.Run("postgres", func(t *testing.T) {
		client, _ := openPostgresStorage(t)
		test(t, client)
	})
}

func sameJSONList(stored string, want ...string) bool {
	var got []string
	return json.Unmarshal([]byte(stored), &got) == nil && slices.Equal(got, want)
}

func registerProductRepo(t *testing.T, client *Client, slug, repoURL, ref string) string {
	t.Helper()
	ctx := context.Background()
	productID, err := client.RegisterProduct(ctx, slug, nil)
	if err != nil {
		t.Fatal(err)
	}
	repoID, err := client.RegisterRepo(ctx, repoURL)
	if err != nil {
		t.Fatal(err)
	}
	id, err := client.RegisterProductRepo(ctx, ProductRepo{ProductID: productID, RepoID: repoID, Ref: ref})
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestRegistryIsIdempotentOnNaturalKeys(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		product, err := client.RegisterProduct(ctx, "volsync", stringPointer("operator"))
		if err != nil {
			t.Fatal(err)
		}
		again, err := client.RegisterProduct(ctx, "volsync", stringPointer("adhoc"))
		if err != nil || again != product {
			t.Fatalf("product = %s, %s, %v", product, again, err)
		}
		var segment string
		if err := sqlDB(client).QueryRow(
			"SELECT segment FROM product WHERE product_id = $1", product,
		).Scan(&segment); err != nil || segment != "adhoc" {
			t.Fatalf("segment = %q, %v", segment, err)
		}
		repo, err := client.RegisterRepo(ctx, "https://github.com/backube/volsync")
		if err != nil {
			t.Fatal(err)
		}
		defaultRef, err := client.RegisterProductRepo(ctx, ProductRepo{ProductID: product, RepoID: repo})
		if err != nil {
			t.Fatal(err)
		}
		release, err := client.RegisterProductRepo(ctx, ProductRepo{ProductID: product, RepoID: repo, Ref: "release-0.14"})
		if err != nil || release == defaultRef {
			t.Fatalf("release = %s, default = %s, %v", release, defaultRef, err)
		}
		repeat, err := client.RegisterProductRepo(ctx, ProductRepo{
			ProductID: product, RepoID: repo, SubService: stringPointer("mover"),
		})
		if err != nil || repeat != defaultRef {
			t.Fatalf("repeat = %s, want %s, %v", repeat, defaultRef, err)
		}
	})
}

func TestFindProductRepoReadsWithoutCreating(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		want := registerProductRepo(t, client, "volsync", "https://github.com/backube/volsync", "release-0.14")
		got, found, err := client.FindProductRepo(ctx, "volsync", "https://github.com/backube/volsync", "release-0.14")
		if err != nil || !found || got != want {
			t.Fatalf("find = %s, %v, %v; want %s", got, found, err, want)
		}
		for _, key := range [][3]string{
			{"volsync", "https://github.com/backube/volsync", ""},
			{"volsync", "https://github.com/Backube/volsync", "release-0.14"},
			{"unknown", "https://github.com/backube/volsync", "release-0.14"},
		} {
			if id, found, err := client.FindProductRepo(ctx, key[0], key[1], key[2]); err != nil || found || id != "" {
				t.Fatalf("find %v = %s, %v, %v; want not found", key, id, found, err)
			}
		}
		var products int
		if err := sqlDB(client).QueryRow("SELECT count(*) FROM product").Scan(&products); err != nil || products != 1 {
			t.Fatalf("products = %d, %v", products, err)
		}
	})
}

func TestRegistryRejectsUnregisteredParentsAndBadText(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		if _, err := client.RegisterProductRepo(ctx, ProductRepo{ProductID: "missing", RepoID: "missing"}); err == nil {
			t.Fatal("product_repo with unregistered parents was accepted")
		}
		if err := client.RegisterRepoOwner(ctx, RepoOwner{ProductRepoID: "missing", Team: "team"}); err == nil {
			t.Fatal("owner of an unregistered product_repo was accepted")
		}
		tests := []struct {
			name string
			call func() error
			want error
		}{
			{"empty slug", func() error { _, err := client.RegisterProduct(ctx, "", nil); return err }, ErrRegistryValueMissing},
			{"empty repo", func() error { _, err := client.RegisterRepo(ctx, ""); return err }, ErrRegistryValueMissing},
			{"NUL slug", func() error { _, err := client.RegisterProduct(ctx, "a\x00b", nil); return err }, ErrInvalidIdentifier},
			{"NUL segment", func() error {
				_, err := client.RegisterProduct(ctx, "a", stringPointer("x\x00"))
				return err
			}, ErrInvalidIdentifier},
			{"empty version", func() error {
				return client.RegisterProductRepoVersion(ctx, ProductRepoVersion{ProductRepoID: "x"})
			}, ErrRegistryValueMissing},
			{"empty image", func() error {
				return client.RegisterProductRepoVersion(ctx, ProductRepoVersion{ProductRepoID: "x", Version: "1", Images: []string{""}})
			}, ErrRegistryValueMissing},
		}
		for _, tt := range tests {
			if err := tt.call(); !errors.Is(err, tt.want) {
				t.Errorf("%s: err = %v, want %v", tt.name, err, tt.want)
			}
		}
	})
}

func TestVersionAndOwnerUpsertEncodeLists(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		owner := registerProductRepo(t, client, "acm", "https://github.com/stolostron/acm-cli", "")
		for _, images := range [][]string{{"quay.io/a:1", "quay.io/b<&>:2"}, {"quay.io/c:3"}} {
			if err := client.RegisterProductRepoVersion(ctx, ProductRepoVersion{
				ProductRepoID: owner, Version: "2.14", Images: images,
			}); err != nil {
				t.Fatal(err)
			}
		}
		var images string
		var operators *string
		if err := sqlDB(client).QueryRow(
			"SELECT images, cluster_operators FROM product_repo_version WHERE product_repo_id = $1", owner,
		).Scan(&images, &operators); err != nil || !sameJSONList(images, "quay.io/c:3") || operators != nil {
			t.Fatalf("version = %s, %v, %v", images, operators, err)
		}
		if err := client.RegisterRepoOwner(ctx, RepoOwner{
			ProductRepoID: owner, Team: "acm-cli", Individuals: []string{"a", "b"},
		}); err != nil {
			t.Fatal(err)
		}
		var individuals string
		if err := sqlDB(client).QueryRow(
			"SELECT individuals FROM repo_owner WHERE product_repo_id = $1", owner,
		).Scan(&individuals); err != nil || !sameJSONList(individuals, "a", "b") {
			t.Fatalf("owner = %s, %v", individuals, err)
		}
	})
}

func TestBindingCarriesProductRepoOutsideIdentity(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		owner := registerProductRepo(t, client, "volsync", "https://github.com/backube/volsync", "")
		other := registerProductRepo(t, client, "volsync", "https://github.com/backube/volsync", "release-0.14")
		artifact, err := types.ParseVulnFindingsArtifact(readFixture(t, "storagetest/testdata/vuln-findings-populated.test.json"))
		if err != nil {
			t.Fatal(err)
		}
		binding := runBinding("local", nil)
		anchored := binding
		anchored.ProductRepoID = stringPointer(owner)
		anchored.CommitSHA = stringPointer("0123456789abcdef0123456789abcdef01234567")

		saved, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: anchored, Artifact: artifact})
		if err != nil {
			t.Fatal(err)
		}
		if saved.BindingID != identifyBinding(saved.Digest, "vuln-findings", binding) {
			t.Fatal("product_repo_id or commit_sha changed the binding identity")
		}
		record, err := client.GetBinding(ctx, saved.BindingID)
		if err != nil || !sameOptional(record.Binding.ProductRepoID, anchored.ProductRepoID) ||
			!sameOptional(record.Binding.CommitSHA, anchored.CommitSHA) {
			t.Fatalf("record = %+v, %v", record.Binding, err)
		}
		retry, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: anchored, Artifact: artifact})
		if err != nil || !retry.AlreadyBound {
			t.Fatalf("retry = %+v, %v", retry, err)
		}
		moved := anchored
		moved.ProductRepoID = stringPointer(other)
		if _, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: moved, Artifact: artifact}); !errors.Is(err, ErrBindingMismatch) {
			t.Fatalf("re-save under another product_repo = %v, want mismatch", err)
		}

		corrected, err := types.ParseVulnFindingsArtifact(replaceJSONField(t, artifact.Payload(), func(document map[string]any) {
			document["findings"].([]any)[0].(map[string]any)["title"] = "Corrected title"
		}))
		if err != nil {
			t.Fatal(err)
		}
		successor := anchored
		successor.SupersedesBindingID = stringPointer(saved.BindingID)
		successor.CommitSHA = stringPointer("fedcba9876543210fedcba9876543210fedcba98")
		crossOwner := successor
		crossOwner.ProductRepoID = stringPointer(other)
		if _, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: crossOwner, Artifact: corrected}); !errors.Is(err, ErrBindingMismatch) {
			t.Fatalf("supersession across product_repos = %v, want mismatch", err)
		}
		if _, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: successor, Artifact: corrected}); err != nil {
			t.Fatalf("supersession with a new commit = %v", err)
		}

		unregistered := binding
		unregistered.RunID = stringPointer("sci:scan-result:8")
		unregistered.ProductRepoID = stringPointer("00000000-0000-4000-8000-000000000000")
		if _, err := client.SaveVulnFindings(ctx, SaveVulnFindingsInput{Binding: unregistered, Artifact: artifact}); err == nil {
			t.Fatal("binding to an unregistered product_repo was accepted")
		}
	})
}
