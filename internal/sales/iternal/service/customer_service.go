package service

// import (
// 	"context"

// 	"github.com/alekseev-bro/ddd/pkg/essrv"
// )

// type CustomerService struct {
// 	Order essrv.Root[Order]
// }

// func (c CustomerService) HandleAllEvents(ctx context.Context, e *essrv.Event[Customer]) error {
// 	switch ev := e.Body.(type) {
// 	case OrderAccepted:
// 		_, err := c.Order.Execute(ctx, ev.OrderID, CloseOrder{}, e.ID.String())
// 		return err

// 	}
// 	return nil
// }
