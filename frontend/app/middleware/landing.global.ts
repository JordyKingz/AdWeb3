export default defineNuxtRouteMiddleware((to) => {
  if (useRuntimeConfig().public.landingOnly && to.path !== '/') {
    return navigateTo('/#pilot', { redirectCode: 302 })
  }
})
