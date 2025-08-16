package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"exampleapp/internal/natsstore"
)

func (app *application) serve() error {
	ctx := context.Background()

	natsSrv, err := app.natsServe()
	if err != nil {
		return err
	}

	if err = app.startSessions(ctx); err != nil {
		return err
	}

	if err = app.startCache(ctx); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		s := <-quit

		app.logger.Info("shutting down servers", slog.String("signal", s.String()))

		ctxCancel, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		app.logger.Info("closing NATS client connection")
		app.natsClient.Close()

		// Shutdown the NATS server
		app.logger.Info("shutting down embedded NATS server")
		natsSrv.Shutdown()
		app.logger.Info("stopped embedded NATS server")

		app.logger.Info("shutting down http server")
		shutdownError <- srv.Shutdown(ctxCancel)
	}()

	app.logger.Info("starting http server", slog.String("addr", srv.Addr))

	if err = srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	if err = <-shutdownError; err != nil {
		return err
	}

	app.logger.Info("stopped http server", slog.String("addr", srv.Addr))

	return nil
}

func (app *application) natsServe() (*server.Server, error) {
	var err error

	opts := &server.Options{
		Host:      app.config.nats.address,
		Port:      app.config.nats.port,
		JetStream: true, // required
		StoreDir:  app.config.nats.storeDir,
		NoSigs:    true, // required
	}

	app.logger.Info(
		"starting embedded NATS server",
		slog.String("addr", opts.Host+":"+strconv.Itoa(opts.Port)),
		slog.String("store-dir", opts.StoreDir),
		slog.Bool("jetstream-enabled", opts.JetStream),
	)

	srv, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	go srv.Start()

	if !srv.ReadyForConnections(5 * time.Second) {
		return nil, errors.New("nats server not ready")
	}

	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", app.config.nats.address, app.config.nats.port))
	if err != nil {
		return nil, err
	}

	app.natsClient = nc

	return srv, nil
}

func (app *application) startSessions(ctx context.Context) error {
	js, err := jetstream.New(app.natsClient)
	if err != nil {
		return err
	}

	sessionStore := natsstore.Must(
		natsstore.New(
			ctx,
			js,
			jetstream.KeyValueConfig{
				Bucket:      app.config.sessions.bucketName,
				Compression: true,
				TTL:         app.config.sessions.TTL,
			},
			natsstore.WithPrefix(app.config.sessions.prefix),
		),
	)

	app.sessions = scs.New()
	app.sessions.Store = sessionStore

	return nil
}

func (app *application) startCache(ctx context.Context) error {
	js, err := jetstream.New(app.natsClient)
	if err != nil {
		return err
	}

	cacheStore, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      app.config.cache.bucketName,
		Compression: true,
		TTL:         app.config.cache.TTL,
	})
	if err != nil {
		return err
	}

	app.cache = cacheStore

	return nil
}
