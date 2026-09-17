package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/openshift/traust-sdk/go/v1/types"
)

type FindingsSummaryRow struct {
	ScopeID      string
	SubjectID    *string
	RunID        *string
	LayerID      *string
	Repository   *string
	Severity     string
	Verdict      *string
	FindingCount int64
}

func (c *Client) GetEvidence(ctx context.Context, digest string) ([]byte, error) {
	if c == nil || c.store == nil {
		return nil, wrap(OperationRead, PhaseInput, ErrNilDatabase)
	}
	if !digestPattern.MatchString(digest) {
		return nil, wrap(OperationRead, PhaseInput, ErrNotFound)
	}
	conn, err := c.store.db.Conn(ctx)
	if err != nil {
		return nil, wrap(OperationRead, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()
	return c.store.readEvidence(ctx, conn, digest)
}

func (c *Client) GetBinding(ctx context.Context, bindingID string) (BindingRecord, error) {
	if c == nil || c.store == nil {
		return BindingRecord{}, wrap(OperationRead, PhaseInput, ErrNilDatabase)
	}
	if !digestPattern.MatchString(bindingID) {
		return BindingRecord{}, wrap(OperationRead, PhaseInput, ErrBindingNotFound)
	}
	conn, err := c.store.db.Conn(ctx)
	if err != nil {
		return BindingRecord{}, wrap(OperationRead, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()
	record, found, err := c.store.getBinding(ctx, conn, bindingID)
	if err != nil {
		return BindingRecord{}, err
	}
	if !found {
		return BindingRecord{}, wrap(OperationRead, PhaseRead, ErrBindingNotFound)
	}
	return record, nil
}

func getTypedArtifact[T any](
	ctx context.Context,
	store *sqlStore,
	name string,
	bindingID string,
	parse func([]byte) (types.Artifact[T], error),
) (types.Artifact[T], error) {
	var zero types.Artifact[T]
	if store == nil {
		return zero, wrap(OperationRead, PhaseInput, ErrNilDatabase)
	}
	if !digestPattern.MatchString(bindingID) {
		return zero, wrap(OperationRead, PhaseInput, ErrBindingNotFound)
	}
	conn, err := store.db.Conn(ctx)
	if err != nil {
		return zero, wrap(OperationRead, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()

	record, found, err := store.getBinding(ctx, conn, bindingID)
	if err != nil {
		return zero, err
	}
	if !found {
		return zero, wrap(OperationRead, PhaseRead, ErrBindingNotFound)
	}
	if record.ArtifactName != name {
		return zero, wrap(OperationRead, PhaseRead, ErrArtifactTypeMismatch)
	}
	payload, err := store.readEvidence(ctx, conn, record.Digest)
	if err != nil {
		return zero, err
	}
	artifact, err := parse(payload)
	if err != nil {
		return zero, wrap(OperationRead, PhaseValidate, err)
	}
	return artifact, nil
}

func (s *sqlStore) readEvidence(ctx context.Context, conn *sql.Conn, digest string) ([]byte, error) {
	var payload []byte
	if err := s.queries.artifactEvidenceGet(
		ctx,
		conn,
		artifactEvidenceGetParams{digest: digest},
	).Scan(&payload); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrap(OperationRead, PhaseRead, ErrNotFound)
		}
		return nil, wrap(OperationRead, PhaseRead, err)
	}
	if storedDigest, _ := identifyArtifact(payload); storedDigest != digest {
		return nil, wrap(OperationRead, PhaseRead, ErrEvidenceCorrupt)
	}
	return payload, nil
}

// QueryFindingsSummary returns severity-and-verdict buckets for the Security Posture dashboard.
// It is narrower than the future compliance posture dashboard.
func (c *Client) QueryFindingsSummary(
	ctx context.Context,
	scopeIDs []string,
) ([]FindingsSummaryRow, error) {
	if c == nil || c.store == nil {
		return nil, wrap(OperationQuery, PhaseInput, ErrNilDatabase)
	}
	scope, err := scopeValue(scopeIDs)
	if err != nil {
		return nil, wrap(OperationQuery, PhaseInput, err)
	}
	conn, err := c.store.db.Conn(ctx)
	if err != nil {
		return nil, wrap(OperationQuery, PhaseConnect, err)
	}
	defer func() { _ = conn.Close() }()

	if err = c.store.prepareConnection(ctx, conn); err != nil {
		return nil, wrap(OperationQuery, PhaseConnect, err)
	}
	statement := "BEGIN"
	if c.store.dialect == dialectPostgres {
		statement = "BEGIN READ ONLY"
	}
	if _, err = conn.ExecContext(ctx, statement); err != nil {
		return nil, wrap(OperationQuery, PhaseBegin, err)
	}
	committed := false
	defer rollbackUnlessCommitted(ctx, conn, &committed)
	if c.store.dialect == dialectPostgres {
		if err = c.store.queries.scopeSet(
			ctx,
			conn,
			scopeSetParams{scopeIds: scope},
		); err != nil {
			return nil, wrap(OperationQuery, PhaseScope, err)
		}
	}

	rows, err := c.store.queries.findingsSummaryList(
		ctx,
		conn,
		findingsSummaryListParams{scopeIds: scope},
	)
	if err != nil {
		return nil, wrap(OperationQuery, PhaseRead, err)
	}
	result, err := scanFindingsSummary(rows)
	if err != nil {
		return nil, wrap(OperationQuery, PhaseRead, err)
	}
	if _, err = conn.ExecContext(ctx, "COMMIT"); err != nil {
		return nil, wrap(OperationQuery, PhaseCommit, err)
	}
	committed = true
	return result, nil
}

func scanFindingsSummary(rows *sql.Rows) (result []FindingsSummaryRow, err error) {
	defer func() { err = errors.Join(err, rows.Close()) }()
	for rows.Next() {
		var row FindingsSummaryRow
		if err := rows.Scan(
			&row.ScopeID,
			&row.SubjectID,
			&row.RunID,
			&row.LayerID,
			&row.Repository,
			&row.Severity,
			&row.Verdict,
			&row.FindingCount,
		); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
