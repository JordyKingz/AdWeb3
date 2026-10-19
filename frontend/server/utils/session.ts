import type { H3Event } from 'h3'
import type { Address } from 'viem'

/**
 * Wallet-link sessions for the mock API. State lives in Nitro storage (in memory by default);
 * the Go identity service replaces this once it exists.
 */

export const SESSION_COOKIE = 'adweb3_session'
/** Binds a sign-in nonce to the browser that requested it. */
export const PRE_SESSION_COOKIE = 'adweb3_pre'

const SESSION_TTL_S = 60 * 60 * 24 * 7

export interface PlayerSession {
  address: Address
  chainId: number
  linkedAt: string
  expiresAt: number
}

export function cookieOptions(maxAge: number) {
  return { httpOnly: true, sameSite: 'lax' as const, secure: !import.meta.dev, path: '/', maxAge }
}

export function randomId(): string {
  return crypto.randomUUID()
}

export async function createPlayerSession(event: H3Event, address: Address, chainId: number): Promise<PlayerSession> {
  const id = randomId()
  const session: PlayerSession = {
    address,
    chainId,
    linkedAt: new Date().toISOString(),
    expiresAt: Date.now() + SESSION_TTL_S * 1000,
  }
  await useStorage().setItem(`sessions:${id}`, session)
  setCookie(event, SESSION_COOKIE, id, cookieOptions(SESSION_TTL_S))
  return session
}

export async function getPlayerSession(event: H3Event): Promise<PlayerSession | null> {
  const id = getCookie(event, SESSION_COOKIE)
  if (!id) return null
  const session = await useStorage().getItem<PlayerSession>(`sessions:${id}`)
  if (!session || session.expiresAt < Date.now()) {
    await useStorage().removeItem(`sessions:${id}`)
    return null
  }
  return session
}

export async function requirePlayerSession(event: H3Event): Promise<PlayerSession> {
  const session = await getPlayerSession(event)
  if (!session) throw createError({ statusCode: 401, statusMessage: 'Wallet not linked' })
  return session
}

export async function destroyPlayerSession(event: H3Event): Promise<void> {
  const id = getCookie(event, SESSION_COOKIE)
  if (id) await useStorage().removeItem(`sessions:${id}`)
  deleteCookie(event, SESSION_COOKIE, { path: '/' })
}
