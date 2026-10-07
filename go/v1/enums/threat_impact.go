// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// ThreatImpact — How bad this threat is if realised, modelled BEFORE reachability is known. Deliberately not the finding severity enum.
type ThreatImpact string

const (
	ThreatImpactLow         ThreatImpact = "low"
	ThreatImpactMedium      ThreatImpact = "medium"
	ThreatImpactHigh        ThreatImpact = "high"
	ThreatImpactCritical    ThreatImpact = "critical"
	ThreatImpactExistential ThreatImpact = "existential"
)

// ThreatImpactValues returns all valid ThreatImpact values.
func ThreatImpactValues() []ThreatImpact {
	return []ThreatImpact{
		ThreatImpactLow,
		ThreatImpactMedium,
		ThreatImpactHigh,
		ThreatImpactCritical,
		ThreatImpactExistential,
	}
}
