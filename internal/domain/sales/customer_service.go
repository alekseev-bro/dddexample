package sales

import (
	"context"
	"ddd/pkg/domain"
)

type CustomerService struct {
	Order domain.Aggregate[Order]
}

func (c *CustomerService) Handle(ctx context.Context, eventID domain.ID[domain.Event[Customer]], e domain.Event[Customer]) error {
	switch ev := e.(type) {
	case *OrderAccepted:
		return c.Order.Command(ctx, eventID.String(), &CloseOrder{OrderID: ev.OrderID})

	}
	return nil
}
