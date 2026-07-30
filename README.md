# idea-to-mvp

A structured, AI-first process framework to take any idea from concept to shippable MVP.

## What Is This

A **process framework** --- like Scrum, Kanban, or Shape Up --- designed for AI agents.
It is not a code scaffold, not a boilerplate, and not tied to any language, stack, or AI
provider. It guides any idea through discovery, architecture, refinement, and delivery
using an 8-phase pipeline with human approval gates at every critical decision point.

Technology agnostic. Provider agnostic. Works with
[any AI tool that supports the Open Agentic Standard](https://agents.md).

**Everything lives in one file: [`AGENTS.md`](AGENTS.md).**

## The Pipeline

| # | Phase | Purpose | Human Gate |
|---|-------|---------|------------|
| 1 | **Capture** | Understand the idea, problem, and audience | MIM: approve problem statement |
| 2 | **Discover** | Research domain, users, constraints | MIM: approve discovery report |
| 3 | **Architect** | Define tech stack, data model, API contracts | MIM: approve architecture |
| 4 | **Refine** | Write user stories with acceptance criteria | MIM: approve backlog |
| 5 | **Plan** | Break stories into tasks with dependency order | MIM: approve sprint plan |
| 6 | **Build** | Implement via TDD with per-task handoffs | PDC after each task |
| 7 | **Verify** | Independent review against specs and DOD | MIM: approve verification |
| 8 | **Accept** | Scrum team review, final human sign-off | MIM: release decision |

**MIM** = Manual Inspection Milestone --- the human decides, the agent waits.

## Key Concepts

- **5 Execution Axioms** --- non-negotiable gates (Handoff, Orchestrator, Echo, TDD, Native)
  that every action must pass before executing
- **Echo System** --- a 4-stage build pipeline (Build, Test, E2E, Audit) that runs
  before every commit
- **Contract Truth Gate** --- tests must trace to typed contract artifacts, not
  ad-hoc mocks or prose descriptions
- **Scaffold Readiness Gate** --- environment must be mechanically verified as
  functional before feature work begins
- **Anti-Rationalization Protocol** --- blocks agents from reasoning around mandatory
  rules ("given the simplicity, we can skip...")
- **Lightweight Mode** --- reduced ceremony for small projects (< 5 stories) while
  preserving quality invariants

## When To Use This

**Good fit**: you have an idea and want an AI agent to build it methodically ---
with architecture decisions, test coverage, and human oversight at critical points.

**Not a fit**: you want to vibe-code a quick prototype without process, or you
need a code template with pre-built components.

## Quick Start

1. Create a new repository from this template
2. Clone your new repo
3. Point your AI agent at the repo --- it reads `AGENTS.md` automatically
4. Describe your idea --- the agent takes it from there

The agent handles everything from here: capturing your idea, researching the
domain, proposing architecture, writing stories, planning sprints, building with
TDD, and presenting deliverables for your review at each MIM gate.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

PostgreSQL License --- see [LICENSE](LICENSE). Covers the framework files only (AGENTS.md,
ARTIFACT-TEMPLATES.md, CONTRIBUTING.md). Projects built with this template are yours ---
license them however you want.
