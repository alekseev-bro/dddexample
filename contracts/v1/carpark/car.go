package carpark

import "github.com/alekseev-bro/ddd/pkg/aggregate"

type Car struct {
	ID  aggregate.ID
	VIN string
	CarModel
	RentState
	MaintananceState
}

type CarModel struct {
	Brand string
	Model string
}

type RentState uint

const (
	NotAvailable RentState = iota
	Available
)

type MaintananceState uint

const (
	InMaintenance MaintananceState = iota
	NotNeeded
	Needed
)
