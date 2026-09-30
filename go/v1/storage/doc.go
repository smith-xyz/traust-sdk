// Package storage records v1 artifact evidence and binds it to caller-owned context.
//
// Callers own the database pool, the object store, and the preprocessing
// boundary. PostgreSQL
// relations use the fixed traust_storage schema; SQLite uses the caller-selected
// database file as its physical namespace. Each named Save operation validates
// a raw-preserving generated artifact, writes the exact bytes supplied to Save
// to the caller's ObjectStore under their sha256 digest, then records the digest
// and byte size, an opaque scope/subject/run/layer binding, and any approved
// projection in one database transaction. Optional enrichment,
// including Ledger fingerprint stamping, happens before Save; storage neither
// computes nor promotes that enrichment as independent authority.
//
// Typed Get operations resolve a binding, read the bytes from the ObjectStore,
// check them against the bound digest, and validate them. GetEvidence returns
// bytes without a type claim. Every artifact
// has a schema-specific SQL projection. Consumers correlate Ledger disposition
// data through the caller-supplied layer ID and artifact-relative finding ID;
// fingerprints remain payload evidence or Ledger-owned state, not storage join
// keys. Stored strings remain untrusted when rendered and require
// context-appropriate output encoding.
package storage
