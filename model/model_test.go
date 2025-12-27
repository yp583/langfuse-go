package model

import "testing"

func TestGenerationSetTraceInheritsEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	gen := &Generation{Name: "test"}
	gen.SetTrace(trace)

	if gen.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", gen.Environment)
	}
}

func TestGenerationSetTracePreservesEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	gen := &Generation{Name: "test", Environment: "staging"}
	gen.SetTrace(trace)

	if gen.Environment != "staging" {
		t.Errorf("expected environment 'staging', got '%s'", gen.Environment)
	}
}

func TestSpanSetTraceInheritsEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	span := &Span{Name: "test"}
	span.SetTrace(trace)

	if span.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", span.Environment)
	}
}

func TestSpanSetTracePreservesEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	span := &Span{Name: "test", Environment: "staging"}
	span.SetTrace(trace)

	if span.Environment != "staging" {
		t.Errorf("expected environment 'staging', got '%s'", span.Environment)
	}
}

func TestEventSetTraceInheritsEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	event := &Event{Name: "test"}
	event.SetTrace(trace)

	if event.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", event.Environment)
	}
}

func TestEventSetTracePreservesEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	event := &Event{Name: "test", Environment: "staging"}
	event.SetTrace(trace)

	if event.Environment != "staging" {
		t.Errorf("expected environment 'staging', got '%s'", event.Environment)
	}
}

func TestScoreSetTraceInheritsEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	score := &Score{Name: "test"}
	score.SetTrace(trace)

	if score.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", score.Environment)
	}
}

func TestScoreSetTracePreservesEnvironment(t *testing.T) {
	trace := &Trace{ID: "trace-1", Environment: "production"}
	score := &Score{Name: "test", Environment: "staging"}
	score.SetTrace(trace)

	if score.Environment != "staging" {
		t.Errorf("expected environment 'staging', got '%s'", score.Environment)
	}
}
