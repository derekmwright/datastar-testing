package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type config struct {
	http struct {
		address string
		port    int
	}
	nats struct {
		address  string
		port     int
		storeDir string
	}
	sessions struct {
		bucketName string
		prefix     string
		TTL        time.Duration
	}
	cache struct {
		bucketName string
		prefix     string
		TTL        time.Duration
	}
	database struct {
		host    string
		port    int
		name    string
		user    string
		pass    string
		sslmode string
	}
}

func setDefaultInt(envVar string, defaultValue int) (int, error) {
	var (
		err   error
		value int
		ok    bool
	)

	if _, ok = os.LookupEnv(envVar); !ok {
		value = defaultValue
	} else {
		value, err = strconv.Atoi(os.Getenv(envVar))
		if err != nil {
			return 0, err
		}
	}

	return value, nil
}

func setDefaultDuration(envVar string, defaultValue time.Duration) (time.Duration, error) {
	var (
		err   error
		value time.Duration
		ok    bool
	)

	if _, ok = os.LookupEnv(envVar); !ok {
		return defaultValue, nil
	} else {
		value, err = time.ParseDuration(os.Getenv(envVar))
		if err != nil {
			return defaultValue, err
		}
	}

	return value, nil
}

func (app *application) setDefaults() error {
	var (
		err error
		ok  bool
	)

	// Check env vars and set defaults
	if app.config.http.address, ok = os.LookupEnv("APP_LISTEN_ADDR"); !ok {
		app.config.http.address = "0.0.0.0"
	}

	if app.config.http.port, err = setDefaultInt("APP_HTTP_PORT", 8080); err != nil {
		app.logger.Error("unable to parse APP_HTTP_PORT", slog.String("error", err.Error()))
		return err
	}

	if app.config.nats.address, ok = os.LookupEnv("APP_NATS_ADDR"); !ok {
		app.config.nats.address = "0.0.0.0"
	}

	if app.config.nats.port, err = setDefaultInt("APP_NATS_PORT", 4222); err != nil {
		app.logger.Error("unable to parse APP_NATS_PORT", slog.String("error", err.Error()))
		return err
	}

	if app.config.nats.storeDir, ok = os.LookupEnv("APP_NATS_STORAGE_DIR"); !ok {
		app.config.nats.storeDir = "./data"
	}

	if app.config.sessions.bucketName, ok = os.LookupEnv("APP_SESSION_BUCKET_NAME"); !ok {
		app.config.sessions.bucketName = "sessions"
	}

	if app.config.sessions.prefix, ok = os.LookupEnv("APP_SESSION_PREFIX"); !ok {
		app.config.sessions.prefix = "scs"
	}

	if app.config.sessions.TTL, err = setDefaultDuration("APP_SESSION_TTL", 24*time.Hour); err != nil {
		app.logger.Error("unable to parse APP_SESSION_TTL", slog.String("error", err.Error()))
		return err
	}

	if app.config.cache.bucketName, ok = os.LookupEnv("APP_CACHE_BUCKET_NAME"); !ok {
		app.config.cache.bucketName = "cache"
	}

	if app.config.cache.prefix, ok = os.LookupEnv("APP_CACHE_PREFIX"); !ok {
		app.config.cache.prefix = ""
	}

	if app.config.cache.TTL, err = setDefaultDuration("APP_CACHE_TTL", 24*time.Hour); err != nil {
		app.logger.Error("unable to parse APP_CACHE_TTL", slog.String("error", err.Error()))
		return err
	}

	if app.config.database.host, ok = os.LookupEnv("APP_DATABASE_HOST"); !ok {
		app.config.database.host = "localhost"
	}

	if app.config.database.port, err = setDefaultInt("APP_DATABASE_PORT", 5432); err != nil {
		app.logger.Error("unable to parse APP_DATABASE_PORT", slog.String("error", err.Error()))
		return err
	}

	if app.config.database.name, ok = os.LookupEnv("APP_DATABASE_NAME"); !ok {
		app.config.database.name = "exampleapp"
	}

	if app.config.database.name, ok = os.LookupEnv("APP_DATABASE_USERNAME"); !ok {
		app.config.database.name = "exampleapp"
	}

	if app.config.database.name, ok = os.LookupEnv("APP_DATABASE_PASSWORD"); !ok {
		app.config.database.name = "exampleapp"
	}

	if app.config.database.name, ok = os.LookupEnv("APP_DATABASE_SSLMODE"); !ok {
		app.config.database.name = "disable"
	}

	return nil
}

func (app *application) parseFlags() {
	flag.StringVar(&app.config.http.address, "http-address", app.config.http.address, "HTTP listen address")
	flag.IntVar(&app.config.http.port, "http-port", app.config.http.port, "HTTP listen port")
	flag.StringVar(&app.config.nats.address, "nats-address", app.config.nats.address, "NATS listen address")
	flag.IntVar(&app.config.nats.port, "nats-port", app.config.nats.port, "NATS listen port")
	flag.StringVar(&app.config.nats.storeDir, "nats-store-dir", app.config.nats.storeDir, "NATS store directory")
	flag.StringVar(&app.config.sessions.bucketName, "sessions-bucket-name", app.config.sessions.bucketName, "Session storage bucket name")
	flag.StringVar(&app.config.sessions.prefix, "sessions-prefix", app.config.sessions.prefix, "Session storage key prefix")
	flag.DurationVar(&app.config.sessions.TTL, "sessions-ttl", app.config.sessions.TTL, "Session storage TTL")
	flag.StringVar(&app.config.cache.bucketName, "cache-bucket-name", app.config.cache.bucketName, "Cache storage bucket name")
	flag.StringVar(&app.config.cache.prefix, "cache-prefix", app.config.cache.prefix, "Cache storage key prefix")
	flag.DurationVar(&app.config.cache.TTL, "cache-ttl", app.config.cache.TTL, "Cache storage TTL")
	flag.StringVar(&app.config.database.host, "database-host", app.config.database.host, "Database host")
	flag.IntVar(&app.config.database.port, "database-port", app.config.database.port, "Database port")
	flag.StringVar(&app.config.database.name, "database-name", app.config.database.name, "Database name")
	flag.StringVar(&app.config.database.user, "database-username", app.config.database.user, "Database user")
	flag.StringVar(&app.config.database.pass, "database-password", app.config.database.pass, "Database password")
	flag.StringVar(&app.config.database.sslmode, "database-sslmode", app.config.database.sslmode, "Database sslmode")
	flag.Parse()

}
