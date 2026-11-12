// Package tracelite is the client SDK: wrap LLM and tool calls into spans and
// report them, batched, to a Tracelite collector with W3C trace-context
// propagation.
package tracelite

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Client reports spans to a Tracelite collector endpoint.
type Client struct {
	endpoint string
	http     *http.Client
}

// NewClient points a Client at a collector (e.g. http://localhost:8080).
func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{Timeout: 5 * time.Second},
	}
}

// Report posts a batch of spans to the collector. Batching is the caller's
// policy for now (TODO Week 1: background flusher with size/time triggers).
func (c *Client) Report(ctx context.Context, spans []Span) error {
	body, err := json.Marshal(spans)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/v1/traces", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return ErrUnexpectedStatus{Code: resp.StatusCode}
	}
	return nil
}

// ErrUnexpectedStatus is returned when the collector rejects a batch.
type ErrUnexpectedStatus struct{ Code int }

func (e ErrUnexpectedStatus) Error() string {
	return "collector returned non-202 status: " + http.StatusText(e.Code)
}
