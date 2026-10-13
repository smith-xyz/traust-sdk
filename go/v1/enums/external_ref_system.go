// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ExternalRefSystem — Values of `system` in layer.schema.json.
type ExternalRefSystem string

const (
	ExternalRefSystemCve      ExternalRefSystem = "cve"
	ExternalRefSystemBugzilla ExternalRefSystem = "bugzilla"
	ExternalRefSystemGhsa     ExternalRefSystem = "ghsa"
	ExternalRefSystemJira     ExternalRefSystem = "jira"
)

// ExternalRefSystemValues returns all valid ExternalRefSystem values.
func ExternalRefSystemValues() []ExternalRefSystem {
	return []ExternalRefSystem{
		ExternalRefSystemCve,
		ExternalRefSystemBugzilla,
		ExternalRefSystemGhsa,
		ExternalRefSystemJira,
	}
}
