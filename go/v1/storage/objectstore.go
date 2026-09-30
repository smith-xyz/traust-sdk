package storage

import "context"

// ObjectStore keeps artifact bytes, addressed by their sha256 digest.
//
// storage/v1 records only an artifact's digest and byte size, so every Client
// needs a place for the bytes themselves. The SDK does not assume one:
// implementations (S3, GCS, Azure Blob, a local directory, ...) live outside
// this module, the same way callers bring their own *sql.DB. The SDK owns digest
// computation, writing bytes before the database commit, and checking bytes
// against their digest on read. An ObjectStore owns only transport.
//
// Implementations must:
//   - accept a repeated PutArtifact for the same digest, because Save retries
//     rely on it;
//   - return an error wrapping ErrNotFound from GetArtifact for an unknown
//     digest;
//   - key objects with ObjectKey, so every SDK language reads what another
//     wrote.
type ObjectStore interface {
	PutArtifact(ctx context.Context, digest string, payload []byte) error
	GetArtifact(ctx context.Context, digest string) ([]byte, error)
}

// ObjectKey is the object key for an artifact digest: "<prefix>sha256/<digest>".
func ObjectKey(prefix, digest string) string {
	return prefix + "sha256/" + digest
}
