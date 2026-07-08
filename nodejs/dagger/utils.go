package main

import "main/internal/dagger"

// Return the current container state
func (n *Nodejs) Container() *dagger.Container {
	return n.Ctr
}

// Return a directory by default the current working directory
func (n *Nodejs) Directory(
	// Define permission on the package in the registry
	// +optional
	path string,
) *dagger.Directory {
	if path == "" {
		return n.Ctr.Directory(workdir)
	}

	return n.Ctr.Directory(path)
}

// Open a shell in the current container or execute a command inside it, like node
func (n *Nodejs) Shell(
	// The command to execute in the terminal
	// +optional
	cmd []string,
) *dagger.Container {
	return n.Ctr.WithDefaultTerminalCmd(cmd).Terminal()
}

// Expose the container as a service
func (n *Nodejs) Serve() *dagger.Service {
	return n.Ctr.AsService()
}

func (n *Nodejs) getCacheKey(cacheKey string) string {
	if n.PipelineID != "" {
		cacheKey = n.PipelineID + "-" + cacheKey
	}

	if n.IsProduction {
		cacheKey = cacheKey + "-prod"
	}

	return cacheKey
}
