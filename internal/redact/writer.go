package redact

import "io"

// Writer wraps an io.Writer and redacts secret values from all writes.
type Writer struct {
	underlying io.Writer
	redactor   *Redactor
}

// NewWriter returns a Writer that scrubs secrets before forwarding to w.
func NewWriter(w io.Writer, r *Redactor) *Writer {
	return &Writer{underlying: w, redactor: r}
}

// Write redacts p before writing to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	clean := w.redactor.Redact(string(p))
	_, err := w.underlying.Write([]byte(clean))
	if err != nil {
		return 0, err
	}
	// Report original length so callers don't get short-write errors.
	return len(p), nil
}
