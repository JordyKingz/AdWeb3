import { erc20Abi, isAddressEqual, zeroAddress, type Hex } from 'viem'
import { campaignEscrowAbi } from '#shared/utils/escrow'
import type { CampaignDraft, CreatedOrder } from '#shared/types/api'

export type FundingStage = 'idle' | 'creating' | 'approving' | 'approved' | 'depositing' | 'confirming' | 'funded'

/**
 * Order → ERC-20 approve → CampaignEscrow.fundOrder, following contracts/src/CampaignEscrow.sol.
 * Without a configured escrow the same steps run in demo mode and send no transactions.
 */
export function useCampaignFunding() {
  const net = useChain()
  const wallet = useWallet()

  const stage = ref<FundingStage>('idle')
  const order = ref<CreatedOrder | null>(null)
  const txHash = ref<Hex | null>(null)
  const error = ref<string | null>(null)
  const busy = computed(() => ['creating', 'approving', 'depositing', 'confirming'].includes(stage.value))

  /** Call when the draft changes after an order exists: the next approval creates a fresh order. */
  function invalidate() {
    if (stage.value === 'funded' || busy.value || !order.value) return
    order.value = null
    stage.value = 'idle'
  }

  async function approve(draft: Omit<CampaignDraft, 'payer'>) {
    error.value = null
    try {
      if (net.fundingMode === 'onchain' && wallet.wrongChain.value) await wallet.switchChain()
      if (!order.value) {
        stage.value = 'creating'
        order.value = await $fetch<CreatedOrder>('/api/campaigns', {
          method: 'POST',
          body: { ...draft, payer: wallet.address.value ?? zeroAddress },
        })
      }
      stage.value = 'approving'
      if (order.value.fundingMode === 'onchain') {
        const cap = BigInt(order.value.capUnits)
        const allowance = await net.publicClient().readContract({
          address: net.tokenAddress!,
          abi: erc20Abi,
          functionName: 'allowance',
          args: [wallet.address.value!, net.escrowAddress!],
        })
        if (allowance < cap) {
          const hash = await wallet.walletClient().writeContract({
            address: net.tokenAddress!,
            abi: erc20Abi,
            functionName: 'approve',
            args: [net.escrowAddress!, cap],
          })
          await net.publicClient().waitForTransactionReceipt({ hash })
        }
      } else {
        await sleep(700)
      }
      stage.value = 'approved'
    } catch (err) {
      stage.value = 'idle'
      error.value = readableError(err)
    }
  }

  async function deposit() {
    const current = order.value
    if (!current) return
    error.value = null
    try {
      if (current.fundingMode === 'onchain') {
        const cap = BigInt(current.capUnits)
        const [payer, orderCap, funded] = await net.publicClient().readContract({
          address: net.escrowAddress!,
          abi: campaignEscrowAbi,
          functionName: 'orders',
          args: [current.orderId],
        })
        if (!isAddressEqual(payer, wallet.address.value!) || orderCap !== cap) {
          throw new Error('This order is not registered on-chain for your wallet yet. Try again once it has been registered.')
        }
        if (funded > 0n) throw new Error('This order is already funded.')
        stage.value = 'depositing'
        const hash = await wallet.walletClient().writeContract({
          address: net.escrowAddress!,
          abi: campaignEscrowAbi,
          functionName: 'fundOrder',
          args: [current.orderId, cap],
        })
        txHash.value = hash
        stage.value = 'confirming'
        const receipt = await net.publicClient().waitForTransactionReceipt({ hash })
        if (receipt.status !== 'success') throw new Error('The deposit transaction reverted.')
      } else {
        stage.value = 'confirming'
        await sleep(1400)
      }
      await $fetch(`/api/campaigns/${current.id}/fund`, { method: 'POST', body: { txHash: txHash.value } })
      await refreshNuxtData('brand-overview')
      stage.value = 'funded'
    } catch (err) {
      stage.value = 'approved'
      error.value = readableError(err)
    }
  }

  return { stage, order, txHash, error, busy, approve, deposit, invalidate }
}
