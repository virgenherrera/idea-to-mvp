# Virgil

End-to-end methodology for AI-assisted software development.

An alternative to a monolithic AGENTS.md: instead of a 1000+ line file
prescribing how an agent operates, Virgil defines what to do — concepts,
definitions, and methodology — in modular documents under `docs/`.

## Current scope

Definitions, concepts, and methodology only. No implementation yet.
The CLI that consumes this methodology will be built once the dogma is
stabilized.

## Documentation

`docs/` contains the complete methodology (Spanish originals; English
versions live under `docs/en/`):

- `docs/overview.md` (English: `docs/en/overview.md`) — general
  overview, governing dogma, progressive adoption
- `docs/execution/` (English: `docs/en/execution/`) — Red/Green/Refactor/Accept
  cycle
- `docs/planning/` (English: `docs/en/planning/`) — phases, artifacts, SM
  behavior
- `docs/operation/` (English: `docs/en/operation/`) — optional post-execution
  facade
- `docs/glossary.md` (English: `docs/en/glossary.md`) — terminology
- `docs/agile-adaptations.md` (English: `docs/en/agile-adaptations.md`) —
  explicit trade-offs vs. classic Agile
- `docs/echo-system.md` (English: `docs/en/echo-system.md`) — hooks and
  verification system
