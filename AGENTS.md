# file-sync

A CLI tool that downloads files from remote URLs to local paths based on a YAML config file. Supports preprocessing (trim/fix lines, append extensions), optional format conversion, and concurrent multi-task execution.

## Project Structure

```
main.go                         # Entry point: flag parsing, config loading, task orchestration
internal/
  config/config.go              # Config structs (SyncConfig, SyncTask, FilePair) + Viper-based YAML loading
  fetcher/fetcher.go            # HTTP GET client with 20s timeout, status code validation, preset headers
  pipeline/pipeline.go          # Business logic: download → preProcess → convert → move
```

## Core Flow

`main()` parses `-c <config>` → `config.InitConfig()` reads YAML via Viper → spawns a goroutine per `SyncTask` → each task spawns goroutines per `FilePair` → each file goes through:

1. **download** — HTTP GET to URL, write to `os.TempDir()`
2. **preProcess** — trim lines, fix `- '+.` prefix → `- '`, skip `payload:` lines, append `Extensions` at EOF
3. **convert** (optional) — strip leading `-` and all `'` characters per line
4. **move** — copy+delete to `path` (skips if file is empty)

## Key Conventions

- All business packages under `internal/` — not importable externally
- Errors returned up the stack, logged at the call site (no silent failures)
- `fetcher.Fetch(url, ua)` returns `(data, error)` — caller checks both
- `config.InitConfig` calls `os.Exit(1)` on fatal errors (config not found, invalid YAML)
- Config file is any format Viper supports (YAML recommended), default path `./config.yaml`

## Config Shape

```yaml
default-ua: string           # fallback User-Agent
sync-tasks:
  - tag: string              # task label (for logging)
    file-pairs:
      - url: string          # download source
        path: string         # destination on disk
        convert: bool        # strip leading '-' and '\''?
        ua: string           # per-file UA override
        extensions: [string] # lines appended during preProcess
```

## Build & Run

```bash
go build -o bin/file-sync .
./bin/file-sync -c config.yaml
```

Cross-compile targets: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`. See `.github/workflows/ci.yml` and `.vscode/tasks.json`.
