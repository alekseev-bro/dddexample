package sales

import (
	"fmt"

	"github.com/alekseev-bro/dddexample/ddd/pkg/domain"
)

type CreateCustomer struct {
	Customer
}

func (c *CreateCustomer) Execute(a *Customer) domain.Event[Customer] {
	if a != nil {
		return &domain.EventError[Customer]{Reason: "customer already exists"}
	}

	return &CustomerCreated{Customer: c.Customer}
}

func (c *CreateCustomer) AggregateID() domain.ID[Customer] {
	return c.Customer.ID
}

type ValidateOrdersError struct {
	etype string
	value int
}

func (e ValidateOrdersError) Error() string {
	return fmt.Sprintf("active orders %s %d", e.etype, e.value)
}

type ValidateAgeError struct {
	age uint
}

func (e ValidateAgeError) Error() string {
	return fmt.Sprintf("age is %d, must be greater than 18", e.age)
}

var ErrMaxOrders = &ValidateOrdersError{">=", 3}
var ErrMinOrders = &ValidateOrdersError{"<", 0}

func NewValidateAgeError(age uint) *ValidateAgeError {
	return &ValidateAgeError{age}
}

type ValidateOrder struct {
	CustomerID domain.ID[Customer]
	OrderID    domain.ID[Order]
}

func (v *ValidateOrder) Execute(c *Customer) domain.Event[Customer] {
	if c.Age <= 18 {
		return &OrderRejected{OrderID: v.OrderID, Error: NewValidateAgeError(c.Age).Error()}
	}
	if c.ActiveOrders >= 3 {
		return &OrderRejected{OrderID: v.OrderID, Error: ErrMaxOrders.Error()}
	}

	return &OrderAccepted{OrderID: v.OrderID}
}

func (v ValidateOrder) AggregateID() domain.ID[Customer] {
	return v.CustomerID
}
