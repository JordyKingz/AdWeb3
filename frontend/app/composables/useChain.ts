import { createPublicClient, http, type Address } from 'viem'
import { resolveChain } from '#shared/utils/chain'
import type { FundingMode } from '#shared/types/api'

let client: ReturnType<typeof createPublicClient> | undefined

/** Network configuration plus a lazily created public client for reads. */
export function useChain() {
  const config = useRuntimeConfig().public
  const chain = resolveChain(config)
  const tokenAddress = (config.tokenAddress || null) as Address | null
  const escrowAddress = (config.escrowAddress || null) as Address | null
  const fundingMode: FundingMode = tokenAddress && escrowAddress ? 'onchain' : 'demo'

  function publicClient() {
    client ??= createPublicClient({ chain, transport: http(config.rpcUrl) })
    return client
  }

  return {
    chain,
    chainName: config.chainName,
    tokenSymbol: config.tokenSymbol,
    tokenDecimals: Number(config.tokenDecimals),
    tokenAddress,
    escrowAddress,
    fundingMode,
    publicClient,
    txUrl: (hash: string) => `${config.explorerUrl}/tx/${hash}`,
    addressUrl: (address: string) => `${config.explorerUrl}/address/${address}`,
  }
}
