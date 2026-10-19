<script setup lang="ts">
const props = defineProps<{ brand: string; headline: string; cta: string }>()

const initials = computed(
  () =>
    props.brand
      .trim()
      .split(/\s+/)
      .map((w) => w.charAt(0))
      .join('')
      .slice(0, 2)
      .toUpperCase() || 'YB',
)
</script>

<template>
  <div class="preview">
    <div class="card-mock">
      <div class="card-mock__image mono">[YOUR CARD IMAGE]</div>
      <div class="card-mock__row">
        <div class="card-mock__logo">{{ initials }}</div>
        <div class="card-mock__copy">
          <div class="card-mock__label">Sponsored · {{ brand || 'Your brand' }}</div>
          <div class="card-mock__headline">{{ headline || 'Your headline' }}</div>
        </div>
      </div>
      <div class="card-mock__cta">{{ cta }}</div>
    </div>
  </div>
</template>

<style scoped>
.preview {
  display: flex;
  justify-content: center;
  padding: 28px 20px;
  border-radius: 22px;
  background: var(--surface);
}

.card-mock {
  width: 100%;
  max-width: 300px;
  padding: 14px;
  border-radius: 22px;
  background: #fff;
  box-shadow: 0 20px 40px -24px rgb(11 11 12 / 0.35);
}

.card-mock__image {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 132px;
  border-radius: 14px;
  background: #ededE9;
  font-size: 12px;
  color: var(--muted);
}

.card-mock__row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 14px;
}

.card-mock__logo {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--ink);
  color: var(--accent);
  font-size: 13px;
  font-weight: 700;
}

.card-mock__copy {
  flex: 1;
  min-width: 0;
}

.card-mock__label {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--muted);
}

.card-mock__headline {
  margin-top: 2px;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.3;
  overflow-wrap: anywhere;
}

.card-mock__cta {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 42px;
  margin-top: 14px;
  border-radius: 12px;
  background: var(--ink);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
}
</style>
