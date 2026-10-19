# Pilot implementation evidence

Plan: [implementation-plan.md](../features/adweb3-platform/implementation-plan.md).

Owner approved pilot Phases 1–5 on 2026-10-19, including planned reviews and Nemesis/Pashov audits. Publisher beta is not authorized.

## Phase 1 — In progress

- Foundry tests first failed for protected balances, replay/authorization and allocation; implementation then passed five tests including a 257-case conservation fuzz run.
- HTTP unauthorized and trailing-JSON tests were observed failing, then passed after implementation/fix.
- Actual Anvil test deployed token/escrow, funded an order, indexed it to a real PostgreSQL ledger and confirmed a partial refund. Initial indexer returned not-implemented; the implemented path passed.
- PostgreSQL tests exercise scoped creation idempotency and eight concurrent refund requests; accepted refund reservations cannot overdraw.
- Review round 1 found canonical-range/reorg, escrow/token-binding and persistent-quarantine defects. Fixes and regression verification are underway; phase is not complete.
- Tool update fetched verified Nemesis/Pashov origins and fast-forwarded clean checkouts (already current). See audit/tools-lock.json. No clean audit report claimed yet.
- Go 1.22's test binary failed macOS loading (`missing LC_UUID`); project toolchain upgraded to Go 1.25.3. Dependency bootstrap is pinned in go.mod/go.sum.

## External verification dependencies

- Android SDK exists at the configured local SDK path; ADB currently reports no connected physical device.
- Claude CLI is authenticated. Actual audit invocation and coverage are still pending.
- Play Integrity app/package setup and a funded Fuji test relayer remain to be verified; no credentials or production-money permissions are inferred.
