package types

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/traust-security/traust-sdk/go/v1/validate"
)

var ErrEmptyArtifact = errors.New("artifact payload is empty")

// Artifact binds schema-valid JSON bytes to their generated Go type.
type Artifact[T any] struct {
	payload []byte
}

// ParseArtifact validates and retains exact bytes for a generated Go type.
func ParseArtifact[T any](schema string, payload []byte) (Artifact[T], error) {
	if err := validate.ValidateBytes(schema, payload); err != nil {
		return Artifact[T]{}, err
	}
	return Artifact[T]{payload: bytes.Clone(payload)}, nil
}

// Payload returns an owned copy of the exact artifact bytes.
func (a Artifact[T]) Payload() []byte {
	return bytes.Clone(a.payload)
}

// Value decodes the exact artifact bytes into their generated Go type.
func (a Artifact[T]) Value() (T, error) {
	var value T
	if len(a.payload) == 0 {
		return value, ErrEmptyArtifact
	}
	decoder := json.NewDecoder(bytes.NewReader(a.payload))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return value, err
	}
	return value, nil
}
