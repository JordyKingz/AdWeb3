# AdWeb3 frontend

Nuxt 4 app for the marketing site and the dApp: brands fund campaigns, players link a wallet and follow their rewards.

## Run it

Use Node 24 (`.nvmrc` pins the locally verified release) and Bun. Dependencies are pinned in `bun.lock`.

```bash
bun install --frozen-lockfile
cp .env.example .env
bun run dev        # http://localhost:3000
```

`bun run typecheck` and `bun run build` must both pass before merging.

## Public landing on Vercel

Import `JordyKingz/AdWeb3` in Vercel and set **Root Directory** to `frontend`, **Node.js Version** to `24.x`, and **Production Branch** to `main`. Keep files outside the Root Directory excluded from the build. `vercel.json` configures:

| Setting | Value |
|---------|-------|
| Framework preset | Other (static output) |
| Install | `bun install --frozen-lockfile` |
| Build | `bun run generate:landing && bun run verify:landing` |
| Output | `.output/public` |
| Required environment variables | None |

The landing build enables `NUXT_PUBLIC_LANDING_ONLY=true`, includes only the `/` page, displays the development notice, and points calls to action to pilot information. Legacy dApp URLs redirect to that information. Wallet reconnection is disabled. Nitro's mock API handlers are not deployed; this build contains static files only. Do not change the build to `nuxt build` to publish the unfinished dApp.

Run `bun run generate:landing` followed by `bun run verify:landing` locally to check the public output. `bun run dev` still provides the full development UI.

The landing command explicitly selects Nitro's `static` preset so local and Vercel builds both write to `.output/public`. Without it, Vercel auto-detection selects `vercel-static` and writes to `.vercel/output/static`, which does not match the verification script or configured output directory. To check the hosted build behavior locally, run `VERCEL=1 CI=1 bun run generate:landing && bun run verify:landing` from a clean output directory.

All `NUXT_PUBLIC_*` values are visible to visitors and included at build time. Leave the public landing's environment-variable list empty. Backend API tokens, database URLs, relayer keys and future session secrets belong only in the backend's private runtime configuration. See [repository security guidance](../SECURITY.md).

## Routes

| Route | What it is |
|-------|------------|
| `/` | Landing page |
| `/brands/campaigns/new` | Four-step campaign setup ending in escrow funding |
| `/brands/campaigns/:id` | Campaign report: delivery funnel, invalid traffic, escrow split, pause/resume |
| `/connect` | Player wallet link with Sign-In with Ethereum |
| `/dashboard` | Player rewards: pending, finalized and paid balances; activity; payouts |

## Web3

- **Wallets**: discovered through EIP-6963 (falls back to `window.ethereum`), driven with viem. No wallet SDK.
- **Wallet linking** (`/connect` → `server/api/auth/*`): single-use nonce bound to the requesting browser, a SIWE message
  checked for domain, URI, chain and expiry, and an EOA signature check. A failed or replayed attempt burns the nonce.
  Contract wallets are rejected, as the pilot plan specifies.
- **Campaign funding** follows `contracts/src/CampaignEscrow.sol`: create the order, ERC-20 `approve(escrow, cap)`,
  then `fundOrder(orderId, cap)`. Before marking an order funded, the server checks the receipt has a matching
  `Funded` event.
  - With `NUXT_PUBLIC_ESCROW_ADDRESS` and `NUXT_PUBLIC_TOKEN_ADDRESS` set, this runs on-chain. The order must already
    be registered by the operator (`registerOrder`), which the Go API will do.
  - With either unset, the same steps run in **demo mode** and no transaction is sent.

## What is mocked

`server/` is a stand-in for the Go API so the UI can be built and tested now:

- State lives in Nitro's in-memory storage and resets when the server restarts.
- The "Spring launch" campaign and every number on the player dashboard are sample data, marked **Sample data** in the UI.
- New campaigns show zero delivery and stay in creative review.

Replace the handlers under `server/api/` with calls to the Go service as its endpoints land. The response shapes are in
`shared/types/api.ts`.

## Design

The design tokens are CSS custom properties in `app/assets/css/main.css`. They are white, ink `#0B0B0C` and one lime
accent `#D4FF3A`, with Geist and Geist Mono. Money is always handled as integer token units (6 decimals) and formatted
only for display.
