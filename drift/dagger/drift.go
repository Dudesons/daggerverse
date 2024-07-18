package main

import (
	"bytes"
	"context"
	"dagger/drift/internal/dagger"
	"github.com/sourcegraph/conc/pool"
	"text/template"
	"time"
)

type report struct {
	StackName    string
	DriftContent string
}

// Trigger the drift detection
func (d *Drift) Detection(
	ctx context.Context,
	// All the terraform/terragrunt code necessary in order to be able to run plan
	src *dagger.Directory,
	// The root path where stack are living
	stackRootPath string,
	// The number of execution in parallel we want to have, 0 mean no limit
	maxParallelization int,
	// Define if the cache burster level is done per day (daily), per hour (hour), per minute (minute), per second (default)
	// +optional
	// +default="hour"
	cacheBursterLevel string,
) (*Drift, error) {
	d.RootStacksPath = stackRootPath
	d.StartTime = time.Now().Format("2006-01-02 3:4:5 PM")
	stacks, err := src.Entries(ctx, dagger.DirectoryEntriesOpts{Path: stackRootPath})
	if err != nil {
		return nil, err
	}

	d.StackLen = len(stacks)
	reportChan := make(chan report, d.StackLen)

	runPool := pool.New()
	if maxParallelization != 0 {
		runPool = runPool.WithMaxGoroutines(maxParallelization)
	}

	driftTemplate, err := dag.CurrentModule().Source().File("templates/drift_detected.tmpl").Contents(ctx)
	if err != nil {
		return nil, err
	}

	templateRenderer, err := template.New("drift").Parse(driftTemplate)
	if err != nil {
		return nil, err
	}

	for _, stack := range stacks {
		runPool.Go(func() {
			internalStackName := stack
			_, err := dag.
				Infrabox().
				Terragrunt().
				WithSource(d.MountPoint, src).
				DisableColor().
				WithCacheBurster(dagger.InfraboxTfWithCacheBursterOpts{CacheBursterLevel: cacheBursterLevel}).
				Plan(d.MountPoint+"/"+d.RootStacksPath+"/"+internalStackName, dagger.InfraboxTfPlanOpts{DetailedExitCode: true}).
				Do(ctx)
			if err != nil {
				reportChan <- report{StackName: internalStackName, DriftContent: err.Error()}
			}
		})

	}

	go func() {
		runPool.Wait()
		close(reportChan)
	}()

	for res := range reportChan {
		buf := new(bytes.Buffer)
		err = templateRenderer.Execute(buf, res)
		if err != nil {
			return nil, err
		}

		d.Reports = append(d.Reports, buf.String())
	}

	d.Endtime = time.Now().Format("2006-01-02 3:4:5 PM")
	d.DriftLen = len(d.Reports)

	return d, nil
}
