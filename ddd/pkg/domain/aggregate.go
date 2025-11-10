package domain

import (
	"context"
	"log/slog"

	reg "ddd/internal/registry"
	"ddd/internal/serde"

	"ddd/pkg/store"

	"errors"
	"fmt"

	"github.com/google/uuid"
)

// var nc *nats.Conn

type messageCount uint

const (
	snapshotSize messageCount = 100
)

type ID[T any] string

func NewID[T any]() ID[T] {
	a := uuid.New()
	return ID[T](a.String())
}

type option[T any] func(a *aggregateRoot[T])

func WithSerde[T any](s serde.Serder) option[T] {
	return func(a *aggregateRoot[T]) {
		a.typeReg = reg.New(s)
	}
}

func NewAggregateRoot[T any](ctx context.Context, es eventStream[T], ss snapshotStore[T], opts ...option[T]) *aggregateRoot[T] {

	//ent := PT(new(T))
	// for _, v := range ent.RegisterEvents() {
	// 	etype := reflect.TypeOf(v)
	// 	if etype.Kind() == reflect.Ptr {
	// 		panic("RegisterEvents return type must be a slice of values")
	// 	}

	// 	eventDefaultRegistry[fmt.Sprintf("%s_%s", aname, etype.Name())] = v
	// 	eventNamesRegistry[v] = etype.Name()

	// }

	aggr := &aggregateRoot[T]{
		eventStream:   es,
		snapshotStore: ss,
		typeReg:       reg.New(serde.NewDefaultSerder()),

		//serder:          &DefaultSerder[T]{},
	}
	for _, o := range opts {
		o(aggr)
	}

	//var ent T
	// ent.Events(func(e Applyable[T]) {

	// 	store.eventRegistry.Add(e)
	// })
	// for _, v := range ent.Events() {

	// }
	//	st, ok := streams[ag.Domain().Type]
	//if !ok {

	//	streams[ag.Domain().Type] = st
	//}

	return aggr
}

type Envelope[T any] struct {
	AggrID  ID[T]
	Version uint64
	Event   Event[T]
}

type commander[T any] interface {
	Command(ctx context.Context, id ID[T], command Command[T]) error
	CommandFunc(ctx context.Context, id ID[T], command func(*T) Event[T]) error
}
type subscriber[T any] interface {
	Subscribe(ctx context.Context, name string, handler func(Event[T]) error, ordered bool)
}

type Aggregate[T any] interface {
	commander[T]
	subscriber[T]
}

type registry interface {
	Register(any)
}

type eventSaver[T any] interface {
	Save(ctx context.Context, events []Envelope[T]) error
}

type eventLoader[T any] interface {
	Load(ctx context.Context, aggrID ID[T], fromSeq uint64, handler func(event Event[T]) error) (uint64, error)
}

type eventStream[T any] interface {
	registry
	eventSaver[T]
	eventLoader[T]
	subscriber[T]
}

type snapshotSaver[T any] interface {
	Save(ctx context.Context, aggrID ID[T], snap []byte) error
}

type snapshotLoader[T any] interface {
	Load(ctx context.Context, aggrID ID[T]) ([]byte, error)
}

type snapshotStore[T any] interface {
	snapshotSaver[T]
	snapshotLoader[T]
}

//	type Registry[T Reducible[T]] interface {
//		Register(Applyable[T])
//	}

type aggregateRoot[T any] struct {
	typeReg reg.TypeRegistry
	eventStream[T]
	snapshotStore[T]
}

type snapshot[T any] struct {
	MsgCount messageCount
	Version  uint64
	Body     *T
}

func (a *aggregateRoot[T]) build(ctx context.Context, id ID[T]) (*snapshot[T], error) {

	var ent T
	//var snap Snapshot[T]
	// rec, err := a.snap.Get(ctx, id)
	// if err != nil {
	// 	if !errors.Is(err, store.ErrNoSnapshot) {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}
	// 	if err := a.serder.Deserialize(rec, &snap); err != nil {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}

	// } else {
	// 	if err := a.serder.Deserialize(rec, &snap); err != nil {
	// 		return nil, fmt.Errorf("build: %w", err)
	// 	}
	// }

	var totalMsgs messageCount

	last, err := a.eventStream.Load(ctx, id, 0, func(e Event[T]) error {

		e.Apply(&ent)
		totalMsgs++

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("buid %w", err)
	}
	sn := &snapshot[T]{Version: last, Body: &ent, MsgCount: totalMsgs}

	return sn, nil
}

type CommandFunc[T any] func(*T) Event[T]

func (f CommandFunc[T]) Execute(t *T) Event[T] {
	return f(t)
}

func (a *aggregateRoot[T]) RegisterEvent(event Event[T]) {

	a.eventStream.Register(event)
}
func (a *aggregateRoot[T]) RegisterCommand(command Command[T]) {
	a.typeReg.Register(command)
}

func (a *aggregateRoot[T]) CommandFunc(ctx context.Context, id ID[T], command func(*T) Event[T]) error {
	return a.Command(ctx, id, CommandFunc[T](command))
}

func (a *aggregateRoot[T]) Command(ctx context.Context, id ID[T], command Command[T]) error {

	var err error

	snap := &snapshot[T]{}

	snap, err = a.build(ctx, id)
	if err != nil {
		if !errors.Is(err, store.ErrNoAggregate) {
			return fmt.Errorf("build aggrigate: %w", err)
		}
		snap = &snapshot[T]{}
	}

	evt := command.Execute(snap.Body)

	if e, ok := evt.(EventError[T]); ok {
		return fmt.Errorf("command: %w", e)
	}

	if err := a.eventStream.Save(ctx, []Envelope[T]{{AggrID: id, Version: snap.Version, Event: evt}}); err != nil {
		return fmt.Errorf("command: %w", err)
	}

	// Save snapshot if aggregate has more than snapshotSize messages
	if snap != nil {
		if snap.MsgCount >= snapshotSize {
			go func() {
				b, err := a.typeReg.Serialize(snap)
				if err != nil {
					slog.Warn(err.Error())
				}
				a.snapshotStore.Save(ctx, id, b)
			}()

		}
	}

	return nil
}
