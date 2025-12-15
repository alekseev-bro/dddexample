package sales

import "github.com/alekseev-bro/ddd/pkg/aggregate"

type Car struct {
	aggregate.ID[Car]
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
