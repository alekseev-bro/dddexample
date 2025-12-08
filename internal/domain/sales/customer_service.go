package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/domain"
)

type CustomerService struct {
	Order domain.Aggregate[Order]
}

func (c *CustomerService) Handle(ctx context.Context, eventID domain.EventID[Customer], e domain.Event[Customer]) error {
	switch ev := e.(type) {
	case *OrderAccepted:
		return c.Order.Execute(ctx, eventID.String(), &CloseOrder{OrderID: ev.OrderID})

	}
	return nil
}
