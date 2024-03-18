package main

import (
	"context"
	"dagger/terrabox/internal/dagger"
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

func (t *Tf) WithSource(path string, src *Directory) *Tf {
	t.RootPath = path
	t.Ctr = t.Ctr.
		WithDirectory(path, src).
		WithWorkdir(path)

	return t
}

func (t *Tf) WithContainer(ctr *Container) *Tf {
	t.Ctr = ctr

	return t
}

func (t *Tf) WithSecretDotEnv(dotEnv *Secret) *Tf {
	return t.WithContainer(dag.Utils().WithDotEnvSecret(t.Ctr, dotEnv))
}

func (t *Tf) DisableColor() *Tf {
	t.NoColor = true
	return t.WithContainer(t.Ctr.WithEnvVariable("TERRAGRUNT_NO_COLOR", "true"))
}

func (t *Tf) Container() *Container {
	return t.Ctr
}

func (t *Tf) run(workDir string, command []string) *dagger.Container {
	return t.Ctr.
		WithWorkdir(workDir).
		WithExec(append([]string{t.Bin}, command...))
}

func (t *Tf) Do(ctx context.Context) (string, error) {
	return t.Ctr.Stdout(ctx)
}

func (t *Tf) Directory() *Directory {
	return t.Ctr.Directory(t.RootPath)
}

func (t *Tf) Shell() *Terminal {
	return t.Ctr.WithDefaultTerminalCmd(nil).Terminal()
}
