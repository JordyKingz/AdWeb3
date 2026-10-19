# AdWeb3 pilot

Android-first advertiser-funded game rewards. The repository is in implementation; see the [plan](docs/features/adweb3-platform/implementation-plan.md) and [execution evidence](docs/evidence/implementation-progress.md). Local Anvil and Fuji test assets only; no mainnet launch.

## Development

Required: Go 1.25.3, Foundry 1.5.1 / Solidity 0.8.30, PostgreSQL 16, and an isolated Anvil chain. Android additionally uses its Gradle wrapper and the SDK. `make verify PHASE=1` runs Foundry, Go, race, real-database/Anvil integration and audit-evidence validator tests. It fails if required services are absent; later phase verification cannot silently pass before it exists.

Start local services with `docker compose up -d`. For native Anvil instead of the compose service: `anvil --host 127.0.0.1 --port 18545 --chain-id 31337`. All known Anvil keys are development-only; never use them for real funds. Tests deploy fresh token/escrow contracts and isolate PostgreSQL schemas. They do not contact mainnet.

Copy `.env.example` to an ignored local configuration and supply actual deployment addresses, an isolated test relayer key and generated advertiser API tokens. Export the environment, then `make api` and `make worker`. The worker registers created orders and processes requested refunds; confirmed chain logs create ledger effects. Open `http://127.0.0.1:8080` for the campaign UI. The wallet payer funds the registered order through the escrow; the UI does not request private keys.

Internal audits use actual [Nemesis/Pashov workflows](docs/runbooks/contract-audit.md). Audits and code-review closure are separate from `make verify`: source-bound original reports, verified findings and reviewer closure are required. The external publisher beta is separate scope.
