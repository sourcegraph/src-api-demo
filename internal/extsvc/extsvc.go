package extsvc

import (
	"github.com/sourcegraph/sourcegraph/lib/errors"

	"github.com/sourcegraph/src-api-demo/srcgql"
)

type Config interface {
	// UpdateToken updates the access token in the configuration.
	UpdateToken(token string) error
	// Raw returns the raw JSONC configuration string.
	Raw() string
}

func New(kind srcgql.ExternalServiceKind, config string) (Config, error) {
	switch kind {
	case srcgql.ExternalServiceKindGithub:
		return NewGitHubConfig(config)
	// TODO: Add support for other external service kinds.
	default:
		return nil, errors.Newf("unsupported external service kind: %q", kind)
	}
}
