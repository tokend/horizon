## Startup

Entry point: cmd/horizon/main.go
Add environment variables (see example config in the config-example.txt)
Create horizon database (make sure `DATABASE_URL` env var contains proper name of DB)
Apply migrations (run horizon with `db history migrate up` command-line arguments)
Build and run horizon
