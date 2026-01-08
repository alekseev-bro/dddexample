package customers

import "github.com/alekseev-bro/ddd/pkg/essrv"

type Order struct{}
type OrderID = essrv.ID[Order]

type CustomerID = essrv.ID[Customer]
