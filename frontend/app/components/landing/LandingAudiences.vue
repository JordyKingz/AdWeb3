<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
const audiences = [
  {
    id: 'brands',
    eyebrow: 'For brands',
    title: 'Buy attention you can audit.',
    bullets: [
      'Pay per 1,000 validated, viewable impressions',
      'Budget held in escrow, unused funds refunded',
      'Delivery reports with invalid-traffic adjustments',
      'Contextual placements. No personal tracking',
    ],
    cta: { to: '/brands/campaigns/new', label: 'Start a campaign', primary: true },
  },
  {
    id: 'players',
    eyebrow: 'For players',
    title: 'Play free. Share the revenue.',
    bullets: [
      'Free games that keep working offline',
      'Earn points while you play online',
      'Points settle into USD₮ and pay out automatically',
      'Connect with a signature. No keys, no approvals',
    ],
    cta: { to: '/connect', label: 'Connect wallet', primary: true },
    filled: true,
  },
  {
    id: 'apps',
    eyebrow: 'For apps',
    title: 'Keep 80%. Pay nothing to start.',
    bullets: [
      'Free Android SDK, written in Kotlin',
      '80% of earned net revenue is yours',
      'No fees when an ad request goes unfilled',
      'Rewards optional. Users never need a wallet',
    ],
    cta: { to: '#apps', label: 'Join the publisher beta', primary: false },
  },
]
</script>

<template>
  <section class="container audiences">
    <SectionHeading :split="false" eyebrow="Who it's for" title="One network. Three ways in." />
    <div class="audiences__grid">
      <article
        v-for="a in audiences"
        :id="a.id"
        :key="a.id"
        class="audience"
        :class="{ 'audience--filled': a.filled }"
      >
        <p class="eyebrow">{{ a.eyebrow }}</p>
        <h3 class="audience__title">{{ a.title }}</h3>
        <ul class="audience__list">
          <li v-for="b in a.bullets" :key="b"><CheckDot />{{ b }}</li>
        </ul>
        <div class="audience__cta">
          <NuxtLink :to="landingOnly ? '#pilot' : a.cta.to" class="btn" :class="a.cta.primary ? 'btn--primary' : 'btn--secondary'">
            {{ landingOnly ? 'Pilot status' : a.cta.label }}
            <AppIcon v-if="a.cta.primary" name="arrow-right" :size="15" :stroke-width="2" />
          </NuxtLink>
        </div>
      </article>
    </div>
  </section>
</template>

<style scoped>
.audiences {
  padding-block: var(--section-y);
}

.audiences__grid {
  margin-top: 56px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.audience {
  flex: 1 1 320px;
  display: flex;
  flex-direction: column;
  padding: 32px;
  border: 1px solid var(--line);
  border-radius: 24px;
  scroll-margin-top: 80px;
}

.audience--filled {
  background: var(--surface);
  border-color: var(--surface);
}

.audience__title {
  margin-top: 18px;
  font-size: 30px;
  line-height: 1.08;
  letter-spacing: -0.03em;
  font-weight: 600;
}

.audience__list {
  list-style: none;
  margin: 28px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.audience__list li {
  display: flex;
  gap: 12px;
  font-size: 15px;
  line-height: 1.45;
  color: var(--ink-2);
}

.audience__cta {
  margin-top: auto;
  padding-top: 36px;
}
</style>
