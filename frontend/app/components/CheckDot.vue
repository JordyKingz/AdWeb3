<script setup lang="ts">
withDefaults(defineProps<{ size?: number }>(), { size: 22 })
</script>

<template>
  <span class="check-dot" :style="{ width: `${size}px`, height: `${size}px` }" aria-hidden="true">
    <AppIcon name="check" :size="Math.round(size / 2)" :stroke-width="3.5" />
  </span>
</template>

<style scoped>
.check-dot {
  flex: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--accent);
  color: var(--ink);
}
</style>
