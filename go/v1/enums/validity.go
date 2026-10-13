// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// Validity — Finding validity disposition state. Indicates whether a finding is confirmed, false positive, etc.
type Validity string

const (
	ValidityConfirmed     Validity = "confirmed"
	ValidityCorrected     Validity = "corrected"
	ValidityFalsePositive Validity = "false_positive"
	ValidityNotVerified   Validity = "not_verified"
	ValidityHardening     Validity = "hardening"
)

// ValidityValues returns all valid Validity values.
func ValidityValues() []Validity {
	return []Validity{
		ValidityConfirmed,
		ValidityCorrected,
		ValidityFalsePositive,
		ValidityNotVerified,
		ValidityHardening,
	}
}
