import { keccak256, stringToHex } from 'viem'
import type { PlayerSummary } from '#shared/types/api'

/** Sample rewards data until the Go settlement service exposes real balances. */
export function buildPlayerSummary(session: PlayerSession): PlayerSummary {
  const today = utcDay(0)
  const day = (offset: number) => new Date(today.getTime() - offset * DAY_MS)
  const range = (from: number, to: number) => `${shortDay(day(from))} – ${shortDay(day(to))}`

  // Validated points for the last 14 days; null marks offline days (progress saved, no earning).
  const points = [62, 75, null, 88, 41, 96, 70, null, 54, 81, 92, null, 88, 46]

  return {
    sample: true,
    address: session.address,
    linkedAt: session.linkedAt,
    game: { name: 'Candy Blocks', platform: 'Android' },
    pendingPoints: 1284,
    finalizedUnits: '421000',
    paidUnits: '1035000',
    payoutThresholdUnits: '1000000',
    payoutCount: 2,
    lastPayoutAt: isoDate(day(7)),
    nextProcessingAt: utcDay(1).toISOString(),
    dailyCap: 120,
    today: { points: 46, onlineMinutes: 18, sponsorCardsSeen: 5 },
    daily: points.map((p, i) => ({ date: isoDate(day(13 - i)), points: p, today: i === 13 })),
    activity: [
      { kind: 'session', title: 'Online session · 18 min', detail: 'Today · Candy Blocks', amount: '+46', status: 'Validated' },
      { kind: 'session', title: 'Online sessions · 32 min', detail: `${shortDay(day(1))} · Candy Blocks`, amount: '+88', status: 'Validated' },
      {
        kind: 'offline',
        title: 'Offline play · 41 min',
        detail: `${shortDay(day(2))} · progress synced, not eligible for rewards`,
        amount: '0',
        status: 'Not eligible',
      },
      { kind: 'rejected', title: '3 events rejected', detail: `${shortDay(day(3))} · duplicate submission`, amount: '0', status: 'Rejected' },
      {
        kind: 'finalized',
        title: `Epoch ${range(14, 8)} finalized`,
        detail: `${shortDay(day(7))} · campaign revenue cleared`,
        amount: '0.5120',
        status: 'Finalized',
      },
    ],
    payouts: [
      {
        id: 'PO-0192',
        epoch: range(14, 8),
        amountUnits: '512000',
        txHash: keccak256(stringToHex('adweb3:sample-payout:PO-0192')),
        status: 'confirmed',
      },
      {
        id: 'PO-0157',
        epoch: range(28, 22),
        amountUnits: '523000',
        txHash: keccak256(stringToHex('adweb3:sample-payout:PO-0157')),
        status: 'confirmed',
      },
    ],
  }
}
