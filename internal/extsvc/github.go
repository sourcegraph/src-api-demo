package extsvc

import (
	"github.com/sourcegraph/sourcegraph/lib/errors"

	"github.com/sourcegraph/src-api-demo/internal/jsonc"
	"github.com/sourcegraph/src-api-demo/srcconf"
)

type GitHubConfig struct {
	s      string
	config srcconf.GitHubConnection
}

func (g *GitHubConfig) UpdateToken(token string) error {
	updated, err := jsonc.Edit(g.s, token, "token")
	if err != nil {
		return err
	}
	g.s = updated
	return nil
}

func (g *GitHubConfig) Raw() string {
	return g.s
}

func NewGitHubConfig(s string) (*GitHubConfig, error) {
	var v srcconf.GitHubConnection
	if err := jsonc.Unmarshal(s, &v); err != nil {
		return nil, errors.Wrap(err, "unmarshal GitHubConnection")

	}
	return &GitHubConfig{
		s:      s,
		config: v,
	}, nil
}
