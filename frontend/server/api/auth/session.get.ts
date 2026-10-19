import type { SessionInfo } from '#shared/types/api'

export default defineEventHandler(async (event): Promise<SessionInfo> => {
  const session = await requirePlayerSession(event)
  return { address: session.address, chainId: session.chainId, linkedAt: session.linkedAt }
})
