package main

import (
	"bytes"
	"context"
	"text/template"
)

func (d *Drift) ReportToSlack(
	ctx context.Context,
	token *Secret,
	channelId string,
	// blabla
	// +optional
	// +default="#9512a6"
	color string,
) error {
	initTemplate, err := dag.CurrentModule().Source().File("templates/init_slack.tmpl").Contents(ctx)
	if err != nil {
		return err
	}

	templateRenderer, err := template.New("init_drift").Parse(initTemplate)
	if err != nil {
		return err
	}

	initReport := new(bytes.Buffer)
	err = templateRenderer.Execute(initReport, d)
	if err != nil {
		return err
	}

	threadId, err := dag.
		Notify().
		Slack().
		SendMessage(
			ctx,
			token,
			color,
			initReport.String(),
			channelId,
		)
	if err != nil {
		return err
	}

	for _, reportFormatted := range d.Reports {
		threadId, err = dag.
			Notify().
			Slack().
			SendMessage(
				ctx,
				token,
				"warning",
				reportFormatted,
				channelId,
				NotifySlackSendMessageOpts{ThreadID: threadId},
			)
		if err != nil {
			return err
		}
	}

	return nil
}
