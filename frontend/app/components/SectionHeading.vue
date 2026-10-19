<script setup lang="ts">
withDefaults(defineProps<{ eyebrow: string; title: string; lede?: string; split?: boolean; dark?: boolean }>(), {
  lede: undefined,
  split: true,
  dark: false,
})
</script>

<template>
  <div class="heading" :class="{ 'heading--split': split && lede, 'heading--dark': dark }">
    <div class="heading__main">
      <p class="eyebrow">{{ eyebrow }}</p>
      <h2 class="display-2 heading__title">{{ title }}</h2>
    </div>
    <p v-if="lede" class="lede heading__lede">{{ lede }}</p>
  </div>
</template>

<style scoped>
.heading {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.heading--split {
  flex-flow: row wrap;
  justify-content: space-between;
  align-items: flex-end;
}

/* Basis widths only make sense in the row layout; in a column they would become heights. */
.heading--split .heading__main {
  flex: 1 1 520px;
}

.heading--split .heading__lede {
  flex: 0 1 420px;
}

.heading__title {
  margin-top: 16px;
  max-width: 860px;
}

.heading__lede {
  max-width: 620px;
}

.heading--dark .eyebrow,
.heading--dark .heading__lede {
  color: var(--on-ink-muted);
}
</style>
