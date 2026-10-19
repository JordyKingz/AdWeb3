import type { CreatedOrder } from '#shared/types/api'

/**
 * Creates a campaign order. In production the Go API registers the order on-chain
 * (CampaignEscrow.registerOrder, operator-only) before the advertiser can fund it.
 */
export default defineEventHandler(async (event): Promise<CreatedOrder> => {
  const fundingMode = currentFundingMode()
  const draft = validateDraft(await readBody(event), fundingMode)
  const campaign = newCampaign(draft, fundingMode)
  await saveCampaign(campaign)
  setResponseStatus(event, 201)
  return {
    id: campaign.id,
    orderId: campaign.orderId,
    capUnits: campaign.budgetUnits,
    terms: campaign.terms,
    fundingMode: campaign.fundingMode,
  }
})
