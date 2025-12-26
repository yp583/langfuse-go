package sampling

import (
	"context"
	"testing"

	"github.com/henomis/langfuse-go/model"
)

const spanIterations = 100

func TestSpanSampling0(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0)

	for i := 0; i < spanIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Span(&model.Span{
			Name:    "test-span",
			TraceID: trace.ID,
			Trace_:  trace,
			Input:   "test input",
			Output:  "test output",
		}, nil)
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0, counter.Get(), spanIterations*2, 0)
}

func TestSpanSampling25(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.25)

	for i := 0; i < spanIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Span(&model.Span{
			Name:    "test-span",
			TraceID: trace.ID,
			Trace_:  trace,
			Input:   "test input",
			Output:  "test output",
		}, nil)
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.25, counter.Get(), spanIterations*2, 0.15)
}

func TestSpanSampling50(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.50)

	for i := 0; i < spanIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Span(&model.Span{
			Name:    "test-span",
			TraceID: trace.ID,
			Trace_:  trace,
			Input:   "test input",
			Output:  "test output",
		}, nil)
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.50, counter.Get(), spanIterations*2, 0.15)
}

func TestSpanSampling75(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.75)

	for i := 0; i < spanIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Span(&model.Span{
			Name:    "test-span",
			TraceID: trace.ID,
			Trace_:  trace,
			Input:   "test input",
			Output:  "test output",
		}, nil)
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.75, counter.Get(), spanIterations*2, 0.15)
}

func TestSpanSampling100(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 1.0)

	for i := 0; i < spanIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Span(&model.Span{
			Name:    "test-span",
			TraceID: trace.ID,
			Trace_:  trace,
			Input:   "test input",
			Output:  "test output",
		}, nil)
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 1.0, counter.Get(), spanIterations*2, 0)
}
