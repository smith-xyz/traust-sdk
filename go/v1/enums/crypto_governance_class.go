// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// CryptoGovernanceClass — Values of `class` in pqc-readiness.schema.json.
type CryptoGovernanceClass string

const (
	CryptoGovernanceClassSelf            CryptoGovernanceClass = "self"
	CryptoGovernanceClassLanguageRuntime CryptoGovernanceClass = "language-runtime"
	CryptoGovernanceClassPlatform        CryptoGovernanceClass = "platform"
	CryptoGovernanceClassInfrastructure  CryptoGovernanceClass = "infrastructure"
)

// CryptoGovernanceClassValues returns all valid CryptoGovernanceClass values.
func CryptoGovernanceClassValues() []CryptoGovernanceClass {
	return []CryptoGovernanceClass{
		CryptoGovernanceClassSelf,
		CryptoGovernanceClassLanguageRuntime,
		CryptoGovernanceClassPlatform,
		CryptoGovernanceClassInfrastructure,
	}
}
