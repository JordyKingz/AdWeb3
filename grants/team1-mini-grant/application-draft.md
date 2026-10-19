# AdWeb3 application answer draft

**Working copy.** These are suggested answer topics, not transcribed form questions. Confirm the decisions in [founder details](founder-details.md), then adapt to the actual field limits. The proposed request is USD 7,500; no funding or commercial commitment is implied.

## Project name

AdWeb3 — working name, pending applicant confirmation.

## Short description

Advertiser-funded game rewards on Avalanche, with transparent campaign escrow and a shared Go backend for mobile games.

## Project overview

AdWeb3 is a proposed advertising and reward platform for mobile games. Companies would fund campaigns in USDT, a shared Go backend would select eligible advertisements and validate delivery, and accepted earned revenue in owned games would be allocated 30% to players, 60% to operations, and 10% to contingency reserves. Refundable advertiser deposits would remain protected until the corresponding campaign obligations are earned.

Android is the confirmed first platform, with iOS deferred. Gameplay remains available offline; monetary eligibility requires online validation. Players would accumulate non-transferable points and, where permitted, receive funded USDT rewards after settlement. A tradable platform token is deferred. The proposed grant project is a testnet demonstration of this funding and reconciliation flow, using a bounded Android adaptation of Candy Blocks as the first candidate game; no Android build is currently verified.

## Problem and intended users

Small game teams need a manageable way to sell direct sponsorships and explain how advertiser money is used. Advertisers need clear delivery terms and records, while players need reward balances that are actually funded rather than promises based on unlimited token emissions.

AdWeb3 would initially serve advertisers buying placements in the applicant's own games. The pilot would test whether unobtrusive sponsor cards and small rewards can coexist with enjoyable gameplay and acceptable advertiser results. The broader product now includes a publisher SDK for other companies' apps. The grant includes only its reusable one-game integration foundation; external publisher onboarding and settlement remain outside the funded pilot.

## Why Avalanche

Avalanche C-Chain offers an EVM environment for campaign escrow and payout contracts, alongside a published native USD₮ deployment. The proposed design keeps frequent gameplay and ad-validation events off-chain while recording funding and settlement on-chain. This lets the project test an application use case for stablecoin transactions without putting every gameplay action on a blockchain.

The grant demonstration would use Fuji testnet and a clearly identified test token, not real USD₮. Mainnet asset integration, payout costs, wallet compatibility, and production security would be verified separately before real funds are accepted. [Tether deployment reference](https://tether.to/en/supported-protocols/), [Avalanche integration reference](https://build.avax.network/docs/primary-network/exchange-integration).

## Current stage and evidence

AdWeb3 is at the researched specification stage. The platform API, escrow contracts, reward ledger, and advertising integration are not built or deployed yet. Existing local mobile projects include Candy Blocks, Train Chaos, Farm Land, Ages of Time, and Photo Drawing. Candy Blocks has a native iOS codebase and deterministic rules engine that can inform an Android adaptation, but cannot be treated as an existing Android integration; this does not establish a live AdWeb3 product.

No verified platform user, revenue, advertiser, or mainnet transaction metrics are available for this application. Public evidence links and founder credentials are still to be supplied. A current gameplay recording and a shareable repository/demo will distinguish existing work from the new grant deliverables.

## Proposed funded deliverables

1. A testnet campaign-funding and ad-decision prototype that attributes deposits once, applies campaign expiry, and protects unused budgets.
2. One bounded playable Android game slice (candidate: Candy Blocks), with local progress and one sponsor placement, integrated through a reusable Android SDK foundation, demonstrating sponsor placement, online activity validation, replay rejection, and wallet-ownership verification while preserving offline gameplay. Include a versioned interface, integration example, and data-flow/privacy notes; no full iOS-game parity, production SDK, or multi-platform support promise. Android implementation and device validation are included work, not claimed prior assets.
3. A complete testnet settlement demonstration, retry-safe payout handling, financial reconciliation report, and documented findings from advertiser/player discovery.

The proposed target is eight weeks after the grant scope is agreed and work begins, subject to Android scope re-estimation and applicant capacity confirmation; the earlier iOS integration estimate does not establish that this target is feasible. [Milestones and budget](milestones-and-budget.md) defines the proposed acceptance evidence; this timeline is not a Team1 requirement or an existing commitment.

## Funding request and use of funds

The proposed request is **USD 7,500**: USD 3,300 for backend/contracts, USD 1,600 provisionally allocated to the bounded Android game slice and SDK foundation, USD 1,200 for validation and bounded security-review work, USD 400 for infrastructure, and USD 1,000 for discovery, documentation, and demo preparation. These are retained draft allocations requiring Android-specific re-estimation, applicant confirmation, and any necessary supplier quotes before submission. The USD 1,600 allowance is not a validated Android port estimate; revise scope, allocations, total, or timeline if needed.

The grant would fund development and validation. It would not be treated as advertiser revenue, token liquidity, or a recurring player payout subsidy. The 30/60/10 advertising allocation is separate from this budget. A full production audit, legal engagement, and any mainnet expansion would need separately scoped funding if they exceed this pilot's resources.

## Expected ecosystem contribution

The project would demonstrate a traceable stablecoin-funded advertising flow and a reusable integration pattern for mobile games on Avalanche. Reported activity would distinguish development/testnet transactions from real users and commercial activity. Success would initially mean a working, reviewable funding-to-settlement demonstration, not a fabricated transaction-volume target.

Selected integration documentation or sample interfaces could be published for other builders, subject to the applicant's IP/licensing decision. No open-source commitment is made until that scope is confirmed.

## Sustainability and commercial validation

The proposed initial offer is a seven-day campaign with a 25-USDT prepaid budget, billed at 5 USDT per 1,000 valid sponsor-card impressions, with unused funds refundable under agreed terms. This price is a hypothesis to test with advertisers, not validated demand.

At illustrative assumptions of six daily opportunities per player, 70% paid fill, 95% accepted delivery, and a 5-USDT net CPM, the player share is about 0.18 USDT per 30 active days. Small rewards must complement an enjoyable game; they are not presented as meaningful income. The pilot would measure advertiser willingness to pay, player response, invalid traffic, operational costs, and payout friction before scaling.

The owner-approved publisher model uses free SDK access with a 20% commission on accepted earned net revenue from partner apps, leaving 80% to the publisher for ordinary ads. Player rewards are optional and require a separately agreed allocation; an illustrative 30% player / 50% publisher / 20% platform model does not charge an additional commission or contingency percentage. The owned-game 30/60/10 stays unchanged. The optional player allocation and remaining contractual details need agreement, and actual publisher demand still needs validation; no publisher commitment is claimed. External onboarding and settlement require separate funding. Adding app inventory does not establish advertiser demand or guaranteed fill.

## Risks and mitigations

Android-first does not establish Google Play approval. Tokenized-asset earning requires the relevant declaration, and potential earnings must not be promoted or glamorized. The exact product also needs assessment against applicable real-money rules and selected markets. The proposed game has no paid entry or wagering. The grant remains a testnet demonstration, without promising Play approval or live cash payouts. iOS is deferred. [Google blockchain policy](https://support.google.com/googleplay/android-developer/answer/13607354?hl=en), [Google real-money games policy](https://support.google.com/googleplay/android-developer/answer/9877032?hl=en).

Other material risks are forged activity, wallet farms, duplicate accounting, and refund liabilities. The proposed controls include online eligibility, bounded sessions, event deduplication, delayed settlement, segregated ledger balances, and unique payout IDs. On-chain payment records do not prove that a human viewed an advertisement; the backend remains responsible for delivery validation.

## Applicant-specific answers still required

See [founder details](founder-details.md) for the applicant identity, relevant experience, team roles, availability, project links, previous funding, launch market, proposed grant amount approval, and publication preferences. Do not submit placeholder values or infer these from local directory names.
