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
	// Define if the plan is saved in a file
	// +optional
	savePlan bool,
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

	if savePlan {
		cmd = append(cmd, "-out="+workDir+"/tfplan")
	}

	ctr := t.run(workDir, cmd)

	if savePlan {
		t.TfPlan = ctr.File(workDir + "/tfplan")
	}

	return t.WithContainer(ctr)
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

	if t.TfPlan != nil {
		t.WithContainer(t.Ctr.WithFile(workDir+"/tfplan", t.TfPlan))
		cmd = append(cmd, "tfplan")
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

// Run a show on a specific state or plan file
func (t *Tf) Show(
	// Define if the output is in machine-readableform
	// +optional
	ojson bool,
	// Define a path to a plan file or state
	// +optional
	path string,
) *Tf {
	cmd := []string{"show"}

	if ojson {
		cmd = append(cmd, "-json")
	}

	if t.NoColor {
		cmd = append(cmd, "-no-color")
	}

	if t.TfPlan != nil {
		t.WithContainer(t.Ctr.WithFile("/tfplan", t.TfPlan))
		cmd = append(cmd, "/tfplan")
	}

	return t.WithContainer(t.run("/", cmd))
}

// Execute the run-all command (only available for terragrunt)
func (t *Tf) RunAll(workDir string, cmd string) *Tf {
	return t.WithContainer(t.run(workDir, []string{"run-all", cmd}))
}
