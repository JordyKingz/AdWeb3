<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
const links = [
  { href: '#how', label: 'How it works' },
  { href: '#brands', label: 'Brands' },
  { href: '#players', label: 'Players' },
  { href: '#apps', label: 'Apps' },
  { href: '#faq', label: 'FAQ' },
]
</script>

<template>
  <header class="nav">
    <nav class="container nav__inner" aria-label="Primary">
      <AppLogo />
      <div class="nav__links">
        <a v-for="link in links" :key="link.href" :href="link.href">{{ link.label }}</a>
      </div>
      <div class="nav__actions">
        <NuxtLink :to="landingOnly ? '#pilot' : '/brands/campaigns/new'" class="nav__text-link">{{ landingOnly ? 'Pilot status' : 'Start a campaign' }}</NuxtLink>
        <NuxtLink :to="landingOnly ? '#how' : '/connect'" class="btn btn--primary btn--sm">{{ landingOnly ? 'Explore the pilot' : 'Launch app' }}</NuxtLink>
      </div>
    </nav>
  </header>
</template>

<style scoped>
.nav {
  position: sticky;
  top: 0;
  z-index: 30;
  border-bottom: 1px solid var(--line);
  background: rgb(255 255 255 / 0.86);
  backdrop-filter: saturate(1.4) blur(14px);
  -webkit-backdrop-filter: saturate(1.4) blur(14px);
}

.nav__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding-block: 12px;
}

.nav__links {
  display: flex;
  align-items: center;
  gap: 32px;
  font-size: 15px;
}

.nav__links a {
  color: var(--muted);
  text-decoration: none;
  transition: color 0.15s ease;
}

.nav__links a:hover {
  color: var(--ink);
}

.nav__actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav__text-link {
  display: inline-flex;
  align-items: center;
  min-height: 44px;
  padding: 0 14px;
  font-size: 15px;
  font-weight: 500;
  text-decoration: none;
}

@media (max-width: 860px) {
  .nav__links,
  .nav__text-link {
    display: none;
  }
}
</style>
