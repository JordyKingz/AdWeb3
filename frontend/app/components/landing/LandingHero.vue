<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
</script>

<template>
  <section class="hero container">
    <div class="hero__copy">
      <p class="hero__badge"><span class="hero__badge-tag mono">PILOT</span>Android first · Settled in USD₮</p>
      <h1 class="hero__title">Ads that <mark>pay</mark> the people who see them.</h1>
      <p class="hero__lede">
        Brands fund campaigns in USD₮. Players earn 30% of the revenue from every validated impression, settled on-chain
        and paid straight to their wallet.
      </p>
      <div class="hero__ctas">
        <NuxtLink :to="landingOnly ? '#pilot' : '/brands/campaigns/new'" class="btn btn--primary btn--lg">
          {{ landingOnly ? 'Explore the pilot' : 'Start a campaign' }} <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
        </NuxtLink>
        <NuxtLink :to="landingOnly ? '#how' : '/connect'" class="btn btn--secondary btn--lg">{{ landingOnly ? 'How it works' : 'Play & earn' }}</NuxtLink>
      </div>
      <p class="hero__note">Test campaigns from 25 USD₮ · Pay only for valid views · Unused budget refunded</p>
    </div>

    <div class="hero__visual" aria-hidden="true">
      <div class="stage">
        <LandingPhone class="stage__phone" />

        <div class="float float--split">
          <div class="float__label">Per 1,000 USD₮ earned</div>
          <div class="float__bar">
            <span style="flex: 30" class="float__seg float__seg--accent" />
            <span style="flex: 60" class="float__seg float__seg--ink" />
            <span style="flex: 10" class="float__seg float__seg--grey" />
          </div>
          <dl class="float__rows">
            <div><dt><span class="float__key float__seg--accent" />Players</dt><dd class="mono">300</dd></div>
            <div><dt><span class="float__key float__seg--ink" />Operations</dt><dd class="mono">600</dd></div>
            <div><dt><span class="float__key float__seg--grey" />Reserve</dt><dd class="mono">100</dd></div>
          </dl>
        </div>

        <div class="float float--payout">
          <div class="float__head"><CheckDot :size="18" />Payout confirmed</div>
          <div class="float__amount mono">+0.5120 <span>USD₮</span></div>
          <div class="float__meta mono">to 0x71C2…9A3f · Avalanche</div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: clamp(40px, 5vw, 72px);
  padding-block: clamp(48px, 8vw, 104px) clamp(56px, 7vw, 104px);
}

.hero__copy {
  flex: 1 1 520px;
  min-width: 0;
}

.hero__copy > * {
  animation: rise 0.7s cubic-bezier(0.2, 0.7, 0.2, 1) both;
}

.hero__copy > :nth-child(2) {
  animation-delay: 0.06s;
}

.hero__copy > :nth-child(3) {
  animation-delay: 0.12s;
}

.hero__copy > :nth-child(n + 4) {
  animation-delay: 0.18s;
}

.hero__badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 5px 14px 5px 5px;
  border: 1px solid var(--line);
  border-radius: 999px;
  font-size: 13px;
  color: var(--ink-2);
}

.hero__badge-tag {
  padding: 4px 9px;
  border-radius: 999px;
  background: var(--ink);
  color: var(--accent);
  font-size: 11px;
  letter-spacing: 0.06em;
}

.hero__title {
  margin-top: 28px;
  font-size: clamp(46px, 6.6vw, 90px);
  line-height: 0.98;
  letter-spacing: -0.045em;
  font-weight: 600;
}

.hero__title mark {
  background: var(--accent);
  color: inherit;
  border-radius: 0.14em;
  padding: 0 0.1em;
  -webkit-box-decoration-break: clone;
  box-decoration-break: clone;
}

.hero__lede {
  margin-top: 28px;
  max-width: 560px;
  font-size: clamp(17px, 1.5vw, 20px);
  line-height: 1.5;
  color: var(--muted);
}

.hero__ctas {
  margin-top: 36px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.hero__note {
  margin-top: 20px;
  font-size: 14px;
  color: var(--muted);
}

.hero__visual {
  flex: 1 1 440px;
  min-width: 0;
  animation: rise 0.9s 0.1s cubic-bezier(0.2, 0.7, 0.2, 1) both;
}

.stage {
  position: relative;
  height: 640px;
  border-radius: 32px;
  background: var(--surface);
  overflow: hidden;
}

.stage__phone {
  position: absolute;
  top: 48px;
  left: 50%;
  transform: translateX(-50%);
}

.float {
  position: absolute;
  width: 214px;
  padding: 16px;
  border-radius: 18px;
}

.float--split {
  top: 64px;
  right: 18px;
  background: #fff;
  border: 1px solid var(--line);
  box-shadow: var(--shadow-float);
  animation: bob 6s ease-in-out infinite;
}

.float--payout {
  top: 318px;
  left: 18px;
  background: var(--ink);
  color: #fff;
  box-shadow: 0 18px 40px -20px rgb(11 11 12 / 0.5);
  animation: bob 6s 1.5s ease-in-out infinite;
}

.float__label {
  font-size: 12px;
  color: var(--muted);
}

.float__bar {
  display: flex;
  gap: 3px;
  height: 10px;
  margin-top: 12px;
}

.float__seg {
  border-radius: 999px;
}

.float__seg--accent {
  background: var(--accent);
  box-shadow: inset 0 0 0 1px rgb(11 11 12 / 0.18);
}

.float__seg--ink {
  background: var(--ink);
}

.float__seg--grey {
  background: #b5b5af;
}

.float__rows {
  margin: 14px 0 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 13px;
}

.float__rows div {
  display: flex;
  justify-content: space-between;
}

.float__rows dt {
  display: flex;
  align-items: center;
  gap: 8px;
}

.float__rows dd {
  margin: 0;
}

.float__key {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.float__head {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--on-ink-muted);
}

.float__amount {
  margin-top: 10px;
  font-size: 24px;
  letter-spacing: -0.02em;
}

.float__amount span,
.float__meta {
  font-size: 11px;
  color: var(--on-ink-muted);
}

.float__amount span {
  font-size: 13px;
}

.float__meta {
  margin-top: 8px;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
}

@keyframes bob {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
}

@media (max-width: 860px) {
  .float {
    display: none;
  }
}
</style>
