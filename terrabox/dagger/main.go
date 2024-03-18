// A generated module for Terrabox functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

type Terrabox struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Terrabox) Terragrunt(
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
