<script setup lang="ts">
import type { WalletOption } from '~/composables/useWallet'

defineProps<{ busy?: string | null }>()
const emit = defineEmits<{ select: [option: WalletOption] }>()

const { options, discover } = useWallet()
const ready = ref(false)

onMounted(async () => {
  await discover()
  ready.value = true
})

function tagFor(option: WalletOption) {
  return option.rdns.includes('core') ? 'Avalanche' : ''
}
</script>

<template>
  <div class="picker">
    <template v-if="options.length">
      <button
        v-for="option in options"
        :key="option.id"
        type="button"
        class="picker__item"
        :disabled="!!busy"
        @click="emit('select', option as WalletOption)"
      >
        <img v-if="option.icon" :src="option.icon" alt="" class="picker__icon" width="36" height="36" />
        <span v-else class="picker__icon picker__icon--mono" aria-hidden="true">{{ option.name.charAt(0) }}</span>
        <span class="picker__name">{{ option.name }}</span>
        <span v-if="busy === option.id" class="spinner" aria-label="Connecting" />
        <span v-else-if="tagFor(option as WalletOption)" class="chip chip--accent">{{ tagFor(option as WalletOption) }}</span>
        <AppIcon name="chevron-right" :size="16" :stroke-width="2" class="picker__chevron" />
      </button>
    </template>

    <div v-else-if="ready" class="picker__empty">
      <p class="picker__empty-title">No browser wallet found</p>
      <p class="picker__empty-text">Install a wallet that supports Avalanche, then reload this page.</p>
      <div class="picker__links">
        <a href="https://core.app" target="_blank" rel="noopener" class="btn btn--secondary btn--sm">
          Get Core <AppIcon name="external" :size="14" />
        </a>
        <a href="https://metamask.io/download/" target="_blank" rel="noopener" class="btn btn--secondary btn--sm">
          Get MetaMask <AppIcon name="external" :size="14" />
        </a>
      </div>
    </div>

    <div v-else class="picker__loading" role="status"><span class="spinner" /> Looking for wallets…</div>
  </div>
</template>

<style scoped>
.picker {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.picker__item {
  display: flex;
  align-items: center;
  gap: 14px;
  min-height: 64px;
  padding: 0 16px;
  border-radius: 16px;
  border: 1px solid var(--line);
  background: #fff;
  text-align: left;
  transition: border-color 0.15s ease;
}

.picker__item:hover:not(:disabled) {
  border-color: var(--ink);
}

.picker__item:disabled {
  cursor: progress;
}

.picker__icon {
  flex: none;
  width: 36px;
  height: 36px;
  border-radius: 10px;
}

.picker__icon--mono {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--surface-2);
  font-weight: 700;
  font-size: 14px;
}

.picker__name {
  flex: 1;
  min-width: 0;
  font-size: 16px;
  font-weight: 500;
}

.picker__chevron {
  color: var(--muted);
}

.picker__empty {
  padding: 20px;
  border-radius: 16px;
  border: 1.5px dashed var(--line-strong);
}

.picker__empty-title {
  font-size: 16px;
  font-weight: 600;
}

.picker__empty-text {
  margin-top: 6px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--muted);
}

.picker__links {
  margin-top: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.picker__loading {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 64px;
  font-size: 14px;
  color: var(--muted);
}
</style>
