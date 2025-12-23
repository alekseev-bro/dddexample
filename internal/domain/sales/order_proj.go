package sales

import (
	"context"
	"sync"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/domain"
)

type DB struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewRamDB() *DB {
	return &DB{
		data: make(map[string]any),
	}
}

func (db *DB) Get(key string) (any, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.data[key]
	return val, ok
}

func (db *DB) Set(key string, val any) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = val
}

type OrderSaga struct {
	cust domain.Aggregate[Customer]
}

func (c *OrderSaga) Handle(ctx context.Context, o *Order, eventID aggregate.EventID[Order]) error {

	if err := c.cust.Update(ctx, o.CustomerID, eventID.String(), func(c *Customer) (*Customer, error) {
		if err := c.AddOrder(); err != nil {
			return nil, err
		}
		return c, nil
	}); err != nil {
		return err
	}
	return nil
}
