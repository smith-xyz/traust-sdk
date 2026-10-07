// Code generated from traust-contracts 6501442e448bb5a78b7e58c6745d0eea37d0e7ac. DO NOT EDIT.

package enums

// CloudProvider — Values of `provider` in cloud-config-audit.schema.json.
type CloudProvider string

const (
	CloudProviderAws        CloudProvider = "aws"
	CloudProviderAzure      CloudProvider = "azure"
	CloudProviderGcp        CloudProvider = "gcp"
	CloudProviderKubernetes CloudProvider = "kubernetes"
	CloudProviderDocker     CloudProvider = "docker"
	CloudProviderOther      CloudProvider = "other"
)

// CloudProviderValues returns all valid CloudProvider values.
func CloudProviderValues() []CloudProvider {
	return []CloudProvider{
		CloudProviderAws,
		CloudProviderAzure,
		CloudProviderGcp,
		CloudProviderKubernetes,
		CloudProviderDocker,
		CloudProviderOther,
	}
}
