package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
)

func (c *Ci) Terrabox(ctx context.Context, testDataSrc *Directory) error {
	var eg errgroup.Group

	eg.Go(func() error {
		_, err := dag.
			Terrabox().
			Terragrunt().
			WithSource("/terraform", testDataSrc.Directory("terraform")).
			DisableColor().
			Plan("/terraform/stacks/dev/europe-west1/staging/qux").
			Apply("/terraform/stacks/dev/europe-west1/staging/qux").
			Plan("/terraform/stacks/dev/europe-west1/staging/qux", TerraboxTfPlanOpts{DetailedExitCode: true}).
			Do(ctx)

		return err
	})

	eg.Go(func() error {
		_, err := dag.
			Terrabox().
			Terragrunt().
			WithSource("/terraform", testDataSrc.Directory("terraform")).
			DisableColor().
			Plan("/terraform/stacks/dev/europe-west1/staging/foo", TerraboxTfPlanOpts{DetailedExitCode: true}).
			Do(ctx)

		if err == nil {
			return errors.New("it should failed because the stack was not applied")
		}

		return nil
	})

	return eg.Wait()
}
