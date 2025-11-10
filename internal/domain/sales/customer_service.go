package sales

import (
	"context"
	"ddd/pkg/domain"
)

type CustomerService struct {
	customer domain.Aggregate[Customer]
	order    domain.Aggregate[Order]
}

func NewCustomerService(ctx context.Context, customer domain.Aggregate[Customer], order domain.Aggregate[Order]) *CustomerService {
	s := &CustomerService{
		customer: customer,
		order:    order,
	}
	s.order.Subscribe(ctx, "sales_customer_service", func(e domain.Event[Order]) error {
		switch ev := e.(type) {
		case *OrderCreated:
			return s.customer.Command(ctx, ev.Order.CustomerID, ValidateOrder{OrderID: ev.Order.ID})
		}

		return nil
	}, false)
	return s
}
