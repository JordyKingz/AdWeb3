// https://nuxt.com/docs/api/configuration/nuxt-config
const landingOnly = process.env.NUXT_PUBLIC_LANDING_ONLY === 'true'

export default defineNuxtConfig({
  hooks: {
    'pages:extend'(pages) {
      if (landingOnly) pages.splice(0, pages.length, ...pages.filter(page => page.path === '/'))
    },
    'prerender:routes'({ routes }) {
      if (landingOnly) for (const route of [...routes]) if (route !== '/') routes.delete(route)
    },
  },
  compatibilityDate: '2026-10-01',
  devtools: { enabled: false },
  css: ['~/assets/css/main.css'],

  app: {
    head: {
      htmlAttrs: { lang: 'en' },
      title: 'AdWeb3 — Ads that pay the people who see them',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            'Brands fund campaigns in USD₮. Players earn 30% of the revenue from every validated impression, settled on-chain.',
        },
        { name: 'theme-color', content: '#0B0B0C' },
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap',
        },
      ],
    },
  },

  // Token amounts are BigInt integer units, so the server bundle must target ES2020+.
  nitro: {
    esbuild: { options: { target: 'es2022' } },
    ...(landingOnly ? { prerender: { crawlLinks: false, routes: ['/'] } } : {}),
  },

  routeRules: {
    '/brands': { redirect: '/brands/campaigns/new' },
    '/brands/campaigns': { redirect: '/brands/campaigns/new' },
  },

  // Every value can be overridden with NUXT_PUBLIC_* env vars (see .env.example).
  runtimeConfig: {
    public: {
      landingOnly,
      chainId: 43113,
      chainName: 'Avalanche Fuji',
      rpcUrl: 'https://api.avax-test.network/ext/bc/C/rpc',
      explorerUrl: 'https://testnet.snowtrace.io',
      tokenAddress: '',
      tokenSymbol: 'USD₮',
      tokenDecimals: 6,
      escrowAddress: '',
    },
  },
})
