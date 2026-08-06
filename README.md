# Virgil

AI-assisted development methodology — distributed as a single Go binary.

Virgil provides three things:

1. **Methodology** — 31 docs defining an end-to-end framework for AI-assisted software development (planning, execution, acceptance)
2. **Mechanisms** — SM behavior, artifact state machine, retrieval patterns, orchestration protocols
3. **Feedback loop** — structured collection of gaps, friction, and wins to improve the methodology over time

## Install

### Download binary (no dependencies)

Go compiles to a single static binary — no runtime, no Go installation
needed.

```bash
# macOS (Apple Silicon)
curl -fsSL https://github.com/virgenherrera/virgil/releases/latest/download/virgil-darwin-arm64 -o virgil
chmod +x virgil
sudo mv virgil /usr/local/bin/

# macOS (Intel)
curl -fsSL https://github.com/virgenherrera/virgil/releases/latest/download/virgil-darwin-amd64 -o virgil
chmod +x virgil
sudo mv virgil /usr/local/bin/

# Linux (amd64)
curl -fsSL https://github.com/virgenherrera/virgil/releases/latest/download/virgil-linux-amd64 -o virgil
chmod +x virgil
sudo mv virgil /usr/local/bin/
```

### Build from source (contributors only)

```bash
# Requires Go 1.19+
git clone https://github.com/virgenherrera/virgil.git
cd virgil
go install ./cmd/virgil/
# Binary at $GOPATH/bin/virgil
```

### Verify

```bash
virgil version
# virgil 0.5.0-dev
```

## Quick Start

### 1. Initialize in your project

```bash
cd your-project
virgil init
```

This creates:

```
your-project/
├── .virgil/
│   ├── config.yaml          # tier, language, feedback settings
│   ├── docs/                # methodology docs (extracted)
│   │   ├── overview.md      # start here
│   │   ├── execution/       # Red/Green/Refactor/Accept
│   │   ├── planning/        # phases, artifacts, roles
│   │   └── ...
│   └── feedback/            # JSONL feedback logs
└── CLAUDE.md                # ← awareness block injected here
```

The key output is the **awareness block** appended to your agent config
(CLAUDE.md, AGENTS.md, or .cursorrules — whichever exists). This teaches
your AI agent what Virgil is and when to use it.

### 2. Talk to your agent naturally

After `virgil init`, your agent knows about Virgil. Just talk:

| You say | Agent does |
|---------|-----------|
| "I have an idea for this project" | Reads overview.md, starts planning flow |
| "Do a takeover of this codebase" | Reads overview.md, runs `virgil scan` |
| "Let's build feature X" | Checks for planning artifacts, starts from the right phase |
| "There's a bug in production" | Reads fast-forward.md, follows interruption protocol |

You don't need to say "virgil" — the agent recognizes the scenario.

### 3. Collect feedback

The awareness block includes a **session observer** that automatically
records methodology feedback at session close. You can also record
feedback manually:

```bash
# Record a gap (methodology didn't cover this)
virgil feedback gap \
  --doc "overview.md" \
  --section "Actors" \
  --scenario "Rate limiting during Green phase" \
  --impact "high" \
  --description "No guidance on handling rate limits mid-test"

# Record a win (methodology worked well here)
virgil feedback win \
  --doc "execution/contracts.md" \
  --section "Contract Types" \
  --scenario "Caught breaking API change" \
  --impact "medium" \
  --description "Contract validation caught it before deployment"

# Record a session summary
virgil feedback session \
  --goal "Implement auth middleware" \
  --docs-used "overview.md,execution/contracts.md" \
  --methodology-grade "B" \
  --notes "fastForward worked well for this scope"

# View aggregated report
virgil feedback report

# Export raw data as JSON
virgil feedback export
```

## Adoption Tiers

```bash
virgil init --tier minimal    # 5 docs — overview, glossary, artifact schemas
virgil init --tier standard   # 24 docs — core + execution + behavior + roles (default)
virgil init --tier full       # 31 docs — everything
```

Start with `standard`. Move to `full` when you want the complete methodology.

## Local Development (Contributing to Virgil)

If you're editing the methodology docs and want to test changes in a
real project immediately:

```bash
# In the virgil source repo:
virgil link --register

# In any consumer project (after virgil init):
virgil link
# .virgil/docs/ is now a symlink to your source repo's docs/en/
# Changes to the source docs are visible immediately

# When done:
virgil unlink
# Re-extracts embedded docs, removes symlink
```

This is the `npm link` equivalent for methodology development.

## Agent Setup

### Claude Code

`virgil init` automatically injects the awareness block into `CLAUDE.md`.
No additional setup needed — just open Claude Code in your project.

### Cursor

If `.cursorrules` exists, `virgil init` appends the awareness block there.
Otherwise it creates `CLAUDE.md` (Cursor reads both).

### Codex / Other Agents

If `AGENTS.md` exists, `virgil init` appends there. The awareness block
is plain markdown — any agent that reads project-level config files will
pick it up.

### Manual Setup

If your agent doesn't read any of those files, copy the awareness block
from the generated file into whatever config your agent reads. The block
is between `<!-- Generated by virgil init -->` and
`<!-- End Virgil auto-configuration -->`.

## Commands

| Command | Description |
|---------|-------------|
| `virgil init [--tier T]` | Bootstrap methodology in current project |
| `virgil link --register [path]` | Register a virgil source repo |
| `virgil link` | Symlink docs to registered source |
| `virgil unlink` | Remove symlink, restore embedded docs |
| `virgil feedback gap` | Record a methodology gap |
| `virgil feedback friction` | Record a friction point |
| `virgil feedback win` | Record a methodology win |
| `virgil feedback session` | Record a session summary |
| `virgil feedback report` | View aggregated feedback report |
| `virgil feedback export` | Export all feedback as JSON |
| `virgil health` | Project health dashboard (coming soon) |
| `virgil scan` | Methodology compliance scan (coming soon) |
| `virgil version` | Print version |

## How It Works

Virgil embeds all 31 methodology docs into the binary via Go's `embed`
package. When you run `virgil init`, it extracts the docs for your tier
into `.virgil/docs/` and injects an awareness block into your agent's
config file.

The awareness block teaches the agent:
- What Virgil is (6 artifacts, 3 actors, execution cycle)
- When to activate it (new project, takeover, feature, bug)
- What CLI commands are available
- How to record feedback at session close (via a non-interfering observer)

The feedback system stores structured JSONL logs in `.virgil/feedback/`.
Each entry is traceable (links to a doc and section), actionable (describes
what was missing), and contextual (includes the scenario). The report
command aggregates feedback by frequency and impact to prioritize
methodology improvements.

## Project Structure

```
cmd/virgil/           # CLI entry point and cobra commands
internal/
  distribution/       # init, link, extract, config
  feedback/           # types, writer, reader, report, export
  awareness/          # awareness block injection
  health/             # health dashboard (stub)
docs/en/              # 31 methodology docs (embedded into binary)
```

## License

MIT
