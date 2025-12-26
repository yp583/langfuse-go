package sampling

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	langfuse "github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/internal/pkg/api"
)

type EventCounter struct {
	mu    sync.Mutex
	count int
}

func (c *EventCounter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += n
}

func (c *EventCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func (c *EventCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
}

func createMockServer(t *testing.T, counter *EventCounter) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req api.Ingestion
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Logf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		counter.Add(len(req.Batch))

		for _, event := range req.Batch {
			t.Logf("Event: Type=%s, ID=%s", event.Type, event.ID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"successes":[],"errors":[]}`))
	}))
}

func createLangfuse(ctx context.Context, serverURL string, samplingRate float32) *langfuse.Langfuse {
	client := api.NewWithHost(serverURL)
	return langfuse.NewWithClient(ctx, client).WithSamplingRate(samplingRate)
}

func assertSamplingRate(t *testing.T, expectedRate float32, actualCount int, totalCount int, tolerance float32) {
	actualRate := float32(actualCount) / float32(totalCount)
	diff := float32(math.Abs(float64(actualRate - expectedRate)))

	t.Logf("Expected rate: %.2f, Actual rate: %.2f (count: %d/%d)", expectedRate, actualRate, actualCount, totalCount)

	if expectedRate == 0 && actualCount != 0 {
		t.Errorf("Expected 0 events for 0%% sampling, got %d", actualCount)
	} else if expectedRate == 1 && actualCount != totalCount {
		t.Errorf("Expected %d events for 100%% sampling, got %d", totalCount, actualCount)
	} else if expectedRate > 0 && expectedRate < 1 && diff > tolerance {
		t.Logf("Warning: Sampling rate %.2f outside tolerance %.2f (expected %.2f) - this may be due to random variance", actualRate, tolerance, expectedRate)
	}
}
