import { getAddress, isHex, verifyMessage } from 'viem'
import { parseSiweMessage, validateSiweMessage } from 'viem/siwe'

/**
 * Verifies a Sign-In with Ethereum message and links the wallet to a session.
 * EOA signatures only, as the pilot plan specifies; contract wallets are rejected.
 */
export default defineEventHandler(async (event) => {
  const body = await readBody<{ message?: unknown; signature?: unknown }>(event)
  const message = body?.message
  const signature = body?.signature
  if (typeof message !== 'string' || message.length > 2000 || typeof signature !== 'string' || !isHex(signature)) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid sign-in request' })
  }

  const fields = parseSiweMessage(message)
  if (!fields.nonce || !fields.address || !fields.expirationTime) {
    throw createError({ statusCode: 400, statusMessage: 'Malformed sign-in message' })
  }

  // Single use: consume the nonce before checking anything else.
  const storage = useStorage()
  const nonceKey = `siwe:${fields.nonce}`
  const record = await storage.getItem<{ preId: string; expiresAt: number }>(nonceKey)
  await storage.removeItem(nonceKey)
  const preId = getCookie(event, PRE_SESSION_COOKIE)
  if (!record || record.expiresAt < Date.now() || !preId || record.preId !== preId) {
    throw createError({ statusCode: 401, statusMessage: 'This sign-in request expired or was started in another browser' })
  }

  const chainId = Number(useRuntimeConfig(event).public.chainId)
  const url = getRequestURL(event)
  const matchesSite =
    validateSiweMessage({ message: fields, domain: url.host, nonce: fields.nonce, time: new Date() }) &&
    fields.chainId === chainId &&
    !!fields.uri?.startsWith(url.origin)
  if (!matchesSite) {
    throw createError({ statusCode: 401, statusMessage: 'The message does not match this site or network' })
  }

  const valid = await verifyMessage({ address: fields.address, message, signature }).catch(() => false)
  if (!valid) throw createError({ statusCode: 401, statusMessage: 'The signature does not match the wallet' })

  deleteCookie(event, PRE_SESSION_COOKIE, { path: '/' })
  const session = await createPlayerSession(event, getAddress(fields.address), chainId)
  return { address: session.address, chainId: session.chainId, linkedAt: session.linkedAt }
})
