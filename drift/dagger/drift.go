package main

import (
	"bytes"
	"context"
	"github.com/sourcegraph/conc/pool"
	"text/template"
	"time"
)

type report struct {
	StackName    string
	DriftContent string
}

func (d *Drift) Detection(ctx context.Context, src *Directory, stackRootPath string, maxParallelization int) (*Drift, error) {
	d.RootStacksPath = stackRootPath
	d.StartTime = time.Now().Format("2006-01-02 3:4:5 PM")
	stacks, err := src.Entries(ctx, DirectoryEntriesOpts{Path: stackRootPath})
	if err != nil {
		return nil, err
	}

	d.StackLen = len(stacks)
	reportChan := make(chan report, len(stacks))

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
				Terrabox().
				Terragrunt().
				WithSource("/terraform", src).
				DisableColor().
				Plan(stackRootPath+"/"+internalStackName, TerraboxTfPlanOpts{DetailedExitCode: true}).
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
