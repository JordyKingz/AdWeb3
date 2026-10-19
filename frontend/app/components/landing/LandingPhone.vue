<script setup lang="ts">
// A Candy Blocks board mid-game; e = empty cell.
const CANDY = { e: '#ededE9', b: '#7aa7ff', o: '#ffb547', p: '#ff8fb1', l: '#b9f23a', v: '#b79cff' } as const
const board: (keyof typeof CANDY)[] = [
  'e', 'e', 'e', 'b', 'e', 'e',
  'e', 'e', 'b', 'b', 'e', 'e',
  'e', 'e', 'e', 'e', 'e', 'o',
  'p', 'e', 'e', 'l', 'o', 'o',
  'p', 'p', 'v', 'l', 'l', 'b',
  'v', 'p', 'v', 'v', 'l', 'b',
]
</script>

<template>
  <div class="phone">
    <div class="phone__screen">
      <div class="phone__top">
        <span class="mono phone__level">LEVEL 12</span>
        <span class="phone__earning"><span class="phone__dot" />Earning</span>
      </div>
      <div class="phone__title">Candy Blocks</div>
      <div class="phone__board">
        <span v-for="(cell, i) in board" :key="i" class="phone__cell" :style="{ background: CANDY[cell] }" />
      </div>
      <div class="phone__play">Play</div>
      <div class="phone__sponsor">
        <div class="phone__sponsor-logo">YB</div>
        <div class="phone__sponsor-copy">
          <div class="phone__sponsor-label">Sponsored</div>
          <div class="phone__sponsor-text">Your brand, seen by real players</div>
        </div>
        <span class="phone__sponsor-cta">Open</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.phone {
  width: 276px;
  height: 600px;
  padding: 10px;
  border-radius: 46px;
  background: var(--ink);
  box-shadow: 0 40px 80px -30px rgb(11 11 12 / 0.45);
}

.phone__screen {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 24px 18px 18px;
  border-radius: 37px;
  background: #fff;
  overflow: hidden;
}

.phone__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.phone__level {
  font-size: 11px;
  color: var(--muted);
}

.phone__earning {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 9px;
  border-radius: 999px;
  background: var(--surface);
  font-size: 11px;
  font-weight: 500;
}

.phone__dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--online);
}

.phone__title {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1;
}

.phone__board {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 5px;
  padding: 10px;
  border-radius: 16px;
  background: #f7f7f4;
}

.phone__cell {
  aspect-ratio: 1;
  border-radius: 7px;
}

.phone__play {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 46px;
  border-radius: 14px;
  background: var(--ink);
  color: #fff;
  font-size: 15px;
  font-weight: 600;
}

.phone__sponsor {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: 16px;
}

.phone__sponsor-logo {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 11px;
  background: var(--ink);
  color: var(--accent);
  font-size: 13px;
  font-weight: 700;
}

.phone__sponsor-copy {
  flex: 1;
  min-width: 0;
}

.phone__sponsor-label {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--muted);
}

.phone__sponsor-text {
  margin-top: 2px;
  font-size: 13px;
  font-weight: 600;
  line-height: 1.25;
}

.phone__sponsor-cta {
  padding: 6px 10px;
  border-radius: 999px;
  background: var(--surface);
  font-size: 12px;
  font-weight: 500;
}
</style>
