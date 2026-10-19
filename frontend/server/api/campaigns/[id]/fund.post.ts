import { decodeEventLog, isAddressEqual, isHex, type Address, type Hex } from 'viem'
import { campaignEscrowAbi } from '#shared/utils/escrow'

/**
 * Marks an order funded. On-chain orders must point at a successful transaction that emitted
 * a matching Funded event from the configured escrow; the client's word is never enough.
 */
export default defineEventHandler(async (event) => {
  const campaign = await getCampaign(getRouterParam(event, 'id') ?? '')
  if (!campaign) throw createError({ statusCode: 404, statusMessage: 'Campaign not found' })
  if (campaign.status !== 'awaiting_funding') throw createError({ statusCode: 409, statusMessage: 'Campaign is already funded' })

  let txHash: Hex | null = null
  if (campaign.fundingMode === 'onchain') {
    const body = await readBody<{ txHash?: unknown }>(event)
    if (typeof body?.txHash !== 'string' || !isHex(body.txHash) || body.txHash.length !== 66) {
      throw createError({ statusCode: 400, statusMessage: 'A funding transaction hash is required' })
    }
    txHash = body.txHash
    const escrow = useRuntimeConfig(event).public.escrowAddress as Address
    const receipt = await usePublicClient().getTransactionReceipt({ hash: txHash })
    const funded =
      receipt.status === 'success' &&
      receipt.logs.some((log) => {
        if (!isAddressEqual(log.address, escrow)) return false
        try {
          const decoded = decodeEventLog({ abi: campaignEscrowAbi, data: log.data, topics: log.topics as [Hex, ...Hex[]] })
          return (
            decoded.eventName === 'Funded' &&
            decoded.args.orderId === campaign.orderId &&
            isAddressEqual(decoded.args.payer, campaign.payer) &&
            decoded.args.amount === BigInt(campaign.budgetUnits)
          )
        } catch {
          return false
        }
      })
    if (!funded) throw createError({ statusCode: 422, statusMessage: 'That transaction did not fund this order' })
  }

  campaign.status = 'in_review'
  campaign.txHash = txHash
  campaign.fundedAt = new Date().toISOString()
  await saveCampaign(campaign)
  return buildReport(campaign)
})
