<script setup lang="ts">
import type { IconName } from '~/components/AppIcon.vue'
import { CPM_UNITS, TEST_BUDGET, maxImpressions } from '#shared/utils/campaign'

const landingOnly = useRuntimeConfig().public.landingOnly
const testBudgetUnits = BigInt(TEST_BUDGET) * 1_000_000n
const rows = [
  { label: 'Prepaid budget', value: `${TEST_BUDGET.toFixed(2)} USD₮` },
  { label: 'Max billed impressions', value: formatInt(maxImpressions(testBudgetUnits)) },
  { label: 'Serving window', value: '168 h from your UTC start' },
  { label: 'Paid in', value: 'USD₮ · Avalanche C-Chain' },
]
const notes: { icon: IconName; text: string }[] = [
  { icon: 'eye', text: 'Every creative is reviewed before it serves' },
  { icon: 'chart', text: 'Live delivery and invalid-traffic report' },
  { icon: 'info', text: 'Incentivized audience disclosed up front' },
]
const cpm = formatToken(CPM_UNITS, 0)
</script>

<template>
  <section id="pricing" class="pricing">
    <div class="container pricing__inner">
      <div class="pricing__copy">
        <p class="eyebrow">Test campaign</p>
        <h2 class="display-2 pricing__title">Start with a 25 USD₮ test.</h2>
        <p class="lede pricing__lede">
          Run a seven-day campaign in Candy Blocks, our first Android game. You're billed only for validated, viewable
          impressions. Anything not delivered comes back to you after the adjustment window.
        </p>
        <ul class="pricing__notes">
          <li v-for="n in notes" :key="n.text"><AppIcon :name="n.icon" :size="18" />{{ n.text }}</li>
        </ul>
      </div>

      <div class="pricing__card">
        <div class="pricing__head">
          <span class="pricing__name">Test campaign</span>
          <span class="chip chip--soft">7 days</span>
        </div>
        <div class="pricing__price">
          <span class="pricing__amount">{{ cpm }} USD₮</span>
          <span class="pricing__unit">per 1,000 validated impressions</span>
        </div>
        <dl class="pricing__rows">
          <div v-for="r in rows" :key="r.label">
            <dt>{{ r.label }}</dt>
            <dd class="mono">{{ r.value }}</dd>
          </div>
        </dl>
        <p class="pricing__example">
          <strong>Example:</strong> 2,000 valid impressions → 10.00 USD₮ billed, 15.00 USD₮ refunded.
        </p>
        <NuxtLink :to="landingOnly ? '#pilot' : '/brands/campaigns/new'" class="btn btn--primary btn--lg btn--block pricing__cta">
          {{ landingOnly ? 'Pilot status' : 'Launch a test campaign' }} <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.pricing {
  background: var(--surface);
}

.pricing__inner {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: clamp(40px, 6vw, 96px);
  padding-block: var(--section-y);
}

.pricing__copy {
  flex: 1 1 440px;
  min-width: 0;
}

.pricing__title {
  margin-top: 16px;
}

.pricing__lede {
  margin-top: 24px;
  max-width: 520px;
}

.pricing__notes {
  list-style: none;
  margin: 32px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
  font-size: 15px;
  color: var(--ink-2);
}

.pricing__notes li {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pricing__card {
  flex: 1 1 420px;
  min-width: 0;
  padding: clamp(24px, 3vw, 36px);
  border: 1px solid var(--line);
  border-radius: 28px;
  background: #fff;
  box-shadow: 0 30px 60px -40px rgb(11 11 12 / 0.3);
}

.pricing__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.pricing__name {
  font-size: 15px;
  font-weight: 600;
}

.pricing__price {
  margin-top: 24px;
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 4px 10px;
}

.pricing__amount {
  font-size: 64px;
  font-weight: 600;
  letter-spacing: -0.05em;
  line-height: 1;
}

.pricing__unit {
  font-size: 15px;
  color: var(--muted);
}

.pricing__rows {
  margin: 28px 0 0;
}

.pricing__rows div {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 0;
  border-top: 1px dashed var(--line-strong);
  font-size: 15px;
}

.pricing__rows dt {
  color: var(--muted);
}

.pricing__rows dd {
  margin: 0;
  text-align: right;
}

.pricing__example {
  margin-top: 8px;
  padding: 16px;
  border-radius: 16px;
  background: var(--surface);
  font-size: 14px;
  line-height: 1.5;
  color: var(--ink-2);
}

.pricing__cta {
  margin-top: 20px;
}
</style>
