// A module for playing on the terraform ecosystem

package main

import "dagger/terrabox/internal/dagger"

type Infrabox struct{}

// Expose a terragrunt runtime
func (m *Infrabox) Terragrunt(
	// The image to use which contain terragrunt ecosystem
	// +optional
	// +default="alpine/terragrunt"
	image string,
	// The version of the image to use
	// +optional
	// +default="1.7.4"
	version string,
	// A container to use as a base
	// +optional
	ctr *dagger.Container,
) *Tf {
	return newTf(image, version, "terragrunt", ctr)
}
