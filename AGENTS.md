# go-telegram-monitor — agent instructions

Small Go service that sends daily CPU/RAM/disk reports to a Telegram chat (07:30 and 17:30) and alerts when
memory or disk usage goes above 80% (checked every 5 minutes). Runs as a systemd service. Public portfolio repo.

## Stack
Go 1.24 (`go.mod`), `gopsutil/v3` (metrics), `gocron/v2` (scheduling), `godotenv` (config). Single package.

## Run
```bash
cp .env.example .env    # fill TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID
go run .                # or: go build -o monitor .
```

## Verify
- Preferred: `go vet ./... && go build -o /dev/null .`
- Without module downloads: `gofmt -l .` (must print nothing).
- There are no tests.

## Layout
- `main.go` — entry point, scheduler, metric collection and Telegram calls.
- `go-telegram-monitor.service` — systemd unit example (`User`, `WorkingDirectory`, `ExecStart` must be adapted).

## Conventions
- Keep it a single, dependency-light binary; config only through environment variables.
- `gofmt` clean; README in English.

## Careful
- This repository is **public**: no real tokens, chat IDs, hostnames, usernames or personal paths in code, docs or commits.
- Running the binary sends real Telegram messages; use a test chat.

## Secrets
`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` — only in a local `.env` (git-ignored) or the service environment. Never commit values.

## When done
- Run the verification above and report the real output and the commit SHA.
- Commit only your files with a clear message; push only when the task asks for it.
