// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// PqcProvenance — Values of `provenance` in pqc-decision-tree.schema.json.
type PqcProvenance string

const (
	PqcProvenanceInheritedPlatform       PqcProvenance = "inherited-platform"
	PqcProvenanceInheritedConstrained    PqcProvenance = "inherited-constrained"
	PqcProvenanceDelegatedDependency     PqcProvenance = "delegated-dependency"
	PqcProvenanceNativeFirstParty        PqcProvenance = "native-first-party"
	PqcProvenanceVendored                PqcProvenance = "vendored"
	PqcProvenanceExternalized            PqcProvenance = "externalized"
	PqcProvenanceNotAssessableFromSource PqcProvenance = "not-assessable-from-source"
)

// PqcProvenanceValues returns all valid PqcProvenance values.
func PqcProvenanceValues() []PqcProvenance {
	return []PqcProvenance{
		PqcProvenanceInheritedPlatform,
		PqcProvenanceInheritedConstrained,
		PqcProvenanceDelegatedDependency,
		PqcProvenanceNativeFirstParty,
		PqcProvenanceVendored,
		PqcProvenanceExternalized,
		PqcProvenanceNotAssessableFromSource,
	}
}
