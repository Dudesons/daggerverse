// Allow to encrypt, decrypt and read secret from SOPS

package main

//type SopsSecretManager struct {
//}
//
//func newsopsSecretManager() *SopsSecretManager {
//	return &SopsSecretManager{}
//}

// Read a secret from secret manager
//func (m *SopsSecretManager) GetSecret(
//	ctx context.Context,
//	// The secret name to read
//	name string,
//	// The GCP project where the secret is stored
//	project string,
//	// The version of the secret to read
//	// +optional
//	// +default="latest"
//	version string,
//	// The path to a credentials json file
//	// +optional
//	filePath *dagger.File,
//	// The path to the gcloud folder
//	// +optional
//	gcloudFolder *dagger.Directory,
//) (*dagger.Secret, error) {
//
//	return dag.SetSecret(name), nil
//}
//
//// Create or update a secret value
//func (m *SopsSecretManager) Encrypt(
//	ctx context.Context,
//	// The secret name to read
//	name string,
//	// The value to set to the secret
//	value string,
//	// The GCP project where the secret is stored
//	project string,
//	// The path to a credentials json file
//	// +optional
//	filePath *dagger.File,
//	// The path to the gcloud folder
//	// +optional
//	gcloudFolder *dagger.Directory,
//) (string, error) {
//
//}
//
//func (m *SopsSecretManager) createSecret(ctx context.Context, name string, project string) error {
//
//}
