package langfuse

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/henomis/langfuse-go/internal/pkg/api"
	"github.com/henomis/langfuse-go/internal/pkg/observer"
	"github.com/henomis/langfuse-go/model"
	"github.com/henomis/langfuse-go/utils"
)

const (
	defaultFlushInterval = 500 * time.Millisecond
)

type Langfuse struct {
	flushInterval time.Duration
	client        *api.Client
	observer      *observer.Observer[model.IngestionEvent]

	samplingRate  float32
}

func New(ctx context.Context) *Langfuse {
	return NewWithClient(ctx, api.New())
}

func NewWithClient(ctx context.Context, client *api.Client) *Langfuse {
	l := &Langfuse{
		flushInterval: defaultFlushInterval,
		client:        client,
		observer: observer.NewObserver(
			ctx,
			func(ctx context.Context, events []model.IngestionEvent) {
				err := ingest(ctx, client, events)
				if err != nil {
					fmt.Println(err)
				}
			},
		),

		samplingRate: 1,
	}

	return l
}

// change to options but keep for now for backward compatability
func (l *Langfuse) WithFlushInterval(d time.Duration) *Langfuse {
	l.flushInterval = d
	return l
}

func (l *Langfuse) WithSamplingRate(s float32) *Langfuse {
	l.samplingRate = s
	return l
}

func ingest(ctx context.Context, client *api.Client, events []model.IngestionEvent) error {
	req := api.Ingestion{
		Batch: events,
	}

	res := api.IngestionResponse{}
	return client.Ingestion(ctx, &req, &res)
}

func (l *Langfuse) Trace(t *model.Trace) (*model.Trace, error) {
	t.ID = utils.BuildID(&t.ID)
	t.ShouldTrace = rand.Float32() < l.samplingRate
	if (t.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeTraceCreate,
				Timestamp: time.Now().UTC(),
				Body:      t,
			},
		)
	}
	return t, nil
}

func (l *Langfuse) Generation(g *model.Generation, parentID *string) (*model.Generation, error) {
	if g.TraceID == "" {
		t, err := l.createTrace(g.Name)
		if err != nil {
			return nil, err
		}

		g.TraceID = t.ID
		g.Trace_ = t
	}

	g.ID = utils.BuildID(&g.ID)

	if parentID != nil {
		g.ParentObservationID = *parentID
	}

	if (g.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeGenerationCreate,
				Timestamp: time.Now().UTC(),
				Body:      g,
			},
		)
	}
	return g, nil
}

func (l *Langfuse) GenerationEnd(g *model.Generation) (*model.Generation, error) {
	if g.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	if (g.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeGenerationUpdate,
				Timestamp: time.Now().UTC(),
				Body:      g,
			},
		)
	}

	return g, nil
}

func (l *Langfuse) Score(s *model.Score) (*model.Score, error) {
	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}
	s.ID = utils.BuildID(&s.ID)

	if (s.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeScoreCreate,
				Timestamp: time.Now().UTC(),
				Body:      s,
			},
		)
	}
	return s, nil
}

func (l *Langfuse) Span(s *model.Span, parentID *string) (*model.Span, error) {
	if s.TraceID == "" {
		t, err := l.createTrace(s.Name)
		if err != nil {
			return nil, err
		}

		s.TraceID = t.ID
		s.Trace_ = t
	}

	s.ID = utils.BuildID(&s.ID)

	if parentID != nil {
		s.ParentObservationID = *parentID
	}

	if (s.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeSpanCreate,
				Timestamp: time.Now().UTC(),
				Body:      s,
			},
		)
	}

	return s, nil
}

func (l *Langfuse) SpanEnd(s *model.Span) (*model.Span, error) {
	if s.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	if (s.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        utils.BuildID(nil),
				Type:      model.IngestionEventTypeSpanUpdate,
				Timestamp: time.Now().UTC(),
				Body:      s,
			},
		)
	}

	return s, nil
}

func (l *Langfuse) Event(e *model.Event, parentID *string) (*model.Event, error) {
	if e.TraceID == "" {
		t, err := l.createTrace(e.Name)
		if err != nil {
			return nil, err
		}

		e.TraceID = t.ID
		e.Trace_ = t
	}

	e.ID = utils.BuildID(&e.ID)

	if parentID != nil {
		e.ParentObservationID = *parentID
	}

	if (e.Trace_.ShouldTrace) {
		l.observer.Dispatch(
			model.IngestionEvent{
				ID:        uuid.New().String(),
				Type:      model.IngestionEventTypeEventCreate,
				Timestamp: time.Now().UTC(),
				Body:      e,
			},
		)
	}

	return e, nil
}

func (l *Langfuse) createTrace(traceName string) (*model.Trace, error) {
	trace, errTrace := l.Trace(
		&model.Trace{
			Name: traceName,
		},
	)
	if errTrace != nil {
		return nil, errTrace
	}

	return trace, nil
}

func (l *Langfuse) Flush(ctx context.Context) {
	l.observer.Wait(ctx)
}

