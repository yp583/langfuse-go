package sampling

import (
	"context"
	"testing"

	"github.com/henomis/langfuse-go/model"
)

const traceIterations = 100

func TestTraceSampling0(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0)

	for i := 0; i < traceIterations; i++ {
		lf.Trace(&model.Trace{Name: "test-trace"})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0, counter.Get(), traceIterations, 0)
}

func TestTraceSampling25(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.25)

	for i := 0; i < traceIterations; i++ {
		lf.Trace(&model.Trace{Name: "test-trace"})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.25, counter.Get(), traceIterations, 0.15)
}

func TestTraceSampling50(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.50)

	for i := 0; i < traceIterations; i++ {
		lf.Trace(&model.Trace{Name: "test-trace"})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.50, counter.Get(), traceIterations, 0.15)
}

func TestTraceSampling75(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.75)

	for i := 0; i < traceIterations; i++ {
		lf.Trace(&model.Trace{Name: "test-trace"})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.75, counter.Get(), traceIterations, 0.15)
}

func TestTraceSampling100(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 1.0)

	for i := 0; i < traceIterations; i++ {
		lf.Trace(&model.Trace{Name: "test-trace"})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 1.0, counter.Get(), traceIterations, 0)
}
