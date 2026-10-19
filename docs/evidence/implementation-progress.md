# Pilot implementation evidence

Plan: [implementation-plan.md](../features/adweb3-platform/implementation-plan.md).

Owner approved pilot Phases 1–5 on 2026-10-19, including planned reviews and Nemesis/Pashov audits. Publisher beta is not authorized.

## Phase 1 — In progress

- Foundry tests first failed for protected balances, replay/authorization and allocation; implementation then passed five tests including a 257-case conservation fuzz run.
- HTTP unauthorized and trailing-JSON tests were observed failing, then passed after implementation/fix.
- Actual Anvil test deployed token/escrow, funded an order, indexed it to a real PostgreSQL ledger and confirmed a partial refund. Initial indexer returned not-implemented; the implemented path passed.
- PostgreSQL tests exercise scoped creation idempotency and eight concurrent refund requests; accepted refund reservations cannot overdraw.
- Review rounds 1–3 found canonical-range/reorg, escrow/token-binding, persistent-quarantine, worker starvation and provisional-registration recovery defects. Fixes and regression verification have run; final review and audit closure remain pending. Phase is not complete.
- Tool update fetched verified Nemesis/Pashov origins and fast-forwarded clean checkouts (already current). See audit/tools-lock.json. No clean audit report claimed yet.
- Go 1.22's test binary failed macOS loading (`missing LC_UUID`); project toolchain upgraded to Go 1.25.3. Dependency bootstrap is pinned in go.mod/go.sum.

## External verification dependencies

- Android SDK exists at the configured local SDK path; ADB currently reports no connected physical device.
- Claude CLI audit invocation failed because its OAuth session expired. Native Pashov review is in progress; the native Nemesis attempt was interrupted by an automatic safety flag. Neither audit is claimed complete or clean.
- Play Integrity app/package setup and a funded Fuji test relayer remain to be verified; no credentials or production-money permissions are inferred.

## Public landing preparation

- The owner requested public GitHub `JordyKingz/AdWeb3` and a Vercel landing deployment. The Git remote is configured. Public publication is a separate website deliverable, not pilot-phase completion.
- `frontend/` preserves the supplied Nuxt design. Its production landing build exports only `/` plus static fallback pages, shows the development notice, disables wallet reconnection and redirects unfinished dApp URLs to pilot information. Mock Nitro handlers are not deployed.
- Supported Node 24, frozen Bun dependencies, type checking and the static export have been verified locally. A build-output verifier rejects exported app/API routes and server functions. The Git history and staged publication files passed redacted Gitleaks scans. See `frontend/README.md` and `SECURITY.md` for Vercel setup and public-variable boundaries.
- Independent publication review completed three rounds. The placeholder contact address and misleading policy links were fixed; the final round also reviewed the local-hook/inactive-CI fallback with zero unresolved MUST-FIX or SHOULD-FIX findings. Final regeneration, static-output verification, staged/history/output secret scans and whitespace validation passed before the public push.
- GitHub rejected the first public push because its CLI OAuth token lacks `workflow` scope. The secret-scanning workflow is retained as an inactive template; a versioned local pre-push hook provides the automatic scan without broadening account permissions. GitHub's native secret scanning and push protection remain enabled. The unpublished publication commit is revised accordingly before retrying.
- The local hook was verified with an isolated temporary repository: a generated fake credential blocked the push with redacted output, deleting the file still left detectable history, and missing Gitleaks also blocked the push. No actual credential or fixture was added to this repository.
- Vercel CLI is logged out. No Vercel account/project has been linked or deployment claimed.
