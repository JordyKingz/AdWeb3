# Implementation Plan Template

Copy this structure to `docs/features/<feature-slug>/implementation-plan.md`.

Every phase block below is mandatory — the test-first list, the binary
definition of done, and the review-gate log in particular. The log is written
during Stage 4 execution, not at planning time, but the empty table must exist
in the plan so there is somewhere for the evidence to go.

---

# Implementation Plan: <Feature Name>

- **Slug:** `<feature-slug>`
- **Spec:** `docs/features/<feature-slug>/feature-spec.md`
- **Status:** Draft | Approved | In progress | Complete
- **Execution policy:** After plan approval, continue automatically through all phases; reconcile implementation with this plan and the spec after each phase and before final completion.
- **Test command:** `<exact command, from spec section 3>`

## Phase overview

| Phase | Deliverable | Stories | Complexity | Depends on | Status |
|-------|-------------|---------|------------|------------|--------|
| 1 | <one line> | US-01 | S / M / L | — | Not started |
| 2 | | US-02, US-03 | | Phase 1 | Not started |

> Ordering rule: dependency first, then risk. The highest-uncertainty work goes
> as early as its dependencies allow — finding out in Phase 5 that the approach
> does not work costs far more than finding out in Phase 1.
>
> Complexity is S/M/L, relative only. No hour estimates.

---

## Phase 1 — <Deliverable name>

**Goal:** <The deliverable in one sentence.>

**Hard-deliverable check:** <If this phase merged alone and nothing else ever
landed, what works and what can be demoed? If the honest answer is "nothing
yet," this is not a phase — reslice it.>

**User stories covered:** US-01 (AC-01.1, AC-01.2)

**Complexity:** S | M | L

**Depends on:** — | Phase N

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `<path>` | <why> |
| Modify | `<path>` | <what changes> |

### Tests to write first

> Write these before any implementation. Each maps to an acceptance criterion
> from the spec. A test with no criterion is either scope creep or a missing
> criterion — resolve it before starting.

| Test | Maps to | Asserts |
|------|---------|---------|
| `<test name / file>` | AC-01.1 | <the observable behaviour> |
| `<test name / file>` | AC-01.2 | <edge case or failure path> |

### Definition of done

> Binary only. Every box is objectively yes or no by inspection or by running
> something. If an item needs a judgement call to check, rewrite it.

- [ ] All tests listed above are written and passing
- [ ] Full suite passes (`<test command>`) — zero failures, zero skips
- [ ] <Specific observable behaviour, e.g. "POST /export with an invalid range returns 422 with the documented error body">
- [ ] <Specific artifact, e.g. "migration `<name>` applies and rolls back cleanly on a copy of the dev database">
- [ ] Code follows the conventions recorded in spec section 3
- [ ] Review gate closed with zero unresolved MUST-FIX and SHOULD-FIX

### Rollback

<How to undo this phase: revert scope, migration down-path, flag to disable,
data implications. "Revert the commit" is only sufficient when nothing stateful
changed — say so explicitly if that is the case.>

### Review-gate log

> Filled in during Stage 4. One row per review round. The phase is not done
> until a round records zero unresolved MUST-FIX and SHOULD-FIX. Disputed findings stay in the table with
> the reasoning, rather than disappearing.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|
| 1 | | | | | | Changes required |
| 2 | | 0 | 0 | | | Clean |

**Round detail**

- **Round 1 — <finding ref>:** <what the reviewer found> → <what was changed, or why it was disputed>

### Completion

Record requirement/acceptance evidence in the verification table, report progress,
and continue to the next phase without requesting confirmation.

`✅ Phase 1 complete — <deliverable> — <files> — [review rounds: N]`

---

## Phase 2 — <Deliverable name>

**Goal:**

**Hard-deliverable check:**

**User stories covered:**

**Complexity:** S | M | L

**Depends on:** Phase 1

### Files

| Action | Path | Purpose |
|--------|------|---------|

### Tests to write first

| Test | Maps to | Asserts |
|------|---------|---------|

### Definition of done

- [ ] All tests listed above are written and passing
- [ ] Full suite passes — zero failures, zero skips
- [ ] <specific>
- [ ] Review gate closed with zero unresolved MUST-FIX and SHOULD-FIX

### Rollback

### Review-gate log

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record requirement/acceptance evidence in the verification table, report progress,
and continue to the next phase without requesting confirmation.

`✅ Phase 2 complete — <deliverable> — <files> — [review rounds: N]`

---

## Story coverage check

> Every story in the spec lands in exactly one phase. Fill this in at planning
> time; a story with no phase will not get built.

| Story | Phase |
|-------|-------|
| US-01 | 1 |
| US-02 | 2 |

**Unassigned stories:** <must be "none" before the plan is approved>

## Deviations from the spec

> Appended during execution. When building reveals the spec was wrong, record
> the change here and update the spec — a plan that quietly drifts from its
> spec stops being a plan.

| Date | Phase | What changed | Why |
|------|-------|--------------|-----|

## Implementation verification

Update after each phase and reconcile the full table before marking Complete.
Every FR/NFR and every acceptance criterion needs implementation and passing
validation evidence. Missing evidence is unfinished work, not an implicit pass.

| Requirement | Story / acceptance criterion | Phase | Implementation evidence | Validation command / result | Status |
|-------------|------------------------------|-------|-------------------------|-----------------------------|--------|

## Final feature review

An independent reviewer checks the complete implementation against the current
spec and plan. Fix and re-review until no unresolved MUST-FIX or SHOULD-FIX
remain. Record disputes and their reviewer resolution.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

- [ ] Every requirement and acceptance criterion has implementation and passing validation evidence
- [ ] Every phase deliverable and definition-of-done item is satisfied
- [ ] Required integrated checks pass (record commands and results)
- [ ] Spec and plan reflect the implemented behavior without silently weakening approved requirements
- [ ] Independent final review is clean; remaining limitations are documented
