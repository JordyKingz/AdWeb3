# Feature Spec Template

Copy this structure to `docs/features/<feature-slug>/feature-spec.md`. Keep the
nine sections in this order and keep the ID formats exactly as shown — Stage 3
and Stage 4 both reference these IDs, and a renamed or missing ID breaks the
traceability chain.

Replace every `<...>` placeholder. Delete the guidance blockquotes; keep the
headings even when a section is thin, and write "None identified" rather than
dropping a section, so a reader can tell the difference between "considered and
empty" and "never considered."

---

# Feature Spec: <Feature Name>

- **Slug:** `<feature-slug>`
- **Status:** Draft | Awaiting answers | Approved
- **Author:** <name>
- **Date:** <YYYY-MM-DD>

## 1. Summary

> Three sentences on what the feature is, then the problem it solves. If it
> takes more than three sentences, the feature is probably two features.

<What it is, in three sentences.>

**Problem it solves:** <The user or business pain that exists today.>

## 2. Requirements

> Number every requirement and phrase it so a test could prove it. "Fast" is not
> testable; "responds in under 300ms at p95 for 10k rows" is. Each requirement
> must appear in the traceability table in section 8.

### Functional

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01 | <System does X when Y.> | Must / Should / Could |
| FR-02 | | |

### Non-functional

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-01 | <Measurable constraint: performance, security, accessibility, observability.> | Must / Should / Could |
| NFR-02 | | |

## 3. Current State

> What exists in this codebase today, with real file paths verified in Stage 1.
> No invented paths. If this section has no paths, investigation was skipped.

- **Architecture:** <How the relevant slice works today.>
- **Key files:**
  - `<path/to/file>` — <what it does and why it matters here>
- **Existing patterns to follow:** <naming, layering, error handling, logging>
- **Test framework:** <framework> — run with `<exact command>`; tests live in `<path>`
- **Relevant dependencies:** <name@version> — <what it is used for>

## 4. Gaps

> Everything missing, ambiguous, or undecided. Never fill a gap by inventing a
> requirement — an invented requirement is indistinguishable from a real one
> once it is written down, and it will get built.

### BLOCKING

> Planning cannot start until these are answered. These become the numbered
> question list presented to the user at the end of Stage 2.

| ID | Question | Why it blocks |
|----|----------|---------------|
| GAP-B01 | <Question.> | <What cannot be planned without the answer.> |

### NON-BLOCKING

> Proceed on a stated assumption, flagged for confirmation. State the assumption
> plainly and say what changes if it is wrong.

| ID | Assumption | Impact if wrong |
|----|------------|-----------------|
| GAP-N01 | <Assuming X.> | <What would need to change.> |

## 5. Impact Analysis

| Dimension | Impact |
|-----------|--------|
| Files and modules changed | <paths> |
| API / contract changes | <new, changed, or removed endpoints, events, signatures; breaking or not> |
| Data model and migrations | <schema changes, migration needed, reversible or not, backfill required> |
| Backwards compatibility | <what breaks for existing callers, clients, or stored data> |
| Performance | <expected effect on latency, memory, query load> |
| Security | <authn/authz, input validation, data exposure, secrets, new attack surface> |
| Testing surface | <what needs new tests, what existing tests change> |
| Affected teams / consumers | <who else depends on this and needs telling> |

## 6. Risks

| ID | Risk | Likelihood | Blast radius | Mitigation |
|----|------|-----------|--------------|------------|
| RISK-01 | <What could go wrong.> | Low / Med / High | <Who and what is affected if it does.> | <Concrete action, not "be careful".> |

## 7. Epics

> Each epic needs a scope boundary as well as a goal. The boundary is what keeps
> the epic from quietly absorbing the next one.

### EP-01: <Epic name>

- **Goal:** <Outcome this epic delivers.>
- **Scope boundary:** <What is inside, and explicitly what is not.>

### EP-02: <Epic name>

- **Goal:**
- **Scope boundary:**

## 8. User Stories

> Every story needs an ID, a parent epic, the As/I want/So that form, and
> Given-When-Then acceptance criteria. The criteria become the test list in
> Stage 3 — vague criteria produce vague tests, so write them as concrete
> observable behaviour.

### US-01: <Story name>

- **Epic:** EP-01
- **As a** <role>
- **I want** <capability>
- **So that** <benefit>

**Acceptance criteria**

- **AC-01.1**
  - **Given** <starting state>
  - **When** <action>
  - **Then** <observable outcome>
- **AC-01.2**
  - **Given** <edge case or failure state>
  - **When** <action>
  - **Then** <observable outcome>

### US-02: <Story name>

- **Epic:** EP-01
- **As a** / **I want** / **So that**

**Acceptance criteria**

- **AC-02.1** — Given / When / Then

### Traceability

> Fill this in and check it. Every requirement from section 2 needs at least one
> story. An unmapped requirement is a requirement nobody will build.

| Requirement | Covered by |
|-------------|------------|
| FR-01 | US-01, US-02 |
| FR-02 | US-03 |
| NFR-01 | US-02 |

**Unmapped requirements:** <must be "none" before this spec is approved>

## 9. Out of Scope

> Explicit non-goals. This is what scope creep bounces off, so name the things
> someone will plausibly ask for.

- <Not doing X — reason.>
- <Not doing Y — deferred to <where>.>
