package command

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type orderMutator eventstore.Mutator[order.Order, *order.Order]
