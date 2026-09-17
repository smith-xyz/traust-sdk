package storage

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Binding associates evidence with opaque caller-owned workflow identities.
type Binding struct {
	ScopeID             string  `json:"scope_id,omitempty"`
	SubjectID           *string `json:"subject_id,omitempty"`
	RunID               *string `json:"run_id,omitempty"`
	LayerID             *string `json:"layer_id,omitempty"`
	SupersedesBindingID *string `json:"supersedes_binding_id,omitempty"`
}

type BindingRecord struct {
	BindingID    string  `json:"binding_id"`
	Digest       string  `json:"digest"`
	ArtifactName string  `json:"artifact_name"`
	Binding      Binding `json:"binding"`
	BoundAt      string  `json:"bound_at"`
}

type SaveResult struct {
	Digest       string `json:"digest"`
	BindingID    string `json:"binding_id"`
	AlreadyBound bool   `json:"already_bound"`
}

type bindingRequirements struct {
	subject bool
	run     bool
	layer   bool
}

var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

func normalizedBinding(binding Binding) Binding {
	if binding.ScopeID == "" {
		binding.ScopeID = defaultScopeID
	}
	return binding
}

func validateBinding(binding Binding, requirements bindingRequirements) error {
	if binding.ScopeID == "" {
		return ErrScopeRequired
	}
	if !validBindingText(binding.ScopeID) {
		return ErrInvalidIdentifier
	}
	for _, value := range []*string{
		binding.SubjectID,
		binding.RunID,
		binding.LayerID,
		binding.SupersedesBindingID,
	} {
		if value != nil && !validBindingText(*value) {
			return ErrInvalidIdentifier
		}
	}
	if requirements.subject && binding.SubjectID == nil {
		return ErrSubjectIDRequired
	}
	if requirements.run && binding.RunID == nil {
		return ErrRunIDRequired
	}
	if requirements.layer && binding.LayerID == nil {
		return ErrLayerIDRequired
	}
	return nil
}

func scopeValue(ids []string) (string, error) {
	if len(ids) == 0 {
		return "", ErrScopeRequired
	}
	for _, id := range ids {
		if id == "" || !validBindingText(id) {
			return "", ErrInvalidIdentifier
		}
	}
	encoded, err := json.Marshal(ids)
	return string(encoded), err
}

func validBindingText(value string) bool {
	return utf8.ValidString(value) && !strings.ContainsRune(value, 0)
}
