import { getAddress, isAddress, keccak256, parseUnits, stringToHex, zeroAddress, type Address, type Hex } from 'viem'
import type { CampaignDraft, CampaignReport, CampaignStatus, FundingMode } from '#shared/types/api'
import {
  CPM_UNITS,
  CTA_OPTIONS,
  DURATIONS,
  FREQUENCY_CAPS,
  HEADLINE_MAX,
  PLACEMENTS,
  TOKEN_DECIMALS,
  billedUnits,
  isValidBudget,
  maxImpressions,
  splitEarned,
  type PlacementId,
} from '#shared/utils/campaign'

export interface StoredCampaign {
  id: string
  name: string
  brand: string
  url: string
  start: string
  durationDays: number
  placements: PlacementId[]
  incentivized: boolean
  frequencyCap: number
  headline: string
  cta: string
  budgetUnits: string
  payer: Address
  orderId: Hex
  terms: Hex
  fundingMode: FundingMode
  status: CampaignStatus
  txHash: Hex | null
  fundedAt: string | null
  createdAt: string
  sample?: boolean
}

export const SAMPLE_CAMPAIGN_ID = 'spring-launch'

/** Delivery for the seeded sample campaign. A real deployment reads delivery from the Go API. */
const SAMPLE_DELIVERY = {
  served: 2540,
  rendered: 2402,
  viewable: 2279,
  daily: [412, 538, 611, 587],
  placements: {
    home: { validated: 1512, clicks: 31 },
    results: { validated: 636, clicks: 9 },
  } as Record<PlacementId, { validated: number; clicks: number }>,
  invalid: [
    { reason: 'Duplicate events', count: 58 },
    { reason: 'App in background', count: 41 },
    { reason: 'Expired ad decision', count: 19 },
    { reason: 'Unverified session', count: 13 },
  ],
}

function sampleCampaign(): StoredCampaign {
  return {
    id: SAMPLE_CAMPAIGN_ID,
    name: 'Spring launch',
    brand: 'Your Brand',
    url: 'https://yourbrand.com/spring',
    start: isoDate(utcDay(-3)),
    durationDays: 7,
    placements: ['home', 'results'],
    incentivized: true,
    frequencyCap: 3,
    headline: 'Your brand, seen by real players',
    cta: 'Learn more',
    budgetUnits: '25000000',
    payer: '0x8a3f6de1b2c4a5970fe3de0c6b1a9a4c7e2d21cb',
    orderId: keccak256(stringToHex('adweb3:sample-order:spring-launch')),
    terms: keccak256(stringToHex('adweb3:sample-terms:v1')),
    fundingMode: 'demo',
    status: 'live',
    txHash: null,
    fundedAt: utcDay(-4).toISOString(),
    createdAt: utcDay(-4).toISOString(),
    sample: true,
  }
}

const storageKey = (id: string) => `campaigns:${id}`

export async function getCampaign(id: string): Promise<StoredCampaign | null> {
  const stored = await useStorage().getItem<StoredCampaign>(storageKey(id))
  if (stored) return stored
  return id === SAMPLE_CAMPAIGN_ID ? sampleCampaign() : null
}

export async function saveCampaign(campaign: StoredCampaign): Promise<void> {
  await useStorage().setItem(storageKey(campaign.id), campaign)
}

export async function listCampaigns(): Promise<StoredCampaign[]> {
  const storage = useStorage()
  const keys = await storage.getKeys('campaigns')
  const items = (await Promise.all(keys.map((k) => storage.getItem<StoredCampaign>(k)))).filter(
    (c): c is StoredCampaign => !!c,
  )
  if (!items.some((c) => c.id === SAMPLE_CAMPAIGN_ID)) items.push(sampleCampaign())
  return items.sort((a, b) => b.createdAt.localeCompare(a.createdAt))
}

export function validateDraft(
  input: Partial<Record<keyof CampaignDraft, unknown>> | null,
  fundingMode: FundingMode,
): CampaignDraft {
  const fail = (message: string): never => {
    throw createError({ statusCode: 400, statusMessage: message })
  }
  const text = (value: unknown, max: number, label: string): string => {
    const v = typeof value === 'string' ? value.trim() : ''
    if (!v) return fail(`${label} is required`)
    if (v.length > max) return fail(`${label} is too long`)
    return v
  }

  const name = text(input?.name, 80, 'Campaign name')
  const brand = text(input?.brand, 60, 'Brand name')
  const url = text(input?.url, 300, 'Landing page')
  let parsedUrl: URL | null = null
  try {
    parsedUrl = new URL(url)
  } catch {
    fail('Landing page must be a valid URL')
  }
  if (parsedUrl?.protocol !== 'https:') fail('Landing page must use https')

  const start = typeof input?.start === 'string' ? input.start : ''
  if (!/^\d{4}-\d{2}-\d{2}$/.test(start) || Number.isNaN(Date.parse(`${start}T00:00:00Z`))) fail('Start date is invalid')
  if (start < isoDate(utcDay(0))) fail('Start date cannot be in the past')

  const durationDays = Number(input?.durationDays)
  if (!(DURATIONS as readonly number[]).includes(durationDays)) fail('Unsupported serving window')

  const placements = Array.isArray(input?.placements) ? [...new Set(input.placements)] : []
  const known = PLACEMENTS.map((p) => p.id) as string[]
  if (!placements.length || placements.some((p) => typeof p !== 'string' || !known.includes(p))) fail('Choose at least one placement')

  if (typeof input?.incentivized !== 'boolean') fail('Audience type is required')

  const frequencyCap = Number(input?.frequencyCap)
  if (!(FREQUENCY_CAPS as readonly number[]).includes(frequencyCap)) fail('Unsupported frequency cap')

  const headline = text(input?.headline, HEADLINE_MAX, 'Headline')
  const cta = typeof input?.cta === 'string' ? input.cta : ''
  if (!(CTA_OPTIONS as readonly string[]).includes(cta)) fail('Unsupported button label')

  const budget = typeof input?.budget === 'string' ? input.budget : ''
  if (!isValidBudget(budget)) fail('Budget must be between 5 and 100,000 USD₮')

  const payer = typeof input?.payer === 'string' && isAddress(input.payer) ? getAddress(input.payer) : fail('Connect a wallet first')
  // Demo orders may be previewed without a wallet; on-chain orders need a real payer.
  if (payer === zeroAddress && fundingMode === 'onchain') fail('Connect a wallet first')

  return {
    name,
    brand,
    url,
    start,
    durationDays: durationDays as CampaignDraft['durationDays'],
    placements: placements as PlacementId[],
    incentivized: input?.incentivized as boolean,
    frequencyCap: frequencyCap as CampaignDraft['frequencyCap'],
    headline,
    cta: cta as CampaignDraft['cta'],
    budget,
    payer,
  }
}

export function newCampaign(draft: CampaignDraft, fundingMode: FundingMode): StoredCampaign {
  const slug = draft.name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '').slice(0, 32) || 'campaign'
  const id = `${slug}-${randomId().slice(0, 6)}`
  const budgetUnits = parseUnits(draft.budget, TOKEN_DECIMALS)
  // The order id and terms hash are what an operator registers on-chain via CampaignEscrow.registerOrder.
  const terms = keccak256(
    stringToHex(
      JSON.stringify({
        v: 1,
        cpmUnits: CPM_UNITS.toString(),
        budgetUnits: budgetUnits.toString(),
        start: draft.start,
        durationDays: draft.durationDays,
        placements: draft.placements,
        incentivized: draft.incentivized,
        frequencyCap: draft.frequencyCap,
        split: '30/60/10',
      }),
    ),
  )
  return {
    id,
    name: draft.name,
    brand: draft.brand,
    url: draft.url,
    start: draft.start,
    durationDays: draft.durationDays,
    placements: draft.placements,
    incentivized: draft.incentivized,
    frequencyCap: draft.frequencyCap,
    headline: draft.headline,
    cta: draft.cta,
    budgetUnits: budgetUnits.toString(),
    payer: draft.payer,
    orderId: keccak256(stringToHex(`adweb3:order:${id}:${randomId()}`)),
    terms,
    fundingMode,
    status: 'awaiting_funding',
    txHash: null,
    fundedAt: null,
    createdAt: new Date().toISOString(),
  }
}

export function buildReport(c: StoredCampaign): CampaignReport {
  const live = c.status === 'live' || c.status === 'paused'
  const delivery = c.sample && live ? SAMPLE_DELIVERY : null

  const placements = c.placements.map((id) => {
    const meta = PLACEMENTS.find((p) => p.id === id)!
    const d = delivery?.placements[id] ?? { validated: 0, clicks: 0 }
    return { name: `${meta.short} card`, validated: d.validated, clicks: d.clicks, billedUnits: billedUnits(d.validated).toString() }
  })
  const validated = placements.reduce((sum, p) => sum + p.validated, 0)
  const billed = billedUnits(validated)
  const budget = BigInt(c.budgetUnits)
  const refundable = c.status === 'awaiting_funding' ? 0n : budget > billed ? budget - billed : 0n
  const split = splitEarned(billed)

  const start = new Date(`${c.start}T00:00:00Z`)
  const end = new Date(start.getTime() + c.durationDays * DAY_MS)
  const dayIndex = Math.floor((utcDay(0).getTime() - start.getTime()) / DAY_MS)
  const daily = Array.from({ length: c.durationDays }, (_, i) => ({
    date: isoDate(new Date(start.getTime() + i * DAY_MS)),
    value: delivery && i < delivery.daily.length ? delivery.daily[i]! : null,
    today: i === dayIndex,
  }))

  return {
    id: c.id,
    name: c.name,
    brand: c.brand,
    status: c.status,
    sample: !!c.sample,
    placementsLabel: placements.map((p) => p.name.replace(' card', '')).join(' + '),
    start: c.start,
    end: isoDate(end),
    durationDays: c.durationDays,
    dayOfWindow: Math.min(Math.max(dayIndex + 1, 0), c.durationDays),
    budgetUnits: c.budgetUnits,
    billedUnits: billed.toString(),
    refundableUnits: refundable.toString(),
    split: {
      players: split.players.toString(),
      operations: split.operations.toString(),
      contingency: split.contingency.toString(),
    },
    delivery: {
      served: delivery?.served ?? 0,
      rendered: delivery?.rendered ?? 0,
      viewable: delivery?.viewable ?? 0,
      validated,
      maxImpressions: maxImpressions(budget),
      daily,
    },
    invalid: delivery?.invalid ?? [],
    placements,
    terms: {
      cpmUnits: CPM_UNITS.toString(),
      frequencyCap: c.frequencyCap,
      incentivized: c.incentivized,
      creativeStatus: live ? 'approved' : 'in_review',
    },
    funding: { mode: c.fundingMode, orderId: c.orderId, txHash: c.txHash, fundedAt: c.fundedAt },
  }
}
