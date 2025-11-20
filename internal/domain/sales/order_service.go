package sales

import (
	"context"
	"ddd/pkg/domain"

	"github.com/google/uuid"
)

type OrderService struct {
	Customer domain.Aggregate[Customer]
}

func (c *OrderService) Handle(ctx context.Context, eventID uuid.UUID, e domain.Event[Order]) error {
	switch ev := e.(type) {
	case *OrderCreated:
		return c.Customer.Command(ctx, eventID.String(), &ValidateOrder{OrderID: ev.Order.ID, CustomerID: ev.Order.CustomerID})

	}
	return nil
}

// func NewCustomerService(ctx context.Context, customer domain.Aggregate[Customer], order domain.Aggregate[Order]) *CustomerService {
// 	s := &CustomerService{
// 		customer: customer,
// 		order:    order,
// 	}
// 	s.order.Subscribe(ctx, "sales_customer_service", func(e domain.Event[Order]) error {
// 		switch ev := e.(type) {
// 		case *OrderCreated:
// 			return s.customer.Command(ctx, ev.Order.CustomerID, ValidateOrder{OrderID: ev.Order.ID})
// 		}

// 		return nil
// 	}, false)
// 	return s
// }
