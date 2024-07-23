package main

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
	"strconv"
)

// Return the Node container with the right base image
func (n *Node) WithVersion(
	// The image name to use
	// +optional
	// +default="node"
	image string,
	// The version of the image to use
	version string,
	// Define if the image to use is an alpine or not
	// +optional
	// +default="true"
	isAlpine bool,
) *Node {
	baseImage := image + ":" + version
	if isAlpine {
		baseImage += "-alpine"
	}
	n.Ctr = dag.Container().From(baseImage)

	n.BaseImageRef = baseImage

	return n
}

// Return the Node container with an environment variable to use in your npmrc file
func (n *Node) WithNpmrcTokenEnv(
	// The name of the environment variable where the npmrc token is stored
	name string,
	// The value of the token
	value *dagger.Secret,
) *Node {
	n.NpmrcTokenName = name
	n.NpmrcToken = value
	n.Ctr = n.Ctr.WithSecretVariable(name, value)

	return n
}

// Return the Node container with npmrc file
func (n *Node) WithNpmrcTokenFile(
	// The npmrc file to inject in the container
	npmrcFile *dagger.Secret,
) *Node {
	n.NpmrcFile = npmrcFile
	n.Ctr = n.Ctr.WithMountedSecret(workdir+"/.npmrc", npmrcFile)

	return n
}

// Return the Node container setup with the right package manager and optionaly cache setup
func (n *Node) WithPackageManager(
	// The package manager to use
	packageManager string,
	// Disable mounting cache volumes.
	// +optional
	disableCache bool,
	// Define a specific version of the package manager.
	// +optional
	version string,
) *Node {
	switch packageManager {
	case "npm":
		return n.WithNpm(disableCache, version)
	case "yarn":
		return n.WithYarn(disableCache, version)
	default:
		return n.WithNpm(disableCache, version)
	}
}

// Return the Node container with npm setup as an entrypoint and npm cache setup
func (n *Node) WithNpm(
	// Disable mounting cache volumes.
	// +optional
	disableCache bool,
	// Define a specific version of npm.
	// +optional
	version string,
) *Node {
	n.PkgMgr = "npm"

	if !disableCache {
		n.Ctr = n.
			Ctr.
			WithMountedCache("/root/.npm", dag.CacheVolume(n.getCacheKey("global-npm-cache")))
	}

	if version != "" {
		n.PkgMgrVersion = version

		n.Ctr = n.
			Ctr.
			WithExec([]string{"npm", "install", "-g", "npm@" + version})
	}

	return n
}

// Return the Node container with yarn setup as an entrypoint and yarn cache setup
func (n *Node) WithYarn(
	// Disable mounting cache volumes.
	// +optional
	disableCache bool,
	// Define a specific version of npm.
	// +optional
	version string,
) *Node {
	n.PkgMgr = "yarn"

	if !disableCache {
		n.Ctr = n.
			Ctr.
			WithMountedCache("/usr/local/share/.cache/yarn", dag.CacheVolume(n.getCacheKey("global-yarn-cache")))
	}

	if version != "" {
		n.PkgMgrVersion = version

		n.Ctr = n.
			Ctr.
			WithExec([]string{"yarn", "set", "version", version})
	}

	return n
}

// Return the Node container with the source code, 'node_modules' cache set up and workdir set
func (n *Node) WithSource(
	// The source code
	src *dagger.Directory,
	// Indicate if the directory is mounted or persisted in the container
	// +optional
	persisted bool,
) *Node {
	if persisted {
		n.Ctr = n.
			Ctr.
			WithDirectory(workdir, src).
			WithWorkdir(workdir)
	} else {
		n.Ctr = n.
			Ctr.
			WithMountedDirectory(workdir, src).
			WithMountedCache(workdir+"/node_modules", dag.CacheVolume(n.getCacheKey("node-modules")))
	}

	n.Ctr = n.Ctr.WithWorkdir(workdir)

	return n
}

// Return the Node container with an additional file in the working dir
func (n *Node) WithFile(
	// The file to use
	file *dagger.File,
	// The path where the file should be mounted
	path string,
	// Indicate if the file is mounted or persisted in the container
	// +optional
	persisted bool,
) *Node {
	if persisted {
		n.Ctr = n.
			Ctr.
			WithFile(workdir+"/"+path, file)
	} else {
		n.Ctr = n.
			Ctr.
			WithMountedFile(workdir+"/"+path, file)
	}

	return n
}

// Return the Node container with an additional directory in the working dir
func (n *Node) WithDirectory(
	// The directory to use
	dir *dagger.Directory,
	// The path where the directory should be mounted
	path string,
	// Indicate if the directory is mounted or persisted in the container
	// +optional
	persisted bool,
) *Node {
	if persisted {
		n.Ctr = n.
			Ctr.
			WithDirectory(workdir+"/"+path, dir)
	} else {
		n.Ctr = n.
			Ctr.
			WithMountedDirectory(workdir+"/"+path, dir)
	}

	return n
}

// Return the Node container with an additional cache volume in the working dir
func (n *Node) WithCache(
	// The cache volume to use
	cache *dagger.CacheVolume,
	// The path where the cache volume should be mounted
	path string,
	// Indicate if the cache volume is mounted or persisted in the container
	// +optional
	persisted bool,
) *Node {
	if persisted {
		tmpPath := "/tmp/" + uuid.New().String() + path
		n.Ctr = n.
			Ctr.
			WithMountedCache(tmpPath, cache).
			WithExec([]string{"cp", "r", tmpPath, workdir + "/" + path})

	} else {
		n.Ctr = n.
			Ctr.
			WithMountedCache(workdir+"/"+path, cache)
	}

	return n
}

// Return a node container with the 'NODE_ENV' set to production
func (n *Node) Production() *Node {
	n.IsProduction = true

	n.Ctr = n.
		Ctr.
		WithEnvVariable("NODE_ENV", "production")
	return n
}

// Execute a command from the package.json
func (n *Node) Run(
	// Command from the package.json to run
	command []string,
) *Node {
	n.Ctr = n.
		Ctr.
		WithExec(append([]string{n.PkgMgr, "run"}, command...))
	return n
}

// Install node modules
func (n *Node) Install() *Node {
	n.Ctr = n.Ctr.WithExec([]string{n.PkgMgr, "install"})
	return n
}

// Execute lint command
func (n *Node) Lint() *Node {
	return n.Run([]string{"lint"})
}

// Execute test command
func (n *Node) Test() *Node {
	return n.Run([]string{"test"})
}

// Execute test commands in parallel
func (n *Node) ParallelTest(
	ctx context.Context,
	cmds [][]string,
) error {
	var eg errgroup.Group

	for _, cmd := range cmds {
		eg.Go(func() error {
			_, err := n.Run(cmd).Do(ctx)
			return err
		})
	}

	return eg.Wait()
}

// Execute clean command
func (n *Node) Clean() *Node {
	return n.Run([]string{"clean"})
}

// Execute the build command
func (n *Node) Build() *Node {
	return n.Run([]string{"build"})
}

// Execute the publish which push a package to a registry
func (n *Node) Publish(
	// Define permission on the package in the registry
	// +optional
	access string,
	// Indicate if the package is publishing as development version
	// +optional
	devTag string,
	// Indicate to dry run the publishing
	// +optional
	dryRun bool,
) *Node {
	publishCmd := []string{n.PkgMgr, "publish"}

	if access != "" {
		publishCmd = append(publishCmd, []string{"--access", access}...)
	}

	if devTag != "" {
		publishCmd = append(publishCmd, []string{"--tag", devTag}...)
	}

	if dryRun {
		publishCmd = append(publishCmd, []string{"--dry-run", strconv.FormatBool(dryRun)}...)
	}

	n.Ctr = n.Ctr.WithExec(publishCmd)
	return n
}

// Bump the package version
func (n *Node) BumpVersion(
	// Define the bump version strategy (major | minor | patch | premajor | preminor | prepatch | prerelease)
	strategy string,
	// The message will use it as a commit message when creating a version commit. If the message config contains %s then that will be replaced with the resulting version number
	// +optional
	message string,
) *Node {
	versionCmd := []string{n.PkgMgr, "version", strategy}

	if message != "" {
		versionCmd = append(versionCmd, []string{"-m", message}...)
	}

	n.Ctr = n.Ctr.WithExec(versionCmd)
	return n
}
