package close_order

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
)

type Command struct {
}

func (cmd Command) Dispatch(o *orders.Order) (essrv.Events[orders.Order], error) {
	return o.CloseOrder(), nil
}
