package sampling

import (
	"context"
	"testing"

	"github.com/henomis/langfuse-go/model"
)

const scoreIterations = 100

func TestScoreSampling0(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0)

	for i := 0; i < scoreIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Score(&model.Score{
			Name:    "test-score",
			TraceID: trace.ID,
			Trace_:  trace,
			Value:   0.95,
			Comment: "test comment",
		})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0, counter.Get(), scoreIterations*2, 0)
}

func TestScoreSampling25(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.25)

	for i := 0; i < scoreIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Score(&model.Score{
			Name:    "test-score",
			TraceID: trace.ID,
			Trace_:  trace,
			Value:   0.95,
			Comment: "test comment",
		})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.25, counter.Get(), scoreIterations*2, 0.15)
}

func TestScoreSampling50(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.50)

	for i := 0; i < scoreIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Score(&model.Score{
			Name:    "test-score",
			TraceID: trace.ID,
			Trace_:  trace,
			Value:   0.95,
			Comment: "test comment",
		})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.50, counter.Get(), scoreIterations*2, 0.15)
}

func TestScoreSampling75(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 0.75)

	for i := 0; i < scoreIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Score(&model.Score{
			Name:    "test-score",
			TraceID: trace.ID,
			Trace_:  trace,
			Value:   0.95,
			Comment: "test comment",
		})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 0.75, counter.Get(), scoreIterations*2, 0.15)
}

func TestScoreSampling100(t *testing.T) {
	counter := &EventCounter{}
	server := createMockServer(t, counter)
	defer server.Close()

	ctx := context.Background()
	lf := createLangfuse(ctx, server.URL, 1.0)

	for i := 0; i < scoreIterations; i++ {
		trace, _ := lf.Trace(&model.Trace{Name: "test-trace"})
		lf.Score(&model.Score{
			Name:    "test-score",
			TraceID: trace.ID,
			Trace_:  trace,
			Value:   0.95,
			Comment: "test comment",
		})
	}
	lf.Flush(ctx)

	assertSamplingRate(t, 1.0, counter.Get(), scoreIterations*2, 0)
}
