package main

import (
	"strconv"
)

// Run a plan on a specific stack
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

// Run an apply on a specific stack
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

// Format the code
func (t *Tf) Format(workDir string, check bool) *Tf {
	checkOptVal := strconv.FormatBool(check)
	if t.Bin == "terragrunt" {
		return t.WithContainer(
			t.run(
				workDir,
				[]string{"hclfmt", "--terragrunt-check=" + checkOptVal},
			).WithExec([]string{
				"terraform",
				"fmt",
				"-recursive",
				"-check=" + checkOptVal,
			}),
		)
	}

	return t.WithContainer(
		t.
			Ctr.
			WithWorkdir(workDir).
			WithExec([]string{
				t.Bin,
				"fmt",
				"-recursive",
				"-check=" + checkOptVal,
			}),
	)
}

// Return the output of a specific stack
func (t *Tf) Output(workDir string, isJson bool) *Tf {
	cmd := []string{"output"}

	if isJson {
		cmd = append(cmd, "-json")
	}

	return t.WithContainer(t.run(workDir, cmd))
}

// Execute the run-all command (only available for terragrunt)
func (t *Tf) RunAll(workDir string, cmd string) *Tf {
	return t.WithContainer(t.run(workDir, []string{"run-all", cmd}))
}

// expose the module catalog (only available for terragrunt)
func (t *Tf) Catalog() *Terminal {
	return t.Ctr.WithDefaultTerminalCmd([]string{t.Bin, "catalog"}).Terminal()
}
