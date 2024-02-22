package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/sourcegraph/src-api-demo/internal/observability"
	"github.com/sourcegraph/src-api-demo/internal/output"
)

func main() {
	liblog := observability.InitLogs("src-api-demo", "dev")
	defer liblog.Sync()

	ctx := context.Background()
	if err := app.RunContext(ctx, os.Args); err != nil {
		_ = output.Render(output.FormatText, err)
		os.Exit(1)
	}
}

var app = &cli.App{
	Name: "src-api-demo",
	Commands: []*cli.Command{
		codeHostConn,
	},
}
