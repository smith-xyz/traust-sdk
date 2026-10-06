// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ValidationApprovalMode — Values of `mode` in validation.schema.json.
type ValidationApprovalMode string

const (
	ValidationApprovalModeInteractive ValidationApprovalMode = "interactive"
	ValidationApprovalModeAuto        ValidationApprovalMode = "auto"
	ValidationApprovalModeDryRun      ValidationApprovalMode = "dry-run"
)

// ValidationApprovalModeValues returns all valid ValidationApprovalMode values.
func ValidationApprovalModeValues() []ValidationApprovalMode {
	return []ValidationApprovalMode{
		ValidationApprovalModeInteractive,
		ValidationApprovalModeAuto,
		ValidationApprovalModeDryRun,
	}
}
