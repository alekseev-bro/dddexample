package orders

import "github.com/alekseev-bro/ddd/pkg/essrv"

type RentOrderStatus uint8

const (
	ProcessingByCustomer RentOrderStatus = iota
	ValidForProcessing
	Closed
)

type Customer struct{}
type Car struct{}

type CustomerID = essrv.ID[Customer]
type OrderID = essrv.ID[Order]
type CarID = essrv.ID[Car]
