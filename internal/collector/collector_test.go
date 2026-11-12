package collector_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/th3m45t3rm1nd/tracelite/internal/collector"
)

func TestIngestHandlerAcceptsSpanBatch(t *testing.T) {
	c := collector.New()

	body := bytes.NewBufferString(`[
		{"trace_id":"t1","span_id":"s1","name":"completion","kind":"llm",
		 "start_time_unix_nano":1000,"end_time_unix_nano":2000,
		 "tokens_in":120,"tokens_out":300,"cost_usd":0.0042},
		{"trace_id":"t1","span_id":"s2","parent_id":"s1","name":"search",
		 "kind":"tool","start_time_unix_nano":1500,"end_time_unix_nano":1800,
		 "error":"upstream timeout"}
	]`)
	req := httptest.NewRequest("POST", "/v1/traces", body)
	rec := httptest.NewRecorder()

	c.Handler()(rec, req)

	if rec.Code != 202 {
		t.Fatalf("expected 202 Accepted, got %d", rec.Code)
	}
	if got := c.Len(); got != 2 {
		t.Fatalf("expected 2 buffered spans, got %d", got)
	}
}

func TestFailedSpansFlagged(t *testing.T) {
	ok := collector.Span{Name: "completion", Kind: collector.KindLLM}
	bad := collector.Span{Name: "search", Kind: collector.KindTool, Error: "timeout"}
	if ok.Failed() || !bad.Failed() {
		t.Fatal("failure detection is inverted")
	}
}
