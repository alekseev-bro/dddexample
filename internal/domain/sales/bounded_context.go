package sales

import (
	"context"
	"ddd/pkg/domain"
	"ddd/pkg/store/esnats"
	"ddd/pkg/store/snapnats"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type boundedContext struct {
	Customer        domain.Aggregate[Customer]
	Order           domain.Aggregate[Order]
	orderService    *OrderService
	customerService *CustomerService
}

// type MySerder struct {
// }

// func (m *MySerder) Serialize(in any) ([]byte, error) {

// 	var buf bytes.Buffer
// 	if err := gob.NewEncoder(&buf).Encode(in); err != nil {

// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func (m *MySerder) Deserialize(data []byte, out any) error {

// 	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(out); err != nil {
// 		return err
// 	}
// 	return nil
// }

func New(ctx context.Context) *boundedContext {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}
	custStream := esnats.NewEventStream[Customer](ctx, js)
	customer := domain.NewAggregateRoot(ctx,
		custStream,
		snapnats.NewSnapshotStore[Customer](ctx, js),
	)

	domain.RegisterEvent[CustomerCreated](customer)
	domain.RegisterEvent[OrderAccepted](customer)

	order := domain.NewAggregateRoot[Order](ctx,
		esnats.NewEventStream[Order](ctx, js),
		snapnats.NewSnapshotStore[Order](ctx, js),
	)

	domain.RegisterEvent[OrderCreated](order)
	domain.RegisterEvent[OrderClosed](order)
	domain.RegisterEvent[OrderVerified](order)

	c := &boundedContext{
		Customer:        customer,
		Order:           order,
		orderService:    NewOrderService(ctx, customer, order),
		customerService: NewCustomerService(ctx, customer, order),
	}
	return c
}
