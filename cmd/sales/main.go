package main

import (
	"context"
	"time"

	"os"
	"os/signal"

	sales "github.com/alekseev-bro/dddexample/internal/sales/iternal"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
	register_customer "github.com/alekseev-bro/dddexample/internal/sales/iternal/features/customer/register_cutomer"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/features/order/post_order"
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

	s := sales.NewModule(ctx, js)

	time.Sleep(time.Second)
	custid := s.CustomerService.NewID()
	cust := &customers.Customer{ID: custid, Name: "Joe", Age: 21}
	_, err = s.CustomerService.ExecuteUnique(ctx, custid, register_customer.Command{Customer: cust})
	if err != nil {
		panic(err)
	}
	custid2 := s.CustomerService.NewID()
	cust2 := &customers.Customer{ID: custid, Name: "Joe", Age: 21}
	_, err = s.CustomerService.ExecuteUnique(ctx, custid2, register_customer.Command{Customer: cust2})
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, custid.String(), &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 20 {

		// ordid := s.Order.NewID()
		// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

		// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		// if err != nil {
		// 	panic(err)
		// }
		ordID := s.OrderService.NewID()
		ord := &orders.Order{ID: ordID, CustomerID: orders.CustomerID(custid)}
		_, err := s.OrderService.ExecuteUnique(ctx, ordID, post_order.Command{Order: ord})
		if err != nil {
			panic(err)
		}
	}
	<-ctx.Done()

}
