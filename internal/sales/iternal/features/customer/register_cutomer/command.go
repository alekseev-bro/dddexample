package register_customer

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/customers"
)

type Command struct {
	*customers.Customer
}

func (cmd Command) Dispatch(c *customers.Customer) (essrv.Events[customers.Customer], error) {
	return c.Register(cmd.Customer), nil
}
