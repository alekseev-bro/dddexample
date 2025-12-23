package main

import (
	"context"
	"time"

	"os"
	"os/signal"

	"github.com/alekseev-bro/dddexample/internal/domain/sales"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	//	slog.SetLogLoggerLevel(slog.LevelWarn)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}

	s := sales.New(ctx, js)
	s.StartOrderCreationSaga(ctx)

	time.Sleep(time.Second)
	custid, err := s.Customer.Create(ctx, &sales.Customer{Name: "joe", Age: 10})

	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 10 {

		// ordid := s.Order.NewID()
		// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

		// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		// if err != nil {
		// 	panic(err)
		// }
		_, err := s.Order.Create(ctx, &sales.Order{CustomerID: custid})
		if err != nil {
			panic(err)
		}
	}

	<-ctx.Done()

}
