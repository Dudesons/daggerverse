// A Nodejs module for managing package, oci image, static website, running run script ...

package main

import (
	"context"
)

const (
	workdir = "/opt/app"
)

type Node struct {
	// +private
	PipelineID string
	// +private
	Ctr *Container
	// +private
	Name string
	// +private
	Version string
	// +private
	DetectTest bool
	// +private
	DetectPackage bool
	// +private
	DetectLint bool
	// +private
	DetectOci bool
	// +private
	PkgMgr string
	// +private
	Platform Platform
	// +private
	IsProduction bool
	// +private
	SystemSetupCmds [][]string
	// +private
	BaseImageRef string
	// +private
	NpmrcTokenName string
	// +private
	NpmrcToken *Secret
	// +private
	NpmrcFile *Secret
	// +private
	DistName string
}

// Define the pipeline id to use
func (n *Node) WithPipelineId(
	// The name to give to the pipeline
	pipelineID string,
) *Node {
	n.PipelineID = pipelineID

	return n
}

// Setup system component like installing packages
func (n *Node) SetupSystem(
	// Indicate attempted system package to install
	// +optional
	systemSetupCmds [][]string,
) *Node {
	n.SystemSetupCmds = append(n.SystemSetupCmds, systemSetupCmds...)

	for _, i := range n.SystemSetupCmds {
		n.Ctr = n.Ctr.WithExec(i)
	}

	return n
}

// Execute all commands
func (n *Node) Do(ctx context.Context) (string, error) {
	return n.Ctr.Stdout(ctx)
}
