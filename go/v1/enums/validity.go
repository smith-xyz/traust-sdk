// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

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
