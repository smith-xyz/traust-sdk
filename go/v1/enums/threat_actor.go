// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// ThreatActor — Weakest-position actor FIRST: scoring reads the first entry.
type ThreatActor string

const (
	ThreatActorRemoteUnauth    ThreatActor = "remote_unauth"
	ThreatActorRemoteAuth      ThreatActor = "remote_auth"
	ThreatActorAdjacentNetwork ThreatActor = "adjacent_network"
	ThreatActorLocalUser       ThreatActor = "local_user"
	ThreatActorLocalAdmin      ThreatActor = "local_admin"
	ThreatActorSupplyChain     ThreatActor = "supply_chain"
	ThreatActorInsider         ThreatActor = "insider"
)

// ThreatActorValues returns all valid ThreatActor values.
func ThreatActorValues() []ThreatActor {
	return []ThreatActor{
		ThreatActorRemoteUnauth,
		ThreatActorRemoteAuth,
		ThreatActorAdjacentNetwork,
		ThreatActorLocalUser,
		ThreatActorLocalAdmin,
		ThreatActorSupplyChain,
		ThreatActorInsider,
	}
}
