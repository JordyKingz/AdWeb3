---
name: feature-reviewer
description: Turns a feature request into a researched Feature Spec (requirements, gaps, impact, risks, epics, Given-When-Then stories), then on approval into a phase-based TDD implementation plan executed continuously through all approved phases with mandatory code-review subagent gates and verification against the spec and plan. Use whenever someone proposes a feature, enhancement, or change request, or asks to spec, plan, scope, or break down work. Not for bug fixes, one-line tweaks, pure refactors, or questions about how existing code works.
---

# Feature Reviewer

Take a feature request from vague ask to merged code without ever guessing. The
value of this skill is that research precedes writing, gaps are surfaced instead
of invented, and no phase is called done until an independent reviewer says so.

## The four stages

Run these in order. Do not skip a stage, merge two stages into one pass, or
start writing code because the request "seems small." Each stage exists because
the next one is unreliable without it: a plan built on an unresearched spec
produces phases that fight the existing architecture, and phases without a
review gate accumulate defects that surface three phases later.

| Stage | Output | Gate before continuing |
|-------|--------|------------------------|
| 1. Investigate | Notes (in context, not a file) | None — flows into Stage 2 |
| 2. Feature Spec | `docs/features/<feature-slug>/feature-spec.md` | **STOP.** User answers BLOCKING gaps |
| 3. Implementation Plan | `docs/features/<feature-slug>/implementation-plan.md` | **STOP.** User approves the plan |
| 4. Execution loop | Code + tests + review log | Review and verify each phase, then continue automatically until the approved plan is complete |

`<feature-slug>` is the feature name in kebab-case (`Bulk CSV Export` →
`bulk-csv-export`). Use the same slug in both documents.

---

## Stage 1 — Codebase investigation (read-only)

Read before you write. The goal is to be able to describe how this codebase
already does things, so the spec proposes work that fits rather than work that
collides.

**Do not write or edit any source file during this stage.** Reading, searching,
and listing only. If the urge to "just quickly add" something appears, that is
the signal that investigation is incomplete.

Find and read:

- The entry points and modules the feature would touch
- The data models, schemas, and migrations in that area
- The existing tests covering that area — these encode the real contract
- Config, feature flags, and environment handling
- Dependency manifests, to know what is already available before proposing anything new

Record, because Stages 2–4 all depend on it:

1. **Current architecture** — how the relevant slice is structured today
2. **Existing patterns and conventions** — naming, layering, error handling,
   logging, how modules talk to each other
3. **Test framework and how tests are run** — the exact command, plus how tests
   are organized and named
4. **Relevant dependencies** — what is available, what versions, what is
   deliberately avoided

If the codebase contradicts your instincts, the codebase wins. Consistency with
surrounding code is worth more than any individual preference, and a reviewer
will flag divergence anyway.

---

## Stage 2 — Feature Spec

Write to `docs/features/<feature-slug>/feature-spec.md` using
`feature-spec-template.md`. Read that template before writing; it
carries the exact section shapes, ID formats, and worked examples.

The nine sections, in this order:

1. **Summary** — the feature in three sentences, plus the problem it solves
2. **Requirements** — functional and non-functional, each numbered (`FR-01`,
   `NFR-01`) and phrased so it can be tested
3. **Current State** — what exists today, with concrete file paths from Stage 1.
   No paths here means Stage 1 was skipped
4. **Gaps** — what is missing, ambiguous, or undecided, split into:
   - `BLOCKING` — planning cannot start until answered
   - `NON-BLOCKING` — proceed on a stated assumption, flagged for confirmation

   **Never silently invent a requirement to fill a gap.** An invented
   requirement reads exactly like a real one and gets built. If something is
   unknown, name it as unknown. Recording an explicit assumption is fine;
   dressing it up as a decision is not.
5. **Impact Analysis** — files and modules changed, API/contract changes, data
   model and migration impact, backwards compatibility, performance, security,
   testing surface, affected teams or consumers
6. **Risks** — each with likelihood, blast radius, and mitigation
7. **Epics** — each with an ID (`EP-01`), a goal, and a scope boundary
8. **User Stories** — each with an ID (`US-01`), parent epic, `As a / I want /
   So that`, and Given-When-Then acceptance criteria. **Every requirement maps
   to at least one story.** Include the traceability table from the template and
   check it; an unmapped requirement is a requirement nobody will build.
9. **Out of Scope** — explicit non-goals, so scope creep has something to bounce off

### Then stop

Present the BLOCKING gaps as a numbered question list and wait. Do not begin
Stage 3, do not draft phases "in the meantime," do not answer the questions on
the user's behalf. If there are no BLOCKING gaps, say so explicitly and ask for
approval to plan.

```
Feature Spec written — docs/features/<feature-slug>/feature-spec.md

Before I can plan, I need answers on N blocking gaps:
1. <question>
2. <question>

I've also assumed the following (non-blocking, tell me if wrong):
- <assumption>
```

---

## Stage 3 — Phase-based implementation plan

Write to `docs/features/<feature-slug>/implementation-plan.md` using
`implementation-plan-template.md`.

Decompose the spec into sequential phases. **Every phase is a hard deliverable:
independently mergeable, demoable, and leaving the codebase green.** A phase
that only adds scaffolding, or wires half a layer, or "sets up the types for
later" is not a phase — it is an unfinished piece of another phase. The test is
simple: if this phase merged alone and nothing else ever landed, would the
codebase still be working and would there be something to show? If no, merge it
into an adjacent phase or reslice.

Each phase specifies:

- **Phase goal** — the deliverable in one sentence
- **User stories covered** — `US-xx` IDs from the spec
- **Files to create/modify** — concrete paths
- **Test cases to write first** — each mapped to a story's Given-When-Then criterion
- **Definition of done** — a binary checklist, every item objectively yes/no.
  "Error handling is solid" is not checkable; "invalid CSV returns 422 with the
  documented error body, covered by a test" is
- **Rollback note** — how to undo this phase, including any migration
- **Dependencies** — which prior phases must land first
- **Complexity** — S / M / L, relative only. Never estimate in hours; hour
  estimates invite schedule pressure onto work whose shape is still uncertain

Order phases by dependency first, then risk: **pull the highest-uncertainty work
early.** Discovering in Phase 5 that the chosen approach cannot work is far more
expensive than discovering it in Phase 1.

Then stop and get the plan approved before executing.

---

## Stage 4 — Phase execution loop

Once the user approves the implementation plan, execute every phase in order
without requesting approval to start or continue individual phases. Phase
boundaries are review and verification checkpoints, not user approval gates.
Continue until the entire approved plan is implemented and verified.

Before each phase, re-read the current feature spec and implementation plan.
Check the phase's requirements, stories, acceptance criteria, dependencies, and
definition of done against the implementation so far. Resolve in-scope gaps and
keep both documents current; do not silently remove or weaken requirements to
make the implementation appear complete.

For each phase, in order:

### 1. TDD red

Write the failing tests from the phase's test list — no implementation yet. Run
them. Confirm each fails **for the right reason**: a test failing on an import
error or a typo proves nothing about the behaviour it claims to cover. Read the
failure output rather than assuming.

### 2. TDD green

Write the minimum code that passes those tests. Then run the **full** suite, not
just the new tests — the fastest way to find a regression is to look for it now.

### 3. Refactor

Clean up with tests staying green. Align naming, structure, and error handling
with the conventions recorded in Stage 1. Re-run the suite after refactoring.

### 4. Review gate

Launch a code-review subagent. Hand it:

- The phase diff
- The phase's definition of done
- The current feature spec and implementation plan, including requirement and story mappings
- `code-review-checklist.md`

**The subagent reviews only and never edits.** A reviewer that fixes what it
finds destroys the independence that makes the gate worth having, and its
findings can no longer be triaged or disputed. Fixes are made by you, in the
main thread, after triage.

### 5. Triage and fix

Sort every finding into:

- `MUST-FIX` — correctness, security, data loss, broken contract, missing or
  false test coverage, definition-of-done item unmet
- `SHOULD-FIX` — convention violations, unclear naming, missing edge-case
  coverage, duplication
- `NOTE` — observations, future work, matters of taste

Fix all `MUST-FIX` and all `SHOULD-FIX`. Re-run the tests. Then re-launch the
review subagent on the fixes.

If a finding looks wrong, say so with reasoning and mark it disputed in the log
rather than silently ignoring it.

### 6. Repeat until clean

Loop steps 4–5 until **zero unresolved `MUST-FIX` and `SHOULD-FIX` remain**.
Disputed findings require a documented reviewer resolution before closure. Log every review round and its
outcome in the plan document's review-gate log for that phase — round number,
counts by severity, what was fixed, and the verdict. The log is the evidence
that the gate actually ran.

### 7. Verify, report, and continue

Reconcile the phase implementation with the feature spec and implementation
plan: every covered requirement and acceptance criterion must have concrete
implementation and passing validation evidence. Record this in the plan's
verification table, update phase status, and resolve omissions through the same
test and review loop before marking the phase complete.

Report progress in commentary:

```
✅ Phase N complete — [deliverable] — [files] — [review rounds: N]
```

Immediately start the next phase. Do not end the turn or ask for confirmation
just because a phase is complete.

### 8. Verify the complete feature

After the last phase, iterate over the entire feature spec and implementation
plan again. Check every FR/NFR, story, acceptance criterion, phase deliverable,
and definition-of-done item against the actual code and validation results.
Run the required integrated checks and have an independent review subagent
review the complete feature against both documents. Record its review rounds
in the plan's final review log. Fix in-scope gaps, re-run affected checks, and
repeat review and reconciliation until there are no unresolved MUST-FIX or
SHOULD-FIX findings and no missing acceptance evidence.

Mark the plan Complete only then. Give a final report with the delivered
behavior, verification evidence, and any remaining non-blocking limitations.

### When to pause

Pause only for a genuine blocker that cannot be resolved within the approved
scope, a material product/scope change requiring the user's decision, a new
action requiring authorization, or an explicit user request to stop. State the
specific blocker and decision needed; continue independent authorized work if
possible. Plan approval does not authorize unapproved destructive operations,
production deployment, or other external actions. Routine fixes, document
reconciliation, reviews, and transitions between approved phases need no new
approval.

---

## Rules that hold across every stage

**Follow existing codebase conventions over your own preferences.** If the
project uses a pattern you would not have chosen, use the project's pattern.
Raise the disagreement as a `NOTE`, separately from the feature work.

**Never mark a phase done with failing or skipped tests.** A skipped test is an
untested requirement wearing a green checkmark. If a test cannot pass yet, the
phase is not sliced correctly — say so and reslice rather than skipping.

**Never invent a requirement, a file path, or an API you have not verified.**
Check it in the codebase or list it as a gap.

**Keep the documents live.** When execution reveals the spec was wrong, update
the spec and note the change — a plan that drifts from its spec stops being a
plan.

---

## Reference files

Read each one at the point it is needed, rather than all up front:

- `feature-spec-template.md` — read before writing Stage 2
- `implementation-plan-template.md` — read before writing Stage 3
- `code-review-checklist.md` — pass to the review subagent in Stage 4
