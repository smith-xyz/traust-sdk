// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// PqcBlastRadius — Values of `blast_radius` in pqc-readiness.schema.json.
type PqcBlastRadius string

const (
	PqcBlastRadiusFleet     PqcBlastRadius = "fleet"
	PqcBlastRadiusProduct   PqcBlastRadius = "product"
	PqcBlastRadiusService   PqcBlastRadius = "service"
	PqcBlastRadiusComponent PqcBlastRadius = "component"
)

// PqcBlastRadiusValues returns all valid PqcBlastRadius values.
func PqcBlastRadiusValues() []PqcBlastRadius {
	return []PqcBlastRadius{
		PqcBlastRadiusFleet,
		PqcBlastRadiusProduct,
		PqcBlastRadiusService,
		PqcBlastRadiusComponent,
	}
}
