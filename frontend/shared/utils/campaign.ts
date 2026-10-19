/**
 * Commercial constants shared by the UI and the mock API.
 * Amounts are integer token units (USD₮ has 6 decimals) to match the escrow contract.
 */

export const TOKEN_DECIMALS = 6

/** Test-campaign rate: 5 USD₮ per 1,000 validated impressions. */
export const CPM_UNITS = 5_000_000n

export const MIN_BUDGET = 5
export const MAX_BUDGET = 100_000
export const TEST_BUDGET = 25

export const DURATIONS = [7, 14, 30] as const
export const FREQUENCY_CAPS = [2, 3, 5] as const
export const CTA_OPTIONS = ['Learn more', 'Shop now', 'Download', 'Play now'] as const
export const HEADLINE_MAX = 40

export const PLACEMENTS = [
  {
    id: 'home',
    name: 'Candy Blocks · Home card',
    short: 'Home',
    format: 'Android · 1200×628',
    description: 'On the home screen between rounds. The first thing players see when they open the game.',
  },
  {
    id: 'results',
    name: 'Candy Blocks · Results card',
    short: 'Results',
    format: 'Android · 1200×628',
    description: 'After a round ends, next to the score. Never over live gameplay.',
  },
] as const

export type PlacementId = (typeof PLACEMENTS)[number]['id']
export type Duration = (typeof DURATIONS)[number]
export type FrequencyCap = (typeof FREQUENCY_CAPS)[number]
export type CtaOption = (typeof CTA_OPTIONS)[number]

/** Basis points of earned revenue; contingency takes the remainder, exactly like CampaignEscrow.releaseEarned. */
export const SPLIT_BPS = { players: 3000n, operations: 6000n } as const

export function splitEarned(earned: bigint) {
  const players = (earned * SPLIT_BPS.players) / 10_000n
  const operations = (earned * SPLIT_BPS.operations) / 10_000n
  return { players, operations, contingency: earned - players - operations }
}

export function billedUnits(validatedImpressions: number): bigint {
  return (BigInt(validatedImpressions) * CPM_UNITS) / 1000n
}

export function maxImpressions(budgetUnits: bigint): number {
  return Number((budgetUnits * 1000n) / CPM_UNITS)
}

/** Budgets are whole or two-decimal USD₮ amounts. */
export const BUDGET_PATTERN = /^\d{1,6}(\.\d{1,2})?$/

export function isValidBudget(value: string): boolean {
  if (!BUDGET_PATTERN.test(value)) return false
  const n = Number(value)
  return n >= MIN_BUDGET && n <= MAX_BUDGET
}
