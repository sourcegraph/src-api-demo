package extsvc

import (
	"github.com/sourcegraph/sourcegraph/lib/errors"

	"github.com/sourcegraph/src-api-demo/internal/jsonc"
	"github.com/sourcegraph/src-api-demo/srcconf"
)

type GitLabConfig struct {
	s      string
	config srcconf.GitLabConnection
}

func (g *GitLabConfig) UpdateToken(token string) error {
	updated, err := jsonc.Edit(g.s, token, "token")
	if err != nil {
		return err
	}
	g.s = updated
	return nil
}

func (g *GitLabConfig) Raw() string {
	return g.s
}

func NewGitLabConfig(s string) (*GitLabConfig, error) {
	var v srcconf.GitLabConnection
	if err := jsonc.Unmarshal(s, &v); err != nil {
		return nil, errors.Wrap(err, "unmarshal GitLabConnection")

	}
	return &GitLabConfig{
		s:      s,
		config: v,
	}, nil
}
