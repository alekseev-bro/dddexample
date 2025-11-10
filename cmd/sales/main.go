package main

import (
	"context"
	"ddd/pkg/domain"
	"os"
	"os/signal"
	"ttt/internal/domain/sales"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	// nc, err := nats.Connect(nats.DefaultURL)
	// if err != nil {
	// 	slog.Error("connect to nats", "error", err)
	// 	panic(err)
	// }

	// _, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: "atest", Subjects: []string{"atest.>"}, AllowAtomicPublish: true})
	// if err != nil {
	// 	slog.Error("create stream", "error", err)
	// 	panic(err)
	// }

	// _, err = js.PublishMsg(ctx, m, jetstream.WithExpectLastSequenceForSubject(uint64(0), "atest.t"))
	// if err != nil {
	// 	slog.Error("publish message", "error", err)
	// 	panic(err)
	// }

	// w.Start()

	s := sales.New(ctx)
	cusid := domain.NewID[sales.Customer]()
	err := s.Customer.Command(ctx, cusid, sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	if err != nil {
		panic(err)
	}

	ordid := domain.NewID[sales.Order]()
	err = s.Order.Command(ctx, ordid, sales.CreateOrder{OrderID: ordid, CustID: cusid})
	if err != nil {
		panic(err)
	}

	// err = s.Order.Command(ctx, ordid, sales.CloseOrder{OrderID: ordid})
	// if err != nil {
	// 	panic(err)
	// }

	<-ctx.Done()
}
