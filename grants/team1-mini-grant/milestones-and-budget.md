# Proposed grant milestones and budget

**Draft ask: USD 7,500.** Applicant confirmation required for scope, capacity, amount, and timeline. Amounts are planning allocations, not supplier quotations, committed costs, or Team1-approved payment terms. The public program materials do not establish the grant's payout schedule or asset.

This proposal funds a bounded testnet pilot, not the full AdWeb3 feature spec or a mainnet launch. Android-first and free SDK access with a 20% AdWeb3 commission are confirmed. The retained eight-week target and USD 7,500 total require Android-specific re-estimation: the original integration allowance assumed an existing iOS client. No Android game build is verified. The target begins after the revised scope is agreed and work starts; no calendar completion date is promised.

## Milestones

| Milestone | Target window | Deliverable | Acceptance evidence | Budget allocation |
|-----------|---------------|-------------|---------------------|-------------------|
| M1 — Fund and schedule a test campaign | Weeks 1–2 | Go campaign API plus Fuji escrow prototype with order/payment attribution, campaign dates, and protected unused budget | Demo funds a campaign with a test token; duplicate funding events credit once; expired campaign returns no paid ad; an unused-budget scenario reconciles; source revision, test results, and testnet links recorded | USD 2,500 |
| M2 — Demonstrate one Android game slice and SDK foundation | Weeks 3–5 | Bounded Android Candy Blocks adaptation (candidate): one playable mode/ruleset, local progress, one labeled sponsor placement, online ingestion, and wallet-signing flow through an Android SDK foundation; versioned interface, integration example, and data-flow notes | Physical Android-device recording shows gameplay → sponsor card → accepted/rejected event results; package/signing identity and integrity-verification outcomes are recorded; replayed events/challenges fail; loss of connectivity preserves gameplay and pauses monetary eligibility; Android lifecycle/app-link tests, SDK/interface version, integration instructions, data-flow notes, and integration revision recorded | USD 3,000 |
| M3 — Demonstrate settlement and evaluate the pilot | Weeks 6–8 | Funded epoch allocation, retry-safe test payout, reconciliation, documented discovery, and final demo | Full test journey traces funding to 30/60/10 accounting and a verified-wallet test payout; no duplicate payout on retry; refundable funds remain unavailable; measured fee estimates and open risks recorded; report distinguishes real feedback from assumptions | USD 2,000 |
| **Total** | **8 weeks provisional; Android re-estimate required** | | | **USD 7,500** |

If payment by milestones is offered, this allocation can inform a discussion; it is not a statement that Team1 pays according to this schedule. If the grant process requires a different format, preserve the scope and reconcile the amounts before submission.

## Cost allocation

| Workstream | M1 | M2 | M3 | Total | Estimate basis |
|------------|----|----|----|-------|----------------|
| Go backend and contracts | 1,500 | 1,200 | 600 | 3,300 | Bounded campaign/ledger/settlement prototype; applicant to confirm development cost basis |
| Android game slice and SDK foundation | 250 | 1,100 | 250 | 1,600 | Retained allowance, not a validated estimate: Android adaptation, one placement, SDK adapter, and device validation; excludes full iOS-game parity, a production SDK, and external publisher integrations |
| Validation and security-review allowance | 400 | 400 | 400 | 1,200 | Foundry invariant/regression testing, updated/pinned Nemesis and Pashov internal reviews, and limited independent review; not a full professional production audit; obtain a quote for any external component |
| Infrastructure | 100 | 200 | 100 | 400 | Pilot hosting, storage, RPC, and observability allowance; confirm provider estimates |
| Discovery, documentation, and demo | 250 | 100 | 650 | 1,000 | Applicant/research effort and reviewable deliverables; no committed paid promotion |
| **Total USD** | **2,500** | **3,000** | **2,000** | **7,500** | |

No engineering hourly rate, available work hours, or supplier quote has been provided. Applicant must confirm whether these costs represent founder compensation, contractor work, or a combination, and check category eligibility in the actual grant terms. Tax inclusion and treatment must also be confirmed; this budget must not silently grow above the requested amount.

## Measurement commitments to confirm

- Demonstrate one complete reproducible end-to-end testnet flow and report failed cases as well as successful ones.
- Record valid/invalid event counts, accounting reconciliation, payout retry behavior, and measured transaction-cost estimates with their measurement date.
- Aim to interview at least three potential advertisers and observe ten voluntary player sessions; these are proposed research targets, not acquired partners/users or guaranteed recruitment outcomes. Report achieved counts and findings honestly.
- Do not equate wallet count, test scripts, or testnet transactions with unique humans, revenue, or market traction.
- Provide an accessible demo and documentation for grant review; public source publication depends on applicant consent and verified grant terms.

## Scope and contingencies

- Use test tokens only; clearly distinguish them from issuer-backed USD₮.
- No tradable token, liquidity pool, wagering, external-publisher onboarding/settlement, public publisher marketplace, production SDK, iOS SDK delivery, or multiple game integrations. A bounded Android game adaptation is included; full parity with all iOS Candy Blocks levels, modes, and polish is excluded. The Android SDK foundation is limited to one integration and its example/docs.
- The Android switch does not automatically increase the draft USD 7,500 total or imply more capacity. Re-estimate the game slice, Android SDK, device testing, and integrity setup before submission. Reconcile any changes across the cost table and milestones; do not promise this scope for the retained allowance without evidence. The external-publisher beta needs its own estimate.
- Grant funds are not advertiser deposits and are not subject to the product's 30/60/10 revenue split. No recurring player rewards are financed from this grant budget.
- The internal-audit workflow adds tool/runtime and review effort to re-estimate; the existing allowance is not proof that all review costs fit.
- Full mainnet audit, formal legal advice, commercial acquisition spend, and production operations beyond the pilot are not represented as covered by this request.
- Google Play reward-policy assessment and launch-market selection remain prerequisites for a real-money product. If the scope must change, agree it with the applicant and grant team before claiming completion against different deliverables.
- If partial funding is offered, negotiate a smaller deliverable set or confirmed co-funding; do not assume the complete scope can be delivered for any award amount.
- Missing physical-device or external-review evidence remains open; simulator results cannot replace it. Google Play approval is not a milestone promised by this application.
