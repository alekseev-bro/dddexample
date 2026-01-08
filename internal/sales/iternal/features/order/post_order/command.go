package post_order

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/iternal/domain/orders"
)

type Command struct {
	*orders.Order
}

func (ce Command) Dispatch(o *orders.Order) (essrv.Events[orders.Order], error) {
	return o.PostOrder(ce.Order), nil
}
