package integration

import (
	"context"
	"log/slog"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/codec"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car"
)

type Publisher interface {
	Publish(ctx context.Context, name string, event []byte) error
}

type carHandler struct {
	publisher Publisher
	codec     codec.Codec
}

func NewCarHandler(publisher Publisher, codec codec.Codec) *carHandler {
	return &carHandler{
		publisher: publisher,
		codec:     codec,
	}
}

func (h *carHandler) HandleEvents(ctx context.Context, event aggregate.Evolver[car.Car]) error {
	switch ev := event.(type) {
	case *car.Arrived:

		b, err := h.codec.Marshal(ev.ToArrivedV1())
		if err != nil {
			slog.Error("can't marshal event", "error", err)
			return err
		}
		h.publisher.Publish(ctx, "carpark.car.arrived", b)

	}
	return nil
}
