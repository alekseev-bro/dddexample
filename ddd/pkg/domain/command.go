package domain

type Command[T any] interface {
	Execute(entity *T) Event[T]
	identer[T]
}

type identer[T any] interface {
	AggregateID() ID[T]
}

func RegisterCommand[E Command[T], T any](root registry) {
	var cmd E
	root.register(cmd)
}
