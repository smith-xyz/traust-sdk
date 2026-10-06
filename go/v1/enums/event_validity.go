// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// EventValidity — Finding validity for disposition events. 4-value set (not_verified is the audit-stage default and is never set by an event). 'hardening' = defense-in-depth gap with no concrete exploit path.
type EventValidity string

const (
	EventValidityConfirmed     EventValidity = "confirmed"
	EventValidityFalsePositive EventValidity = "false_positive"
	EventValidityCorrected     EventValidity = "corrected"
	EventValidityHardening     EventValidity = "hardening"
)

// EventValidityValues returns all valid EventValidity values.
func EventValidityValues() []EventValidity {
	return []EventValidity{
		EventValidityConfirmed,
		EventValidityFalsePositive,
		EventValidityCorrected,
		EventValidityHardening,
	}
}
