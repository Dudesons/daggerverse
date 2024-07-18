package main

import (
	"bytes"
	"context"
	"dagger/drift/internal/dagger"
	"text/template"
)

// Send the report formated to slack
func (d *Drift) ReportToSlack(
	ctx context.Context,
	// The slack token to use
	token *dagger.Secret,
	// The channel  id where the report will be posted
	channelId string,
	// Define the sidebar color of the message in slack
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
				dagger.NotifySlackSendMessageOpts{ThreadID: threadId},
			)
		if err != nil {
			return err
		}
	}

	return nil
}
