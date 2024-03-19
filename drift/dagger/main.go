// A Drift detection module around terraform/terragrunt which allow to send a report

package main

func New(
	// Define where the code is mounted, this could impact for absolute module path
	// +optional
	// +default="/terraform"
	mountPoint string,
) *Drift {
	return &Drift{
		MountPoint: mountPoint,
	}
}

type Drift struct {
	// +private
	Reports []string
	// +private
	StartTime string
	// +private
	Endtime string
	// +private
	StackLen int
	// +private
	DriftLen int
	// +private
	RootStacksPath string
	// +private
	MountPoint string
}
