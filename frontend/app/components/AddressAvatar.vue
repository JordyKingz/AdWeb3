<script setup lang="ts">
withDefaults(defineProps<{ size?: number }>(), { size: 24 })
</script>

<template>
  <span
    class="avatar"
    :style="{ width: `${size}px`, height: `${size}px`, boxShadow: `inset 0 0 0 ${Math.round(size / 4)}px var(--accent)` }"
    aria-hidden="true"
  />
</template>

<style scoped>
.avatar {
  flex: none;
  display: inline-block;
  border-radius: 50%;
  background: var(--ink);
}
</style>
