package sales

import (
	"ddd/pkg/domain"
)

type CreateOrder struct {
	OrderID domain.ID[Order]
	CustID  domain.ID[Customer]
}

func (c CreateOrder) Execute(o *Order) domain.Event[Order] {
	event := &OrderCreated{
		Order{ID: c.OrderID, CustomerID: c.CustID,
			Cars: make(map[domain.ID[Car]]struct{}), Status: ProcessingByCustomer,
		}}
	return event
}

type CloseOrder struct {
	OrderID domain.ID[Order]
	CustID  domain.ID[Customer]
}

func (c CloseOrder) Execute(o *Order) domain.Event[Order] {
	event := &OrderClosed{
		OrderID: c.OrderID,
		CustID:  o.CustomerID,
	}
	return event
}

type AddCarToOrder struct {
	OrderID domain.ID[Order]
	CarID   domain.ID[Car]
}

func (c AddCarToOrder) Execute(o *Order) domain.Event[Order] {
	event := &CarAddedToOrder{
		OrderID: c.OrderID,
		CarID:   c.CarID,
	}
	return event
}
