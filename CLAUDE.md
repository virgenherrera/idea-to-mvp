# Virgil — Development Mode

You are working on **Virgil itself** — the methodology source repo.
This is NOT a consumer project. Do not run `virgil init` here.

## Project Identity

| Key | Value |
|-----|-------|
| Remote | `git@github.com:virgenherrera/virgil.git` |
| Go module | `github.com/virgenherrera/virgil` |
| Owner | Hugo Virgen (`virgenherrera`) |
| Binary | `cmd/virgil/` → single static binary |

**BEFORE writing any module path, import, URL, or package reference**:
run `git remote -v` and use that namespace. NEVER use another org's
namespace. NEVER guess. The remote is the source of truth.

## Repository Structure

```
docs/       ← Spanish methodology docs (original)
docs/en/    ← English methodology docs (canonical, embedded in binary)
cmd/virgil/ ← CLI entry point (cobra)
internal/
  distribution/ ← init, link, extract, config
  feedback/     ← types, writer, reader, report, export
  awareness/    ← awareness block injection for consumers
  health/       ← health dashboard (stub)
embedded.go     ← go:embed all:docs/en
```

## Conventions

- Go module: `github.com/virgenherrera/virgil`
- All imports use this module path — verify before writing
- Methodology docs: English in `docs/en/` is canonical (embedded in binary)
- Spanish in `docs/` is the original — keep in sync
- Commits: conventional commits, no AI attribution
- Tests: zero tokens — no LLM calls in tests, ever
- Binary distribution: static binary download, `go install` for contributors only

## Build

```bash
go build ./cmd/virgil/        # local binary
go install ./cmd/virgil/      # install to $GOPATH/bin
go vet ./...                  # lint
go test ./...                 # tests (when they exist)
```

## What Virgil IS (three responsibilities, nothing more)

1. **Methodology provider** — source of truth (RAG-able docs)
2. **Mechanism provider** — SM behavior, artifactStore, orchestration protocols
3. **Orchestrator** — SM orchestrates specialist agents as needed

Virgil does NOT generate project code. Virgil provides HOW to work,
not WHAT to build.

## Awareness Block (for consumers, not this repo)

`virgil init` injects an awareness block into consumer projects'
CLAUDE.md/AGENTS.md. That block teaches the agent to recognize Virgil
scenarios and act accordingly. See `internal/awareness/awareness.go`.

## Feedback System

Feedback is JSONL in `.virgil/feedback/` (consumer side). Three types:
- `gap` — methodology didn't cover this scenario
- `friction` — covered but unclear/slow
- `win` — worked well, preserve

`virgil feedback report` aggregates by frequency × impact.
