package output

import (
	"text/template"

	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/urfave/cli/v2"
)

type withTemplate struct {
	template string
	data     any
}

func (w *withTemplate) Apply(opts *options) error {
	if w.template == "" {
		return nil
	}
	template, err := template.New("").Parse(w.template)
	if err != nil {
		return errors.Wrap(err, "parse template")
	}
	opts.template = template
	opts.data = w.data
	return nil
}

func WithTemplate(cli *cli.Context, d any) Option {
	return &withTemplate{
		template: cli.String("template"),
		data:     d,
	}
}
