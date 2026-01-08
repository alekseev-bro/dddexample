package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/features/saga/order_verification"
	"github.com/nats-io/nats.go/jetstream"
)

type Module struct {
	CustomerService essrv.Root[customers.Customer]
	OrderService    essrv.Root[orders.Order]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {
	customer := essrv.New(ctx,
		esnats.NewEventStream[customers.Customer](ctx, js, esnats.EventStreamConfig{
			StoreType: esnats.Memory,
		}),
		snapnats.NewSnapshotStore[customers.Customer](ctx, js, snapnats.SnapshotStoreConfig{
			StoreType: snapnats.Memory,
		}),
		essrv.AggregateConfig{
			SnapthotMsgThreshold: 5,
		},
		essrv.WithEvent[customers.OrderRejected](),
		essrv.WithEvent[customers.OrderAccepted](),
		essrv.WithEvent[customers.CustomerRegistered](),
	)

	order := natsstore.NewAggregate(ctx, js,
		natsstore.NatsAggregateConfig{
			AggregateConfig: essrv.AggregateConfig{
				SnapthotMsgThreshold: 5,
			},
			StoreType: natsstore.Memory,
		},
		essrv.WithEvent[orders.OrderClosed](),
		essrv.WithEvent[orders.OrderPosted](),
		essrv.WithEvent[orders.OrderVerified](),
	)
	order_verification.StartSaga(ctx)

	return &Module{
		CustomerService: customer,
		OrderService:    order,
	}
}
