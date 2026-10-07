// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// RemediationEffort — Values of `remediation_effort` in pqc-blockers.schema.json.
type RemediationEffort string

const (
	RemediationEffortTrivial         RemediationEffort = "trivial"
	RemediationEffortModerate        RemediationEffort = "moderate"
	RemediationEffortSignificant     RemediationEffort = "significant"
	RemediationEffortBlockedExternal RemediationEffort = "blocked-external"
	RemediationEffortXs              RemediationEffort = "xs"
	RemediationEffortS               RemediationEffort = "s"
	RemediationEffortM               RemediationEffort = "m"
	RemediationEffortL               RemediationEffort = "l"
	RemediationEffortXl              RemediationEffort = "xl"
)

// RemediationEffortValues returns all valid RemediationEffort values.
func RemediationEffortValues() []RemediationEffort {
	return []RemediationEffort{
		RemediationEffortTrivial,
		RemediationEffortModerate,
		RemediationEffortSignificant,
		RemediationEffortBlockedExternal,
		RemediationEffortXs,
		RemediationEffortS,
		RemediationEffortM,
		RemediationEffortL,
		RemediationEffortXl,
	}
}
