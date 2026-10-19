<script setup lang="ts">
const parts = [
  { key: 'players', pct: 30, label: 'Players', text: 'The reward pool. Settled in USD₮ and paid to verified wallets.' },
  { key: 'ops', pct: 60, label: 'Operations', text: 'Servers, payout gas, fraud review and the team that runs it all.' },
  { key: 'reserve', pct: 10, label: 'Contingency', text: 'A reserve for disputes and delivery adjustments.' },
]
</script>

<template>
  <section class="split">
    <div class="container split__inner">
      <SectionHeading
        dark
        :split="false"
        eyebrow="The split"
        title="A fixed share, written into every campaign."
        lede="The split is locked in the moment a campaign is bought, so it can't change underneath anyone. Refundable deposits never fund payouts. Only revenue that's actually been earned does."
      />
      <div class="split__bar" role="img" aria-label="Revenue split: 30 percent players, 60 percent operations, 10 percent contingency">
        <div v-for="p in parts" :key="p.key" class="split__seg mono" :class="`split__seg--${p.key}`" :style="{ flex: p.pct }">
          {{ p.pct }}%
        </div>
      </div>
      <div class="split__legend">
        <div v-for="p in parts" :key="p.key" class="split__legend-item" :style="{ flex: `${p.pct} 1 220px` }">
          <div class="split__legend-title"><span class="split__key" :class="`split__seg--${p.key}`" />{{ p.label }}</div>
          <p class="split__legend-text">{{ p.text }}</p>
        </div>
      </div>
      <div class="split__example mono">
        <span>1,000 USD₮ earned</span>
        <span class="split__example-strong">→ 300 to players</span>
        <span>→ 600 to operations</span>
        <span>→ 100 to reserve</span>
      </div>
    </div>
  </section>
</template>

<style scoped>
.split {
  background: var(--ink);
  color: #fff;
}

.split__inner {
  padding-block: var(--section-y);
}

.split__bar {
  display: flex;
  gap: 6px;
  height: 96px;
  margin-top: 64px;
}

.split__seg {
  min-width: 72px;
  display: flex;
  align-items: flex-end;
  padding: 14px 16px;
  border-radius: 18px;
  font-size: clamp(18px, 2.4vw, 30px);
  font-weight: 500;
}

.split__seg--players {
  background: var(--accent);
  color: var(--ink);
}

.split__seg--ops {
  background: #fff;
  color: var(--ink);
}

.split__seg--reserve {
  background: #3a3a38;
  color: #fff;
}

.split__legend {
  margin-top: 32px;
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
}

.split__legend-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 17px;
  font-weight: 600;
}

.split__key {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  padding: 0;
  min-width: 0;
}

.split__key.split__seg--reserve {
  background: #6b6b66;
}

.split__legend-text {
  margin-top: 8px;
  max-width: 320px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--on-ink-muted);
}

.split__example {
  margin-top: 56px;
  padding-top: 28px;
  border-top: 1px solid #2a2a2a;
  display: flex;
  flex-wrap: wrap;
  gap: 12px 40px;
  font-size: 14px;
  color: var(--on-ink-muted);
}

.split__example-strong {
  color: #fff;
}
</style>
