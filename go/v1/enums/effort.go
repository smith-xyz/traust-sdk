// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// Effort — Canonical effort-size labels for explicitly selected future-write contracts.
type Effort string

const (
	EffortXs Effort = "xs"
	EffortS  Effort = "s"
	EffortM  Effort = "m"
	EffortL  Effort = "l"
	EffortXl Effort = "xl"
)

// EffortValues returns all valid Effort values.
func EffortValues() []Effort {
	return []Effort{
		EffortXs,
		EffortS,
		EffortM,
		EffortL,
		EffortXl,
	}
}
