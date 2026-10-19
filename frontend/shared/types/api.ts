import type { Address, Hex } from 'viem'
import type { CtaOption, Duration, FrequencyCap, PlacementId } from '../utils/campaign'

export type FundingMode = 'onchain' | 'demo'

export interface CampaignDraft {
  name: string
  brand: string
  url: string
  /** UTC start date, YYYY-MM-DD. */
  start: string
  durationDays: Duration
  placements: PlacementId[]
  incentivized: boolean
  frequencyCap: FrequencyCap
  headline: string
  cta: CtaOption
  /** Decimal USD₮ amount, e.g. "25" or "25.50". */
  budget: string
  payer: Address
}

export interface CreatedOrder {
  id: string
  orderId: Hex
  capUnits: string
  terms: Hex
  fundingMode: FundingMode
}

export type CampaignStatus = 'awaiting_funding' | 'in_review' | 'live' | 'paused'

export interface CampaignReport {
  id: string
  name: string
  brand: string
  status: CampaignStatus
  sample: boolean
  placementsLabel: string
  start: string
  end: string
  durationDays: number
  dayOfWindow: number
  budgetUnits: string
  billedUnits: string
  refundableUnits: string
  split: { players: string; operations: string; contingency: string }
  delivery: {
    served: number
    rendered: number
    viewable: number
    validated: number
    maxImpressions: number
    /** One entry per UTC day of the window; null for days not served yet. */
    daily: { date: string; value: number | null; today: boolean }[]
  }
  invalid: { reason: string; count: number }[]
  placements: { name: string; validated: number; clicks: number; billedUnits: string }[]
  terms: { cpmUnits: string; frequencyCap: number; incentivized: boolean; creativeStatus: 'approved' | 'in_review' }
  funding: { mode: FundingMode; orderId: Hex; txHash: Hex | null; fundedAt: string | null }
}

export interface BrandOverview {
  inEscrowUnits: string
  campaigns: { id: string; name: string; status: CampaignStatus }[]
}

export interface SessionInfo {
  address: Address
  chainId: number
  linkedAt: string
}

export type ActivityKind = 'session' | 'offline' | 'rejected' | 'finalized'

export interface PlayerSummary {
  sample: boolean
  address: Address
  linkedAt: string
  game: { name: string; platform: string }
  pendingPoints: number
  finalizedUnits: string
  paidUnits: string
  payoutThresholdUnits: string
  payoutCount: number
  lastPayoutAt: string | null
  nextProcessingAt: string
  dailyCap: number
  today: { points: number; onlineMinutes: number; sponsorCardsSeen: number }
  daily: { date: string; points: number | null; today: boolean }[]
  activity: { kind: ActivityKind; title: string; detail: string; amount: string; status: string }[]
  payouts: { id: string; epoch: string; amountUnits: string; txHash: Hex; status: 'confirmed' | 'submitted' }[]
}
