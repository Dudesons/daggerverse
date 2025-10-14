// Allow to create or read secret from GCP secret manager

package main

type SecretManager struct {
}

func (m *SecretManager) Gcp(
	// Used to overwrite the default image used for internal action (mainly used to avoid rate limit with dockerhub)
	// +optional
	// +default="alpine:latest"
	InternalImage string,
) *GcpSecretManager {
	return newGcpSecretManager(InternalImage)
}

func (m *SecretManager) Aws(
	// Used to overwrite the default image used for internal action (mainly used to avoid rate limit with dockerhub)
	// +optional
	// +default="alpine:latest"
	InternalImage string,
) *AwsSecretManager {
	return newAwsSecretManager(InternalImage)
}
