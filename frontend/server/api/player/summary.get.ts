export default defineEventHandler(async (event) => {
  const session = await requirePlayerSession(event)
  return buildPlayerSummary(session)
})
