package gotrail

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HTTPOTLPConnector struct {
	// HTTP Client
	Client *http.Client
	// Endpoint for publishing traces
	URL string

	// Modify request before sending it
	PrepareRequest func(*http.Request) error
}

// Sends spans to specified endpoint that supports OTLP over HTTP
// with JSON encoded body. Sets Content-Type to application/json
func (c *HTTPOTLPConnector) Send(spans []*Span) error {
	if len(spans) == 0 {
		return nil
	}

	data, err := SpansToOTLPJSON(spans)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	err = c.PrepareRequest(req)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// respBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func NewHTTPOTLPConnector(client *http.Client, url string, prepareReq func(*http.Request) error) *HTTPOTLPConnector {
	return &HTTPOTLPConnector{
		client,
		url,
		prepareReq,
	}
}
