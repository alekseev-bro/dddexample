package sales

import "github.com/alekseev-bro/dddexample/ddd/pkg/domain"

type Car struct {
	domain.ID[Car]
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
