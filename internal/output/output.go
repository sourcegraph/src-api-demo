package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/sourcegraph/sourcegraph/lib/errors"
	liboutput "github.com/sourcegraph/sourcegraph/lib/output"
	"sigs.k8s.io/yaml"
)

type Format string

const (
	// FormatJSON renders plain, unformatted JSON.
	FormatJSON Format = "json"
	// FormatYAML renders plain, unformatted YAML.
	FormatYAML Format = "yaml"
	// FormatPretty renders pretty, human-readable content, by default indented,
	// color-coded JSON.
	FormatPretty Format = "pretty"
	// FormatText renders plain-text content, by default the '%v' directive formatting.
	FormatText Format = "text"
	// FormatNone renders nothing.
	FormatNone Format = "none"

	FormatGoTemplate Format = "go-template"
)

var (
	// Formats is a slice of all supported default formats.
	Formats = []Format{
		FormatJSON,
		FormatYAML,
		FormatPretty,
		FormatText,
		FormatNone,
	}
)

// ErrFormatUnimplemented can be used by Renderer implementations for
// unhandled format cases.
var ErrFormatUnimplemented = errors.New("custom output.Renderer does not implement this format")

// Renderer is an interface that can be implemented by types that want to
// support custom render formats.
//
// The implementation should return ErrFormatUnimplemented as the fallback
// case - this tells the top-level output.Render implementation to use its
// default behaviour if the type does not implement it.
type Renderer interface{ Render(io.Writer, Format) error }

// Render outputs v in the requested format. For custom formatting, v may
// implement output.Renderer to the default override formatting behaviour.
func Render(format Format, v any, opts ...Option) error {
	var options options
	for _, opt := range opts {
		if err := opt.Apply(&options); err != nil {
			return err
		}
	}

	if r, ok := (v).(Renderer); ok {
		err := r.Render(os.Stdout, format)
		// If no issue occurred, we are done.
		if err == nil {
			return nil
		}
		// If unimplemented, continue with default behaviour - otherwise
		// the custom render implementation has failed.
		if !errors.Is(err, ErrFormatUnimplemented) {
			return err
		}
	}

	switch format {
	case FormatJSON:
		return json.NewEncoder(os.Stdout).Encode(v)
	case FormatYAML:
		b, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		_, err = fmt.Print(string(b))
		return err
	case FormatPretty:
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		return liboutput.NewOutput(os.Stdout, liboutput.OutputOpts{}).
			WriteCode("json", string(data))
	case FormatText:
		data := strings.TrimSpace(fmt.Sprintf("%v", v))
		_, err := fmt.Println(data)
		return err
	case FormatNone:
		return nil
	case FormatGoTemplate:
		if options.template == nil {
			return errors.New("--template is required for go-template format")
		}
		var b bytes.Buffer
		if err := options.template.Execute(&b, options.data); err != nil {
			return errors.Wrap(err, "execute go-template")
		}
		_, err := fmt.Println(b.String())
		return err
	default:
		return fmt.Errorf("unknown format %q", format)
	}
}

type options struct {
	// for WithTemplate
	template *template.Template
	data     any
}

type Option interface {
	Apply(*options) error
}
