import { createPublicClient, http } from 'viem'
import { resolveChain } from '#shared/utils/chain'
import type { FundingMode } from '#shared/types/api'

export function usePublicClient() {
  const config = useRuntimeConfig().public
  return createPublicClient({ chain: resolveChain(config), transport: http(config.rpcUrl) })
}

/** On-chain funding needs both a deployed escrow and its token; otherwise orders run in demo mode. */
export function currentFundingMode(): FundingMode {
  const config = useRuntimeConfig().public
  return config.escrowAddress && config.tokenAddress ? 'onchain' : 'demo'
}
