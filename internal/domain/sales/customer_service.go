package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type CustomerService struct {
	Order aggregate.Aggregate[Order]
}

func (c *CustomerService) Handle(ctx context.Context, eventID aggregate.EventID[Customer], e aggregate.Event[Customer]) error {
	switch ev := e.(type) {
	case *OrderAccepted:
		_, err := c.Order.Execute(ctx, eventID.String(), &CloseOrder{OrderID: ev.OrderID})
		return err

	}
	return nil
}
