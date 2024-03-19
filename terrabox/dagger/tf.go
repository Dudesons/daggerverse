package main

import (
	"context"
	"dagger/terrabox/internal/dagger"
	"fmt"
	"strconv"
	"time"
)

type Tf struct {
	// +private
	Ctr *Container
	// +private
	Bin string
	// +private
	RootPath string
	// +private
	NoColor bool
}

func newTf(
	image string,
	version string,
	binary string,
) *Tf {
	return &Tf{
		Bin: binary,
		Ctr: dag.
			Container().
			From(image+":"+version).
			WithMountedCache("/root/.terraform.d/plugin-cache", dag.CacheVolume("terraform-plugins")),
	}
}

// Mount the source code at the given path
func (t *Tf) WithSource(path string, src *Directory) *Tf {
	t.RootPath = path
	t.Ctr = t.Ctr.
		WithDirectory(path, src).
		WithWorkdir(path)

	return t
}

// Use a new container
func (t *Tf) WithContainer(ctr *Container) *Tf {
	t.Ctr = ctr

	return t
}

// Convert a dotfile format to secret environment variables in the container (could be use to configure providers)
func (t *Tf) WithSecretDotEnv(dotEnv *Secret) *Tf {
	return t.WithContainer(dag.Utils().WithDotEnvSecret(t.Ctr, dotEnv))
}

// Indicate to disable the the color in the output
func (t *Tf) DisableColor() *Tf {
	t.NoColor = true
	return t.WithContainer(t.Ctr.WithEnvVariable("TERRAGRUNT_NO_COLOR", "true"))
}

// Expose the container
func (t *Tf) Container() *Container {
	return t.Ctr
}

func (t *Tf) run(workDir string, command []string) *dagger.Container {
	return t.Ctr.
		WithWorkdir(workDir).
		WithExec(append([]string{t.Bin}, command...))
}

// Execute the call chain
func (t *Tf) Do(ctx context.Context) (string, error) {
	return t.Ctr.Stdout(ctx)
}

// Return the source directory
func (t *Tf) Directory() *Directory {
	return t.Ctr.Directory(t.RootPath)
}

// Open a shell
func (t *Tf) Shell() *Terminal {
	return t.Ctr.WithDefaultTerminalCmd(nil).Terminal()
}

func (t *Tf) WithCacheBurster(
	// Define if the cache burster level is done per day (daily), per hour (hour), per minute (minute), per second (default)
	// +optional
	cacheBursterLevel string,
) *Tf {
	utcNow := time.Now().UTC()
	cacheBursterKey := fmt.Sprintf("%d%d%d", utcNow.Year(), utcNow.Month(), utcNow.Day())

	switch cacheBursterLevel {
	case "daily":
	case "hour":
		cacheBursterKey += strconv.Itoa(utcNow.Hour())
	case "minute":
		cacheBursterKey += fmt.Sprintf("%d%d", utcNow.Hour(), utcNow.Minute())
	default:
		cacheBursterKey += fmt.Sprintf("%d%d%d", utcNow.Hour(), utcNow.Minute(), utcNow.Second())
	}

	return t.WithContainer(t.Ctr.WithEnvVariable("CACHE_BURSTER", cacheBursterKey))
}
