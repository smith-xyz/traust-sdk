// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package types

type OrgParameters struct {
	DeclaredBy string                                  `json:"declared_by"`
	DeclaredOn *string                                 `json:"declared_on,omitempty"`
	Note       *string                                 `json:"note,omitempty"`
	Parameters map[string]OrgParametersParametersEntry `json:"parameters"`
	Version    int                                     `json:"version"`
}

type OrgParametersParametersEntry struct {
	Controls []string    `json:"controls"`
	Note     *string     `json:"note,omitempty"`
	Value    interface{} `json:"value"`
}
