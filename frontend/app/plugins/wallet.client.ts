export default defineNuxtPlugin(() => {
  if (useRuntimeConfig().public.landingOnly) return
  // Restore a previously connected wallet in the background; never blocks rendering.
  useWallet()
    .reconnect()
    .catch(() => {})
})
