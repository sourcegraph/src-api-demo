package main

import (
	"github.com/Khan/genqlient/graphql"
	"github.com/sourcegraph/src-api-demo/srcgql"
	"github.com/urfave/cli/v2"
)

func getClient(c *cli.Context) graphql.Client {
	endpoint := flagSrcEndpoint.Get(c)
	if endpoint == "" {
		panic("-endpoint is required")
	}
	token := flagSrcAccessToken.Get(c)
	if token == "" {
		panic("-access-token is required")
	}
	return srcgql.NewGraphQLClient(endpoint, token)
}
