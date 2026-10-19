<script setup lang="ts">
import type { IconName } from '~/components/AppIcon.vue'

const route = useRoute()
const { data: overview } = await useFetch('/api/brand/overview', { key: 'brand-overview' })

const latest = computed(() => overview.value?.campaigns.find((c) => c.status !== 'awaiting_funding'))

const nav = computed<{ label: string; icon: IconName; to?: string; active?: boolean }[]>(() => [
  {
    label: 'Overview',
    icon: 'grid',
    to: latest.value ? `/brands/campaigns/${latest.value.id}` : undefined,
    active: route.path.startsWith('/brands/campaigns/') && route.path !== '/brands/campaigns/new',
  },
  { label: 'New campaign', icon: 'megaphone', to: '/brands/campaigns/new', active: route.path === '/brands/campaigns/new' },
  { label: 'Creatives', icon: 'image' },
  { label: 'Reports', icon: 'chart' },
  { label: 'Billing', icon: 'card' },
])
</script>

<template>
  <div class="app">
    <header class="app__header">
      <div class="app__brand">
        <AppLogo />
        <span class="chip chip--soft">Brands</span>
      </div>
      <div class="app__header-actions">
        <NetworkPill class="app__network" />
        <ClientOnly>
          <WalletChip show-balance />
          <template #fallback><span class="app__wallet-placeholder" /></template>
        </ClientOnly>
      </div>
    </header>

    <div class="app__body">
      <aside class="app__side">
        <nav class="app__nav" aria-label="Advertiser">
          <template v-for="item in nav" :key="item.label">
            <NuxtLink
              v-if="item.to"
              :to="item.to"
              class="app__nav-item"
              :class="{ 'app__nav-item--active': item.active }"
              :aria-current="item.active ? 'page' : undefined"
            >
              <AppIcon :name="item.icon" />{{ item.label }}
            </NuxtLink>
            <span v-else class="app__nav-item app__nav-item--soon" aria-disabled="true">
              <AppIcon :name="item.icon" />{{ item.label }}<span class="app__soon">Soon</span>
            </span>
          </template>
        </nav>
        <div class="app__escrow">
          <div class="app__escrow-label">In escrow</div>
          <div class="app__escrow-value mono">
            {{ overview ? formatToken(overview.inEscrowUnits) : '—' }} <span>USD₮</span>
          </div>
          <p class="app__escrow-note">Refundable until delivery earns it.</p>
        </div>
      </aside>

      <main class="app__main">
        <div class="app__content">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
  background: var(--bg-app);
}

.app__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  padding: 10px clamp(16px, 3vw, 28px);
  border-bottom: 1px solid var(--line);
  background: #fff;
}

.app__brand,
.app__header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app__header-actions {
  gap: 8px;
}

.app__wallet-placeholder {
  display: inline-block;
  width: 150px;
  height: 44px;
}

.app__body {
  display: flex;
  flex-wrap: wrap;
  align-items: stretch;
}

.app__side {
  flex: 1 1 232px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 20px 14px;
  border-right: 1px solid var(--line);
  background: #fff;
}

.app__nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.app__nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 44px;
  padding: 0 12px;
  border-radius: 10px;
  color: var(--ink-2);
  font-size: 15px;
  text-decoration: none;
  white-space: nowrap;
}

a.app__nav-item:hover {
  background: var(--surface);
}

.app__nav-item--active {
  background: var(--surface-2);
  color: var(--ink);
  font-weight: 500;
}

.app__nav-item--soon {
  color: var(--muted);
  cursor: default;
}

.app__soon {
  margin-left: auto;
  font-size: 11px;
  color: var(--muted);
}

.app__escrow {
  margin-top: auto;
  padding: 16px;
  border-radius: 14px;
  background: var(--surface);
}

.app__escrow-label {
  font-size: 12px;
  color: var(--muted);
}

.app__escrow-value {
  margin-top: 4px;
  font-size: 20px;
  letter-spacing: -0.02em;
}

.app__escrow-value span {
  font-size: 12px;
  color: var(--muted);
}

.app__escrow-note {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.45;
  color: var(--muted);
}

.app__main {
  flex: 999 1 560px;
  min-width: 0;
  padding: clamp(20px, 3vw, 40px);
}

.app__content {
  max-width: 1120px;
  margin: 0 auto;
}

@media (max-width: 860px) {
  .app__network {
    display: none;
  }

  .app__side {
    border-right: 0;
    border-bottom: 1px solid var(--line);
    padding-block: 8px;
  }

  .app__nav {
    flex-direction: row;
    overflow-x: auto;
  }

  .app__nav-item--soon,
  .app__escrow {
    display: none;
  }
}
</style>
