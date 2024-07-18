package main

import "main/internal/dagger"

// Return the current container state
func (n *Node) Container() *dagger.Container {
	return n.Ctr
}

// Return the current working directory
func (n *Node) Directory() *dagger.Directory {
	return n.Ctr.Directory(workdir)
}

// Open a shell in the current container or execute a command inside it, like node
func (n *Node) Shell(
	// The command to execute in the terminal
	// +optional
	cmd []string,
) *dagger.Container {
	return n.Ctr.WithDefaultTerminalCmd(cmd).Terminal()
}

// Expose the container as a service
func (n *Node) Serve() *dagger.Service {
	return n.Ctr.AsService()
}

func (n *Node) getCacheKey(cacheKey string) string {
	if n.PipelineID != "" {
		cacheKey = n.PipelineID + "-" + cacheKey
	}

	if n.IsProduction {
		cacheKey = cacheKey + "-prod"
	}

	return cacheKey
}
