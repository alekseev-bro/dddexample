package order_verification

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/features/customer/verify_order"
)

type CustomerID = essrv.ID[customers.Customer]

func StartSaga(ctx context.Context) {
	essrv.SagaStep(ctx, func(event orders.OrderPosted) (verify_order.Command, CustomerID) {
		return verify_order.Command{OrderID: event.ID}, CustomerID(event.CustomerID)
	})

}
