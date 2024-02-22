package main

import "github.com/urfave/cli/v2"

var (
	globalFlags = []cli.Flag{flagSrcEndpoint, flagSrcAccessToken}

	flagSrcEndpoint = &cli.StringFlag{
		Name:     "endpoint",
		EnvVars:  []string{"SRC_ENDPOINT"},
		Usage:    "Sourcegraph instance GraphQL endpoint, e.g., https://sourcegraph.acme.com/.api/graphql, https://acme.sourcegraphcloud.com/.api/graphql",
		Required: true,
	}
	flagSrcAccessToken = &cli.StringFlag{
		Name:     "access-token",
		Aliases:  []string{"token"},
		EnvVars:  []string{"SRC_ACCESS_TOKEN", "SRC_TOKEN"},
		Usage:    "Sourcegraph access token",
		Required: true,
	}
)

// mergeFlagSets combines multiple sets of flags into a single slice of flags.
func mergeFlagSets(flagSets ...[]cli.Flag) (flags []cli.Flag) {
	for _, fs := range flagSets {
		flags = append(flags, fs...)
	}
	return flags
}
