export const DAY_MS = 86_400_000

/** Midnight UTC, offset by whole days from today. */
export function utcDay(offsetDays = 0): Date {
  const d = new Date()
  d.setUTCHours(0, 0, 0, 0)
  return new Date(d.getTime() + offsetDays * DAY_MS)
}

export function isoDate(d: Date): string {
  return d.toISOString().slice(0, 10)
}

export function shortDay(d: Date): string {
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', timeZone: 'UTC' })
}
