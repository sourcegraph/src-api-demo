package main

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/urfave/cli/v2"

	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/src-api-demo/srcgql"
)

var (
	defaultCodeIntelligencePolicies = []srcgql.CodeIntelligenceConfigurationPolicyFields{
		{
			Name:                   "[Sourcegraph Managed] michael test policy",
			Type:                   "GIT_COMMIT",
			Pattern:                "HEAD",
			RetentionEnabled:       false,
			IndexingEnabled:        true,
			IndexCommitMaxAgeHours: func(i int) *int { return &i }(8064),
		},
	}
)

var (
	codeIntel = &cli.Command{
		Name:  "code-intel",
		Flags: mergeFlagSets(globalFlags),
		Action: func(c *cli.Context) error {
			ctx := c.Context
			client := getClient(c)
			err := createCodeIntelligencePolicy(ctx, client, defaultCodeIntelligencePolicies[0])
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func createCodeIntelligencePolicy(ctx context.Context, client graphql.Client, policy srcgql.CodeIntelligenceConfigurationPolicyFields) error {
	_, err := srcgql.CreateCodeIntelligenceConfigurationPolicy(ctx, client,
		policy.Name,
		policy.Type,
		policy.Pattern,
		&policy.RetentionEnabled,
		policy.RetentionDurationHours,
		&policy.RetainIntermediateCommits,
		&policy.IndexingEnabled,
		policy.IndexCommitMaxAgeHours,
		&policy.IndexIntermediateCommits,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create code intelligence configuration policy")
	}

	return nil
}
