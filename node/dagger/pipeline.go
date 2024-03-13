package main

import (
	"context"
	"strings"
)

// Execute the whole pipeline in general used with the function 'with-auto-setup'
func (n *Node) Pipeline(
	ctx context.Context,
	// Define hooks to execute before all
	// +optional
	preHooks [][]string,
	// Define hooks to execute after tests and before build
	// +optional
	postHooks [][]string,
	// Indicate if the artifact is an oci build or not
	// +optional
	isOci bool,
	// Indicate to dry run the publishing
	// +optional
	// +default="false"
	dryRun bool,
	// Define permission on the package in the registry
	// +optional
	// +default="true"
	packageAccess string,
	// Indicate if the package is publishing as development version
	// +optional
	packageDevTag string,
	// Define path to file to fetch from the build container
	// +optional
	fileContainerArtifacts []string,
	// Define path to directories to fetch from the build container
	// +optional
	directoryContainerArtifacts []string,
	// Define registries where to push the image
	// +optional
	ociRegistries []string,
	// Define the ttl registry to use
	// +optional
	// +default="ttl.sh"
	ttlRegistry string,
	// Define the ttl in the ttl registry
	// +optional
	// +default="60m"
	ttl string,
) (string, error) {
	pipeline := n.Install()

	for _, hook := range preHooks {
		pipeline = pipeline.Run(hook)
	}

	if n.DetectLint {
		pipeline = pipeline.Lint()
	}

	if n.DetectTest {
		pipeline = pipeline.Test()
	}

	// TODO(Move it at the end)
	for _, hook := range postHooks {
		pipeline = pipeline.Run(hook)
	}

	pipeline = pipeline.Build()

	if n.DetectPackage {
		return pipeline.
			Publish(packageAccess, packageDevTag, dryRun).
			Do(ctx)
	}

	if n.DetectOci || isOci {
		refs, err := pipeline.
			OciBuild(
				ctx,
				fileContainerArtifacts,
				directoryContainerArtifacts,
				ociRegistries,
				dryRun,
				ttlRegistry,
				ttl,
			)

		return strings.Join(refs, "\n"), err
	}

	return pipeline.Do(ctx)
}
