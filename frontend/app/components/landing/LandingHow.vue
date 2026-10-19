<script setup lang="ts">
import type { IconName } from '~/components/AppIcon.vue'

const steps: { title: string; text: string; icon: IconName }[] = [
  { title: 'Fund', icon: 'card', text: 'A brand deposits USD₮ into campaign escrow. It stays refundable until delivery earns it.' },
  { title: 'Serve', icon: 'phone', text: 'Sponsor cards appear on game home and results screens. Never over live gameplay.' },
  {
    title: 'Validate',
    icon: 'shield',
    text: 'Every impression is checked server-side: viewable, unique, from a verified install. Bots and duplicates earn nothing.',
  },
  {
    title: 'Settle',
    icon: 'coin',
    text: 'Earned revenue splits 30 / 60 / 10. The player share settles in USD₮ and pays out to verified wallets.',
  },
]
</script>

<template>
  <section id="how" class="container how">
    <SectionHeading
      eyebrow="How it works"
      title="How the money moves."
      lede="Every USD₮ is accounted for, from a brand's deposit to a player's wallet. Nothing is paid out until delivery has earned it."
    />
    <ol class="how__steps">
      <li v-for="(step, i) in steps" :key="step.title" class="how__step" :class="{ 'how__step--final': i === steps.length - 1 }">
        <div class="how__top">
          <span class="mono how__num">0{{ i + 1 }}</span>
          <span class="how__icon"><AppIcon :name="step.icon" :size="20" /></span>
        </div>
        <h3 class="how__title">{{ step.title }}</h3>
        <p class="how__text">{{ step.text }}</p>
      </li>
    </ol>
  </section>
</template>

<style scoped>
.how {
  padding-block: var(--section-y);
  scroll-margin-top: 64px;
}

.how__steps {
  list-style: none;
  margin: 56px 0 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.how__step {
  flex: 1 1 250px;
  min-height: 280px;
  display: flex;
  flex-direction: column;
  padding: 28px;
  border: 1px solid var(--line);
  border-radius: 24px;
  transition:
    transform 0.2s ease,
    border-color 0.2s ease;
}

.how__step:hover {
  transform: translateY(-3px);
  border-color: var(--line-strong);
}

.how__step--final {
  background: var(--ink);
  border-color: var(--ink);
  color: #fff;
}

.how__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.how__num {
  font-size: 13px;
  color: var(--muted);
}

.how__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--surface);
}

.how__title {
  margin-top: auto;
  padding-top: 56px;
  font-size: 22px;
  letter-spacing: -0.02em;
  font-weight: 600;
}

.how__text {
  margin-top: 10px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}

.how__step--final .how__num,
.how__step--final .how__text {
  color: var(--on-ink-muted);
}

.how__step--final .how__icon {
  background: var(--accent);
  color: var(--ink);
}
</style>
