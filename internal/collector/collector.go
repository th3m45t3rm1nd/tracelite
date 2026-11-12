// Package collector ingests trace spans from instrumented applications and
// buffers them for batched writes to storage (PostgreSQL for trace trees,
// Meridian for numeric series).
package collector

import (
	"encoding/json"
	"net/http"
	"sync"
)

// SpanKind classifies what a span measures.
type SpanKind string

const (
	KindLLM    SpanKind = "llm"
	KindTool   SpanKind = "tool"
	KindAgent  SpanKind = "agent"
	KindCustom SpanKind = "custom"
)

// Span is a single unit of work in a trace (LLM call, tool call, agent step).
// ParentID links spans into a tree; an empty ParentID marks a root span.
type Span struct {
	TraceID   string            `json:"trace_id"`
	SpanID    string            `json:"span_id"`
	ParentID  string            `json:"parent_id,omitempty"`
	Name      string            `json:"name"`
	Kind      SpanKind          `json:"kind"`
	Model     string            `json:"model,omitempty"`
	StartTime int64             `json:"start_time_unix_nano"`
	EndTime   int64             `json:"end_time_unix_nano"`
	TokensIn  int               `json:"tokens_in,omitempty"`
	TokensOut int               `json:"tokens_out,omitempty"`
	CostUSD   float64           `json:"cost_usd,omitempty"`
	Error     string            `json:"error,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
}

// Failed reports whether the span represents an error — the primary tail
// sampling signal (failed traces are always retained).
func (s Span) Failed() bool { return s.Error != "" }

// Collector buffers spans and flushes them in batches to storage backends.
// Flush/sampling pipelines land in Week 1–2 (see README roadmap).
type Collector struct {
	mu    sync.Mutex
	spans []Span
}

// New returns an empty Collector.
func New() *Collector { return &Collector{} }

// Ingest buffers a single span.
func (c *Collector) Ingest(s Span) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spans = append(c.spans, s)
}

// Len reports how many spans are currently buffered.
func (c *Collector) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.spans)
}

// Handler returns the HTTP ingest endpoint (POST /v1/traces).
func (c *Collector) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var spans []Span
		if err := json.NewDecoder(r.Body).Decode(&spans); err != nil {
			http.Error(w, "bad payload: "+err.Error(), http.StatusBadRequest)
			return
		}
		for _, s := range spans {
			c.Ingest(s)
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
