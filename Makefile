SHELL := /bin/bash
PHASE ?= 1
DATABASE_URL ?= postgres://postgres:adweb3-local-only@127.0.0.1:55432/adweb3?sslmode=disable
EVM_RPC_URL ?= http://127.0.0.1:18545
export DATABASE_URL EVM_RPC_URL
.PHONY: verify test contracts integration audit-tools-update audit-contracts verify-contract-audit api worker
verify: contracts test integration
	python3 -m unittest discover -s tests/audit -v
	@if [ "$(PHASE)" -gt 1 ]; then echo 'Later phase verification is not implemented yet; cannot pass.' >&2; exit 1; fi
test:
	go test ./...
	go test -race ./...
contracts:
	forge test --root contracts
integration: contracts
	go test -tags=integration ./tests/integration/... -count=1
audit-tools-update:
	./scripts/audit/update-tools.sh
audit-contracts:
	./scripts/audit/run-contract-audits.sh $(PHASE)
verify-contract-audit:
	@test -n "$(AUDIT_MANIFEST)" || (echo 'AUDIT_MANIFEST required' >&2; exit 1)
	python3 scripts/audit/verify-reports.py "$(AUDIT_MANIFEST)"
api:
	go run ./cmd/api
worker:
	go run ./cmd/worker
