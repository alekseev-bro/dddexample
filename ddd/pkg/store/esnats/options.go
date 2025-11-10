package esnats

import (
	"ddd/internal/registry"
	"ddd/internal/serde"
)

type option[T any] func(*eventStream[T]) error

func WithSerder[T any](serder serde.Serder) option[T] {
	return func(es *eventStream[T]) error {
		es.TypeRegistry = registry.New(serder)
		return nil
	}
}
