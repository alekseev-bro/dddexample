package main

import (
	"context"
	"ddd/pkg/domain"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	"ttt/internal/domain/sales"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	//	slog.SetLogLoggerLevel(slog.LevelError)
	// nc, err := nats.Connect(nats.DefaultURL)
	// if err != nil {
	// 	slog.Error("connect to nats", "error", err)
	// 	panic(err)
	// }

	// _, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: "atest", Subjects: []string{"atest.>"}, AllowAtomicPublish: true})
	// if err != nil {
	// 	slog.Error("create stream", "error", err)
	// 	panic(err)
	// }

	// _, err = js.PublishMsg(ctx, m, jetstream.WithExpectLastSequenceForSubject(uint64(0), "atest.t"))
	// if err != nil {
	// 	slog.Error("publish message", "error", err)
	// 	panic(err)
	// }

	// w.Start()

	s := sales.New(ctx)

	go func() {
		for {

			cusid := domain.NewID[sales.Customer]()
			idempc := domain.NewIdempotencyKey(cusid, "CreateCustomer")

			err := s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
			if err != nil {
				panic(err)
			}
			// for range 1 {

			ordid := domain.NewID[sales.Order]()
			idempo := domain.NewIdempotencyKey(ordid, "CreateOrder")

			err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
			if err != nil {
				panic(err)
			}
			<-time.After(1 * time.Second)
		}

	}()

	mux := http.NewServeMux()

	serv := http.Server{
		Addr:    ":8086",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("shutting down...")
		if err := serv.Shutdown(sctx); err != nil {
			// Only log an error if s.Shutdown returned one.
			// This handles cases where the timeout was reached or another error occurred.
			slog.Error("server shutdown failed", "error", err)
			return
		}
		slog.Info("server shutdown complete")
	}()
	switch err := serv.ListenAndServe(); err {
	case http.ErrServerClosed:
		// This is the expected, successful shutdown exit.
		slog.Info("server stopped gracefully")
	case nil:
		// ListenAndServe shouldn't return nil, but as a safeguard.
		slog.Info("server exited without error")
	default:
		// Log any other non-nil, unexpected error.
		slog.Error("server failed to start or stopped unexpectedly", "error", err)
	}

}
