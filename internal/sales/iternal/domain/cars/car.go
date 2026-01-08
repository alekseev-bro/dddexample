package cars

import "github.com/alekseev-bro/ddd/pkg/essrv"

type Car struct {
	essrv.ID[Car]
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
