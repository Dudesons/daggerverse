package main

import (
	"strconv"
)

func (t *Tf) Plan(
	// Define the path where to execute the command
	workDir string,
	// Define if we are executing the plan in destroy mode or not
	// +optional
	destroyMode bool,
	// Define if the exit code is in detailed mode or not (0 - Succeeded, diff is empty (no changes) | 1 - Errored | 2 - Succeeded, there is a diff)
	// +optional
	detailedExitCode bool,
) *Tf {
	cmd := []string{"plan", "-input=false"}

	if destroyMode {
		cmd = append(cmd, "-destroy")
	}

	if detailedExitCode {
		cmd = append(cmd, "-detailed-exitcode")
	}

	if t.NoColor {
		cmd = append(cmd, "-no-color")
	}

	return t.WithContainer(t.run(workDir, cmd))
}

func (t *Tf) Apply(
	// Define the path where to execute the command
	workDir string,
	// Define if we are executing the plan in destroy mode or not
	// +optional
	destroyMode bool) *Tf {
	cmd := []string{"apply", "-input=false", "-auto-approve"}

	if destroyMode {
		cmd = append(cmd, "-destroy")
	}

	if t.NoColor {
		cmd = append(cmd, "-no-color")
	}

	return t.WithContainer(t.run(workDir, cmd))
}

func (t *Tf) Format(workDir string, check bool) *Tf {
	checkOptVal := strconv.FormatBool(check)
	if t.Bin == "terragrunt" {
		t = t.WithContainer(
			t.run(
				workDir,
				[]string{"hclfmt", "--terragrunt-check=" + checkOptVal},
			),
		)
	}

	// TODO(Find a better way to handle that in particular if it's opentofu)
	return t.WithContainer(
		t.
			Ctr.
			WithWorkdir(workDir).
			WithExec([]string{
				"terraform",
				"fmt",
				"-recursive",
				"-check=" + checkOptVal,
			}),
	)
}

func (t *Tf) Output(workDir string, isJson bool) *Tf {
	cmd := []string{"output"}

	if isJson {
		cmd = append(cmd, "-json")
	}

	return t.WithContainer(t.run(workDir, cmd))
}

func (t *Tf) RunAll(workDir string, cmd string) *Tf {
	return t.WithContainer(t.run(workDir, []string{"run-all", cmd}))
}

func (t *Tf) Catalog() *Terminal {
	return t.Ctr.WithDefaultTerminalCmd([]string{t.Bin, "catalog"}).Terminal()
}
