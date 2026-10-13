// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// CvssPrivilegesRequired — Privileges the probe identity ACTUALLY held at success time.
type CvssPrivilegesRequired string

const (
	CvssPrivilegesRequiredNone CvssPrivilegesRequired = "none"
	CvssPrivilegesRequiredLow  CvssPrivilegesRequired = "low"
	CvssPrivilegesRequiredHigh CvssPrivilegesRequired = "high"
)

// CvssPrivilegesRequiredValues returns all valid CvssPrivilegesRequired values.
func CvssPrivilegesRequiredValues() []CvssPrivilegesRequired {
	return []CvssPrivilegesRequired{
		CvssPrivilegesRequiredNone,
		CvssPrivilegesRequiredLow,
		CvssPrivilegesRequiredHigh,
	}
}
