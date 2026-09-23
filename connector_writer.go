package gotrail

import (
	"encoding/json"
	"io"
)

type WriterConnector struct {
	Writer io.Writer
}

func NewWriterConnector(w io.Writer) *WriterConnector {
	return &WriterConnector{
		w,
	}
}

func (c *WriterConnector) Send(spans []*Span) error {
	if len(spans) == 0 {
		return nil
	}

	otlpSpans, err := SpansToOTLPJSON(spans)
	if err != nil {
		return err
	}

	data, err := json.Marshal(otlpSpans)
	if err != nil {
		return err
	}

	if _, err := c.Writer.Write(data); err != nil {
		return err
	}

	return nil
}
