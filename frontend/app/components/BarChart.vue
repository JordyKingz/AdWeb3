<script setup lang="ts">
export interface Bar {
  label: string
  /** null draws a dashed placeholder (not served yet, or offline). */
  value: number | null
  highlight?: boolean
}

const props = withDefaults(
  defineProps<{
    bars: Bar[]
    max: number
    ticks: number[]
    /** Text alternative describing the data. */
    summary: string
    height?: number
    valueLabels?: boolean
    /** Draws the top tick as a dashed limit line with this label. */
    capLabel?: string
  }>(),
  { height: 210, valueLabels: false, capLabel: undefined },
)

const sortedTicks = computed(() => [...props.ticks].sort((a, b) => b - a))

function barHeight(bar: Bar) {
  if (bar.value === null) return '14%'
  return `${Math.max((bar.value / props.max) * 100, bar.value > 0 ? 2 : 0)}%`
}
</script>

<template>
  <figure class="chart" role="img" :aria-label="summary">
    <div class="chart__plot" :style="{ height: `${height}px` }">
      <div class="chart__grid" aria-hidden="true">
        <div v-for="tick in sortedTicks" :key="tick" class="chart__tick" :style="{ bottom: `${(tick / max) * 100}%` }">
          <span class="chart__tick-label mono">{{ tick }}</span>
          <span
            class="chart__tick-line"
            :class="{ 'chart__tick-line--base': tick === 0, 'chart__tick-line--cap': capLabel && tick === max }"
          />
          <span v-if="capLabel && tick === max" class="chart__cap">{{ capLabel }}</span>
        </div>
      </div>
      <div class="chart__bars" :class="{ 'chart__bars--capped': capLabel }" aria-hidden="true">
        <div v-for="(bar, i) in bars" :key="i" class="chart__col">
          <span v-if="valueLabels && bar.value !== null" class="chart__value mono">{{ formatInt(bar.value) }}</span>
          <div
            class="chart__bar"
            :class="{ 'chart__bar--empty': bar.value === null, 'chart__bar--highlight': bar.highlight }"
            :style="{ height: barHeight(bar) }"
          />
        </div>
      </div>
    </div>
    <div class="chart__labels mono" :class="{ 'chart__labels--capped': capLabel }" aria-hidden="true">
      <span v-for="(bar, i) in bars" :key="i" :class="{ 'chart__label--strong': bar.highlight }">{{ bar.label }}</span>
    </div>
  </figure>
</template>

<style scoped>
.chart {
  margin: 0;
}

.chart__plot {
  position: relative;
}

.chart__grid {
  position: absolute;
  inset: 0;
}

.chart__tick {
  position: absolute;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  transform: translateY(50%);
}

.chart__tick-label {
  width: 30px;
  text-align: right;
  font-size: 11px;
  color: var(--muted);
}

.chart__tick-line {
  flex: 1;
  height: 1px;
  background: #f0f0ec;
}

.chart__tick-line--base {
  background: #d6d6d1;
}

.chart__tick-line--cap {
  height: 0;
  background: none;
  border-top: 1.5px dashed var(--ink);
}

.chart__cap {
  width: 56px;
  font-size: 11px;
  color: var(--ink-2);
}

.chart__bars {
  position: absolute;
  inset: 0 0 0 40px;
  display: flex;
  align-items: flex-end;
  gap: clamp(4px, 1.4vw, 18px);
  padding: 0 4px;
}

.chart__bars--capped {
  right: 64px;
}

.chart__labels--capped {
  margin-right: 64px;
}

.chart__col {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
}

.chart__value {
  margin-bottom: 6px;
  font-size: 11px;
}

.chart__bar {
  width: 100%;
  max-width: 56px;
  border-radius: 7px 7px 2px 2px;
  background: var(--ink);
  transition: height 0.4s ease;
}

.chart__bar--empty {
  background: none;
  border: 1.5px dashed var(--line-dashed);
}

.chart__bar--highlight {
  background: var(--accent);
  box-shadow: inset 0 0 0 1.5px var(--ink);
}

.chart__labels {
  display: flex;
  gap: clamp(4px, 1.4vw, 18px);
  margin-top: 10px;
  padding: 0 4px 0 44px;
  font-size: 11px;
  color: var(--muted);
}

.chart__labels span {
  flex: 1;
  text-align: center;
}

.chart__label--strong {
  color: var(--ink);
  font-weight: 500;
}
</style>
