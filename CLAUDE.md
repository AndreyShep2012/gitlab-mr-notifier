# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # Compile binary for local OS → ./gitlab-mr-notifier
make build-linux    # Cross-compile for Linux (CGO disabled) → ./bootstrap
make aws-zip        # Build Linux binary + create Lambda deployment ZIP
make test           # Run all tests with race detection and coverage
make check-escape   # Run Go escape analysis
```

Run a single test:
```bash
go test ./internal/slack/... -run TestGetBody -v
```

## Architecture

A Go service that polls GitLab for open merge requests and sends Slack notifications. Supports three deployment modes: local binary with scheduling, Docker container, and AWS Lambda.

**Execution flow:**

```
cmd/main.go
  └─ Detects environment (AWS_LAMBDA_RUNTIME_API env var → internal/utils/aws.go)
      ├─ Local: internal/runner/local.go → loads config, optionally schedules with cron
      └─ Lambda: internal/runner/awslambda.go → starts Lambda runtime handler

internal/runner/common.go check()
  ├─ Instantiates notifier (Slack or Log) and message formatter
  ├─ Calls GitLab API for each configured group/project ID
  └─ Sends notifications
```

**Key design pattern:** All major components are coded against interfaces in `internal/interfaces/` (Runner, Notifier, MessageFormatter, Cron, GitlabApi), enabling substitution (e.g., `logNotifier` instead of Slack).

**Configuration** (`internal/config/config.go`): Environment-based via `cleanenv`. Optional file loading via `CONFIG_PATH` env var pointing to a `.env` file. AWS Lambda variant (`internal/config/awslambda.go`) decrypts KMS-encrypted environment variables.

| Variable | Required | Notes |
|----------|----------|-------|
| `GITLAB_TOKEN` | Yes | GitLab API read token |
| `GITLAB_GROUP_IDS` | Yes* | Comma-separated group IDs |
| `GITLAB_PROJECT_IDS` | No | Comma-separated project IDs |
| `SLACK_WEBHOOK_URL` | Yes | Slack incoming webhook |
| `CRON_PERIOD` | No | e.g. `10m`, `1d`, `1w`; empty = run once |
| `CRON_TIME` | No | Required with `1d`/`1w`, e.g. `10:30` |
| `NOTIFIER` | No | `slack` (default) or `log` |
| `MESSAGE_DESCRIPTION_LIMIT` | No | Default 500 chars |

*At least one of `GITLAB_GROUP_IDS` or `GITLAB_PROJECT_IDS` required.

**Scheduling** (`internal/cron/cron.go`): Wraps `go-co-op/gocron`. Period formats: `10s`, `10m` (every N duration), `1d`/`1w` (daily/weekly at a specific time).

**Slack notifications** (`internal/slack/`): MRs are partitioned into success and failed pipeline groups. Each successful-pipeline MR gets its own message; failed-pipeline MRs are batched into a single message. Message formatting lives in `slack/messageFormatter.go` and respects `MESSAGE_DESCRIPTION_LIMIT`.

**GitLab API** (`internal/gitlabapi/gitlabapi.go`): Uses `xanzy/go-gitlab`. Fetches open, non-WIP MRs; enriches each with full MR details and unresolved discussion thread counts.

## Deployment

```bash
# Docker
docker run -it --rm --env-file=.env gitlab-mr-notifier

# AWS Lambda: deploy gitlab-mr-notifier.zip, trigger via EventBridge Scheduler
# KMS-encrypt env vars using internal/config/awslambda.go decryption conventions
```
