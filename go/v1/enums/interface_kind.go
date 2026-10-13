// Code generated from traust-contracts a150b2b56cee34857c7bfd12aa5c39fdcf033a8e. DO NOT EDIT.

package enums

// InterfaceKind — Values of `kind` in isolation-review.schema.json.
type InterfaceKind string

const (
	InterfaceKindApi       InterfaceKind = "api"
	InterfaceKindDataStore InterfaceKind = "data-store"
	InterfaceKindQueue     InterfaceKind = "queue"
	InterfaceKindIngress   InterfaceKind = "ingress"
	InterfaceKindWebhook   InterfaceKind = "webhook"
	InterfaceKindCli       InterfaceKind = "cli"
	InterfaceKindOther     InterfaceKind = "other"
)

// InterfaceKindValues returns all valid InterfaceKind values.
func InterfaceKindValues() []InterfaceKind {
	return []InterfaceKind{
		InterfaceKindApi,
		InterfaceKindDataStore,
		InterfaceKindQueue,
		InterfaceKindIngress,
		InterfaceKindWebhook,
		InterfaceKindCli,
		InterfaceKindOther,
	}
}
