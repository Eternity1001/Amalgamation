# Amalgamation

## Run with Docker

Start the Go server and PostgreSQL:

```sh
docker compose up --build
```

The server listens on <http://localhost:8080> and PostgreSQL is available on
port `5432`. The database defaults are for local development; set
`POSTGRES_DB`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` in the environment or a
`.env` file to override them. PostgreSQL data is stored in the `postgres_data`
Docker volume.

The Go server serves `html/index.html` at `/`, files in `js/` under `/js/`, and
files in `photo/` under `/photo/`. The `html/`, `js/`, and `photo/` directories
are copied into the application image. The server currently does not connect
to or query PostgreSQL; the database service and connection environment
variables are provided for future use.

#To Do

# Make a socket using GOO
# Return Some HTML