# Internal smart-contract review: Nemesis and Pashov

**Status:** Required development workflow; no AdWeb3 contract audit has run yet because the contracts are not implemented.

**Requested by the owner:** Update and use the Nemesis and Pashov skills from the existing Capyfi audit-tool checkouts during smart-contract development, contract review and internal audit. Reports concern AdWeb3, not the Capyfi contracts. This is an internal security process; independent production security review before mainnet remains required.

## Tool sources verified locally

Shared tool root: `/Users/homer/development/audits/capyfi-smart-contracts/audit/tools`.

| Tool | Checkout | Entry point |
|------|----------|-------------|
| Nemesis | `nemesis-auditor` | `.claude/skills/nemesis-auditor/SKILL.md` |
| Nemesis dependencies | Same checkout | `.claude/skills/feynman-auditor/SKILL.md` and `.claude/skills/state-inconsistency-auditor/SKILL.md` |
| Pashov Solidity Auditor | `pashov-skills` | `solidity-auditor/SKILL.md`, `VERSION`, and its `references/` resources |

Inspection baseline on 2026-10-19, **not a claim of latest upstream revisions**:

- Nemesis origin: `https://github.com/0xiehnnkta/nemesis-auditor.git`; local HEAD `75cecc6dbd798f82ed8928d1a906078be9c575de`.
- Pashov origin: `https://github.com/pashov/skills.git`; local HEAD `92b0ea6786e6c210bc9a839e4813d50876c1c550`; local Solidity Auditor version `4.1`.
- Both working trees were clean at inspection. The Claude CLI exists at `/Users/homer/.local/bin/claude`; account authentication, capacity and a successful audit invocation have not been verified.

These are shared checkouts. Do not reset, clean, force-pull, overwrite instructions, or install AdWeb3 source into the Capyfi repository. No global skill installation is needed for this workflow.

## When both tools are required

- Phase 1: escrow funding, release, refund and authorization logic, plus deploy/setup scripts.
- Phase 4: distributor, budget/entitlement commitments and payout logic, rechecking their interaction with escrow.
- Phase 5: integrated pilot internal audit of the final candidate revision, including deployment roles and Go-to-contract settlement assumptions.
- Phase 7: publisher-related contract changes, then the complete feature candidate.
- Any other phase or later change that affects Solidity, deployment scripts, constructor/role configuration, upgrade assumptions, token behavior or backend transaction/funding semantics.

During development, use focused runs on changed production contracts **plus their callers, dependencies, inherited behavior and cross-contract state**, not just changed lines. Final pilot/full-feature runs cover all production contracts and deployment scripts, with supporting interfaces and token/library behavior supplied for context. Tests, mocks, dependencies and generated artifacts are not mislabeled as production contracts; include an explicitly named dependency if its behavior needs review.

## 1. Update safely and pin the audit cycle

At the start of a contract audit cycle, inspect each checkout's status, origin URL, current branch and tracking branch. Confirm it is still the intended upstream. For a clean, normally tracking checkout, the update sequence is:

```bash
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor status --short
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor remote get-url origin
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor fetch origin
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor merge --ff-only '@{upstream}'
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor rev-parse HEAD

git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills status --short
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills remote get-url origin
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills fetch origin
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills merge --ff-only '@{upstream}'
git -C /Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills rev-parse HEAD
```

Do not run a merge if the status is dirty, fetch fails, the checkout is detached, tracking is absent, the upstream is unexpected or the branch has diverged. Preserve local work and report the exact update condition. A separate clean checkout of the verified upstream can be used for AdWeb3 if needed; record its actual path and revision. Do not claim an updated-tool gate passed after an unsuccessful update. Continue unrelated implementation while resolving the tool issue.

Read the updated skill entry points and required referenced resources before running them; do not assume the inspected instructions remain unchanged. Record before/after SHAs, origin, branch, fetch/update result, Pashov VERSION, runtime and model settings. A clean fast-forward update within this approved workflow needs no extra routine permission.

Freeze the tool versions for a findings/fix/retest cycle. Re-check upstream at the next audit cycle and record any newly adopted versions. If tools change during a cycle, record the new baseline and rerun affected coverage; previous reports remain historical evidence.

## 2. Freeze the target and prepare isolated audit workspaces

Commit the candidate AdWeb3 changes in the local development repository and record the source SHA. A dirty or uncommitted target requires an explicit content manifest/hash and cannot be mislabeled as a reviewed commit. Include:

- `contracts/src/`, deployment/setup scripts and constructor/role settings.
- Foundry configuration, pinned dependency/compiler versions, relevant interfaces/libraries, and tests as supporting evidence.
- Current feature spec and implementation plan, allocation/escrow invariants, trust boundaries, intended chains/tokens and known open findings.
- Relevant Go settlement/indexing/payout code where its assumptions influence on-chain safety; Nemesis is language-agnostic, while Pashov's primary audited source is Solidity.

Create independent, content-identical snapshots for the two audit tools outside the live source tree. Resolve dependencies at their recorded revisions. Copy only the required skill directories into each snapshot's isolated `.claude/skills/` directory, preserving references and Nemesis's two dependency skills; never replace the project's existing `.claude/` directory. Check source hashes before and after every audit. Allow reports and skill metadata in the audit workspace, but audit agents must not modify contract/source files.

Store evidence under `audit/reports/<source-sha>/<cycle-id>/` in AdWeb3. Ignore disposable bundles/snapshots while retaining reports, manifests and regression-test references. `manifest.json` plus a readable `manifest.md` summary records source/tool SHAs, scope and exclusions, snapshot hashes, tool/runtime versions, model, pass settings, timestamps, completion state and output paths. No secrets, private keys, credentials or production environment files belong in bundles or reports.

## 3. Run Nemesis

Load the pinned full Nemesis skill and its Feynman/state-inconsistency dependencies. Its documented Claude Code entry point is `/nemesis`, invoked from the isolated AdWeb3 audit workspace. It is an agent skill, not a shell executable called `nemesis`.

Use the complete alternating workflow:

1. Recon and full Feynman pass.
2. Full state-inconsistency pass enriched by the first pass.
3. Targeted alternating passes until convergence, with the skill's maximum of six total passes.
4. Consolidation and verification, retaining unresolved suspects and coverage limits.

Preserve all pass outputs and `.audit/findings/nemesis-verified.md`. Record actual pass count, convergence status and any interrupted or incomplete analysis. Reaching a pass limit is not automatically convergence or a clean review.

Review focus: deposit versus earned balances; refund/release ordering; allocation conservation; immutable agreement/cohort/entitlement binding; unauthorized role paths; external-token calls and reentrancy; decimal/rounding behavior; payout retries; Go/chain reconciliation and reorg assumptions.

## 4. Run Pashov Solidity Auditor

Load the pinned `pashov-skills/solidity-auditor/SKILL.md` and its required reference resources. Invoke the `solidity-auditor` skill with explicit **`--loop 3 --file-output`** settings on the AdWeb3 scope. The explicit three-pass choice is part of this workflow; do not repeat its pass-count question at every gate. Keep the current runtime's inherited model unless the owner specifies otherwise.

Follow that installed version's discovery, specialist review, finding validation, deduplication and report assembly. Its current workflow uses twelve specialist reviewers per pass and memory for multi-pass scans. Use a runtime that supports the workflow; respect the host's concurrency limits. Do not replace the specialists with one generic reviewer and call that a complete Pashov run. Record any missing specialty, failed pass, runtime adaptation or partial coverage.

Include production deployment scripts; they set authorities and constructor values. Preserve `.solidity-auditor/runs/<stamp>/full-report.md`, scope/run metadata, the output copy, and the multi-pass memory ledger. The skill's own report format and severity/confidence rules govern its original report; any AdWeb3 summary must link back to the unmodified report.

Nemesis and Pashov should assess the same candidate independently before their findings are merged. Prior-cycle findings are valid retest context; do not feed one tool's first-pass report to the other as a substitute for independent discovery.

## 5. Triage, fix and retest

Create `findings.md` with a stable AdWeb3 ID, source-tool IDs, affected source revision/files, severity/confidence, root cause, exploit/trigger, impact, evidence, regression-test location and status. Deduplicate by root cause without losing either tool's original finding.

For each issue:

1. Validate the report against the actual code and intended trust model.
2. For a reproducible defect, write a failing Foundry regression/PoC in the implementation tree or a clearly isolated test harness. Capture its failure on the reviewed revision; reviewers remain read-only.
3. Implement the fix in the development tree; run unit, fuzz, invariant and relevant Go/chain integration suites.
4. Re-run both auditors against the new source revision and affected scope, then repeat until the gate closes. Final candidate runs cover the integrated production-contract scope.
5. Disputed or false-positive findings retain the reasoning, proof and reviewer resolution; they do not disappear from the record.

Map all confirmed exploitable issues to MUST-FIX. Classify other engineering findings using the feature-reviewer MUST-FIX/SHOULD-FIX/NOTE rules, preserving the original audit severity. Close with zero unresolved MUST-FIX/SHOULD-FIX and no unresolved confirmed Critical/High/Medium security issue. Low-severity findings must be fixed or documented with reviewed reasoning under the same gate; a severity label alone is not permission to ignore a defect.

## 6. Completion evidence

An audit gate is complete only when:

- [ ] Both tool updates succeeded, or a documented clean-upstream alternative was used; exact tool SHAs are pinned.
- [ ] Both requested skill workflows completed on the candidate; no silent coverage loss, failed-pass-as-clean or nonexistent output.
- [ ] Scope, exclusions, source hashes and runtime/pass settings are recorded.
- [ ] Both original reports and the consolidated finding ledger are retained.
- [ ] Confirmed issues have fixes, passing regression tests and retest evidence at the final source revision.
- [ ] Foundry unit/fuzz/invariant and affected Go/chain checks pass; final-candidate snapshot hashes match the implementation revision being approved.
- [ ] Independent review resolves all MUST-FIX/SHOULD-FIX findings and disputes.
- [ ] The implementation plan's phase and final review logs link the evidence and accurately state internal-audit scope/limitations.

If an auditor cannot run or a report is incomplete, record the gate as incomplete. It is never a clean result. Review and evidence collection do not submit transactions or deploy contracts.

## Future execution wrappers

Phase 1 will implement `make audit-tools-update`, `make audit-contracts PHASE=N` and `make verify-contract-audit PHASE=N` with documented implementations in `scripts/audit/`. These are future orchestration targets, not commands already available in this documentation-only repository. The first target performs the guarded update and writes tool provenance; the second prepares the immutable target and orchestrates both skill invocations using the supported agent runtime; the verifier checks evidence completeness and source/tool revision binding. A successful wrapper exit cannot itself prove absence of vulnerabilities; independent report triage and reviewer closure remain mandatory.
