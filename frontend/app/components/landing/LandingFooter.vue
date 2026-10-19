<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
const columns = [
  {
    title: 'Product',
    links: [
      { label: 'For brands', to: '#brands' },
      { label: 'For players', to: '#players' },
      { label: 'For apps', to: '#apps' },
    ],
  },
  {
    title: 'App',
    links: [
      { label: 'Start a campaign', to: '/brands/campaigns/new' },
      { label: 'Connect wallet', to: '/connect' },
      { label: 'Rewards dashboard', to: '/dashboard' },
    ],
  },
  {
    title: 'Legal',
    links: [
      { label: 'Terms', to: '#top' },
      { label: 'Privacy', to: '#top' },
      { label: 'Reward policy', to: '#top' },
    ],
  },
]
</script>

<template>
  <footer class="footer">
    <div class="container footer__inner">
      <div class="footer__top">
        <div class="footer__brand">
          <AppLogo />
          <p class="footer__tagline">Advertising that pays the players, settled in USD₮.</p>
        </div>
        <div class="footer__cols">
          <div v-for="col in columns" :key="col.title">
            <h3 class="footer__heading">{{ col.title }}</h3>
            <ul class="footer__links">
              <li v-for="link in col.links" :key="link.label">
                <span v-if="landingOnly && col.title === 'Legal'">{{ link.label }} · In preparation</span>
                <NuxtLink v-else :to="landingOnly && link.to.startsWith('/') ? '#pilot' : link.to">{{ link.label }}</NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </div>
      <div class="footer__bottom">
        <span>© {{ new Date().getFullYear() }} AdWeb3</span>
        <span class="footer__disclaimer">
          Rewards depend on campaign funding, validated activity and availability in your country. Incentivized audiences
          are disclosed to advertisers.
        </span>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  border-top: 1px solid var(--line);
}

.footer__inner {
  padding-block: 56px 40px;
}

.footer__top {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 40px;
}

.footer__brand {
  flex: 1 1 280px;
}

.footer__tagline {
  margin-top: 12px;
  max-width: 300px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}

.footer__cols {
  display: flex;
  flex-wrap: wrap;
  gap: 56px;
}

.footer__heading {
  font-size: 13px;
  font-weight: 500;
  color: var(--muted);
}

.footer__links {
  list-style: none;
  margin: 16px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  font-size: 15px;
}

.footer__links a {
  text-decoration: none;
}

.footer__links a:hover {
  text-decoration: underline;
}

.footer__bottom {
  margin-top: 56px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 16px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--muted);
}

.footer__disclaimer {
  max-width: 640px;
}
</style>
