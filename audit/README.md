# Internal audit evidence

Run `make audit-tools-update`, commit the clean contract candidate, then `make audit-contracts PHASE=1` (or 4/5). This invokes actual Nemesis and Pashov agent skills in isolated Git snapshots. Runtime success is not audit closure: inspect native pass coverage, set completion truthfully, triage/fix findings and retest. `make verify-contract-audit AUDIT_MANIFEST=...` requires explicit review closure and evidence. See [the runbook](../docs/runbooks/contract-audit.md).

Retain `reports/` and `tools-lock.json`; disposable `work/` snapshots are ignored. Never add secrets. Agent credentials are supplied by the configured runtime, not committed to the repository.
