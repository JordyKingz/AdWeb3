<script setup lang="ts">
withDefaults(defineProps<{ to?: string; inverted?: boolean; label?: boolean }>(), {
  to: '/',
  inverted: false,
  label: true,
})
</script>

<template>
  <NuxtLink :to="to" class="logo" :class="{ 'logo--inverted': inverted }" aria-label="AdWeb3 home">
    <svg width="26" height="26" viewBox="0 0 26 26" aria-hidden="true">
      <rect width="26" height="26" rx="7" :fill="inverted ? '#FFFFFF' : '#0B0B0C'" />
      <path
        d="M7 18.5 13 7l6 11.5"
        :stroke="inverted ? '#0B0B0C' : '#FFFFFF'"
        stroke-width="2.2"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <circle cx="19.2" cy="18.6" r="2.7" fill="#D4FF3A" />
    </svg>
    <span v-if="label" class="logo__word">AdWeb3</span>
  </NuxtLink>
</template>

<style scoped>
.logo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  color: var(--ink);
  text-decoration: none;
}

.logo--inverted {
  color: #fff;
}

.logo__word {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.02em;
}
</style>
