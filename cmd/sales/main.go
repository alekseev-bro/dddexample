package main

import (
	"context"
	"ddd/pkg/domain"
	"log/slog"

	"dddexample/internal/domain/sales"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelInfo)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}
	s := sales.New(ctx, js)

	go func() {
		for {

			cusid := domain.NewID[sales.Customer]()
			idempc := domain.NewIdempotencyKey(cusid, "CreateCustomer")

			err := s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
			if err != nil {
				panic(err)
			}

			ordid := domain.NewID[sales.Order]()
			idempo := domain.NewIdempotencyKey(ordid, "CreateOrder")

			err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
			if err != nil {
				panic(err)
			}
			<-time.After(1 * time.Second)
		}

	}()

	<-ctx.Done()

}
