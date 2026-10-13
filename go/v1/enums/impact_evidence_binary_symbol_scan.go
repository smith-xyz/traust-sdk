// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ImpactEvidenceBinarySymbolScan — Vulnerable symbol names present in the shipped binary's strings (.dynstr) — the binary references the vulnerable function, not just the library.
type ImpactEvidenceBinarySymbolScan string

const (
	ImpactEvidenceBinarySymbolScanSymbolsPresent ImpactEvidenceBinarySymbolScan = "symbols_present"
	ImpactEvidenceBinarySymbolScanSymbolsAbsent  ImpactEvidenceBinarySymbolScan = "symbols_absent"
)

// ImpactEvidenceBinarySymbolScanValues returns all valid ImpactEvidenceBinarySymbolScan values.
func ImpactEvidenceBinarySymbolScanValues() []ImpactEvidenceBinarySymbolScan {
	return []ImpactEvidenceBinarySymbolScan{
		ImpactEvidenceBinarySymbolScanSymbolsPresent,
		ImpactEvidenceBinarySymbolScanSymbolsAbsent,
	}
}
