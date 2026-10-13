// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package types

type AttackMapping struct {
	AttackVersion  string                 `json:"attack_version"`
	Attribution    string                 `json:"attribution"`
	CapabilityMap  map[string]interface{} `json:"capability_map"`
	CategoryMap    map[string]interface{} `json:"category_map"`
	Documentation  *string                `json:"documentation,omitempty"`
	MappingVersion string                 `json:"mapping_version"`
	Schema         *string                `json:"schema,omitempty"`
	Source         string                 `json:"source"`
}
