package sales

import (
	"context"
	"ddd/pkg/domain"
	"ddd/pkg/store/natsstore/esnats"
	"ddd/pkg/store/natsstore/snapnats"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type boundedContext struct {
	Customer domain.Aggregate[Customer]
	Order    domain.Aggregate[Order]
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

	customer := domain.NewAggregate[Customer](ctx,
		esnats.NewEventStream[Customer](ctx, js),
		snapnats.NewSnapshotStore[Customer](ctx, js),
	)
	domain.RegisterEvent[*OrderRejected](customer)
	domain.RegisterEvent[*CustomerCreated](customer)
	domain.RegisterEvent[*OrderAccepted](customer)
	oes := esnats.NewEventStream[Order](ctx, js)
	order := domain.NewAggregate[Order](ctx, oes,

		snapnats.NewSnapshotStore[Order](ctx, js),
	)

	domain.RegisterEvent[*OrderCreated](order)
	domain.RegisterEvent[*OrderClosed](order)
	domain.RegisterEvent[*OrderVerified](order)

	// order.Project(ctx, &OrderService{
	// 	Customer: customer,
	// })
	customer.Project(ctx, &CustomerService{
		Order: order,
	}, domain.WithEventFilter[*OrderAccepted]())

	domain.Saga(ctx, order, customer, func(e *OrderCreated) *ValidateOrder {
		return &ValidateOrder{CustomerID: e.Order.CustomerID, OrderID: e.Order.ID}
	})

	// 	return nil
	// })
	// domain.Subscribe(ctx, order, func(e *OrderClosed) error {
	// 	return nil
	// })

	return &boundedContext{
		Customer: customer,
		Order:    order,
	}
}
