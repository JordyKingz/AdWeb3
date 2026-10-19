<script setup lang="ts">
const props = defineProps<{ steps: string[]; current: number; done: (step: number) => boolean }>()
const emit = defineEmits<{ go: [step: number] }>()
</script>

<template>
  <ol class="stepper">
    <li v-for="(label, i) in props.steps" :key="label" class="stepper__item">
      <button
        type="button"
        class="stepper__btn"
        :class="{ 'stepper__btn--active': current === i + 1, 'stepper__btn--done': done(i + 1) }"
        :aria-current="current === i + 1 ? 'step' : undefined"
        @click="emit('go', i + 1)"
      >
        <span class="stepper__num mono">
          <AppIcon v-if="done(i + 1)" name="check" :size="12" :stroke-width="3.5" />
          <template v-else>{{ i + 1 }}</template>
        </span>
        {{ label }}
      </button>
      <span
        v-if="i < props.steps.length - 1"
        class="stepper__line"
        :class="{ 'stepper__line--done': current > i + 1 }"
        aria-hidden="true"
      />
    </li>
  </ol>
</template>

<style scoped>
.stepper {
  list-style: none;
  margin: 28px 0 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px 12px;
}

.stepper__item {
  flex: 1 1 150px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.stepper__btn {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  padding: 0 2px;
  border: 0;
  background: none;
  color: var(--muted);
  font-size: 15px;
  font-weight: 500;
  white-space: nowrap;
}

.stepper__btn--active,
.stepper__btn--done {
  color: var(--ink);
}

.stepper__btn--active {
  font-weight: 600;
}

.stepper__num {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1.5px solid #cfcfca;
  background: #fff;
  font-size: 13px;
}

.stepper__btn--active .stepper__num {
  border-color: var(--ink);
  background: var(--ink);
  color: #fff;
}

.stepper__btn--done .stepper__num {
  border-color: var(--accent);
  background: var(--accent);
  color: var(--ink);
}

.stepper__line {
  flex: 1;
  min-width: 16px;
  height: 1.5px;
  background: #e2e2de;
}

.stepper__line--done {
  background: var(--ink);
}
</style>
