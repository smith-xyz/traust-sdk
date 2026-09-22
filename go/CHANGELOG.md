# Changelog

All notable changes to the Go SDK are documented here.

## [0.3.0]

## Changes

- **Storage tracks traust-contracts 0.34.0 (storage `REVISION` 16).** Four
  new fan-out tables (`impact_repo`, `threat_boundary`, `boundary_threat`,
  `doc_variance_record`), the widened `layer_event` (all nine actor fields)
  and `subject_ownership` (`report_kind`), and the rewritten
  `advisory_exposure` reach the generated bootstrap and query bindings.
- **Four new readers**, so every consumption view the contract declares has a
  Go reader (the parity gate says so): `QueryCensusDistinct`,
  `QueryCensusBranch`, `QueryBoundaryCurrent`, `QueryDocVarianceCurrent`.
  `policy_report_current` is declared intermediate in both languages.
- The hand-written layer and registry projectors write the new columns:
  `projectLayer` fills the eight actor identity fields, `projectCorpusRegistry`
  fills `report_kind`.
- The generator escapes Go keywords in generated identifiers
  (`threat_boundary.interface` produced invalid Go) and mirrors the
  contract's view creation order for the new intermediate view, which
  PostgreSQL resolves at CREATE time.

### Breaking

`CensusPopulationRow` and `CensusExposureRow` gain a trailing `ReportKind
*string`. `AdvisoryExposureRow.L1DependsOn`, `.L1VersionInRange`,
`.L4PackageImported` and `.FeaturePatternMatches` are `*int64` (the columns
are INTEGER now), no longer `*string`.

## [0.2.0]

## Changes

- Storage uses the fixed PostgreSQL `traust_storage` schema and fully qualified
  canonical SQL, preventing application-table collisions and `search_path`
  redirection. SQLite continues to use the caller-selected database file as its
  physical namespace.

- **Two write-path verbs the Python client already had are now on the Go
  client**, reaching the ledger's new REST endpoints (traust-ledger >= 0.3.0):
  - `ledger.Client.StampEventIdentities(ctx, layerID, StampInput{Fingerprints})`
    — POST `/v1/ledger/layers/{layer_id}/stamp`. Backfills event fingerprints
    from a `finding_ref -> fingerprint` map and re-signs the layer; the ledger
    never overwrites an existing fingerprint. Returns the new Merkle root and
    the count stamped.
  - `ledger.Client.Whoami(ctx)` — GET `/v1/ledger/whoami`. Returns the
    token-verified `types.Actor` without recording anything.

## [0.1.1]

## Changes

- `occurred_at` no longer gets a midnight suffix appended to a value that
  already carries a time. Three duplicate `*OccurredAt` helpers collapse into
  one `eventOccurredAt`; a report date with a time component passes through
  unchanged, and only a bare date is padded. Previously a timestamped report
  date produced `"2026-01-16T00:00:00ZT00:00:00+00:00"`, which the ledger's
  dict-only write path stores without complaint and then cannot read back.

- The schema compiler now calls `AssertFormat()`. Without it `format` is
  annotation-only, so `format: date` and `format: date-time` accepted any
  string and report validation passed values the Python side rejects.

### Note

`traust-ledger` >= 0.1.1 enforces RFC 3339 on event timestamps when reading a
layer, but its write path does not validate, so a malformed `occurred_at`
submitted by an older SDK is accepted and only fails on the next read. Upgrade
any Go producer before it writes.

## [0.1.0]

**First public release.**

Typed Go client for the traust-ledger service: convert and submit
triage/validation/verification reports as disposition events, request
layer merkle signing, and query layers/findings/events. Generated
`v1/types`, `v1/enums`, and `v1/validate` track `traust-contracts` schemas.
