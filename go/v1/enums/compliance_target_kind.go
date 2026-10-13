// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ComplianceTargetKind — Values of `kind` in compliance-assessment.schema.json.
type ComplianceTargetKind string

const (
	ComplianceTargetKindProduct     ComplianceTargetKind = "product"
	ComplianceTargetKindEnvironment ComplianceTargetKind = "environment"
	ComplianceTargetKindBoth        ComplianceTargetKind = "both"
)

// ComplianceTargetKindValues returns all valid ComplianceTargetKind values.
func ComplianceTargetKindValues() []ComplianceTargetKind {
	return []ComplianceTargetKind{
		ComplianceTargetKindProduct,
		ComplianceTargetKindEnvironment,
		ComplianceTargetKindBoth,
	}
}
