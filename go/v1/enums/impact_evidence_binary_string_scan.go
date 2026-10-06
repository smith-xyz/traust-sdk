// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ImpactEvidenceBinaryStringScan — ELF binary string scan result.
type ImpactEvidenceBinaryStringScan string

const (
	ImpactEvidenceBinaryStringScanPackagePathPresent ImpactEvidenceBinaryStringScan = "package_path_present"
	ImpactEvidenceBinaryStringScanPackagePathAbsent  ImpactEvidenceBinaryStringScan = "package_path_absent"
	ImpactEvidenceBinaryStringScanNotScanned         ImpactEvidenceBinaryStringScan = "not_scanned"
)

// ImpactEvidenceBinaryStringScanValues returns all valid ImpactEvidenceBinaryStringScan values.
func ImpactEvidenceBinaryStringScanValues() []ImpactEvidenceBinaryStringScan {
	return []ImpactEvidenceBinaryStringScan{
		ImpactEvidenceBinaryStringScanPackagePathPresent,
		ImpactEvidenceBinaryStringScanPackagePathAbsent,
		ImpactEvidenceBinaryStringScanNotScanned,
	}
}
