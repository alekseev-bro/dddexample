package esnats

type option[T any] func(*eventStream[T]) error
