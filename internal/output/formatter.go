// Package output provides formatting utilities for displaying secret
// metadata and diagnostic information to the user without exposing
// sensitive values.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format controls how output is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes structured diagnostic output.
type Formatter struct {
	w      io.Writer
	format Format
}

// NewFormatter returns a Formatter writing to w in the given format.
// If w is nil, os.Stdout is used.
func NewFormatter(w io.Writer, format Format) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Formatter{w: w, format: format}
}

// PrintSecretKeys writes the resolved key names (not values) to the writer.
func (f *Formatter) PrintSecretKeys(path string, keys []string) error {
	switch f.format {
	case FormatJSON:
		quoted := make([]string, len(keys))
		for i, k := range keys {
			quoted[i] = fmt.Sprintf("%q", k)
		}
		_, err := fmt.Fprintf(f.w, `{"path":%q,"keys":[%s]}\n`, path, strings.Join(quoted, ","))
		return err
	default:
		tw := tabwriter.NewWriter(f.w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "PATH\tKEYS\n")
		fmt.Fprintf(tw, "%s\t%s\n", path, strings.Join(keys, ", "))
		return tw.Flush()
	}
}

// PrintEnvPreview writes a masked preview of the environment variables
// that would be injected, without revealing their values.
func (f *Formatter) PrintEnvPreview(entries map[string]string) error {
	switch f.format {
	case FormatJSON:
		pairs := make([]string, 0, len(entries))
		for k := range entries {
			pairs = append(pairs, fmt.Sprintf("%q:%q", k, "***"))
		}
		_, err := fmt.Fprintf(f.w, "{%s}\n", strings.Join(pairs, ","))
		return err
	default:
		tw := tabwriter.NewWriter(f.w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "KEY\tVALUE\n")
		for k := range entries {
			fmt.Fprintf(tw, "%s\t***\n", k)
		}
		return tw.Flush()
	}
}
