// Allow to create or read secret from GCP secret manager

package main

type SecretManager struct {
}

func (m *SecretManager) Gcp() *GcpSecretManager {
	return newGcpSecretManager()
}

func (m *SecretManager) Aws() *AwsSecretManager {
	return newAwsSecretManager()
}
