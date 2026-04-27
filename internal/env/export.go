package env

import (
	"fmt"
	"sort"
	"strings"
)

// ExportFormat controls how exported variables are rendered.
type ExportFormat int

const (
	FormatShell ExportFormat = iota // export KEY="VALUE"
	FormatDotenv                    // KEY=VALUE
	FormatJSON                      // {"KEY": "VALUE"}
)

// Exporter serialises an env map to a target format.
type Exporter struct {
	format ExportFormat
	quote  bool
}

// NewExporter creates an Exporter for the given format.
// quote controls whether values are double-quoted (ignored for JSON).
func NewExporter(format ExportFormat, quote bool) *Exporter {
	return &Exporter{format: format, quote: quote}
}

// Export converts the provided map to the configured string format.
func (e *Exporter) Export(env map[string]string) (string, error) {
	if env == nil {
		return "", fmt.Errorf("env map must not be nil")
	}
	switch e.format {
	case FormatShell:
		return e.toShell(env), nil
	case FormatDotenv:
		return e.toDotenv(env), nil
	case FormatJSON:
		return e.toJSON(env), nil
	default:
		return "", fmt.Errorf("unknown export format: %d", e.format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (e *Exporter) toShell(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		v := env[k]
		if e.quote {
			fmt.Fprintf(&sb, "export %s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "export %s=%s\n", k, v)
		}
	}
	return sb.String()
}

func (e *Exporter) toDotenv(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		v := env[k]
		if e.quote {
			fmt.Fprintf(&sb, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}

func (e *Exporter) toJSON(env map[string]string) string {
	var sb strings.Builder
	keys := sortedKeys(env)
	sb.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		fmt.Fprintf(&sb, "%q:%q", k, env[k])
	}
	sb.WriteString("}")
	return sb.String()
}
