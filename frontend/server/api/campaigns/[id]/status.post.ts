export default defineEventHandler(async (event) => {
  const campaign = await getCampaign(getRouterParam(event, 'id') ?? '')
  if (!campaign) throw createError({ statusCode: 404, statusMessage: 'Campaign not found' })
  const body = await readBody<{ paused?: unknown }>(event)
  if (typeof body?.paused !== 'boolean') throw createError({ statusCode: 400, statusMessage: 'paused must be a boolean' })
  if (campaign.status !== 'live' && campaign.status !== 'paused') {
    throw createError({ statusCode: 409, statusMessage: 'Only live campaigns can be paused or resumed' })
  }
  campaign.status = body.paused ? 'paused' : 'live'
  await saveCampaign(campaign)
  return buildReport(campaign)
})
