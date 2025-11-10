package esnats

import (
	"context"
	"ddd/internal/registry"
	"ddd/internal/serde"
	"ddd/pkg/domain"
	"ddd/pkg/store"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"math"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-io/orbit.go/jetstreamext"
)

// const (
// 	eventTypeHeader string = "ev_type"
// )

type eventStream[T any] struct {
	registry.TypeRegistry
	tname      string
	boundedCtx string
	js         jetstream.JetStream
}

func NewEventStream[T any](ctx context.Context, js jetstream.JetStream, opts ...option[T]) *eventStream[T] {
	aname, bcname := registry.MetaFromType[T]()

	stream := &eventStream[T]{js: js, tname: aname, boundedCtx: bcname, TypeRegistry: registry.New(serde.NewDefaultSerder())}

	for _, opt := range opts {
		opt(stream)
	}

	_, err := stream.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Subjects:    []string{stream.allSubjects()},
		Name:        stream.streamName(),
		Storage:     jetstream.MemoryStorage,
		AllowDirect: true,
	})
	if err != nil {
		panic(err)
	}
	return stream
}

func (s *eventStream[T]) subjectNameForID(agrid string) string {
	return fmt.Sprintf("%s:%s.%s", s.boundedCtx, s.tname, agrid)
}

func (s *eventStream[T]) streamName() string {
	return fmt.Sprintf("%s:%s", s.boundedCtx, s.tname)
}

func (s *eventStream[T]) allSubjects() string {
	return fmt.Sprintf("%s.*", s.streamName())
}

func (s *eventStream[T]) Save(ctx context.Context, events []domain.Envelope[T]) error {
	for _, envel := range events {
		msg := nats.NewMsg(s.subjectNameForID(string(envel.AggrID)))
		//msg.Header.Add(eventTypeHeader, event.Type)
		msg.Header.Add(jetstream.MsgIDHeader, uuid.New().String())
		b, err := s.Serialize(envel.Event)
		if err != nil {
			return fmt.Errorf("command: %w", err)
		}

		tname := registry.TypeNameFrom(envel.Event)
		if !s.TypeExists(tname) {
			slog.Error("event type not registered", "event", tname)
			panic(fmt.Errorf("event type not registered: %s", tname))
		}
		rec := StoreRecord{Body: b, Kind: tname}

		r, err := s.Serialize(rec)
		if err != nil {
			return fmt.Errorf("command: %w", err)
		}

		msg.Data = r
		retries := 0
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				_, err := s.js.PublishMsg(ctx, msg, jetstream.WithExpectLastSequencePerSubject(envel.Version))
				if err != nil {

					var seqerr *jetstream.APIError

					if errors.As(err, &seqerr) {
						if seqerr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
							slog.Warn("OCC", "version", envel.Version, "name", s.subjectNameForID(string(envel.AggrID)))
							retries++
							if retries > 50 {
								panic("OCC DeadLock")
							}
							continue
						}
					}
					return fmt.Errorf("store event func: %w", err)

				}
				slog.Info("event stored", "kind", tname, "subject", s.subjectNameForID(string(envel.AggrID)), "stream", s.streamName())
				return nil
			}
		}
	}
	return nil
}

type StoreRecord struct {
	Body json.RawMessage
	Kind string
}

func (s *eventStream[T]) Load(ctx context.Context, id domain.ID[T], version uint64, handler func(event domain.Event[T]) error) (uint64, error) {

	subj := s.subjectNameForID(string(id))
	msgs, err := jetstreamext.GetBatch(ctx,
		s.js, s.streamName(), math.MaxInt, jetstreamext.GetBatchSubject(subj),
		jetstreamext.GetBatchSeq(version+1))
	//fmt.Println(time.Since(start))

	if err != nil {
		return 0, fmt.Errorf("get events: %w", err)
	}

	var lastevent uint64
	for msg, err := range msgs {
		if err != nil {
			if errors.Is(err, jetstreamext.ErrNoMessages) {
				return 0, store.ErrNoAggregate
			}
			return 0, fmt.Errorf("build func can't get msg batch: %w", err)
		}

		lastevent = msg.Sequence
		var rec StoreRecord

		if err := s.Deserialize(msg.Data, &rec); err != nil {

			return 0, fmt.Errorf("build : %w", err)
		}

		ev, err := s.GetType(rec.Kind, rec.Body)
		if err != nil {
			panic(fmt.Sprintf("event not registered: %s", rec.Kind))
		}

		if err := handler(ev.(domain.Event[T])); err != nil {
			return 0, fmt.Errorf("get events: %w", err)
		}
	}
	return lastevent, nil
}

func (e *eventStream[T]) Subscribe(ctx context.Context, name string, handler func(event domain.Event[T]) error, ordered bool) {
	maxpend := 1000
	if ordered {
		maxpend = 1
	}

	cons, err := e.js.CreateOrUpdateConsumer(ctx, e.streamName(), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("subscription-%s", name),
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: maxpend,
	})
	if err != nil {
		panic(fmt.Errorf("subscription create consumer: %w", err))
	}
	ct, err := cons.Consume(func(msg jetstream.Msg) {

		var rec StoreRecord
		e.Deserialize(msg.Data(), &rec)
		ev, err := e.GetType(rec.Kind, rec.Body)
		if err != nil {
			slog.Error("subscribe type error", "error", err)
			panic(err)
		}

		if err := handler(ev.(domain.Event[T])); err != nil {
			slog.Warn("redelivering", "error", err)
			msg.NakWithDelay(1 * time.Second)
			return
		}
		msg.Ack()

	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {}))
	if err != nil {
		panic(fmt.Errorf("subscription consume: %w", err))
	}
	go func() {
		<-ctx.Done()
		ct.Drain()
		fmt.Println("CLOSED")
	}()
}
