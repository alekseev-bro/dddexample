package orders

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
)

type Order struct {
	ID         OrderID
	CustomerID CustomerID
	Cars       map[CarID]struct{}
	Status     RentOrderStatus
	Deleted    bool
}

func (o *Order) PostOrder(ord *Order) essrv.Events[Order] {
	if o.ID.IsZero() {
		return essrv.NewEvents(OrderPosted{Order: ord})
	}
	return nil
}

func (o *Order) CloseOrder() essrv.Events[Order] {
	if o.Status != Closed {
		return essrv.NewEvents(OrderClosed{})
	}
	return nil
}
