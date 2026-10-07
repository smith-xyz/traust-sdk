// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ComplianceVerdictSource — deterministic controls: 'check' (code-computed) or 'human_override' — never 'agent' (amendment 1)
type ComplianceVerdictSource string

const (
	ComplianceVerdictSourceCheck         ComplianceVerdictSource = "check"
	ComplianceVerdictSourceAgent         ComplianceVerdictSource = "agent"
	ComplianceVerdictSourceHumanOverride ComplianceVerdictSource = "human_override"
)

// ComplianceVerdictSourceValues returns all valid ComplianceVerdictSource values.
func ComplianceVerdictSourceValues() []ComplianceVerdictSource {
	return []ComplianceVerdictSource{
		ComplianceVerdictSourceCheck,
		ComplianceVerdictSourceAgent,
		ComplianceVerdictSourceHumanOverride,
	}
}
