import { generateSiweNonce } from 'viem/siwe'

const NONCE_TTL_MS = 5 * 60 * 1000

/** Issues a single-use sign-in nonce, bound to this browser via a short-lived cookie. */
export default defineEventHandler(async (event) => {
  let preId = getCookie(event, PRE_SESSION_COOKIE)
  if (!preId) {
    preId = randomId()
    setCookie(event, PRE_SESSION_COOKIE, preId, cookieOptions(30 * 60))
  }

  const nonce = generateSiweNonce()
  const issuedAt = new Date()
  const expirationTime = new Date(issuedAt.getTime() + NONCE_TTL_MS)
  await useStorage().setItem(`siwe:${nonce}`, { preId, expiresAt: expirationTime.getTime() })

  setResponseHeader(event, 'cache-control', 'no-store')
  return { nonce, issuedAt: issuedAt.toISOString(), expirationTime: expirationTime.toISOString() }
})
