package main

import (
	"context"
	"main/internal/dagger"
	"strings"
)

// Allow to let the pipeline to be setup automatically based on the package.json aka lazy mode
func (n *Node) WithAutoSetup(
	ctx context.Context,
	// A name to use in the pipeline and injected in cache keys
	pipelineId string,
	// The code to mount in the node container
	src *dagger.Directory,
	// Patterns to exclude in the analysis for the auto-detection
	// +optional
	patternExclusions []string,
	// The image name to use
	// +optional
	// +default="node"
	image string,
	// Define if the image to use is an alpine or not
	// +optional
	// +default="true"
	isAlpine bool,
	// Container options
	// +optional
	// +default="linux/amd64"
	containerPlatform dagger.Platform,
	// Indicate attempted system package to install
	// +optional
	systemSetupCmds [][]string,
	// Define a specific version of the package manager.
	// +optional
	packageManagerVersion string,
	// Node workspaces to use during the pipeline
	// +optional
	workspaces []string,
) (*Node, error) {
	var err error
	nodeAutoSetup := &Node{
		PipelineID:      pipelineId,
		PkgMgr:          "npm",
		Platform:        containerPlatform,
		SystemSetupCmds: systemSetupCmds,
		Workspaces:      workspaces,
		Ctr: dag.
			Container(dagger.ContainerOpts{
				Platform: containerPlatform,
			}),
	}

	nodeAnalyzer := dag.
		Autodetection().
		Node(
			src,
			dagger.AutodetectionNodeOpts{
				PatternExclusions: append(
					[]string{"node_modules"},
					patternExclusions...,
				),
			},
		)
	n.DetectOci, err = dag.
		Autodetection().
		Oci(
			src,
			dagger.AutodetectionOciOpts{
				PatternExclusions: append(
					[]string{"node_modules"},
					patternExclusions...,
				),
			},
		).
		IsOci(ctx)
	if err != nil {
		return nil, err
	}

	engineVersion, err := nodeAnalyzer.GetEngineVersion(ctx)
	if err != nil {
		return nil, err
	}

	isYarn, err := nodeAnalyzer.IsYarn(ctx)
	if err != nil {
		return nil, err
	}
	if isYarn {
		nodeAutoSetup.PkgMgr = "yarn"
	}

	appVersion, err := nodeAnalyzer.GetVersion(ctx)
	if err != nil {
		return nil, err
	}
	nodeAutoSetup.Version = appVersion

	appName, err := nodeAnalyzer.GetName(ctx)
	if err != nil {
		return nil, err
	}
	nodeAutoSetup.Name = appName

	nodeAutoSetup.DetectTest, err = nodeAnalyzer.IsTest(ctx)
	if err != nil {
		return nil, err
	}

	nodeAutoSetup.DetectPackage, err = nodeAnalyzer.IsPackage(ctx)
	if err != nil {
		return nil, err
	}

	nodeAutoSetup.DetectLint, err = nodeAnalyzer.Is(ctx, "lint")
	if err != nil {
		return nil, err
	}

	rootWorkspacePaths, err := nodeAnalyzer.GetWorkspaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, i := range rootWorkspacePaths {
		nodeAutoSetup.RootWorkspacePaths = append(nodeAutoSetup.RootWorkspacePaths, strings.ReplaceAll(i, "*", ""))
	}

	return nodeAutoSetup.
			WithVersion(image, engineVersion, isAlpine).
			WithSource(src, false).
			WithPackageManager(n.PkgMgr, false, packageManagerVersion),
		nil
}
