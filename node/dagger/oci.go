package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"main/internal/dagger"
	"strings"
)

// Build a production image and push to one or more registries
func (n *Node) OciBuild(
	ctx context.Context,
	// Define folder names to map in the working directory, these names match 'folder-artifacts'
	// +optional
	folderArtifactNames []string,
	// Define folders to map in the working directory, these folders match 'folder-artifact-names'
	// +optional
	folderArtifacts []*Directory,
	// Define files to mount in the working directory, these names match 'file-artifact-names'
	// +optional
	fileArtifactNames []string,
	// Define file names to map in the working directory, these names match 'file-artifacts'
	// +optional
	fileArtifacts []*File,
	// Define path to fo file to fetch from the build container
	// +optional
	fileContainerArtifacts []string,
	// Define path to fo directories to fetch from the build container
	// +optional
	directoryContainerArtifacts []string,
	// Define registries where to push the image
	registries []string,
	// Define the ttl registry to use
	// +optional
	isTtl bool,
	// Define the ttl registry to use
	// +optional
	// +default="ttl.sh"
	ttlRegistry string,
	// Define the ttl in the ttl registry
	// +optional
	// +default="60m"
	ttl string,
) ([]string, error) {
	var err error
	var eg errgroup.Group
	var fullyQualifiedImageNames []string

	if n.DistName == "" {
		n.DistName = "dist"
	}

	result := make(chan string)
	productionBuild := &Node{
		PipelineID:      n.PipelineID,
		PkgMgr:          n.PkgMgr,
		Platform:        n.Platform,
		SystemSetupCmds: n.SystemSetupCmds,
		DistName:        n.DistName,
		Ctr: dag.
			Pipeline(n.PipelineID + "oci-build").
			Container(dagger.ContainerOpts{
				Platform: n.Platform,
			}),
	}

	baseImageRefParts := strings.Split(n.BaseImageRef, ":")
	productionBuild = productionBuild.
		WithVersion(baseImageRefParts[0], baseImageRefParts[1], false)

	ctrDirArtifacts := []string{
		n.DistName,
	}

	ctrFileArtifacts := []string{
		"package.json",
	}

	if directoryContainerArtifacts != nil && len(directoryContainerArtifacts) > 0 {
		ctrDirArtifacts = append(ctrDirArtifacts, directoryContainerArtifacts...)
	}

	if fileContainerArtifacts != nil && len(fileContainerArtifacts) > 0 {
		ctrFileArtifacts = append(fileContainerArtifacts, fileContainerArtifacts...)
	}

	if n.NpmrcToken != nil {
		ctrFileArtifacts = append(ctrFileArtifacts, ".npmrc")
		productionBuild = productionBuild.WithNpmrcTokenEnv(n.NpmrcTokenName, n.NpmrcToken)
	}

	if n.NpmrcFile != nil {
		productionBuild = productionBuild.WithNpmrcTokenFile(n.NpmrcFile)
	}

	switch n.PkgMgr {
	case "npm":
		ctrFileArtifacts = append(ctrFileArtifacts, "package-lock.json")
	case "yarn":
		ctrFileArtifacts = append(ctrFileArtifacts, "yarn.lock")
	default:
		ctrFileArtifacts = append(ctrFileArtifacts, "package-lock.json")
	}

	if len(folderArtifacts) != len(folderArtifactNames) {
		return nil, fmt.Errorf("error mapping in folder artifacts")
	}

	if len(fileArtifacts) != len(fileArtifactNames) {
		return nil, fmt.Errorf("error in mapping in file artifacts")
	}

	productionBuild.Ctr = productionBuild.Ctr.WithWorkdir(workdir)

	productionBuild = productionBuild.
		SetupSystem(nil).
		Production().
		WithPackageManager(n.PkgMgr, true)

	for i, name := range folderArtifactNames {
		productionBuild.Ctr = productionBuild.
			Ctr.
			WithDirectory(workdir+"/"+name, folderArtifacts[i])
	}

	for _, name := range ctrDirArtifacts {
		path := workdir + "/" + name
		productionBuild.Ctr = productionBuild.
			Ctr.
			WithDirectory(path, n.Ctr.Directory(path))
	}

	for i, name := range fileArtifactNames {
		productionBuild.Ctr = productionBuild.
			Ctr.
			WithFile(workdir+"/"+name, fileArtifacts[i])
	}

	for _, name := range ctrFileArtifacts {
		path := workdir + "/" + name
		productionBuild.Ctr = productionBuild.
			Ctr.
			WithFile(path, n.Ctr.File(path))
	}

	productionBuild = productionBuild.Install()

	if isTtl {
		registries = []string{ttlRegistry}
	}

	for _, registry := range registries {
		eg.Go(func() error {
			ref := fmt.Sprintf("%s/%s:%s", registry, n.Name, n.Version)
			if isTtl {
				ref = fmt.Sprintf("%s/%s:%s", ttlRegistry, uuid.New().String(), ttl)
			}

			ref, err := productionBuild.Ctr.Publish(ctx, ref)
			result <- ref

			return err
		})
	}

	go func() {
		err = eg.Wait()
		close(result)
	}()

	for res := range result {
		fullyQualifiedImageNames = append(fullyQualifiedImageNames, res)
	}

	return fullyQualifiedImageNames, err
}
