export default defineEventHandler(async (event) => {
  await destroyPlayerSession(event)
  return { ok: true }
})
