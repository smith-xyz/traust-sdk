package storage

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// ProductRepo is a repo as a product ships it. Ref "" is the default branch.
type ProductRepo struct {
	ProductID    string
	RepoID       string
	Ref          string
	SubService   *string
	ResourceType *string
}

// ProductRepoVersion records that a product version ships a product_repo.
type ProductRepoVersion struct {
	ProductRepoID    string
	Version          string
	Category         *string
	ClusterOperators []string
	Images           []string
}

// RepoOwner records a team that owns a product_repo.
type RepoOwner struct {
	ProductRepoID string
	Team          string
	Manager       *string
	Individuals   []string
	Source        *string
	JiraProject   *string
	JiraComponent *string
}

// RegisterProduct registers a product by slug, idempotently, and returns its
// product_id. Segment is descriptive: re-registering sets it.
func (c *Client) RegisterProduct(ctx context.Context, slug string, segment *string) (string, error) {
	if err := registryText(slug, true, segment); err != nil {
		return "", err
	}
	id := newRegistryID()
	return c.register(ctx, func(conn *sql.Conn, s *sqlStore, at string) (string, error) {
		if err := s.queries.productUpsert(ctx, conn, productUpsertParams{
			productId: id, slug: slug, segment: segment, registeredAt: at,
		}); err != nil {
			return "", err
		}
		return storedID(s.queries.productGet(ctx, conn, productGetParams{slug: slug}))
	})
}

// RegisterRepo registers a repo by exact URL, idempotently, and returns its repo_id.
func (c *Client) RegisterRepo(ctx context.Context, repoURL string) (string, error) {
	if err := registryText(repoURL, true); err != nil {
		return "", err
	}
	id := newRegistryID()
	return c.register(ctx, func(conn *sql.Conn, s *sqlStore, at string) (string, error) {
		if err := s.queries.repoUpsert(ctx, conn, repoUpsertParams{
			repoId: id, repoUrl: repoURL, registeredAt: at,
		}); err != nil {
			return "", err
		}
		return storedID(s.queries.repoGet(ctx, conn, repoGetParams{repoUrl: repoURL}))
	})
}

// RegisterProductRepo registers a product's repo at a ref, idempotently, and
// returns its product_repo_id. Product and repo must already be registered.
func (c *Client) RegisterProductRepo(ctx context.Context, pr ProductRepo) (string, error) {
	if err := registryText(pr.ProductID, true, pr.SubService, pr.ResourceType); err != nil {
		return "", err
	}
	if err := registryText(pr.RepoID, true); err != nil {
		return "", err
	}
	if err := registryText(pr.Ref, false); err != nil {
		return "", err
	}
	id := newRegistryID()
	return c.register(ctx, func(conn *sql.Conn, s *sqlStore, at string) (string, error) {
		if err := s.queries.productRepoUpsert(ctx, conn, productRepoUpsertParams{
			productRepoId: id,
			productId:     pr.ProductID,
			repoId:        pr.RepoID,
			ref:           pr.Ref,
			subService:    pr.SubService,
			resourceType:  pr.ResourceType,
			registeredAt:  at,
		}); err != nil {
			return "", err
		}
		return storedID(s.queries.productRepoGet(ctx, conn, productRepoGetParams{
			productId: pr.ProductID, repoId: pr.RepoID, ref: pr.Ref,
		}))
	})
}

// RegisterProductRepoVersion records a version of a product_repo; re-registering refreshes it.
func (c *Client) RegisterProductRepoVersion(ctx context.Context, v ProductRepoVersion) error {
	if err := registryText(v.ProductRepoID, true, v.Category); err != nil {
		return err
	}
	if err := registryText(v.Version, true); err != nil {
		return err
	}
	operators, err := registryList(v.ClusterOperators)
	if err != nil {
		return err
	}
	images, err := registryList(v.Images)
	if err != nil {
		return err
	}
	_, err = c.register(ctx, func(conn *sql.Conn, s *sqlStore, at string) (string, error) {
		return "", s.queries.productRepoVersionUpsert(ctx, conn, productRepoVersionUpsertParams{
			productRepoId:    v.ProductRepoID,
			version:          v.Version,
			category:         v.Category,
			clusterOperators: operators,
			images:           images,
			registeredAt:     at,
		})
	})
	return err
}

// RegisterRepoOwner records a team that owns a product_repo; re-registering refreshes it.
func (c *Client) RegisterRepoOwner(ctx context.Context, o RepoOwner) error {
	if err := registryText(o.ProductRepoID, true, o.Manager, o.Source, o.JiraProject, o.JiraComponent); err != nil {
		return err
	}
	if err := registryText(o.Team, true); err != nil {
		return err
	}
	individuals, err := registryList(o.Individuals)
	if err != nil {
		return err
	}
	_, err = c.register(ctx, func(conn *sql.Conn, s *sqlStore, at string) (string, error) {
		return "", s.queries.repoOwnerUpsert(ctx, conn, repoOwnerUpsertParams{
			productRepoId: o.ProductRepoID,
			team:          o.Team,
			manager:       o.Manager,
			individuals:   individuals,
			source:        o.Source,
			jiraProject:   o.JiraProject,
			jiraComponent: o.JiraComponent,
			registeredAt:  at,
		})
	})
	return err
}

// FindProductRepo returns the product_repo_id for (product slug, repo URL, ref)
// without creating anything; found is false when the registry has no such row.
func (c *Client) FindProductRepo(ctx context.Context, slug, repoURL, ref string) (id string, found bool, err error) {
	if c == nil || c.store == nil {
		return "", false, wrap(OperationRead, PhaseInput, ErrNilDatabase)
	}
	conn, err := c.store.db.Conn(ctx)
	if err != nil {
		return "", false, wrap(OperationRead, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()
	err = c.store.queries.productRepoFind(ctx, conn, productRepoFindParams{
		slug: slug, repoUrl: repoURL, ref: ref,
	}).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, wrap(OperationRead, PhaseRead, err)
	}
	return id, true, nil
}

type registryWrite func(conn *sql.Conn, store *sqlStore, registeredAt string) (string, error)

func (c *Client) register(ctx context.Context, write registryWrite) (string, error) {
	if c == nil || c.store == nil {
		return "", wrap(OperationRegister, PhaseInput, ErrNilDatabase)
	}
	conn, err := c.store.db.Conn(ctx)
	if err != nil {
		return "", wrap(OperationRegister, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()
	if err := c.store.prepareConnection(ctx, conn); err != nil {
		return "", wrap(OperationRegister, PhaseConnect, err)
	}
	if err := begin(ctx, conn, c.store.dialect); err != nil {
		return "", wrap(OperationRegister, PhaseBegin, err)
	}
	committed := false
	defer rollbackUnlessCommitted(ctx, conn, &committed)

	id, err := write(conn, c.store, nowUTC())
	if err != nil {
		return "", wrap(OperationRegister, PhaseRegistry, err)
	}
	if _, err := conn.ExecContext(ctx, "COMMIT"); err != nil {
		return "", wrap(OperationRegister, PhaseCommit, err)
	}
	committed = true
	return id, nil
}

func storedID(row *sql.Row) (string, error) {
	var id string
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrRegistryRowMissing
		}
		return "", err
	}
	return id, nil
}

// registryText validates a key value (required: non-empty) and optional
// attributes: NUL-free UTF-8, matching contracts' text rules.
func registryText(value string, required bool, optional ...*string) error {
	if required && value == "" {
		return wrap(OperationRegister, PhaseInput, ErrRegistryValueMissing)
	}
	if !validBindingText(value) {
		return wrap(OperationRegister, PhaseInput, ErrInvalidIdentifier)
	}
	for _, attribute := range optional {
		if attribute != nil && !validBindingText(*attribute) {
			return wrap(OperationRegister, PhaseInput, ErrInvalidIdentifier)
		}
	}
	return nil
}

// registryList encodes a list attribute as a compact JSON array; empty stays NULL.
func registryList(values []string) (*string, error) {
	if len(values) == 0 {
		return nil, nil
	}
	for _, value := range values {
		if err := registryText(value, true); err != nil {
			return nil, err
		}
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(values); err != nil {
		return nil, wrap(OperationRegister, PhaseInput, err)
	}
	text := string(bytes.TrimRight(buf.Bytes(), "\n"))
	return &text, nil
}

func newRegistryID() string {
	return uuid.NewString()
}
