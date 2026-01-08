package service

// import (
// 	"context"

// 	"github.com/alekseev-bro/ddd/pkg/essrv"
// )

// type OrderService struct {
// 	CustomerStore essrv.Root[Customer]
// }

// func (c OrderService) Handle(ctx context.Context, e *essrv.Event[Order]) error {
// 	switch ev := e.Body.(type) {
// 	case OrderPosted:
// 		_, err := c.CustomerStore.Execute(ctx, ev.CustomerID, VerifyOrder{OrderID: ev.ID}, e.ID.String())
// 		return err

// 	}
// 	return nil
// }
