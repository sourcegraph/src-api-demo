package main

import (
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/urfave/cli/v2"

	"github.com/sourcegraph/src-api-demo/internal/extsvc"
	"github.com/sourcegraph/src-api-demo/internal/output"
	"github.com/sourcegraph/src-api-demo/srcgql"
)

var (
	codeHostConn = &cli.Command{
		Name:    "code-host-conn",
		Aliases: []string{"extsvc"},
		Usage:   "Manage code host connections",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List code host connections",
				Flags: mergeFlagSets(globalFlags),
				Action: func(c *cli.Context) error {
					ctx := c.Context

					client := getClient(c)
					resp, err := srcgql.GetExternalServices(ctx, client)
					if err != nil {
						return err
					}
					return output.Render(output.FormatJSON, resp.GetExternalServices())
				},
			},
			{
				Name:      "update-token",
				ArgsUsage: "<token>",
				Usage:     "Update the access token for a code host connection",
				Flags: mergeFlagSets(globalFlags, []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "The ID of the code host connection to update",
						Required: true,
					},
				}),
				Action: func(c *cli.Context) error {
					ctx := c.Context

					if c.Args().Len() != 1 {
						return errors.New("token is required")
					}
					newToken := c.Args().First()

					client := getClient(c)
					resp, err := srcgql.GetExternalServices(ctx, client)
					if err != nil {
						return err
					}

					// found the extsvc
					var got srcgql.GetExternalServicesExternalServicesExternalServiceConnectionNodesExternalService
					for _, svc := range resp.GetExternalServices().Nodes {
						if svc.Id == c.String("id") {
							got = svc
							break
						}
					}

					// parse and update token
					config, err := extsvc.New(got.Kind, got.Config)
					if err != nil {
						return errors.Wrapf(err, "parse extsvc config of kind %q", got.Kind)
					}
					config.UpdateToken(newToken)

					// update
					updated, err := srcgql.UpdateExternalService(ctx, client, got.Id, got.DisplayName, string(config.Raw()))
					if err != nil {
						return errors.Wrapf(err, "update extsvc %q", got.Id)
					}

					return output.Render(output.FormatJSON, updated)
				},
			},
		},
	}
)
