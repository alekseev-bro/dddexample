package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alekseev-bro/dddexample/internal/carpark"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car/command"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type jsPublisher struct {
	js jetstream.JetStream
}

func (p *jsPublisher) Publish(ctx context.Context, kind string, event []byte) error {
	_, err := p.js.Publish(ctx, kind, event)
	if err != nil {
		return err
	}
	return nil
}

func NewJSPublisher(js jetstream.JetStream) *jsPublisher {
	return &jsPublisher{js: js}
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelInfo)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}

	// Initialize the carpark module
	cp := carpark.NewModule(ctx, js, NewJSPublisher(js))
	cmd := command.RegisterCar{
		Car: car.New("X5", "BMW"),
	}
	_, err = cp.RegisterCarHandler.HandleCommand(ctx, cmd)
	if err != nil {
		panic(err)
	}
	// Start the carpark module

	<-ctx.Done()
}
