export default defineEventHandler(async (event) => {
  const campaign = await getCampaign(getRouterParam(event, 'id') ?? '')
  if (!campaign) throw createError({ statusCode: 404, statusMessage: 'Campaign not found' })
  return buildReport(campaign)
})
