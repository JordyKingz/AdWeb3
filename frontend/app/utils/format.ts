import { formatUnits } from 'viem'

export function formatToken(units: bigint | string, fractionDigits = 2, decimals = 6): string {
  const value = Number(formatUnits(BigInt(units), decimals))
  return value.toLocaleString('en-US', { minimumFractionDigits: fractionDigits, maximumFractionDigits: fractionDigits })
}

export function formatInt(n: number): string {
  return n.toLocaleString('en-US')
}

export function formatPercent(part: number, whole: number, fractionDigits = 1): string {
  if (!whole) return '0%'
  return `${((part / whole) * 100).toFixed(fractionDigits)}%`
}

export function shortAddress(address: string): string {
  return `${address.slice(0, 6)}…${address.slice(-4)}`
}

export function shortHash(hash: string): string {
  return `${hash.slice(0, 6)}…${hash.slice(-4)}`
}

/** "Oct 20" from an ISO date or timestamp, in UTC. */
export function shortDate(iso: string, withYear = false): string {
  const d = new Date(iso.length === 10 ? `${iso}T00:00:00Z` : iso)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: withYear ? 'numeric' : undefined, timeZone: 'UTC' })
}

export function utcTime(iso: string): string {
  return new Date(iso).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })
}
