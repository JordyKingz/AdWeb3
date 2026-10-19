#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."
exec python3 scripts/audit/run_contract_audits.py "${1:?phase required}" "${2:-both}"
