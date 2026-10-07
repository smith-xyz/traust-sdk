// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// AuditProfile — Report profile discriminator: 'code' = secure-code-audit (source/manifest audit), 'rpm' = secure-rpm-audit (dist-git packaging audit), 'container' = secure-container-audit (registry image audit via skopeo/syft/grype).
type AuditProfile string

const (
	AuditProfileCode      AuditProfile = "code"
	AuditProfileRpm       AuditProfile = "rpm"
	AuditProfileContainer AuditProfile = "container"
)

// AuditProfileValues returns all valid AuditProfile values.
func AuditProfileValues() []AuditProfile {
	return []AuditProfile{
		AuditProfileCode,
		AuditProfileRpm,
		AuditProfileContainer,
	}
}
