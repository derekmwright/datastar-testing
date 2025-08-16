# Go + Datastar + NATS.io + PostgreSQL Example Site

Just building out a basic example hypermedia web framework that I can clone and use for other projects.

Feel free to use it if you want...

## Database Setup

```postgresql
CREATE USER exampleapp WITH PASSWORD 'exampleapp';

CREATE DATABASE exampleapp OWNER exampleapp;

GRANT ALL PRIVILEGES ON DATABASE exampleapp TO exampleapp;
```

Creating new migrations

```shell
migrate create -dir ./db/migrations -seq -digits 3 -ext sql create_items_table
```

Running migrations

```shell
migrate -source file://db/migrations -database $APP_DATABASE_URI up
```

# License

[MIT](https://mit-license.org/)
