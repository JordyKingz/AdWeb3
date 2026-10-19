# AdWeb3 Implementation Plan

> **For agentic workers:** Use `superpowers:executing-plans` for implementation, with the feature-reviewer skill's independent review after every phase and final reconciliation. Read this plan and the spec before each phase. Approval authorizes continuous work through the approved scope; do not request approval at routine phase boundaries.

- **Slug:** `adweb3-platform`
- **Spec:** [feature-spec.md](feature-spec.md)
- **Status:** In progress — owner approved pilot Phases 1–5 with reviews/internal audits; publisher beta remains excluded
- **Date:** 2026-10-19
- **Execution scope:** Owner approved pilot Phases 1–5 on 2026-10-19, including reviews and internal audits. Phases 6–7 describe the separately budgeted publisher beta and require explicit inclusion in the approved scope. Completing the pilot does not complete the full publisher feature.
- **Execution policy:** TDD → implementation → refactor → independent review → fix/re-review → spec/plan reconciliation → next approved phase. Zero unresolved MUST-FIX or SHOULD-FIX findings before proceeding. Pause only for a genuine dependency/blocker, material scope decision, new external authorization, or explicit user stop.

**Goal:** Demonstrate one Android game receiving funded campaign ads and earning validated points that settle into test-token payouts, then extend the same SDK/backend to approved third-party publishers.

**Architecture:** One modular Go backend with API and worker processes, PostgreSQL ledger and transactional outbox, object storage/CDN for approved creatives, and bounded EVM escrow/distributor contracts. A Kotlin Android SDK separates ad delivery from optional rewards; one small Android game is its first consumer. Game progress remains local/offline; money, ad eligibility, and settlement authority remain server-side.

**Tech stack — proposed implementation choices:** Go standard HTTP routing, PostgreSQL with pgx, versioned SQL migrations, OpenAPI, Solidity with Foundry (Forge/Anvil/Cast), Kotlin Android library modules and a Compose sample/game client. Use server-rendered Go templates for minimal advertiser/admin screens and the wallet-signing page; no separate SPA or microservice estate. Pin compatible stable toolchain/dependency versions during Phase 1/2 bootstrap and record them; this plan does not pretend an uninstalled version has been validated.

## Confirmed constraints and open decisions

**Verification amendment accepted 2026-10-19:** Owner has a pending Play Console account and is waiting for the physical Android device. Use one emulator for the current pilot development/verification cycle. Physical-device validation and genuine Play Integrity evidence are deferred to a follow-up when available; do not label simulation as real-device/integrity evidence or enable unverified monetary activity. Emulator/device clauses below refer to this amendment where a physical check cannot yet run.

- Android first; iOS and cross-engine wrappers are deferred.
- Free publisher SDK access; 20% AdWeb3 commission and 80% publisher allocation for ordinary partner ads, applied once to accepted earned net revenue.
- Owned-game allocation: 30% players / 60% operations / 10% contingency. Partner reserves come from the platform's 20%, not an additional charge.
- Points plus funded USDT is the intended product; the implementation scope here uses test tokens only. No tradable platform token, wagering, paid earning boosts, or guaranteed income.
- Offline play does not earn cash entitlement. No new paid ads offline; expired decisions cannot bill.
- Refundable advertiser money, player liabilities, publisher liabilities, operations, and reserves remain segregated and backed. Integer amounts only; do not use floating point for balances or pricing.
- Existing Candy Blocks Swift code is reference material. No Android game, SDK, backend, contract, migration, or test harness currently exists in this workspace. No Git metadata was found. Paths below marked Create are proposed new paths, not discovered implementation.
- No source modifications in the external iOS portfolio are required by this plan. Do not copy assets without confirming rights. Keep grant spending separate from ad-revenue allocations.

The owner requested this draft while spec GAP-B01–B04 remain open. The following assumptions make the draft concrete; they are **not** silently recorded as user decisions. Plan approval must accept or revise the relevant assumptions before dependent implementation begins.

| Decision | Draft basis | Boundary / when required |
|----------|-------------|---------------------------|
| Game | Candy Blocks candidate, one campaign-style level with bounded rules and Home/results; reuse verified rules fixtures, not a complete 30-level/Endless port | Confirm chosen level/rule subset, reference licence, and asset rights before Phase 2. Alternative game changes Phase 2/3 scope. |
| Chain | Local EVM first, then Avalanche Fuji test token; intended later mainnet is Avalanche C-Chain native USDt | Local chain 31337 and Fuji 43113 only; registry rejects other chains/tokens. Owner must accept network before Phase 1 execution. Never label the demo token real USDt. |
| Campaign | Proposed 5 test-token units CPM, 25-unit prepaid cap, 168-hour window, actual delivery and unused-budget refund | Fixture terms, not commercial offers. Real price, tax/deduction rules and dispute window remain GAP-B02. |
| Audience/market | Invited adult testers; no paid public launch or claimed DAU | Actual geography/age eligibility, operating costs and acquisition assumptions remain GAP-B01/B03; no real-money activation until resolved. |
| Grant capacity | USD 7,500 and eight weeks remain provisional | This plan is not an estimate validating those numbers. Re-cost Android work before submitting the grant. |
| Publisher rewards | Ordinary ads first; optional 30/50/20 is a disabled test fixture only | 20% platform fee is confirmed; optional player share and partner terms must be approved before enabling reward inventory. |

**Proposed test-profile constants:** 30-second decision/reward lease renewed every 15 seconds; 120-second bounded retry window for events from previously valid online periods; static card at least 50% visible for 1 continuous foreground second; one accepted view per decision; at most 6 eligible views per player per UTC day and one per campaign per player per 10 minutes. One validated eligible view supplies one point, never clicks or raw minutes. These are deterministic test assumptions, not market standards or production policy. Bind each rule to an immutable version. Freeze a cohort only after the retry/adjustment windows close; late arrival cannot rewrite finalized entitlements. Changed values require versioned tests and terms.

Use a 24-hour test-profile dispute hold after campaign close, daily processing, and a 1-test-token payout threshold. Advance an injected clock in tests instead of waiting days. Retain sub-threshold balances. Zero-weight allocations stay reserved for the same campaign/app pending disclosed disposition; do not move them to unrelated players. No production values are implicitly approved by these test fixtures.

## Repository and interface boundaries

All paths below are relative to `/Users/homer/development/evm/AdWeb3`. Reference source, read only: `/Volumes/Dock 4TB/Development - Dock/hobby/personal/apps/candy_tetris/Packages/CandyBlocksCore/Sources/CandyBlocksCore/Hybrid/GameEngine.swift` and adjacent rules/fixtures. Its `start(level:mode:seed:)` and tick-indexed `advance(state:commands:)` define the candidate reference behavior; fixture extraction must record the source revision and supported subset.

| Proposed area | Responsibility |
|---------------|----------------|
| `cmd/api/main.go`, `cmd/worker/main.go` | Separate processes from one Go module; dependency wiring and graceful shutdown |
| `internal/{identity,campaigns,delivery,ledger,settlement,chain,payouts,publishers,operations}/` | Domain ownership; no HTTP handler may bypass ledger/tenant authorization |
| `api/openapi.yaml`, `api/events.schema.json` | Versioned transport contract, identity audiences, error codes, idempotency and event schemas |
| `migrations/`, `internal/store/` | PostgreSQL schema, transactions, migration harness and outbox |
| `contracts/src/`, `contracts/test/` | Escrow/distributor contracts and isolated test token; Solidity tests/invariants |
| `android/sdk-core/`, `android/sdk-rewards/` | Core ads/identity versus optional wallet/reward module; no end-user wallet requirement for ordinary ads |
| `android/game-core/`, `android/pilot/`, `android/sample/` | Deterministic rules; playable first app; later ordinary-ad publisher example |
| `web/templates/`, `web/static/` | Minimal protected advertiser/operator views and bound wallet-auth page |
| `tests/integration/`, `tests/e2e/`, `tests/load/`, `docs/runbooks/` | Real database/chain tests, end-to-end scenarios, performance evidence and operator procedures |

**Transport decisions to encode in Phase 1/2:** UUID resource IDs, UTC timestamps, integer minor-unit amounts serialized as decimal strings, basis-point allocation integers summing to 10,000, explicit chain ID/token address, and immutable policy versions. `api/openapi.yaml` owns every HTTP request/response type named below. Domain-only structs such as `Journal`, `ChainEvent`, `Epoch` and `TxIntent` live with their owning Go modules; use explicit schema conversion rather than transport types inside accounting. An idempotency key is scoped to actor and operation; the same key with a different body is a conflict. Never accept tenant, price, allocation, or wallet destination as authoritative merely because a client sent it.

- Campaign API: `POST /v1/orders` → `Order`; `GET /v1/orders/{id}`; `POST /v1/orders/{id}/refunds` → `RefundIntent`; `GET /v1/campaigns/{id}/report`. Actor scope is authenticated advertiser/operator.
- Identity: `POST /v1/installations/challenges` → `Challenge`; `POST /v1/installations` → `Session`; `POST /v1/sessions/refresh`; `POST /v1/wallet-links/challenges` → `WalletChallenge`; `POST /v1/wallet-links/verify`; bound one-time app completion endpoint. Credentials never appear in redirect query parameters or logs.
- Delivery: `POST /v1/ad-decisions` → `AdDecision or NoFill`; `POST /v1/events:batch` → per-event `accepted | duplicate | rejected | pending`; `POST /v1/progress:sync` is independent of monetary evidence; `POST /v1/reward-sessions` issues the bounded online lease.
- Rewards: `GET /v1/rewards/balance` distinguishes pending points, finalized amount, below-threshold amount, submitted payouts, confirmed payouts, and rejections. It never reports a transaction submission as payment confirmation.
- Publisher beta: management routes under `/v1/publishers/{id}/` for apps, placements, statements, and payout setup; organization membership is verified on each operation, not trusted from route IDs.

**Contract trust boundary:** The backend is trusted to validate delivery; contracts enforce amounts, roles, budgets and replay protection, not human attention. The proposed escrow interface is `fundOrder(bytes32 orderId,uint256 amount)`, `releaseEarned(bytes32 settlementId,bytes32 orderId,uint256 amount,bytes32 allocationHash)`, and `refundOrder(bytes32 refundId,bytes32 orderId,address recipient,uint256 amount)`. Register order token, payer, purchase/allocation hash and amount limits before funding. A privileged role cannot silently alter those terms or refund to another payer. Check received-token balance delta, supported token behavior, and all role restrictions. No unbounded client-supplied loop.

`RewardDistributor` uses `registerBudget(bytes32 epochId,uint256 amount)` callable only by escrow after an atomic funding transfer, and `commitEntitlement(bytes32 epochId,bytes32 entitlementId,address recipient,uint256 amount)` callable by the restricted settlement role. Registration is unique and checks actual backing; total committed amounts cannot exceed the budget. Entitlement ID, epoch, recipient and amount are immutable. It then exposes `pay(bytes32 payoutId,bytes32 entitlementId,address recipient,uint256 amount)` with on-chain recipient/amount binding and paid-ID checks. The contract rejects uncommitted entitlements, overspending and changed destinations. Encode these signatures and role restrictions in the ABI before generating bindings; every committed entitlement retains its epoch provenance. Reconcile a durable intent → submitted → confirmed state machine across ledger and chain; only confirmed backing becomes spendable. A failed recipient is retried individually so it cannot block unrelated payouts. Production privileged-key governance and independent audit are separate launch requirements; use separate isolated test signers now.

## Verification commands and execution rules

No commands below have passed yet: the files/toolchains are future deliverables. Phase 1 creates the documented wrappers and pins tools. Contract-bearing phase completion additionally requires both internal audit workflows; `make verify` alone cannot establish that an agent audit occurred. From the repository root:

| Command | Contract |
|---------|----------|
| `go test ./...` and `go test -race ./...` | All Go unit tests; race pass with deterministic clocks/IDs |
| `go test -tags=integration ./tests/integration/...` | Real disposable PostgreSQL/local-EVM migrations, concurrency and crash recovery; missing services fail rather than skip |
| `make audit-tools-update` | Guarded fast-forward update of both shared audit tools; pin provenance for the cycle (future Phase 1 wrapper) |
| `make audit-contracts PHASE=N` | Prepare identical immutable AdWeb3 snapshots and run both installed skills with the settings above; archive reports |
| `make verify-contract-audit PHASE=N` | Validate report completeness, source/tool bindings, triage and retest references; independent reviewer still judges closure |
| `forge test --root contracts` | Contract unit, fuzz and invariant suites with deterministic repro seeds retained |
| `./android/gradlew -p android testDebugUnitTest lintDebug` | Android library/game JVM tests and lint after Android modules exist |
| `ANDROID_SERIAL=<selected-device> ./android/gradlew -p android connectedDebugAndroidTest` | Exactly one selected physical Android test device; record device/OS/build. Missing device is unmet evidence, not a pass |
| `make verify PHASE=N` | All suites introduced through N, including real integration tests; no placeholder/empty pass or test skipping |
| `make verify-device` | One selected emulator now; physical-device and genuine Play Integrity evidence tracked as owner-approved follow-up |
| `make verify-e2e` | Full local funded-ad-to-payout flow, migrations, rollback/restore rehearsal and restart recovery from Phase 4 onward |
| `make verify-load` | Phase 5: 100 ad decisions/s, p95 <=250ms, recorded dataset/hardware, excluding creative downloads |
| `make verify-fuji` | Phase 5: test-token deployment/funding/refund/payout on Fuji with recorded addresses, revision and RPC/finality evidence |

Target suites grow with the implementation; absence of future-phase code is not a skipped test. Once introduced, a suite remains required. Real Play Integrity/device checks require an applicant-controlled package/signing setup, Play project and credentials; unit-test fakes are permitted only in isolated tests, never as fallback for network-facing eligibility. Missing account/device dependencies can block a phase's completion, but independent backend work may continue without claiming it finished. Do not upload/publish a build, submit a grant, or launch mainnet merely because this plan is approved; external release actions use the applicable authorization.

Each task below uses this checklist: **(1)** write named failing assertions; **(2)** run their focused suite and confirm a behavioral failure, not a broken import; **(3)** implement the listed interface; **(4)** run focused checks and `make verify PHASE=N`; **(5)** review the diff and commit the verified task in the initialized local repository. Select a repository-safe module/package name at bootstrap; public release namespace is an owner-controlled publishing choice. Never commit secrets or device signing material.

## Required Nemesis and Pashov internal audit gates

Owner-requested workflow: [contract-audit runbook](../../runbooks/contract-audit.md). Source tools are `/Users/homer/development/audits/capyfi-smart-contracts/audit/tools/nemesis-auditor` and `/Users/homer/development/audits/capyfi-smart-contracts/audit/tools/pashov-skills`. Update clean tracking checkouts with fetch/fast-forward only before each audit cycle; preserve local changes and pin actual tool SHAs. Read updated skill instructions, prepare isolated snapshots of the same AdWeb3 candidate, and run Nemesis's full alternating workflow plus Pashov `solidity-auditor --loop 3 --file-output` using a supported agent runtime. Pass count is specified here; no repeated routine question is required.

Keep both original reports, scope/runtime/tool/source provenance, findings triage, Foundry regressions, fixes and both-tool retest results under `audit/reports/<source-sha>/<cycle-id>/`. Reports target AdWeb3, never the unrelated Capyfi source tree. Auditors remain read-only with respect to production source. No contract audit has run yet.

Required at Phases 1, 4, 5 and 7, and any other phase whose changes affect contracts, deploy scripts, roles, token semantics or Go/contract money flows. Each gate requires zero unresolved MUST-FIX/SHOULD-FIX and reviewed disposition of security findings. Partial runs, absent specialist coverage or a six-pass Nemesis limit without convergence remain explicit incomplete evidence, not a clean result. Pashov's tool/runtime requirements must be satisfied rather than silently replacing its specialists with a generic review.

These reviews complement the existing per-phase independent review. The final pre-mainnet independent security review remains outside this internal-audit assurance claim. An approval of this plan includes the safe tool-update/review cycle without new approval at each iteration.

## Review focus

1. Two API instances race to consume the last budget units: Phase 3 must reserve/settle under database transactions without overselling.
2. Server crashes between a chain submission and saving its hash: Phase 4 must reconcile durable nonce/intent and paid IDs without duplicate payment.
3. Android process dies, changes wall clock, or returns from a wallet with an injected redirect: Phases 2–3 preserve offline saves and reject forged monetary authority.
4. Fraud adjustment or late clearance occurs after the activity day: Phases 3–4 retain original cohort weights and never take unrelated refund/reward backing.
5. A legitimate publisher substitutes another tenant's IDs or enables rewards client-side: Phases 6–7 reject it at authorization, delivery and settlement boundaries.

## Phase overview

| Phase | Deliverable | Story completion owner | Complexity | Depends on | Status |
|-------|-------------|-----------------------|------------|------------|--------|
| 1 | Fund and refund a test campaign with protected ledger backing | US-01, US-08 | L | Accepted testnet/network assumptions | In progress |
| 2 | Play an Android game and establish verified installation/wallet sessions | US-04, US-05 | L | 1; selected game slice and device/account access | Not started |
| 3 | Serve and measure valid expiring ads in that game; safe offline transitions | US-02, US-03, US-06, US-15 | L | 1–2 | Not started |
| 4 | Settle original cohorts and make retry-safe test payouts | US-07, US-09 | L | 1–3 | Not started |
| 5 | Deliver a measurable, recoverable Android/Fuji pilot and grant evidence | US-10 | M | 1–4 | Not started |
| 6 | Integrate an approved external publisher with the Android ads SDK | US-11, US-12 | L | 1–5; separate beta scope/terms approval | Not started — separate beta |
| 7 | Settle partner earnings at 80/20 and enforce inventory/reward terms | US-13, US-14 | L | 6 | Not started — separate beta |

A story has one completion owner. Earlier phases may implement supporting criteria, but cannot claim the story complete until its owner phase validates all criteria. In particular, offline game work begins in Phase 2 while US-15 closes only after actual ad lifecycle integration in Phase 3.

## Phase 1 — Fund and refund a protected test campaign

**Goal:** An authenticated advertiser creates a test order, funds it once and recovers unused funds while all balances remain backed.

**Hard-deliverable check:** A local web/API demonstration shows order → test-token funding → indexed confirmation → refund, with an independently balanced ledger.

**User stories covered:** US-01 (AC-01.1–01.2), US-08 (AC-08.1–08.2)

**Complexity:** L · **Depends on:** Accepted local/Fuji test-token assumption; no existing product source.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `go.mod, go.sum, Makefile, compose.yaml, .env.example` | Pinned bootstrap and disposable PostgreSQL/EVM verification environment. |
| Create | `cmd/api/main.go, cmd/worker/main.go, internal/store/store.go` | Service composition and transactional persistence. |
| Create | `api/openapi.yaml, api/events.schema.json, migrations/0001_campaign_ledger.sql` | Orders, campaign terms, immutable ledger, chain events, outbox and scoped idempotency. |
| Create | `internal/campaigns/orders.go, internal/ledger/journal.go, internal/chain/indexer.go` | Funding, refunds, amount invariants and finalized indexing. |
| Create | `scripts/audit/update-tools.sh, scripts/audit/run-contract-audits.sh, scripts/audit/verify-reports.py, audit/README.md` | Guarded tool update, isolated snapshot/agent orchestration, source-bound evidence validation and retained-report policy. |
| Modify | `docs/runbooks/contract-audit.md` | Keep verified tool paths, supported invocation and actual evidence schema current. |
| Create | `contracts/foundry.toml, contracts/src/CampaignEscrow.sol, contracts/test/TestToken.sol` | Minimal allowlisted escrow and test-only token. |
| Create | `web/templates/orders.html, internal/identity/operator.go` | Protected advertiser/order UI and role-scoped bootstrap accounts. |
| Create | `internal/campaigns/orders_test.go, internal/ledger/journal_test.go, contracts/test/CampaignEscrow.t.sol, tests/integration/funding_test.go` | Unit, contract and real-store funding/refund checks. |

### Implementation tasks

- [ ] **1.1 — Money and order contract.** Interface: `CreateOrder(ctx, actor, OrderInput) (Order, error)`; `PostJournal(ctx, tx, Journal) error`. Encode decimal-string minor units and allocation versions; create schema/migration and bootstrap auth with owned tenant boundaries. Reject unbalanced or duplicate journals and changed idempotency bodies. Run the test-first cycle described above before marking this task complete.
- [ ] **1.2 — Escrow and indexer.** Interface: `ApplyFinalizedEvent(ctx, ChainEvent) error`; escrow methods defined above. Implement contract checks, confirmed-event deduplication and reorg quarantine. Indexing never treats a pending transaction as available money. Run the test-first cycle described above before marking this task complete.
- [ ] **1.3 — Refund journey.** Interface: `RequestRefund(ctx, actor, orderID, key) (RefundIntent, error)`. Protect payer/destination, serialize refund/release races and expose order/ledger state through minimal authenticated UI. Demo receipt and refund with local test assets. Run the test-first cycle described above before marking this task complete.

- [ ] **Contract internal audit.** Run `make audit-tools-update`, pin tool revisions, then `make audit-contracts PHASE=1` and `make verify-contract-audit PHASE=1` per the runbook. Scope: Escrow, deploy/setup scripts, refund/release role paths and Go/indexer money assumptions. Verify findings, add failing Foundry regression tests, fix in the implementation tree and rerun both skills against the fixed revision. Archive raw reports and independent reviewer closure; no contract-bearing phase is complete without both tools.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Duplicate final funding / wrong token, chain, short amount | AC-01.1, AC-01.2 | One credit for replay; mismatches cannot activate a campaign. Test fee-on-transfer or unexpected received amount rejection. |
| Refund backing and allocation version invariants | AC-08.1, AC-08.2 | Operations cannot withdraw refundable funds; historical 3000/6000/1000 basis points remain unchanged; journals sum to zero. |
| Concurrent refund/release, unauthorized call, reorg | AC-01.1, AC-08.1 | Only one serialized consumption of remaining funds; role/recipient changes fail; provisional/reorged events do not release backing. |
| Migration and restart on real services | AC-01.1, AC-08.2 | Fresh/repeated migration, crash/replay and backup restore preserve unique event IDs and balances. |
| Internal-audit gate evidence and regression checks | AC-10.5 | Both updated/pinned skill workflows complete on the candidate; absent reports, source mismatch, failed passes or unresolved confirmed defects prevent closure. Verified bugs receive failing-then-passing Foundry regressions and both-tool retest. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=1` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] An actual local test-token transfer and refund reconcile to the API ledger; no fixture pretends to be live ad revenue.
- [ ] All transport schemas, migration and test commands exist with pinned dependencies; invalid production chain/token configuration fails startup.
- [ ] Protected order UI cannot access another advertiser/order by substituting an ID.
- [ ] Both Nemesis and Pashov internal audits complete on the final candidate with pinned tool/source revisions, retained reports, passing fix/regression evidence, and no unresolved MUST-FIX/SHOULD-FIX.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Stop funding/indexing and disable new orders. Preserve chain events and journal history; deploy the previous API only if schema-compatible. A disposable local database can be recreated; a shared database requires backup/forward repair, never dropping financial records. Deployed test contracts are immutable: replace and explicitly reconcile their balances/addresses, not pretend a code revert undoes transactions.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 1 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Do not request routine phase-start permission.

## Phase 2 — Playable Android client with verified installation and wallet linking

**Goal:** A small offline Android game works, and an online user can authenticate the installation and prove ownership of a supported wallet.

**Hard-deliverable check:** Play/save/resume a bounded game on Android; complete a wallet signature and see server-confirmed ownership without any transfer authorization.

**User stories covered:** US-04 (AC-04.1–04.2), US-05 (AC-05.1–05.2); supports US-15 (AC-15.1–15.2), completed in Phase 3

**Complexity:** L · **Depends on:** Phase 1; chosen game/rule subset and rights; Android device, package/signing and Play project access for real integrity evidence.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `android/settings.gradle.kts, android/build.gradle.kts, android/gradle/libs.versions.toml, android/gradlew, android/sdk-core/build.gradle.kts, android/sdk-rewards/build.gradle.kts, android/game-core/build.gradle.kts, android/pilot/build.gradle.kts` | Pinned Android project and wrapper; module namespace is private until release ownership verified. |
| Create | `android/game-core/src/main/kotlin/adweb3/game/GameEngine.kt, android/game-core/src/test/kotlin/adweb3/game/GameEngineTest.kt, tests/fixtures/game/replay.json` | Selected deterministic rules translation and reference fixture evidence. |
| Create | `android/pilot/src/main/kotlin/adweb3/pilot/GameScreen.kt, android/pilot/src/main/kotlin/adweb3/pilot/ProgressStore.kt` | One playable mode/level, local save and clear unavailable-reward state. |
| Create | `android/sdk-core/src/main/kotlin/adweb3/sdk/InstallationClient.kt, android/sdk-rewards/src/main/kotlin/adweb3/rewards/WalletLinkClient.kt` | Per-install key, bounded sessions, optional reward auth and bound wallet return. |
| Create | `internal/identity/installation.go, internal/identity/wallet.go, internal/identity/playintegrity.go, migrations/0002_identity.sql` | Nonce/key/audience binding, integrity verification, refresh rotation, wallet proofs. |
| Create | `web/templates/wallet-link.html, web/static/wallet-link.js, internal/identity/identity_test.go, tests/integration/identity_test.go` | Human-readable domain-bound signing page and replay/session tests. |
| Create | `android/pilot/src/androidTest/kotlin/adweb3/pilot/IdentityLifecycleTest.kt` | Device integrity, app-link tampering, save/restore and wallet return evidence. |

### Implementation tasks

- [ ] **2.1 — Playable rules and progress.** Interface: `GameEngine.start(level, seed): GameState`; `advance(state, commands): Transition`; `ProgressStore.save/load`. Import only the approved rules subset and licensed content. Record reference revision and seeded expected states. Save progress atomically; corrupt saves produce recoverable UI, no reward claims. Run the test-first cycle described above before marking this task complete.
- [ ] **2.2 — Installation sessions.** Interface: `Enroll(ctx, Enrollment) (Session, error)`; Kotlin `InstallationClient.enroll()`. Verify integrity request content/package/certificate and freshness server-side; bind generated device key and rotating credentials. No fallback trust when verification fails. Play remains available. Run the test-first cycle described above before marking this task complete.
- [ ] **2.3 — Wallet proof.** Interface: `VerifyWalletLink(ctx, session, SignedChallenge) (WalletLink, error)`. Use SIWE parsing/verification with a single-use nonce; pilot supports explicitly documented EOA wallets only, rejecting unsupported contract-wallet signatures. Bind browser request, app completion verifier and destination; never ask for approvals/transfers. Reauthentication and new proof required for a change. Run the test-first cycle described above before marking this task complete.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Enrollment replay, wrong audience/certificate, expired token | AC-04.1, AC-04.2 | Reject authority without breaking offline game; requestHash/nonce binds the intended request, not just device attestation. |
| SIWE field tampering and wallet-change hijack | AC-05.1, AC-05.2 | Wrong domain/chain/nonce/time/signature and cross-session completion fail; a copied signing URL cannot silently redirect another account. |
| Seeded rules, malformed commands, process death | AC-15.1, AC-15.2 | Fixture states match; bad input cannot mutate accepted state; atomic local progress survives restart without wallet/network. |
| Real device and injected redirect | AC-04.1, AC-05.1 | Real verification evidence recorded separately from test fakes; deep-link flag alone cannot mark a wallet verified. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=2` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] The selected game is playable and resumes offline on a physical Android device.
- [ ] Wallet proof works for at least one documented supported Android EOA-wallet flow; no test-token or asset transfer is requested by signing.
- [ ] A real signed-build integrity check succeeds and recorded failure cases remain gameplay-only; mocked verdicts cannot be configured on the network-facing testnet service.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Disable enrollment/wallet mutations, revoke affected refresh credentials and retain audit records. Preserve local save format backward compatibility or provide a verified migration. Revert the SDK/client build only after checking saved progress and session compatibility; linked wallets and historical proofs must not be silently rewritten.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 2 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Do not request routine phase-start permission.

## Phase 3 — Expiring sponsor cards, validated delivery and offline-safe sync

**Goal:** The Android game shows an eligible funded sponsor card and records only validated delivery, with expiry and offline transitions enforced.

**Hard-deliverable check:** Fund a short test campaign, view its card on Home/results, see accepted delivery, disconnect while playing, then reconnect without credit for offline-only play.

**User stories covered:** US-02 (AC-02.1–02.3), US-03 (AC-03.1–03.2), US-06 (AC-06.1–06.2), US-15 (AC-15.1–15.3)

**Complexity:** L · **Depends on:** Phases 1–2; versioned test measurement, caps, late-event and hold policies accepted as test assumptions.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `internal/delivery/decision.go, internal/delivery/events.go, internal/delivery/validation.go, migrations/0003_delivery.sql` | Persistent pacing, spend reservations, viewability/sequence validation and original cohort attribution. |
| Create | `android/sdk-core/src/main/kotlin/adweb3/sdk/AdClient.kt, android/sdk-core/src/main/kotlin/adweb3/sdk/SponsorCard.kt, android/sdk-core/src/main/kotlin/adweb3/sdk/EventQueue.kt` | Lifecycle-aware rendering, lease expiry and bounded durable telemetry queue. |
| Create | `android/sdk-rewards/src/main/kotlin/adweb3/rewards/RewardSession.kt, internal/delivery/progress.go` | Online lease and separate progress synchronization. |
| Modify | `android/pilot/src/main/kotlin/adweb3/pilot/GameScreen.kt, api/openapi.yaml, api/events.schema.json, web/templates/orders.html` | Integrate actual sponsor slot and delivery/refund report. |
| Create | `internal/delivery/delivery_test.go, tests/integration/delivery_test.go, android/pilot/src/androidTest/kotlin/adweb3/pilot/AdLifecycleTest.kt` | Expiry, race, offline, queue and physical visibility assertions. |

### Implementation tasks

- [ ] **3.1 — Authorized ad decisions.** Interface: `Decide(ctx, session, DecisionRequest) (AdDecision, error)`. Select only approved/funded eligible campaigns at server time. Reserve potential billable spend transactionally, expire unused reservations, and issue decision IDs with min(campaign end, lease expiry). CDN serves approved assets; no arbitrary executable advertiser code. Run the test-first cycle described above before marking this task complete.
- [ ] **3.2 — Measurement and accounting.** Interface: `IngestBatch(ctx, session, EventBatch) (BatchResult, error)`. Validate app/decision binding, ordering, foreground visibility and caps. Use rational CPM arithmetic with carried integer remainder so rounding individual views cannot overcharge. One impression per decision; server processing and terminal state decide credit. Store per-campaign/app/activity-day weights and pending charges. Run the test-first cycle described above before marking this task complete.
- [ ] **3.3 — Actual Android delivery/offline behavior.** Interface: Kotlin `AdClient.request(placementId)` and `EventQueue.enqueue/flush`; `SyncProgress(ctx, session, ProgressBatch)`. Integrate one Home/results slot. Expiry uses server-relative monotonic time; app resume revalidates. Bound queue to 500 events/1 MiB, dropping oldest unsubmitted telemetry with an observable count; never synthesize credit for dropped events. Local progress has its own durable store. On clock uncertainty use no-fill. Run the test-first cycle described above before marking this task complete.
- [ ] **3.4 — Close campaign and report.** Interface: `CloseCampaign(ctx, campaignID, clock) (SettlementProposal, error)`. After configured report/adjustment windows, show accepted net charge and unused refund. Do not release disputed backing. Preserve all original weights for later reward settlement. Support equal-share scheduling only as an explicitly configured test policy; the public pilot offer remains capped actual-delivery CPM. Run the test-first cycle described above before marking this task complete.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Campaign T-1/T boundary with expiry worker disabled | AC-02.2, AC-02.3 | At T the decision is ineligible and official SDK clears/rejects it; pause and no-fill never block gameplay. |
| Ten equal shares and concurrent last-budget requests | AC-02.1 | In deterministic 1,000-opportunity equal-share fixture each gets 100; multi-instance decisions cannot oversell reserved inventory or the prepaid cap. |
| Visibility, duplicate IDs, caps, reordered/late/offline events | AC-03.1, AC-06.2 | No background/playtime/click credit; accepted online reports within retry policy dedupe; offline or after-cutoff reports cannot add weights; reject malformed batches without corrupting accepted ones. |
| 2,000 impressions at 5 CPM from a 25-unit deposit | AC-03.2 | 10 units charge, 15 refundable; no lost integer remainder, double deduction or release of contested funds. |
| Process death, wall-clock jump, outage, queue overflow | AC-06.1, AC-06.2, AC-15.2, AC-15.3 | Durable progress survives; earning pauses; paid card cannot outlive lease; bounded telemetry never blocks the game. |
| Replay fixtures and actual Android sponsor placement | AC-15.1, AC-15.3 | Selected rules remain correct after SDK integration; no sponsor over active gameplay or retry gate; record SDK/package revision. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=3` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] A funded test campaign produces actual Android delivery events and a reconciled accepted/unused-budget report.
- [ ] Expiry correctness holds with the scheduler stopped and with two API instances racing.
- [ ] All US-15 fixture, device and save/resume criteria pass against the integrated app, not only an isolated game module.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Pause new decisions, preserve raw events/accepted weights and let existing leases expire. Keep ingestion compatible long enough to record allowed late reports before closing cohorts. Do not re-bill previously accepted decision IDs. Only remove new schema on disposable environments; shared event/financial data gets forward repair.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 3 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Do not request routine phase-start permission.

## Phase 4 — Funded cohort settlement and retry-safe payouts

**Goal:** Validated activity turns into bounded test-token entitlements for the original cohort and verified wallets are paid at most once.

**Hard-deliverable check:** Show pending points, clear the campaign hold, finalize 30/60/10 backing, submit a test payout, kill/restart the worker and demonstrate one confirmed payment.

**User stories covered:** US-07 (AC-07.1–07.3), US-09 (AC-09.1–09.2)

**Complexity:** L · **Depends on:** Phases 1–3; immutable eligibility/hold/rounding/payout profile.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `internal/settlement/epochs.go, internal/payouts/worker.go, internal/payouts/reconcile.go, migrations/0004_rewards.sql` | Funded original cohorts, immutable entitlements and durable nonce/transaction intent state. |
| Create | `contracts/src/RewardDistributor.sol, contracts/test/RewardDistributor.t.sol` | Funded budget/entitlement commitments, destination binding and paid-ID protection. |
| Modify | `contracts/src/CampaignEscrow.sol, internal/ledger/journal.go, internal/chain/indexer.go, cmd/worker/main.go` | Confirmed release/distributor backing and cross-layer reconciliation. |
| Create | `android/sdk-rewards/src/main/kotlin/adweb3/rewards/BalanceClient.kt, web/templates/reconciliation.html` | Pending/finalized/submitted/paid state and operator audit trail. |
| Create | `internal/settlement/epochs_test.go, tests/integration/payout_test.go, tests/e2e/reward_journey_test.go` | Cohort, concurrency, rounding, failure and exactly-once effect scenarios. |

### Implementation tasks

- [ ] **4.1 — Immutable allocations.** Interface: `FinalizeEpoch(ctx, cohortID) (Epoch, error)`. Snapshot cleared backing, validation version, original participants and weights. Reward is floor(budget*weight/totalWeight); keep dust reserved. Freeze once; zero weights never mint claims. Record journal and outbox atomically. Run the test-first cycle described above before marking this task complete.
- [ ] **4.2 — Funded commitments and payments.** Interface: `SubmitPayout(ctx, payoutID) (TxIntent, error)`. Register immutable contract entitlements against confirmed funded budget. Aggregate finalized eligible sub-threshold amounts only under explicit destination/cohort provenance. Allocate nonce durably before submission and reconcile with chain state; no unlimited privileged treasury withdrawal path. Run the test-first cycle described above before marking this task complete.
- [ ] **4.3 — Recovery and balance journey.** Interface: `ReconcilePayout(ctx, payoutID) (PayoutStatus, error)`. Handle timeout, dropped/replaced transactions, revert, reorg, insufficient gas and blocked recipients. Confirm payment only after finality; keep failed entitlements payable and allow unrelated recipients to proceed. Run the test-first cycle described above before marking this task complete.

- [ ] **Contract internal audit.** Run `make audit-tools-update`, pin tool revisions, then `make audit-contracts PHASE=4` and `make verify-contract-audit PHASE=4` per the runbook. Scope: Distributor plus escrow integration, entitlement/recipient binding, token behavior and Go payout recovery. Verify findings, add failing Foundry regression tests, fix in the implementation tree and rerun both skills against the fixed revision. Archive raw reports and independent reviewer closure; no contract-bearing phase is complete without both tools.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Concurrent finalize, unequal weights and rounding dust | AC-07.1 | Exactly one entitlement set; sums never exceed confirmed budget; all ledger accounts balance in smallest units. |
| Zero weights/unfunded epoch and delayed clearance | AC-07.2, AC-07.3 | No unfunded claim; funds remain reserved; clearance days later uses original app/campaign cohort and cannot reward later players. |
| Crash at each intent/submission/confirmation boundary | AC-09.1 | One effective on-chain payout even if hash persistence fails; replacement/retry cannot double-spend entitlement. |
| Wrong destination, unverified wallet, failing recipient, gas loss | AC-09.2 | Fail closed and preserve liability; another recipient still pays; no paid status before finality. |
| End-to-end refundable-versus-earned funding | AC-07.1, AC-09.1 | Only cleared accepted revenue funds entitlements; refunded or disputed deposits cannot fund a payout. |
| Internal-audit gate evidence and regression checks | AC-10.5 | Both updated/pinned skill workflows complete on the candidate; absent reports, source mismatch, failed passes or unresolved confirmed defects prevent closure. Verified bugs receive failing-then-passing Foundry regressions and both-tool retest. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=4` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] The entire local campaign → eligible view → original cohort → verified-wallet test payout journey is reproducible.
- [ ] Ledger backing, contract budgets, finalized entitlements and confirmed payments reconcile exactly after injected crashes/reorgs.
- [ ] Android/server reports distinguish pending points, finalized balance and paid transactions; below-threshold money is retained visibly.
- [ ] Both Nemesis and Pashov internal audits complete on the final candidate with pinned tool/source revisions, retained reports, passing fix/regression evidence, and no unresolved MUST-FIX/SHOULD-FIX.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Pause new finalization/payment submissions and preserve immutable entitlements, nonce records and signed transaction intents. Reconcile already submitted transactions before any replacement worker starts. Never delete a paid ID, reuse a nonce blindly, or reverse an on-chain payment by editing the database. Redeploy contracts only with explicit test-fund reconciliation.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 4 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Do not request routine phase-start permission.

## Phase 5 — Verified Android/Fuji pilot and grant evidence

**Goal:** Produce a reviewable testnet pilot with operational controls, privacy handling, measured performance and accurate evidence.

**Hard-deliverable check:** One physical Android device runs a Fuji test campaign through payment, expiry, refund and payout; operators can pause and reconcile it, and evidence distinguishes test activity from customers.

**User stories covered:** US-10 (AC-10.1–10.4)

**Complexity:** M · **Depends on:** Phases 1–4; applicant-controlled test environment/device; no mainnet or store-publication authorization inferred.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `internal/operations/audit.go, internal/operations/policy.go, internal/operations/retention.go, migrations/0005_operations.sql` | Role-audited controls, deny-by-default real-money switches and retention/deletion jobs. |
| Create | `web/templates/operations.html, tests/integration/operations_test.go, tests/load/ad_decisions.js` | Operational UI, permission/privacy tests and recorded load test. |
| Create | `docs/runbooks/pilot.md, docs/runbooks/recovery.md, docs/runbooks/data-retention.md, docs/evidence/pilot-report.md` | Reproducible deployments, incident restore, data inventory and genuine evidence. |
| Modify | `Makefile, grants/team1-mini-grant/evidence-checklist.md, grants/team1-mini-grant/milestones-and-budget.md` | Verification commands and actual milestone evidence/cost revisions; never mark unperformed research complete. |

### Implementation tasks

- [ ] **5.1 — Operable and restricted pilot.** Interface: `CanEnableMoney(ctx, appID, market) error`; audited admin handlers. Default all real-money flags off. Limit test environment to its chain/token registry. Show terms, permissions and rejection reasons without secrets. Add pause controls and reconcile owner-only fund actions. Run the test-first cycle described above before marking this task complete.
- [ ] **5.2 — Privacy and recovery.** Interface: `ApplyDeletion(ctx, playerID, RetentionPolicy) (DeletionReport, error)`. Document data inventory, lawful financial retention decision dependency, access minimization and test-data cleanup. Exercise restore and deletion with retained financial evidence access-restricted. Test policy values are not a legal retention decision. Run the test-first cycle described above before marking this task complete.
- [ ] **5.3 — Performance/device/Fuji evidence.** Interface: `make verify-load`, `make verify-fuji`, `make verify-device`. Record dataset, environment, source revision, contract addresses, exact tool versions, device/OS and results. Load for 10 minutes at 100 decisions/s with funded/expired/no-fill mix and p95 <=250ms. Report API errors separately; never improve latency by silently failing requests. Run the test-first cycle described above before marking this task complete.
- [ ] **5.4 — Pilot reconciliation and handover.** Interface: `docs/evidence/pilot-report.md` plus grant evidence register. Map every pilot criterion to real implementation/test evidence. Record actual interviews/playtests or leave recruitment explicitly unmet. Re-cost Android deliverables before grant submission; independent final pilot review must be clean. Run the test-first cycle described above before marking this task complete.

- [ ] **Contract internal audit.** Run `make audit-tools-update`, pin tool revisions, then `make audit-contracts PHASE=5` and `make verify-contract-audit PHASE=5` per the runbook. Scope: The complete pilot contract/deployment candidate and its Go-to-chain trust boundaries. Verify findings, add failing Foundry regression tests, fix in the implementation tree and rerun both skills against the fixed revision. Archive raw reports and independent reviewer closure; no contract-bearing phase is complete without both tools.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Audit trace and privilege/log inspection | AC-10.1 | Funding through payout references are traversable; private keys, authentication tokens and sensitive signing material do not appear in logs. |
| Unapproved market/app enable attempt | AC-10.2 | Real-money activation denied without recorded launch determination; testnet remains clearly labeled; admin requests are role-checked and audited. |
| Load and bounded resource test | AC-10.3 | Required recorded p95/throughput met on stated hardware; data sizes and error rate reported. |
| Deletion, retention restriction and restore | AC-10.4 | Deletable data disappears; required financial evidence remains restricted/audited; restore does not silently resurrect deleted personal data outside the approved policy. |
| Internal-audit gate evidence and regression checks | AC-10.5 | Both updated/pinned skill workflows complete on the candidate; absent reports, source mismatch, failed passes or unresolved confirmed defects prevent closure. Verified bugs receive failing-then-passing Foundry regressions and both-tool retest. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=5` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] All pilot FR/NFR criteria have passing evidence; backend, contract, Android device and Fuji checks are real and recorded.
- [ ] `docs/evidence/pilot-report.md` explains remaining mainnet/store/market and audience/cost uncertainties without presenting test wallets as users.
- [ ] Independent integrated review of Phases 1–5 is clean; complete the pilot verification table, not the separate publisher-beta rows.
- [ ] Grant milestone evidence is mapped honestly; no automatic application submission, release upload or mainnet deployment occurs.
- [ ] Both Nemesis and Pashov internal audits complete on the final candidate with pinned tool/source revisions, retained reports, passing fix/regression evidence, and no unresolved MUST-FIX/SHOULD-FIX.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Pause pilot campaigns/rewards, preserve financial evidence, revoke exposed service credentials if necessary and use the rehearsed restore procedure. Disable privacy cleanup jobs if defective before data loss; restoration must reapply deletion tombstones and retention restrictions. Roll back observability/UI only where schema-compatible.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 5 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. If approval covers the pilot only, deliver the pilot report and stop at the authorized scope boundary; do not silently start the publisher beta.

## Phase 6 — Approved external publisher with reusable Android ads SDK

**Goal:** One approved external app can integrate ordinary contextual ads without requiring end-user wallets, with strict tenant isolation.

**Hard-deliverable check:** Two separate publisher fixtures show isolated apps/reports; one real authorized beta integration uses the versioned Android SDK and handles no-fill/offline safely.

**User stories covered:** US-11 (AC-11.1–11.3), US-12 (AC-12.1–12.3)

**Complexity:** L · **Depends on:** Phases 1–5; separately approved/budgeted beta; actual partner app/signing ownership and permission to integrate. Local fixtures do not establish a real partner.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `internal/publishers/onboarding.go, internal/publishers/authorization.go, migrations/0006_publishers.sql` | Organizations, roles, app ownership, agreement registration and isolation controls. |
| Modify | `internal/delivery/decision.go, internal/delivery/events.go, internal/identity/installation.go, api/openapi.yaml` | Publisher-bound sessions/decisions; advertiser eligibility; block paid external traffic until Phase 7 settlement is enabled. |
| Create | `android/sample/build.gradle.kts, android/sample/src/main/kotlin/adweb3/sample/MainActivity.kt, docs/sdk/android-integration.md, docs/sdk/data-disclosures.md` | Wallet-free ads sample, API/version support policy and integration/privacy guidance. |
| Modify | `android/sdk-core/build.gradle.kts, android/sdk-core/src/main/kotlin/adweb3/sdk/AdClient.kt` | Versioned local Maven package, consumer rules, bounded callback/queue behavior and core-only dependency graph. |
| Create | `tests/integration/tenant_isolation_test.go, android/sample/src/androidTest/kotlin/adweb3/sample/AdsOnlyTest.kt` | Negative cross-tenant checks and independent package-consumer tests. |

### Implementation tasks

- [ ] **6.1 — Verified organization and app.** Interface: `ApproveApp(ctx, operator, PublisherApp) error`; `AuthorizePublisher(ctx, actor, publisherID, action) error`. Verify organization membership and app control, register signing identity/placements and scoped credentials. Apply authorization to all protected read/write/export/asset routes, not just dashboard routes. Run the test-first cycle described above before marking this task complete.
- [ ] **6.2 — Ads-only package.** Interface: Kotlin `AdWeb3Ads.initialize(config, consent)`; `requestAd(placementId): AdResult`. Publish to a local test Maven repository and consume from the independent sample; document semantic/API versions and supported SDK range. Rewards remain a separate optional module; no player login/wallet needed for ordinary contextual ads. Run the test-first cycle described above before marking this task complete.
- [ ] **6.3 — Controlled beta validation.** Interface: Publisher/app suspension and SDK integration checklist. Onboard one authorized partner using approved test campaigns. No external paid traffic until Phase 7 financial controls pass. Verify consent-denied behavior, queue limits, expiry and account isolation with the partner; publishing the SDK externally needs its separate release authorization. Run the test-first cycle described above before marking this task complete.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Copied app ID/mismatched certificate and suspension | AC-11.1, AC-11.3 | Neither app IDs nor copied SDK keys prove ownership; suspension prevents new decisions and prior decisions cannot bill beyond lease. |
| Cross-tenant route/query/event/payout substitutions | AC-11.2 | No protected data leaks or balance/attribution changes through direct IDs, exports, assets or replayed events. |
| Core-only sample without reward dependency | AC-12.1 | Actual sample renders approved ads and reports evidence without wallet calls/account UI. |
| Denied consent/outage/full queue and client version compatibility | AC-12.2, AC-12.3 | Permitted contextual behavior only; app remains usable; memory/disk queue bounds hold; documented older client behavior stays compatible. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=6` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] Versioned Android SDK artifacts build and integrate in a separate consumer using only published package interfaces.
- [ ] Isolation tests cover protected data and monetary effects; no external paid traffic is enabled ahead of settlement.
- [ ] One authorized external integration is evidenced, or phase stays incomplete; two synthetic tenants alone do not meet that claim.
- [ ] SDK documentation matches observed collection/consent/reward dependencies and supported versions.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Disable external publisher serving and revoke management/session credentials without disturbing owned-game campaigns. Mark incompatible SDK release unavailable and preserve supported clients until their leases expire. Never merge/delete publisher financial histories; migrate tenant records forward. Local package deletion is safe only for unpublished artifacts.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 6 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Do not request routine phase-start permission.

## Phase 7 — Publisher settlement and inventory controls

**Goal:** Ordinary partner ads allocate revenue exactly 80/20 with transparent statements, while inventory and optional reward policies are enforced server-side.

**Hard-deliverable check:** A partner test campaign earns 1,000 units: 800 becomes publisher liability and 200 platform allocation; statements, verified payout and retry recovery reconcile. Excluded/incentivized inventory is denied.

**User stories covered:** US-13 (AC-13.0–13.3), US-14 (AC-14.1–14.3)

**Complexity:** L · **Depends on:** Phase 6; publisher settlement asset, hold/minimum/cadence, net deductions and dispute terms agreed before external money. Reward-enabled allocation remains disabled unless separately accepted.

### Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `internal/publishers/settlement.go, internal/publishers/statements.go, internal/publishers/payout_setup.go, migrations/0007_publisher_settlement.sql` | Agreement-version allocations, distinct liabilities, statements and audited wallet change. |
| Modify | `internal/settlement/epochs.go, internal/payouts/worker.go, internal/delivery/decision.go, contracts/src/RewardDistributor.sol` | Generalized funded beneficiary types; no migration of existing player commitments/destinations. |
| Create | `web/templates/publisher-statement.html, tests/integration/publisher_settlement_test.go, tests/e2e/publisher_journey_test.go` | Partner statement and complete approval-to-settlement checks. |
| Create | `docs/evidence/publisher-report.md, docs/runbooks/publisher-disputes.md` | Full feature evidence and explicit late-dispute/hold policy. |

### Implementation tasks

- [ ] **7.1 — Agreement-backed allocations.** Interface: `SettlePublisherDelivery(ctx, deliveryID) (Allocation, error)`. Snapshot agreement and inventory mode per delivery. Ordinary rows allocate 8000/2000 basis points exactly once; no mandatory subscription/no-fill charge. Any platform reserve is internal to its allocation. Run the test-first cycle described above before marking this task complete.
- [ ] **7.2 — Publisher liability and payment.** Interface: `BuildStatement(ctx, publisherID, period) (Statement, error)`; `VerifyPayoutDestination(ctx, actor, SignedChallenge)`. Reuse funded payout discipline with distinct publisher accounts and immutable approved destination/commitment. Show gross, deductions, holds, net split, payable and paid state. Chargeback after final payment follows separately funded risk policy. Run the test-first cycle described above before marking this task complete.
- [ ] **7.3 — Inventory and optional rewards.** Interface: `ValidateInventory(ctx, campaign, app, incentiveMode) error`. Advertiser-approved app/mode cannot be overridden by client flags. Test 3000/5000/2000 allocation as a disabled fixture; permit activation only with approved terms, game validation contract and market decision. Preserve original cohort when money clears later. Run the test-first cycle described above before marking this task complete.
- [ ] **7.4 — Full feature review.** Interface: Complete verification and final review tables. Re-run all pilot/beta journeys after publisher changes. Confirm no changes to owned 30/60/10 or old payouts. Mainnet, broad public SDK release, and store submission remain distinct launch decisions. Run the test-first cycle described above before marking this task complete.

- [ ] **Contract internal audit.** Run `make audit-tools-update`, pin tool revisions, then `make audit-contracts PHASE=7` and `make verify-contract-audit PHASE=7` per the runbook. Scope: Publisher contract changes and the complete integrated feature candidate. Verify findings, add failing Foundry regression tests, fix in the implementation tree and rerun both skills against the fixed revision. Archive raw reports and independent reviewer closure; no contract-bearing phase is complete without both tools.

### Tests to write first

Use the phase test files listed above; convert each named scenario into executable assertions before its implementation.

| Test / scenario | Maps to | Required assertions |
|-----------------|---------|---------------------|
| Free access and 1,000-unit ordinary allocation | AC-13.0, AC-13.1 | No subscription/no-fill fee; 800 publisher, 200 platform, no additional commission or reserve deduction. |
| Refund/dispute/payout retries and agreement/destination changes | AC-13.2, AC-13.3 | Backing protected, previous terms/destinations immutable, at-most-once effect and unrelated publisher/player balances untouched. |
| Inventory exclusion and client-incentive override | AC-14.1, AC-14.3 | Excluded or unapproved reward apps cannot get decisions or monetary authority; ordinary mode remains independently configurable. |
| Optional 30/50/20 fixture with late clearance | AC-14.2 | 300 original player cohort, 500 publisher, 200 platform; no cross-app/cohort transfer. Passing fixture does not authorize enabling this product. |
| Mixed owned/partner campaign regression | AC-13.1, AC-14.1, AC-14.2 | Each accepted delivery uses its own recorded mode/terms; no double allocation and total accepted charges stay within escrow. |
| Internal-audit gate evidence and regression checks | AC-10.5 | Both updated/pinned skill workflows complete on the candidate; absent reports, source mismatch, failed passes or unresolved confirmed defects prevent closure. Verified bugs receive failing-then-passing Foundry regressions and both-tool retest. |

### Definition of done

- [ ] All listed tests were observed failing for the intended reason and now pass.
- [ ] `make verify PHASE=7` passes all introduced suites with zero failures/skips; applicable physical-device/Fuji evidence is recorded, not replaced by mocks.
- [ ] Partner statements and confirmed test payouts reconcile by agreement, campaign, original cohort and beneficiary.
- [ ] The production feature switch for unapproved optional reward terms remains off; test coverage is not policy approval.
- [ ] Every spec requirement and acceptance criterion has actual implementation/validation evidence and the complete independent review is clean.
- [ ] Final report distinguishes implemented testnet beta from real partner commercial traction or mainnet launch.
- [ ] Both Nemesis and Pashov internal audits complete on the final candidate with pinned tool/source revisions, retained reports, passing fix/regression evidence, and no unresolved MUST-FIX/SHOULD-FIX.
- [ ] Independent code review closed with zero unresolved MUST-FIX/SHOULD-FIX; disputes include reviewer resolution.
- [ ] Spec/plan criteria are reconciled to code and evidence; update the verification table before marking complete.

### Rollback

Pause external settlement/serving, retain owed publisher/player balances and reconcile all in-flight transactions before restart. Undo code only with compatible immutable agreement/destination semantics; never migrate 80/20 earnings into owned-game accounts. Use compensating entries, not destructive history edits, for validated corrections.

### Review-gate log

Not run. Fill during implementation; no placeholder clean verdicts.

| Round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|-------|------|----------|------------|------|------------------|---------|

### Completion

Record evidence and report `Phase 7 complete — deliverable — files — review rounds`. Continue directly to the next approved phase. Mark full implementation Complete only after the final feature checklist passes.

## Grant milestone mapping and scope limits

| Grant deliverable | Plan evidence required | Scope note |
|-------------------|------------------------|------------|
| M1 — Fund/schedule a test campaign | Phase 1 funding/refund plus Phase 3 decision-time expiry; Fuji evidence in Phase 5 | A ledger-only Phase 1 is not all of M1. The grant weeks are provisional and must be reconciled with this dependency order before submission. |
| M2 — Android slice and SDK foundation | Phases 2–3 and physical-device evidence | Includes actual Android development, not a renamed Swift integration. |
| M3 — Settlement and evaluation | Phases 4–5 and documented actual discovery outcomes | Does not guarantee recruited users/advertisers, store acceptance or mainnet readiness. |
| External publisher SDK beta | Phases 6–7 | Outside the current USD 7,500 pilot ask; needs separate budget and partner access. |

This plan neither validates eight weeks nor spends the grant. Update the grant budget/timing with an actual Android implementation estimate before submission. No audit/compliance budget is represented as a full production assurance package.

## Story coverage check

| Story | Completion phase |
|-------|------------------|
| US-01 | 1 |
| US-02 | 3 |
| US-03 | 3 |
| US-04 | 2 |
| US-05 | 2 |
| US-06 | 3 |
| US-07 | 4 |
| US-08 | 1 |
| US-09 | 4 |
| US-10 | 5 |
| US-11 | 6 |
| US-12 | 6 |
| US-13 | 7 |
| US-14 | 7 |
| US-15 | 3 |

**Unassigned stories:** None. Pilot approval does not authorize the separately scoped beta stories.

## Deviations from the spec / planning decisions

| Date | Phase | What changed | Why |
|------|-------|--------------|-----|
| 2026-10-19 | Planning | Added FR-17 / US-15 and concrete Android acceptance criteria to the spec | Existing Android game deliverable needed explicit traceability; no additional game or full-port scope added. |
| 2026-10-19 | Planning | Added owner-requested Nemesis/Pashov update, review and retest gates with NFR-09 / AC-10.5 | Contract development requires both internal-audit skills from the existing shared tool checkouts. |
| 2026-10-19 | Planning | Conditional draft written while GAP-B01–B04 remain open | Owner explicitly requested a plan; test assumptions are exposed for approval, not treated as resolved commercial/launch decisions. |

## Implementation verification

No implementation evidence exists yet. Populate each acceptance criterion individually when executing; a row listing a story is not permission to omit one of its criteria. Commands must record revision, environment and actual result. Pilot completion leaves beta rows explicitly Not started; full completion requires all rows.

| Requirement | Stories / acceptance coverage | Completion phases | Implementation evidence | Validation result | Status |
|-------------|-------------------------------|-------------------|-------------------------|-------------------|--------|
| FR-01 | US-01 — all listed ACs | 1 | None yet | Not run | Not started |
| FR-02 | US-01, US-02 — all listed ACs | 1, 3 | None yet | Not run | Not started |
| FR-03 | US-02 — all listed ACs | 3 | None yet | Not run | Not started |
| FR-04 | US-03 — all listed ACs | 3 | None yet | Not run | Not started |
| FR-05 | US-06 — all listed ACs | 3 | None yet | Not run | Not started |
| FR-06 | US-04 — all listed ACs | 2 | None yet | Not run | Not started |
| FR-07 | US-05 — all listed ACs | 2 | None yet | Not run | Not started |
| FR-08 | US-07 — all listed ACs | 4 | None yet | Not run | Not started |
| FR-09 | US-08 — all listed ACs | 1 | None yet | Not run | Not started |
| FR-10 | US-09 — all listed ACs | 4 | None yet | Not run | Not started |
| FR-11 | US-02, US-03, US-10 — all listed ACs | 3, 5 | None yet | Not run | Not started |
| FR-12 | US-10, US-14 — all listed ACs | 5, 7 | None yet | Not run | Not started |
| FR-13 | US-11, US-13 — all listed ACs | 6, 7 | None yet | Not run | Not started |
| FR-14 | US-11, US-12 — all listed ACs | 6 | None yet | Not run | Not started |
| FR-15 | US-13 — all listed ACs | 7 | None yet | Not run | Not started |
| FR-16 | US-14 — all listed ACs | 7 | None yet | Not run | Not started |
| FR-17 | US-15 — all listed ACs | 3 | None yet | Not run | Not started |
| NFR-01 | US-01, US-03, US-07, US-08, US-09 — all listed ACs | 1, 3, 4 | None yet | Not run | Not started |
| NFR-02 | US-02, US-04, US-06 — all listed ACs | 2, 3 | None yet | Not run | Not started |
| NFR-03 | US-04, US-05, US-10 — all listed ACs | 2, 5 | None yet | Not run | Not started |
| NFR-04 | US-03, US-10 — all listed ACs | 3, 5 | None yet | Not run | Not started |
| NFR-05 | US-10 — all listed ACs | 5 | None yet | Not run | Not started |
| NFR-06 | US-10 — all listed ACs | 5 | None yet | Not run | Not started |
| NFR-07 | US-11 — all listed ACs | 6 | None yet | Not run | Not started |
| NFR-08 | US-12 — all listed ACs | 6 | None yet | Not run | Not started |
| NFR-09 | US-10 / AC-10.5; repeated contract-change gates | 1, 4, 5, 7 | None yet | Not run | Not started |

### Acceptance-criterion evidence register

Each criterion below has a planned owner, not a passing result. Expand the evidence column with exact file/test references and observed output during implementation.

| Criterion | Owner phase | Evidence / result | Status |
|-----------|-------------|-------------------|--------|
| AC-01.1 | 1 | None / not run | Not started |
| AC-01.2 | 1 | None / not run | Not started |
| AC-02.1 | 3 | None / not run | Not started |
| AC-02.2 | 3 | None / not run | Not started |
| AC-02.3 | 3 | None / not run | Not started |
| AC-03.1 | 3 | None / not run | Not started |
| AC-03.2 | 3 | None / not run | Not started |
| AC-04.1 | 2 | None / not run | Not started |
| AC-04.2 | 2 | None / not run | Not started |
| AC-05.1 | 2 | None / not run | Not started |
| AC-05.2 | 2 | None / not run | Not started |
| AC-06.1 | 3 | None / not run | Not started |
| AC-06.2 | 3 | None / not run | Not started |
| AC-07.1 | 4 | None / not run | Not started |
| AC-07.2 | 4 | None / not run | Not started |
| AC-07.3 | 4 | None / not run | Not started |
| AC-08.1 | 1 | None / not run | Not started |
| AC-08.2 | 1 | None / not run | Not started |
| AC-09.1 | 4 | None / not run | Not started |
| AC-09.2 | 4 | None / not run | Not started |
| AC-10.1 | 5 | None / not run | Not started |
| AC-10.2 | 5 | None / not run | Not started |
| AC-10.3 | 5 | None / not run | Not started |
| AC-10.4 | 5 | None / not run | Not started |
| AC-10.5 | 5 (also required at 1, 4 and 7) | None / not run | Not started |
| AC-11.1 | 6 | None / not run | Not started |
| AC-11.2 | 6 | None / not run | Not started |
| AC-11.3 | 6 | None / not run | Not started |
| AC-12.1 | 6 | None / not run | Not started |
| AC-12.2 | 6 | None / not run | Not started |
| AC-12.3 | 6 | None / not run | Not started |
| AC-13.0 | 7 | None / not run | Not started |
| AC-13.1 | 7 | None / not run | Not started |
| AC-13.2 | 7 | None / not run | Not started |
| AC-13.3 | 7 | None / not run | Not started |
| AC-14.1 | 7 | None / not run | Not started |
| AC-14.2 | 7 | None / not run | Not started |
| AC-14.3 | 7 | None / not run | Not started |
| AC-15.1 | 3 | None / not run | Not started |
| AC-15.2 | 3 | None / not run | Not started |
| AC-15.3 | 3 | None / not run | Not started |

## Final feature review

Use an independent read-only reviewer with the current spec, this plan, full diff, evidence register, and `code-review-checklist.md` from the feature-request skill. Review the pilot at Phase 5 and the complete feature at Phase 7; label which scope was reviewed. Fix/retest/re-review until no unresolved MUST-FIX/SHOULD-FIX remains. Never fill a clean result in advance.

| Scope / round | Date | MUST-FIX | SHOULD-FIX | NOTE | Fixed / disputed | Verdict |
|---------------|------|----------|------------|------|------------------|---------|

- [ ] Every requirement and criterion in the approved scope has implementation and passing validation evidence.
- [ ] Every approved phase deliverable and binary definition of done is satisfied.
- [ ] Integrated checks, contract invariants, concurrency/recovery and physical-device checks pass with recorded revisions.
- [ ] Updated/pinned Nemesis and Pashov audits cover the final production-contract/deployment candidate; reports, verified fixes, regressions and both-tool retest evidence are linked and complete.
- [ ] Spec and plan describe actual behavior without silently weakening requirements or treating deferred work as complete.
- [ ] Independent final review is clean; mainnet, public distribution, market approval and unverified business assumptions remain clearly distinguished.
- [ ] If only pilot scope was approved, pilot is marked complete and publisher beta stays Not started; the overall plan is not marked Complete.

## Planning validation and reference material

Document validation is separate from product verification. Planning self-review on 2026-10-19 checked all 26 requirement mappings, 15 unique story completion owners, 41 acceptance-criterion evidence rows, mandatory sections for seven phases, file-action consistency, scope boundaries, and local links. Corrected the Android SDK build-file creation dependency; product checks remain Not run. No Go, Android or contract suite can run in the current documentation-only workspace.

References checked while planning: [Play Integrity overview and request binding](https://developer.android.com/google/play/integrity/overview), [Android library publication](https://developer.android.com/build/publish-library), [Solidity security considerations](https://docs.soliditylang.org/en/latest/security-considerations.html). Use current stable implementation documentation when pinning versions. Official distribution/integrity setup and library metadata matter; an emulator or mock verifier is not a substitute for the required physical-device evidence.
