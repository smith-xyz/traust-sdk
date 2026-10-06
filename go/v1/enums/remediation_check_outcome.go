// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// RemediationCheckOutcome — Values of `outcome` in remediation.schema.json.
type RemediationCheckOutcome string

const (
	RemediationCheckOutcomePass  RemediationCheckOutcome = "pass"
	RemediationCheckOutcomeFail  RemediationCheckOutcome = "fail"
	RemediationCheckOutcomeSkip  RemediationCheckOutcome = "skip"
	RemediationCheckOutcomeError RemediationCheckOutcome = "error"
)

// RemediationCheckOutcomeValues returns all valid RemediationCheckOutcome values.
func RemediationCheckOutcomeValues() []RemediationCheckOutcome {
	return []RemediationCheckOutcome{
		RemediationCheckOutcomePass,
		RemediationCheckOutcomeFail,
		RemediationCheckOutcomeSkip,
		RemediationCheckOutcomeError,
	}
}
