// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ValidationTechnique — Values of `technique` in validation.schema.json.
type ValidationTechnique string

const (
	ValidationTechniqueReplay  ValidationTechnique = "replay"
	ValidationTechniqueAdapted ValidationTechnique = "adapted"
	ValidationTechniqueChained ValidationTechnique = "chained"
	ValidationTechniqueNovel   ValidationTechnique = "novel"
	ValidationTechniqueRecon   ValidationTechnique = "recon"
	ValidationTechniqueSkip    ValidationTechnique = "skip"
)

// ValidationTechniqueValues returns all valid ValidationTechnique values.
func ValidationTechniqueValues() []ValidationTechnique {
	return []ValidationTechnique{
		ValidationTechniqueReplay,
		ValidationTechniqueAdapted,
		ValidationTechniqueChained,
		ValidationTechniqueNovel,
		ValidationTechniqueRecon,
		ValidationTechniqueSkip,
	}
}
