# AdWeb3 project brief

**Proposed Team1 Mini Grant project · USD 7,500 draft request · Eight-week target requiring Android re-estimation and applicant confirmation**

AdWeb3 would connect direct advertising in mobile games to auditable funding and reward settlement on Avalanche. Android is the confirmed first platform, with a shared Go backend and Android SDK foundation. A bounded Android adaptation of Candy Blocks is the proposed first game; the existing native iOS code is reference material, not a ready Android build.

## Product

Advertisers prepay campaign budgets in USDT. The platform serves approved sponsor cards within the purchased dates, validates delivery, and bills only the accepted earned amount. Unused or refundable budgets remain protected. In owned games, earned net advertising revenue is allocated 30% to players, 60% to operations, and 10% to contingency reserves.

Games remain playable offline. Monetary eligibility requires online validation, and a wallet signature proves ownership before payout. The first reward model uses non-transferable points and funded USDT settlement; a tradable platform token is not part of this grant scope.

## Why this project

The proposed value is a clear connection between advertiser budgets, delivered placements, and funded rewards across a portfolio of games. The design avoids depending on token inflation to fund player promises. It does not claim that tiny payments alone will attract or retain players: game quality, advertiser demand, and reward economics are explicit validation questions.

## Avalanche use

The proposed production architecture uses Avalanche C-Chain and native USD₮ for campaign escrow and reward settlement. High-volume gameplay and advertising events stay off-chain. The Go service supplies validated settlement instructions; contracts constrain funding and payment, but do not independently verify ad views.

The proposed grant demonstration runs on Fuji with a test token that has no monetary value. No mainnet funding or production security claim is included in the demonstration.

## Existing work and proposed outcome

Existing work comprises a researched AdWeb3 feature spec and a separate mobile-app portfolio. The advertising platform itself has not been implemented. Candy Blocks supplies an existing iOS game design and deterministic rules reference; its Android adaptation remains unbuilt and must be scoped and estimated; current gameplay/demo evidence still needs packaging for external review.

The grant would produce a reviewable testnet journey: campaign funding → eligible sponsor card → validated activity → earned revenue allocation → verified-wallet test payout → reconciliation report. Evidence would include source references, automated verification, a demo recording, and transaction links. The proposed Android client covers one playable mode with a bounded ruleset, local progress, and one sponsor placement, without full parity with the iOS game. The integration would provide a reusable Android SDK foundation, a versioned interface, an integration example, and data-flow documentation. Scope and estimates are detailed in the accompanying milestone proposal.

## Commercial and launch constraints

The initial pricing hypothesis is 5 USDT per 1,000 valid sponsor-card impressions with a small prepaid campaign budget. No advertiser commitments or platform traction are asserted. The pilot would seek evidence of willingness to pay and assess whether the player share improves retention enough to justify its cost and friction.

Android is the confirmed first platform; Google Play declarations and the applicability of its blockchain and real-money policies still require assessment for the exact product. The funded outcome is proposed as a testnet validation exercise, not a guarantee of Google Play acceptance or real-money launch. Production requires an agreed permitted product, market selection, and separately scoped security/compliance work.

## Publisher SDK expansion

The longer-term product lets other app developers show AdWeb3 ads through an SDK, including ordinary apps without player wallets or rewards. The owner-confirmed commercial model is free SDK access and an 80% publisher / 20% AdWeb3 share of accepted earned net revenue. An optional reward-enabled example allocates 30% to players, 50% to the publisher, and 20% to AdWeb3, subject to app/market eligibility and advertiser agreement. Only the optional player allocation remains a proposal. Owner approval of free access and the 20% commission is not a signed publisher contract or a market benchmark. They do not replace the owned-game 30/60/10 split.

The grant supports only the reusable one-game SDK foundation. A production SDK, external-publisher onboarding, tenant isolation verification, publisher statements/payouts, and additional platforms require a separately scoped and funded beta. No partners or SDK adoption are asserted.
