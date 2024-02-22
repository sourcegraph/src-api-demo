//go:build tools
// +build tools

// This package imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
package tools

import (
	// https://github.com/Khan/genqlient
	_ "github.com/Khan/genqlient"
)
