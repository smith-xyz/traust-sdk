// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// RevalidationMethod — Values of `method` in remediation.schema.json.
type RevalidationMethod string

const (
	RevalidationMethodValidateLive     RevalidationMethod = "validate-live"
	RevalidationMethodValidatePlatform RevalidationMethod = "validate-platform"
	RevalidationMethodValidateFindings RevalidationMethod = "validate-findings"
	RevalidationMethodManual           RevalidationMethod = "manual"
	RevalidationMethodNone             RevalidationMethod = "none"
)

// RevalidationMethodValues returns all valid RevalidationMethod values.
func RevalidationMethodValues() []RevalidationMethod {
	return []RevalidationMethod{
		RevalidationMethodValidateLive,
		RevalidationMethodValidatePlatform,
		RevalidationMethodValidateFindings,
		RevalidationMethodManual,
		RevalidationMethodNone,
	}
}
