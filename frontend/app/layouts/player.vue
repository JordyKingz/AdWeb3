<script setup lang="ts">
import type { SessionInfo } from '#shared/types/api'

const { data: session } = await useFetch<SessionInfo>('/api/auth/session', { key: 'player-session' })
const online = useOnline()

const tabs = [
  { label: 'Dashboard', href: '#overview' },
  { label: 'Activity', href: '#activity' },
  { label: 'Payouts', href: '#payouts' },
  { label: 'Wallet', href: '#wallet' },
]
</script>

<template>
  <div class="player">
    <header class="player__header">
      <div class="container player__header-inner">
        <div class="player__left">
          <AppLogo />
          <nav v-if="session" class="player__tabs" aria-label="Player">
            <a v-for="(tab, i) in tabs" :key="tab.href" :href="tab.href" :aria-current="i === 0 ? 'page' : undefined">{{ tab.label }}</a>
          </nav>
        </div>
        <div class="player__right">
          <span v-if="session" class="player__status" role="status">
            <span class="player__dot" :class="{ 'player__dot--off': !online }" />
            {{ online ? 'Earning · online' : 'Earning paused · offline' }}
          </span>
          <NuxtLink v-if="session" to="/connect?change=1" class="player__wallet mono" :aria-label="`Payout wallet ${session.address}`">
            <AddressAvatar />{{ shortAddress(session.address) }}
          </NuxtLink>
          <NuxtLink v-else to="/connect" class="btn btn--primary btn--sm">Connect wallet</NuxtLink>
        </div>
      </div>
    </header>
    <slot />
  </div>
</template>

<style scoped>
.player {
  min-height: 100vh;
  background: var(--bg-app);
}

.player__header {
  border-bottom: 1px solid var(--line);
  background: #fff;
}

.player__header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  padding-block: 10px;
}

.player__left {
  display: flex;
  align-items: center;
  gap: 28px;
  min-width: 0;
}

.player__tabs {
  display: flex;
  gap: 4px;
}

.player__tabs a {
  display: inline-flex;
  align-items: center;
  min-height: 40px;
  padding: 0 14px;
  border-radius: 999px;
  font-size: 14px;
  color: var(--muted);
  text-decoration: none;
}

.player__tabs a:hover {
  color: var(--ink);
}

.player__tabs a[aria-current='page'] {
  background: var(--surface-2);
  color: var(--ink);
  font-weight: 500;
}

.player__right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.player__status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 40px;
  padding: 0 14px;
  border-radius: 999px;
  border: 1px solid var(--line);
  font-size: 13px;
  color: var(--ink-2);
  white-space: nowrap;
}

.player__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--online);
}

.player__dot--off {
  background: var(--line-dashed);
}

.player__wallet {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 0 14px 0 8px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
  font-size: 13px;
  text-decoration: none;
}

@media (max-width: 860px) {
  .player__tabs,
  .player__status {
    display: none;
  }
}
</style>
