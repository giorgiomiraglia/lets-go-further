package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(app.logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// Listens in background for termination signals
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		// Blocks until SIGTERM or SIGINT is received
		s := <-quit

		app.logger.PrintInfo(
			"Termination signal received",
			map[string]string{
				"signal": s.String(),
			})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Gives the server up to 5s to clean-up,
		// if everything goes well, nil is sent.
		// Otherwise, sends a timeout error or any other error
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo(
		"starting %s server on %s",
		map[string]string{
			"addr": srv.Addr,
			"env":  app.config.env,
		})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Blocks until the channel receives data.
	// If an error ocurred during the shutdown or it exceeded the 5s,
	// returns the error sent
	err = <-shutdownError
	if err != nil {
		return err
	}

	// Otherwise, the shutdown went ok
	app.logger.PrintInfo(
		"stopped server",
		map[string]string{"addr": srv.Addr},
	)

	return nil
}
