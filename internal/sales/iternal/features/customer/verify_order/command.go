package verify_order

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
)

type Command struct {
	OrderID orders.OrderID
}

func (v Command) Dispatch(c *customers.Customer) (essrv.Events[customers.Customer], error) {
	return c.VerifyOrder(customers.OrderID(v.OrderID)), nil
}
