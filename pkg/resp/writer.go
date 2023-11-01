package resp

import (
	"io"
)

// RespWriter represents a RESP value writer.
type RespWriter struct {
	writer io.Writer
}

// NewRespWriter creates a new RESP writer.
func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{writer: w}
}

// WriteValue writes a RESP value to the writer.
func (rw *RespWriter) WriteValue(v Value) error {
	_, err := rw.writer.Write(v.Marshal())
	return err
}
