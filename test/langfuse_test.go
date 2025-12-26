package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	langfuse "github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/internal/pkg/api"
	"github.com/henomis/langfuse-go/model"
)

func TestIngestionWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req api.Ingestion
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Logf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Logf("=== Received %d events ===", len(req.Batch))
		for i, event := range req.Batch {
			t.Logf("Event %d:", i+1)
			t.Logf("  Event ID: %s", event.ID)
			t.Logf("  Type: %s", event.Type)
			t.Logf("  Timestamp: %s", event.Timestamp)

			bodyBytes, _ := json.MarshalIndent(event.Body, "  ", "  ")
			t.Logf("  Body: %s", string(bodyBytes))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"successes":[],"errors":[]}`))
	}))
	defer server.Close()

	ctx := context.Background()
	client := api.NewWithHost(server.URL)
	lf := langfuse.NewWithClient(ctx, client)

	trace, err := lf.Trace(&model.Trace{Name: "test-trace"})
	if err != nil {
		t.Fatalf("Failed to create trace: %v", err)
	}
	t.Logf("Created trace with ID: %s, ShouldTrace: %v", trace.ID, trace.ShouldTrace)

	gen, err := lf.Generation(&model.Generation{
		Name:    "test-generation",
		TraceID: trace.ID,
		Trace_:  trace,
		Input:   "What is 2+2?",
		Output:  "4",
	}, nil)
	if err != nil {
		t.Fatalf("Failed to create generation: %v", err)
	}
	t.Logf("Created generation with ID: %s", gen.ID)

	lf.Flush(ctx)
	t.Logf("Flushed events to mock server")
}

func TestSamplingRate(t *testing.T) {
	eventCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req api.Ingestion
		json.NewDecoder(r.Body).Decode(&req)
		eventCount += len(req.Batch)

		for _, event := range req.Batch {
			t.Logf("Received event: Type=%s, ID=%s", event.Type, event.ID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"successes":[],"errors":[]}`))
	}))
	defer server.Close()

	ctx := context.Background()
	client := api.NewWithHost(server.URL)
	lf := langfuse.NewWithClient(ctx, client).WithSamplingRate(0.5)

	tracedCount := 0
	for i := 0; i < 100; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "sample-trace"})
		if trace.ShouldTrace {
			tracedCount++
		}
	}

	lf.Flush(ctx)

	t.Logf("Out of 100 traces: %d were sampled (ShouldTrace=true)", tracedCount)
	t.Logf("Total events sent to server: %d", eventCount)

	if tracedCount < 30 || tracedCount > 70 {
		t.Logf("Warning: Sampling rate seems off (expected ~50, got %d)", tracedCount)
	}
}
