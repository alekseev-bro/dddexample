package sales

import (
	"ddd/pkg/domain"
)

type CustomerCreated struct {
	Customer Customer
}

func (cc CustomerCreated) Apply(c *Customer) {
	*c = cc.Customer
}

type CustomerOrderClosed struct {
	OrderID domain.ID[Order]
}

func (CustomerOrderClosed) Apply(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	OrderID domain.ID[Order]
}

func (OrderAccepted) Apply(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID domain.ID[Order]
	Error   error
}

func (OrderRejected) Apply(c *Customer) {

}
