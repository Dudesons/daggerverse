// Allow to create or read secret from GCP secret manager

package main

import (
	"context"
	"fmt"
	"main/internal/dagger"
	"regexp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

const (
	gcpCredentialsFilePath   = "/root/gcp_creds.json"
	gcpCredentialsGcloudPath = "/root/.config/gcloud"
)

type GcpSecretManager struct {
	secretManagerClient *secretmanager.Client
}

func newGcpSecretManager() *GcpSecretManager {
	return &GcpSecretManager{}
}

func (m *GcpSecretManager) withCredentials(
	ctx context.Context,
	filePath *dagger.File,
	gcloudFolder *dagger.Directory,
) error {
	syncFileToSandboxCtr := dag.Container().From("alpine:latest")

	if filePath != nil {
		_, err := syncFileToSandboxCtr.
			WithMountedFile(gcpCredentialsFilePath, filePath).
			File(gcpCredentialsFilePath).
			Export(ctx, gcpCredentialsFilePath)

		return err
	}

	_, err := syncFileToSandboxCtr.
		WithMountedDirectory(gcpCredentialsGcloudPath, gcloudFolder).
		Directory(gcpCredentialsGcloudPath).
		Export(ctx, gcpCredentialsGcloudPath)

	return err
}

func (m *GcpSecretManager) auth(ctx context.Context, filePath *dagger.File, gcloudFolder *dagger.Directory) error {
	var gcpOptions []option.ClientOption

	if filePath != nil || gcloudFolder != nil {
		err := m.withCredentials(ctx, filePath, gcloudFolder)
		if err != nil {
			return err
		}
	}

	if filePath != nil {
		gcpOptions = append(gcpOptions, option.WithCredentialsFile(gcpCredentialsFilePath))
	}

	client, err := secretmanager.NewClient(ctx, gcpOptions...)
	if err != nil {
		return err
	}

	m.secretManagerClient = client

	return nil
}

// Read a secret from secret manager
func (m *GcpSecretManager) GetSecret(
	ctx context.Context,
	// The secret name to read
	name string,
	// The GCP project where the secret is stored
	project string,
	// The version of the secret to read
	// +optional
	// +default="latest"
	version string,
	// The path to a credentials json file
	// +optional
	filePath *dagger.File,
	// The path to the gcloud folder
	// +optional
	gcloudFolder *dagger.Directory,
) (*dagger.Secret, error) {
	err := m.auth(ctx, filePath, gcloudFolder)
	if err != nil {
		return nil, err
	}

	response, err := m.secretManagerClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", project, name, version),
	})
	if err != nil {
		return nil, err
	}

	return dag.SetSecret(name, string(response.Payload.GetData())), nil
}

// Create or update a secret value
func (m *GcpSecretManager) SetSecret(
	ctx context.Context,
	// The secret name to read
	name string,
	// The value to set to the secret
	value string,
	// The GCP project where the secret is stored
	project string,
	// The path to a credentials json file
	// +optional
	filePath *dagger.File,
	// The path to the gcloud folder
	// +optional
	gcloudFolder *dagger.Directory,
) (string, error) {
	notFoundRegexp, err := regexp.Compile("rpc error: code = NotFound desc = Secret \\[projects/\\w+/secrets/\\w+] not found")
	if err != nil {
		return "", err
	}

	err = m.auth(ctx, filePath, gcloudFolder)
	if err != nil {
		return "", err
	}

	secretName := fmt.Sprintf("projects/%s/secrets/%s", project, name)
	_, err = m.secretManagerClient.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})

	if err != nil {
		if !notFoundRegexp.MatchString(err.Error()) {
			return "", err
		}

		err = m.createSecret(ctx, name, project)
		if err != nil {
			return "", err
		}
	}

	resp, err := m.secretManagerClient.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretName,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(value),
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Name, nil
}

func (m *GcpSecretManager) createSecret(ctx context.Context, name string, project string) error {
	_, err := m.secretManagerClient.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		SecretId: name,
		Parent:   fmt.Sprintf("projects/%s", project),
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	})

	return err
}
