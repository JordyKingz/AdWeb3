import { defineChain, type Chain } from 'viem'
import { avalanche, avalancheFuji } from 'viem/chains'

export interface ChainConfig {
  chainId: number | string
  chainName: string
  rpcUrl: string
  explorerUrl: string
}

/** Builds the viem chain for the configured network, keeping the configured RPC and explorer. */
export function resolveChain(config: ChainConfig): Chain {
  const id = Number(config.chainId)
  const known = [avalanche, avalancheFuji].find((c) => c.id === id)
  return defineChain({
    id,
    name: config.chainName || known?.name || `Chain ${id}`,
    nativeCurrency: known?.nativeCurrency ?? { name: 'Avalanche', symbol: 'AVAX', decimals: 18 },
    rpcUrls: { default: { http: [config.rpcUrl] } },
    blockExplorers: { default: { name: 'Explorer', url: config.explorerUrl } },
    testnet: known?.testnet ?? id !== avalanche.id,
  })
}
