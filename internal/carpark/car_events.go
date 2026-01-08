package carpark

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
)

type CarRentRejected struct {
	OrderID essrv.ID[Car]
}

func (ce CarRentRejected) Apply(c *Car) {

}

type CarRented struct {
	OrderID essrv.ID[Car]
}

func (ce CarRented) Apply(c *Car) {
	c.RentState = NotAvailable
}
