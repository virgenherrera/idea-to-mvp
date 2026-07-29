# Contributing to idea-to-mvp

This framework improves through empirical testing. If you use it on a real project
and find failure modes, report them through your platform's issue tracker.

## Reporting a Framework Defect

Include these fields:

- **What the agent did wrong**: describe the incorrect behavior observed
- **Framework rule violated or missing**: which section of AGENTS.md was violated,
  or what rule should exist but does not
- **Evidence**: logs, screenshots, agent output, or postmortem excerpts
- **Severity**: CRITICAL (survived to "done"), HIGH (caught late), MEDIUM (caught
  at next gate), LOW (agent self-corrected)
- **Pipeline phase**: where in the 8-phase pipeline the defect occurred
- **Proposed amendment** (optional): suggested change to AGENTS.md

## Proposing an Enhancement

Include these fields:

- **Problem or gap**: what failure mode or gap does this address
- **Proposed change**: the rule, gate, or protocol to add or modify
- **Consolidation candidate**: which existing rule does this subsume or render
  redundant? If none, explain why this is genuinely net-new
- **Evidence** (optional): empirical testing results or prior art

## Change Checklist

When modifying the framework, verify:

- [ ] AGENTS.md changes preserve agent/technology/platform agnosticism
- [ ] No platform-specific references introduced
- [ ] Cross-references (`[Section](#anchor)`) verified
- [ ] Lightweight Mode invariant gates reviewed (if new gate added)
- [ ] Consolidation review performed (META-CONSOLIDATION)
