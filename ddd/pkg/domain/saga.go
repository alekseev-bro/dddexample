package domain

import (
	"context"
	"fmt"
)

type sagaHandlerFunc[E Event[T], C Command[U], T any, U any] func(event E) C

type sagaHandler[E Event[T], C Command[U], T any, U any] struct {
	sub     projector[T]
	cmd     commander[U]
	handler sagaHandlerFunc[E, C, T, U]
}

func (sf *sagaHandler[E, C, T, U]) Handle(ctx context.Context, eventID EventID[T], event Event[T]) error {

	return sf.cmd.Command(ctx, eventID.String(), sf.handler(event.(E)))

}

func Saga[E Event[T], C Command[U], T any, U any](ctx context.Context, sub projector[T], cmd commander[U], shf sagaHandlerFunc[E, C, T, U]) {

	sh := &sagaHandler[E, C, T, U]{
		sub:     sub,
		cmd:     cmd,
		handler: shf,
	}
	var (
		ee E
		cc C
		uu U
		tt T
	)

	ename := typeNameFrom(ee)
	cname := typeNameFrom(cc)
	sname := typeNameFrom(tt)
	cmname := typeNameFrom(uu)
	durname := fmt.Sprintf("%s:%s|%s:%s", sname, ename, cmname, cname)

	sub.Project(ctx, sh, WithName(durname), WithOrder(false), WithEventFilter[E]())

}
