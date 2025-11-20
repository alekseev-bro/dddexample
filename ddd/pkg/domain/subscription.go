package domain

import (
	"context"
)

type SubscribeParams struct {
	DurableName string
	Ordered     bool
	Kind        []string
}

func WithOrder(ordered bool) SubOption {
	return func(p *SubscribeParams) {
		p.Ordered = ordered
	}
}

func WithEventFilter[E Event[T], T any]() SubOption {
	return func(p *SubscribeParams) {
		var ev E
		p.Kind = append(p.Kind, typeNameFrom(ev))
	}
}

func WithName(name string) SubOption {
	return func(p *SubscribeParams) {
		p.DurableName = name
	}
}

type EventHandler[T any] interface {
	Handle(ctx context.Context, eventID ID[Event[T]], event Event[T]) error
}

func (a *aggregate[T]) Project(ctx context.Context, h EventHandler[T], opts ...SubOption) {
	params := &SubscribeParams{
		DurableName: typeNameFrom(h),
		Ordered:     true,
		Kind:        nil,
	}
	for _, opt := range opts {
		opt(params)
	}

	a.es.Subscribe(ctx, func(envel *Envelope) error {

		ev := a.getType(envel.Kind, envel.Payload)
		return h.Handle(ctx, ID[Event[T]](envel.ID.String()), ev.(Event[T]))
	}, params)
}
