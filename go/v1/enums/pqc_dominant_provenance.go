// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// PqcDominantProvenance — Values of `dominant` in pqc-readiness.schema.json.
type PqcDominantProvenance string

const (
	PqcDominantProvenanceInheritedPlatform       PqcDominantProvenance = "inherited-platform"
	PqcDominantProvenanceInheritedConstrained    PqcDominantProvenance = "inherited-constrained"
	PqcDominantProvenanceDelegatedDependency     PqcDominantProvenance = "delegated-dependency"
	PqcDominantProvenanceNativeFirstParty        PqcDominantProvenance = "native-first-party"
	PqcDominantProvenanceVendored                PqcDominantProvenance = "vendored"
	PqcDominantProvenanceExternalized            PqcDominantProvenance = "externalized"
	PqcDominantProvenanceNotAssessableFromSource PqcDominantProvenance = "not-assessable-from-source"
	PqcDominantProvenanceNone                    PqcDominantProvenance = "none"
)

// PqcDominantProvenanceValues returns all valid PqcDominantProvenance values.
func PqcDominantProvenanceValues() []PqcDominantProvenance {
	return []PqcDominantProvenance{
		PqcDominantProvenanceInheritedPlatform,
		PqcDominantProvenanceInheritedConstrained,
		PqcDominantProvenanceDelegatedDependency,
		PqcDominantProvenanceNativeFirstParty,
		PqcDominantProvenanceVendored,
		PqcDominantProvenanceExternalized,
		PqcDominantProvenanceNotAssessableFromSource,
		PqcDominantProvenanceNone,
	}
}
