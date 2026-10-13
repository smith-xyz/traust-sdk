package storage

import (
	"context"
	"errors"
	"testing"
)

func TestQueryProductsAndProductReposEnumerate(t *testing.T) {
	eachDialect(t, func(t *testing.T, client *Client) {
		ctx := context.Background()
		volsyncDefault := registerProductRepo(t, client, "volsync", "https://github.com/backube/volsync", "")
		volsyncRelease := registerProductRepo(t, client, "volsync", "https://github.com/backube/volsync", "release-0.14")
		acm := registerProductRepo(t, client, "acm", "https://github.com/stolostron/acm-cli", "")

		products, err := client.QueryProducts(ctx, RegistryFilter{Limit: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(products) != 2 || products[0].Slug != "acm" || products[1].Slug != "volsync" {
			t.Fatalf("products = %+v", products)
		}
		if products[0].ProductID == "" || products[0].RegisteredAt == "" || products[0].Segment != nil {
			t.Fatalf("product fields = %+v", products[0])
		}

		volsyncOnly := stringPointer("volsync")
		filtered, err := client.QueryProducts(ctx, RegistryFilter{ProductSlug: volsyncOnly, Limit: 10})
		if err != nil || len(filtered) != 1 || filtered[0].Slug != "volsync" {
			t.Fatalf("filtered products = %+v, %v", filtered, err)
		}

		page, err := client.QueryProducts(ctx, RegistryFilter{Limit: 1})
		if err != nil || len(page) != 1 || page[0].Slug != "acm" {
			t.Fatalf("paged products = %+v, %v", page, err)
		}
		next, err := client.QueryProducts(ctx, RegistryFilter{Limit: 1, Offset: 1})
		if err != nil || len(next) != 1 || next[0].Slug != "volsync" {
			t.Fatalf("next products = %+v, %v", next, err)
		}

		repos, err := client.QueryProductRepos(ctx, RegistryFilter{Limit: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(repos) != 3 {
			t.Fatalf("repos = %+v", repos)
		}
		want := []struct {
			id, slug, url, ref string
		}{
			{acm, "acm", "https://github.com/stolostron/acm-cli", ""},
			{volsyncDefault, "volsync", "https://github.com/backube/volsync", ""},
			{volsyncRelease, "volsync", "https://github.com/backube/volsync", "release-0.14"},
		}
		for i, expected := range want {
			row := repos[i]
			if row.ProductRepoID != expected.id || row.ProductSlug != expected.slug ||
				row.RepoURL != expected.url || row.Ref != expected.ref {
				t.Fatalf("repos[%d] = %+v; want %v", i, row, expected)
			}
			if row.RegisteredAt == "" {
				t.Fatalf("repos[%d] missing registered_at", i)
			}
		}

		bySlug, err := client.QueryProductRepos(ctx, RegistryFilter{ProductSlug: volsyncOnly, Limit: 10})
		if err != nil || len(bySlug) != 2 || bySlug[0].Ref != "" || bySlug[1].Ref != "release-0.14" {
			t.Fatalf("repos by slug = %+v, %v", bySlug, err)
		}

		pattern := stringPointer("%acm%")
		byPattern, err := client.QueryProductRepos(ctx, RegistryFilter{RepoURLPattern: pattern, Limit: 10})
		if err != nil || len(byPattern) != 1 || byPattern[0].ProductRepoID != acm {
			t.Fatalf("repos by pattern = %+v, %v", byPattern, err)
		}

		if _, err := client.QueryProductRepos(ctx, RegistryFilter{Limit: 0}); !errors.Is(err, ErrRegistryLimitInvalid) {
			t.Fatalf("limit 0: err = %v, want %v", err, ErrRegistryLimitInvalid)
		}
		if _, err := client.QueryProducts(ctx, RegistryFilter{Limit: 5, Offset: -1}); !errors.Is(err, ErrRegistryLimitInvalid) {
			t.Fatalf("negative offset: err = %v, want %v", err, ErrRegistryLimitInvalid)
		}
	})
}
