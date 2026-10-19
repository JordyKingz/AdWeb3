# Public repository and configuration

This repository is public. Commit only placeholders and public configuration. `.env.example` files are examples; `.env`, local overrides, Vercel project metadata, signing keys and service-account files are ignored. Never force-add an ignored secret file.

The Vercel landing is a static export and needs no secrets. Anything named `NUXT_PUBLIC_*` is public, including values entered through Vercel's environment settings. Never use that prefix for API tokens, database credentials, session secrets, relayer private keys or seed phrases. Keep backend secrets in its separate runtime configuration or secret store.

The fixed local PostgreSQL password and standard Anvil keys used by integration tests are public development fixtures. They must never secure a shared environment or hold real assets. Generate distinct credentials and wallets for other environments.

Before each public push, inspect staged files and run a secret scan. Install [Gitleaks](https://github.com/gitleaks/gitleaks) (locally verified: 8.30.1), then enable the versioned pre-push hook with `git config --local core.hooksPath .githooks`. The hook scans Git history and blocks pushes when Gitleaks is missing or reports findings. The local setup can also use the checksum-verified binary under ignored `.tools/gitleaks/`.

GitHub secret scanning and push protection are enabled for this repository; these checks cannot recognize every kind of secret. A pinned CI workflow template is saved at `scripts/security/github-actions-secrets.yml`. To enable it, authorize GitHub CLI's `workflow` scope (for example, `gh auth refresh --hostname github.com --scopes workflow`), move the template to `.github/workflows/public-repo.yml`, and commit it. The current CLI login cannot publish Actions workflows, so the template is intentionally inactive.

If a real credential is exposed, revoke or rotate it immediately, including when it was later removed from a file. Remove it from history where appropriate after rotation. Do not paste credentials into issues, audit reports or application documents. Complete private grant applicant details outside this public repository.
