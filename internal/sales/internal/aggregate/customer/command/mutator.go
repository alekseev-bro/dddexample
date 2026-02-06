package command

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type customerMutator eventstore.Mutator[customer.Customer, *customer.Customer]
