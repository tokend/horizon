## Startup

- Entry point: cmd/horizon/main.go
- Add environment variables (see example config in the config-example.txt)
- Create Horizon database (make sure `DATABASE_URL` env var contains proper name of DB)
- Apply migrations (run Horizon with `db history migrate up` command-line arguments)
- Build and run Horizon
