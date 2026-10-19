# Feature Request

A Claude Code skill that takes a feature request from a vague ask to merged code
without guessing along the way.

Research happens before anything is written, unknowns are surfaced as gaps
instead of being invented, and no phase is called done until an independent
review subagent says so.

---

## What it does

Four stages, run in order. Each one exists because the next is unreliable
without it — a plan built on an unresearched spec produces phases that fight the
existing architecture, and phases without a review gate accumulate defects that
only surface later.

| Stage | Output | Gate before continuing |
|-------|--------|------------------------|
| 1. Investigate | Notes, held in context | None — flows into Stage 2 |
| 2. Feature Spec | `docs/features/<slug>/feature-spec.md` | **Stops.** You answer the BLOCKING gaps |
| 3. Implementation Plan | `docs/features/<slug>/implementation-plan.md` | **Stops.** You approve the plan |
| 4. Execution loop | Code, tests, review log | Continues automatically through all approved phases, with review and verification gates |

**Stage 1 — Investigation.** Read-only. Maps the current architecture, existing
conventions, the test framework and its exact run command, and the available
dependencies. No source file is edited during this stage.

**Stage 2 — Feature Spec.** Nine sections: summary, numbered requirements
(`FR-01` / `NFR-01`), current state with real file paths, gaps split into
`BLOCKING` and `NON-BLOCKING`, impact analysis, risks, epics (`EP-01`), user
stories (`US-01`) with Given-When-Then criteria, and explicit non-goals. Every
requirement maps to at least one story, checked against a traceability table.

**Stage 3 — Implementation Plan.** Decomposes the spec into phases where each
phase is independently mergeable, demoable, and leaves the codebase green.
Phases are ordered by dependency first, then by risk, so the highest-uncertainty
work lands early. Complexity is S/M/L only — never hours.

**Stage 4 — Execution loop.** Per phase: TDD red, TDD green against the full
suite, refactor, then a mandatory code-review subagent gate. Findings are triaged
`MUST-FIX` / `SHOULD-FIX` / `NOTE`, fixed in the main thread, and the reviewer is
re-run until zero unresolved `MUST-FIX` and `SHOULD-FIX` remain. Every round is logged as evidence the gate
actually ran. It checks implementation against the spec and plan, then immediately
continues to the next phase. A final independent review and requirements-to-evidence
check must pass before the full plan is marked complete. Genuine blockers, material
scope decisions, and actions requiring new authorization still pause execution.

## When it triggers

Use it when someone proposes a feature, an enhancement, or a change request, or
asks to spec, plan, scope, or break down work.

It deliberately does **not** apply to bug fixes, one-line tweaks, pure
refactors, or questions about how existing code already works.

## What you get on disk

```
docs/features/<feature-slug>/
├── feature-spec.md           # Stage 2
└── implementation-plan.md    # Stage 3, plus the review-gate log per phase
```

`<feature-slug>` is the feature name in kebab-case — `Bulk CSV Export` becomes
`bulk-csv-export` — and stays the same across both documents.

---

## Claude Installation

Claude Code discovers skills in two places. Pick one.

### Personal — available in every project

```bash
git clone --depth 1 https://github.com/JordyKingz/jordykingz.git /tmp/jdk-skills
mkdir -p ~/.claude/skills
cp -R /tmp/jdk-skills/skills/feature-request ~/.claude/skills/
rm -rf /tmp/jdk-skills
```

### Project-scoped — committed alongside the repo it serves

Run from the root of the target repository:

```bash
git clone --depth 1 https://github.com/JordyKingz/jordykingz.git /tmp/jdk-skills
mkdir -p .claude/skills
cp -R /tmp/jdk-skills/skills/feature-request .claude/skills/
rm -rf /tmp/jdk-skills
```

### Verify

```bash
ls ~/.claude/skills/feature-request
```

You should see `SKILL.md` alongside the three reference files. Start a new
Claude Code session — skills are read at startup — and confirm it is listed:

```
/skills
```

## Codex Installation

Codex discovers skills in two places. Pick one.

### Personal — available in every project

```bash
git clone --depth 1 https://github.com/JordyKingz/jordykingz.git /tmp/jdk-skills
mkdir -p ~/.agents/skills
cp -R /tmp/jdk-skills/skills/feature-request ~/.agents/skills/
rm -rf /tmp/jdk-skills
```

### Project-scoped — committed alongside the repo it serves

Run from the root of the target repository:

```bash
git clone --depth 1 https://github.com/JordyKingz/jordykingz.git /tmp/jdk-skills
mkdir -p .agents/skills
cp -R /tmp/jdk-skills/skills/feature-request .agents/skills/
rm -rf /tmp/jdk-skills
```

### Verify

```bash
ls ~/.agents/skills/feature-request
```

You should see `SKILL.md` alongside the three reference files. Start a new
Codex session — skills are read at startup — and confirm it is listed:

```
/feature-reviewer 
```

## Usage

Invoke it by name, or just describe the work and let it trigger:

```
Spec out bulk CSV export for the reporting module.
```

```
I want to add SSO. Scope it before we build anything.
```

Expect decision gates after the spec (blocking questions) and after the plan
(approval). Once approved, all phases run sequentially without further phase
approvals, with reviews and checks against the spec and plan throughout.

## Files

| File | Purpose |
|------|---------|
| `SKILL.md` | The skill itself — stages, gates, and the rules that hold across all of them |
| `feature-spec-template.md` | Section shapes, ID formats and worked examples for Stage 2 |
| `implementation-plan-template.md` | Phase block structure and the review-gate log table for Stage 3 |
| `code-review-checklist.md` | Handed to the review subagent at each Stage 4 gate, with the severity triage rules |

## Requirements

- Claude Code or Codex
- A repository with a working test command — Stage 4 is test-first and runs the
  full suite at every phase
- Nothing else; the skill has no dependencies of its own

## Design notes

Three rules do most of the work, and they are worth knowing before you use it:

- **Never invent a requirement to fill a gap.** An invented requirement reads
  exactly like a real one and gets built. Unknowns are named as unknowns;
  assumptions are recorded as assumptions.
- **The review subagent reviews only, never edits.** A reviewer that fixes what
  it finds destroys the independence that makes the gate worth having, and
  leaves nothing to triage or dispute.
- **A skipped test is an untested requirement wearing a green checkmark.** If a
  test cannot pass yet, the phase is sliced wrong — reslice rather than skip.

## License

See the repository root.
