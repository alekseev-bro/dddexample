package sales

import (
	"context"
	"time"

	"github.com/alekseev-bro/ddd/pkg/domain"
	"github.com/alekseev-bro/ddd/pkg/domain/event"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

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
//
//

func New(ctx context.Context, js jetstream.JetStream) *boundedContext {

	customer := domain.NewAggregate(ctx,
		esnats.NewEventStream(ctx, js, esnats.WithInMemory[Customer]()),
		snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Customer]()),
		aggregate.WithSnapshotThreshold[Customer](10, time.Second),
		domain.WithEvent[*OrderRejected](),
		domain.WithEvent[*OrderAccepted](),
	)

	order := natsstore.NewAggregate(ctx, js,
		natsstore.WithSnapshotThreshold[Order](10, time.Second),
		natsstore.WithInMemory[Order](),
		natsstore.WithEvent[*OrderClosed](),
	)

	bc := &boundedContext{
		Customer: customer,
		Order:    order,
	}

	return bc
}

func (b *boundedContext) StartOrderCreationSaga(ctx context.Context) {

	aggregate.SagaStep(ctx, b.Order.(aggregate.Aggregate[Order]), b.Customer.(aggregate.Aggregate[Customer]), func(e *event.Created[Order]) *ValidateOrder {
		return &ValidateOrder{CustomerID: e.Body.CustomerID, OrderID: e.ID}
	})

	b.Order.StartService(ctx, &OrderSaga{b.Customer})

}
