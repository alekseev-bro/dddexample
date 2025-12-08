package esnats

type option[T any] func(*eventStream[T]) error

func WithPartitions[T any](partitions byte) option[T] {
	return func(es *eventStream[T]) error {
		es.partnum = partitions
		return nil
	}
}
