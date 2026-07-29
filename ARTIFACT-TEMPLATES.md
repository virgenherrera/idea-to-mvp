# Artifact Templates --- idea-to-mvp Framework

These templates are referenced by [AGENTS.md](AGENTS.md). Load only the template
for the artifact being created.

This file defines the format for every artifact the pipeline produces. The orchestrator
and workers MUST follow these formats. Each artifact begins with a breadcrumb linking back
to the index: `> [INDEX](path/to/INDEX.md) / [Parent] / Document Title`.

## INDEX.md --- Project Dashboard

File: `docs/INDEX.md`. Created during Capture phase. This is the agent's primary
reference for project state --- read it to understand what exists and what remains.

Sections:

1. **Header** --- project name, one-line description
2. **Document Map** --- Mermaid flowchart showing relationships between all artifacts
3. **Project Overview** --- checkboxes for project-brief.md and domain-glossary.md
4. **Epics** --- checkboxes for each epic, nested with their user stories
5. **Architecture** --- checkboxes for tech-stack, testing-strategy, and other arch docs
6. **Sprint Planning** --- checkboxes for batch plans and handoff files
7. **Navigation Notes** --- how documents cross-reference each other

Checkbox convention: `- [ ]` not started, `- [x]` complete. The orchestrator updates
INDEX.md after every story completion (see Story Completion protocol).

## Project Brief

File: `docs/project-brief.md`. Created during Capture phase.

Sections:

1. **Vision** --- one paragraph describing the desired future state
2. **Problem Statement** --- what pain or gap this project addresses
3. **Deliverable Map** --- Mermaid flowchart of high-level deliverable dependencies
4. **Deliverables** --- checkbox list, each linking to its epic
5. **Constraints** --- time, budget, technology, compliance, or other limitations
6. **Evaluation Criteria** --- table: Criterion | Weight (High/Medium/Low) | Description
7. **Related Documents** --- links to domain glossary, architecture, epics

## Domain Glossary

File: `docs/domain-glossary.md`. Created during Discover phase.

Sections:

1. **Entities** --- table: Entity | Description | Key Attributes | Relationships | Rationale
2. **Value Objects** --- table: Name | Description | Constraints
3. **Domain Events** --- table: Event | Trigger | Outcome
4. **Status/State Values** --- allowed states for entities with lifecycles (Mermaid state diagram)
5. **Entity Relationships** --- Mermaid ER diagram
6. **Conventions** --- naming, casing, date/time formats

**Competitive Landscape** (populated during Discover):

| Competitor/Alternative | Approach | Strengths | Weaknesses | Relevance |
| ---------------------- | -------- | --------- | ---------- | --------- |

## Epic

File: `docs/epics/EP{NN}-{slug}.md`. Created during Discover phase, refined during Refine.

Sections:

1. **Summary** --- 2-3 sentences describing the epic's scope
2. **Business Value** --- why this epic matters to the user/stakeholder
3. **Domain Flow** --- Mermaid diagram showing the domain-level flow this epic enables
4. **User Stories** --- checkbox list with priority labels (`Must Have`, `Should Have`, `Could Have`)
5. **Acceptance Boundaries** --- constraints that apply to ALL stories in this epic
6. **Related Architecture** --- links to relevant architecture docs
7. **Related Documents**

## User Story

File: `docs/user-stories/US-{NNN}-{slug}.md`. Created during Refine phase. This is the
most detailed artifact --- it is the contract between planning and implementation.

Sections (all required unless marked optional):

1. **Metadata** --- epic link, priority (`Must Have` | `Should Have` | `Could Have`), estimation
   (`S` | `M` | `L`), status
2. **Story** --- `As a {persona}, I want {action}, so that {value}`
3. **Definition of Ready** --- checkbox list. Items:
   - Domain entity contract frozen (fields, types, nullability, constraints)
   - Interface or API contract frozen (request/response shapes, error formats)
   - Input validation rules enumerated with exact boundaries
   - Edge cases identified with boundary behavior defined
   - Dependencies identified and resolved or deferred
   - Test plan exists with test names mapped to ACs
   - Out-of-scope items listed

   > Items referencing API contracts or request/response shapes apply only when the story
   > touches an API boundary. For non-API stories, mark these as N/A with justification.

4. **Acceptance Criteria** --- numbered with story prefix (AC-{NNN}.1, AC-{NNN}.2, etc.).
   Each criterion MUST use Given/When/Then format:

   ```text
   - [ ] **AC-{NNN}.1: {title}**
     - **Given** {precondition}
     - **When** {action}
     - **Then** {expected outcome}
   ```

5. **Definition of Done** --- checkbox list. Items:
   - All ACs pass with automated test evidence
   - Unit tests green for domain logic and validation
   - Integration tests green against real dependencies (when applicable)
   - No regressions in existing test suite
   - Error responses conform to agreed shape
   - Code reviewed
   - INDEX.md updated
6. **Deliverables** --- two tables:
   - Files to Create: File Path | Contents
   - Files to Modify: File Path | Change
7. **Test Plan** --- table: Test Name | AC | Assertion
   Every AC must be covered by at least one test. Each test states the specific assertion.
8. **Validation Rules** --- enumerated constraints per input field with exact boundaries
   (min, max, required, format, edge cases)
9. **Risks** --- table: Severity | Risk | Mitigation
   Severities: `CRITICAL` | `HIGH` | `MEDIUM` | `LOW`
10. **Out of Scope** --- bullet list of explicit exclusions to prevent scope creep
11. **Notes** (optional) --- implementation hints, open questions, team decisions
12. **Related Documents** --- links to epic, architecture, related stories
13. **Handoff Files** (populated during Plan) --- links to handoff files implementing this story
14. **Change Log** --- date, change, reason for each modification after initial creation

## Batch Plan

File: `docs/subtasks/{epic}/batch-{N}-plan.md`. Created during Plan phase.

Sections:

1. **Scope** --- what this batch delivers, referencing the engineering addenda
2. **Task List** --- table: Task ID | Task Name | Persona | Model Tier | Depends On
3. **Dependency Graph** --- Mermaid DAG showing task dependencies. Solid arrows for hard
   dependencies, dashed arrows for soft dependencies
4. **Execution Order** --- numbered waves (parallel wave 1, sequential step 2, etc.)
5. **Definition of Done --- Batch** --- checkboxes for batch-level completion criteria
6. **Related Documents**

## Handoff File

File: `docs/subtasks/{epic}/{task-id}-{slug}.md`. Created during Plan phase. The
12-section structure is defined in the [Delegation Contract](AGENTS.md#delegation-contract)
section of AGENTS.md.

Additional rules:

- Quality gates MUST be copy-pasteable shell commands. Use the project's actual commands
  from `docs/architecture/tech-stack.md` (not placeholders)
- Context bundle lists EXACT files with line ranges and justification for each
- Boundaries section lists at least 3 things explicitly OUT OF SCOPE
- Compact rules are pasted as inline text, never referenced by file path
- Line budget: see [Handoff Structure](AGENTS.md#delegation-contract) in AGENTS.md for the
  canonical 300-line rule. If exceeded, split the task

## Engineering Addenda

File: `docs/epics/{epic}-batch-{N}-refinement.md`. Created during Refine phase as the
output of the grooming/refinement ceremony.

Sections:

1. **Batch Info** --- batch number, scope, who refined it
2. **Prerequisites** --- checkbox list of items that must be resolved before implementation
3. **Dependencies** --- table: Package/Tool | Version | Target Project/File
4. **Implementation Order** --- numbered list of files/components in dependency order. When
   TDD is active, this becomes the test-first order
5. **Risks** --- subsections by severity (CRITICAL, HIGH, MEDIUM, LOW). Each risk states the
   problem, why it matters, and the mitigation
6. **Related Documents**

## Batch Progress

File: `docs/subtasks/{epic}/batch-{N}-progress.md`. Created when the first task in a
batch is delegated. Updated atomically before and after each delegation.

Sections:

1. **Batch Reference** --- link to batch plan
2. **Task Status** --- table: Task ID | Status (NOT_STARTED/IN_PROGRESS/BLOCKED/DONE/FAILED) | Worker | Evidence
3. **Active Wave** --- which wave is currently executing
4. **Blockers** --- any active blockers with context
5. **Last Updated** --- timestamp of last status change

This artifact solves crash recovery: a new session reads this file to know exactly
where execution stopped and which tasks need re-delegation.

## Deferred Backlog

File: `docs/deferred-backlog.md`

Persistent artifact capturing ideas and stories excluded from the current iteration with
rationale. Initialized during Capture, updated during Discover and Refine, consumed during
the next iteration's Capture phase.

| Column | Description |
| ------ | ----------- |
| ID | Sequential identifier (DF-001, DF-002) |
| Item | Story or feature title |
| Originating Phase | Phase where the deferral decision was made |
| Reason | Human-readable explanation |
| Category | `time` · `dependency` · `risk` · `low-value` · `scope-creep` |
| Reconsider When | Condition that would make this item a candidate again |
| Related Deps | Links to dependency DAG nodes (if category is `dependency`) |

Categories:

- `time` --- valuable but does not fit in this iteration
- `dependency` --- requires something that does not exist yet
- `risk` --- requires investigation or spike first
- `low-value` --- does not justify the cost now
- `scope-creep` --- not part of the original problem statement

When a blocking dependency is resolved, items with `dependency` category and matching
`Related Deps` become candidates for the next iteration automatically.
