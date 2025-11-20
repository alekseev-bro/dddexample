package sales

import (
	"ddd/pkg/domain"
)

type OrderCreated struct {
	Order Order
}

func (ce *OrderCreated) Apply(c *Order) {
	*c = ce.Order
}

// func (ce OrderCreated) String() string {
// 	return "ORDER_CREATED"
// }

type CarAddedToOrder struct {
	OrderID domain.ID[Order]
	CarID   domain.ID[Car]
}

func (ce *CarAddedToOrder) Apply(c *Order) {
	c.Cars[ce.CarID] = struct{}{}
}

type CarRemovedFromOrder struct {
	OrderID domain.ID[Order]
	CarID   domain.ID[Car]
}

func (ce *CarRemovedFromOrder) Apply(c *Order) {
	delete(c.Cars, ce.CarID)
}

type OrderVerified struct {
	OrderID domain.ID[Order]
}

func (ce *OrderVerified) Apply(c *Order) {
	c.Status = ValidForProcessing
}

type OrderClosed struct {
	OrderID domain.ID[Order]
	CustID  domain.ID[Customer]
}

func (ce *OrderClosed) Apply(c *Order) {
	c.Status = Closed
}
