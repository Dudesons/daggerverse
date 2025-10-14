// Allow to create or read secret from AWS secret manager

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awssm "github.com/aws/aws-sdk-go/service/secretsmanager"
	"main/internal/dagger"
)

const (
	awsCredentialsPath = "/root/.aws"
)

type AwsSecretManager struct {
	secretManagerClient *awssm.SecretsManager
	// +private
	AccessKey string
	// +private
	SecretKey string
	// +private
	Region string
	// +private
	Profile string
	// +private
	AwsFolder *dagger.Directory
	// +private
	InternalImage string
}

func newAwsSecretManager(internalImage string) *AwsSecretManager {
	return &AwsSecretManager{InternalImage: internalImage}
}

// Authenticate to AWS using access and secret key
func (m *AwsSecretManager) WithCredentialsKeys(
	accessKey string,
	secretKey string,
) *AwsSecretManager {
	m.AccessKey = accessKey
	m.SecretKey = secretKey

	return m
}

// Authenticate to AWS using .aws folder
func (m *AwsSecretManager) WithCredentialsFolder(
	awsFolder *dagger.Directory,
) *AwsSecretManager {
	m.AwsFolder = awsFolder
	return m
}

// Authenticate to AWS using access and secret key
func (m *AwsSecretManager) WithRegion(name string) *AwsSecretManager {
	m.Region = name

	return m
}

func (m *AwsSecretManager) WithProfile(name string) *AwsSecretManager {
	m.Profile = name

	return m
}

func (m *AwsSecretManager) auth(ctx context.Context) error {
	config := &aws.Config{}

	if m.Region != "" {
		config.Region = aws.String(m.Region)
	}

	if m.AwsFolder != nil {
		// Sync folder to sandbox
		_, err := dag.
			Container().
			From(m.InternalImage).
			WithMountedDirectory(awsCredentialsPath, m.AwsFolder).
			Directory(awsCredentialsPath).
			Export(ctx, awsCredentialsPath)
		if err != nil {
			return err
		}
	} else if m.AccessKey != "" && m.SecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(m.AccessKey, m.SecretKey, "")
	}

	sessionOpts := session.Options{
		Config: *config,
	}

	if m.Profile != "" {
		sessionOpts.Profile = m.Profile
		sessionOpts.SharedConfigState = session.SharedConfigEnable
	}

	sess, err := session.NewSessionWithOptions(sessionOpts)
	if err != nil {
		return err
	}
	m.secretManagerClient = awssm.New(sess)

	return nil
}

// Retrieve a secret from SecretsManager
func (m *AwsSecretManager) GetSecret(
	ctx context.Context,
	name string,
) (*dagger.Secret, error) {
	err := m.auth(ctx)
	if err != nil {
		return nil, err
	}

	input := &awssm.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	value, err := m.secretManagerClient.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	return dag.SetSecret(name, *(value.SecretString)), nil
}

// Create or update a secret value
func (m *AwsSecretManager) SetSecret(
	ctx context.Context,
	name string,
	value string,
) (string, error) {
	err := m.auth(ctx)
	if err != nil {
		return "", err
	}

	input := &awssm.PutSecretValueInput{
		SecretId:     aws.String(name),
		SecretString: aws.String(value),
	}

	resp, err := m.secretManagerClient.PutSecretValue(input)
	if err != nil {
		return "", err
	}

	return *resp.ARN, err
}
