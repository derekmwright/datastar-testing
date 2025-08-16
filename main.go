package main

import (
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var Version = "0.0.1"

type application struct {
	config config

	sessions   *scs.SessionManager
	logger     *slog.Logger
	natsClient *nats.Conn
	ready      bool
	cache      jetstream.KeyValue
	db         *pgxpool.Pool
}

func main() {
	// Register any types serialized into session data here, example:
	// gob.Register(time.Time{})

	app := &application{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	if err := app.setDefaults(); err != nil {
		os.Exit(1)
	}
	app.parseFlags()

	if err := app.serve(); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
