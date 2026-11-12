// Package tracelite — wire-format types shared with the collector.
package tracelite

// Span mirrors the collector's wire format (see internal/collector).
// TODO Week 1: unify with a shared model package and add span-start helpers
// (StartSpan/EndSpan with automatic timing) and traceparent generation.
type Span struct {
	TraceID   string            `json:"trace_id"`
	SpanID    string            `json:"span_id"`
	ParentID  string            `json:"parent_id,omitempty"`
	Name      string            `json:"name"`
	Kind      string            `json:"kind"`
	Model     string            `json:"model,omitempty"`
	StartTime int64             `json:"start_time_unix_nano"`
	EndTime   int64             `json:"end_time_unix_nano"`
	TokensIn  int               `json:"tokens_in,omitempty"`
	TokensOut int               `json:"tokens_out,omitempty"`
	CostUSD   float64           `json:"cost_usd,omitempty"`
	Error     string            `json:"error,omitempty"`
	Attrs     map[string]string `json:"attrs,omitempty"`
}
