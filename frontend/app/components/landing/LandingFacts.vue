<script setup lang="ts">
const facts = [
  { value: '30%', label: 'of earned ad revenue funds the player reward pool' },
  { value: '80%', label: 'of earned revenue stays with partner apps' },
  { value: '5 USD₮', label: 'per 1,000 validated impressions on test campaigns' },
  { value: '0', label: 'fees on unfilled ad requests. The SDK is free' },
]
</script>

<template>
  <section class="facts" aria-label="Key facts">
    <div class="container facts__inner">
      <div v-for="fact in facts" :key="fact.value" class="facts__item">
        <div class="facts__value">{{ fact.value }}</div>
        <p class="facts__label">{{ fact.label }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.facts {
  border-block: 1px solid var(--line);
}

.facts__inner {
  display: flex;
  flex-wrap: wrap;
  column-gap: 32px;
}

.facts__item {
  flex: 1 1 220px;
  padding-block: 36px;
}

.facts__value {
  font-size: 46px;
  font-weight: 600;
  letter-spacing: -0.045em;
  line-height: 1;
}

.facts__label {
  margin-top: 10px;
  max-width: 240px;
  font-size: 15px;
  line-height: 1.45;
  color: var(--muted);
}
</style>
