package sales

import (
	"github.com/alekseev-bro/ddd/pkg/domain"
)

type RentOrderStatus uint8

const (
	ProcessingByCustomer RentOrderStatus = iota
	ValidForProcessing
	Closed
)

type Order struct {
	ID         domain.ID[Order]
	CustomerID domain.ID[Customer]
	Cars       map[domain.ID[Car]]struct{}
	Status     RentOrderStatus
}
