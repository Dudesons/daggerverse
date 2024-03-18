package main

import (
	"bytes"
	"context"
	"github.com/joho/godotenv"
)

func (m *Utils) WithDotEnvSecret(ctx context.Context, ctr *Container, data *Secret) (*Container, error) {

	clearData, err := data.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	d, err := godotenv.Parse(bytes.NewReader([]byte(clearData)))
	if err != nil {
		return nil, err
	}

	for k, v := range d {
		ctr = ctr.WithSecretVariable(k, dag.Secret(v))
	}

	return ctr, nil
}
