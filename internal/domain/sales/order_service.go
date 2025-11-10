package sales

import (
	"context"
	"ddd/pkg/domain"
)

type OrderService struct {
	customer domain.Aggregate[Customer]
	order    domain.Aggregate[Order]
}

func NewOrderService(ctx context.Context, customer domain.Aggregate[Customer], order domain.Aggregate[Order]) *OrderService {
	s := &OrderService{
		customer: customer,
		order:    order,
	}
	s.customer.Subscribe(ctx, "sales_order_service", func(e domain.Event[Customer]) error {
		switch ev := e.(type) {
		case *OrderAccepted:
			return s.order.Command(ctx, ev.OrderID, CloseOrder{OrderID: ev.OrderID})

		}
		return nil
	}, false)
	return s
}
