<script setup lang="ts">
const { chainName } = useChain()
</script>

<template>
  <span class="network"><span class="network__dot" />{{ chainName }}</span>
</template>

<style scoped>
.network {
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

.network__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--avax);
}
</style>
