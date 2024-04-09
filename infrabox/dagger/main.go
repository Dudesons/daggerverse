// A module for playing on the terraform ecosystem

package main

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
) *Tf {
	return newTf(image, version, "terragrunt")
}
