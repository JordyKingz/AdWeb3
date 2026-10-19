<script setup lang="ts">
import type { IconName } from '~/components/AppIcon.vue'

const states = [
  {
    chip: 'Pending',
    tone: 'chip--outline',
    title: "Points you've earned",
    text: "Validated play waiting for campaign revenue to clear. No fixed value yet, and we won't pretend otherwise.",
  },
  {
    chip: 'Finalized',
    tone: 'chip--ink',
    title: "USD₮ you're owed",
    text: "Your share of a funded epoch, locked in and fully backed. It can't be spent on anyone else.",
  },
  {
    chip: 'Paid',
    tone: 'chip--accent',
    title: 'USD₮ in your wallet',
    text: 'Sent automatically with a unique payout ID and a transaction hash you can look up.',
  },
]

const features: { icon: IconName; title: string; text: string }[] = [
  { icon: 'lock', title: 'Wallet-safe sign-in', text: 'Link a wallet by signing a message. No seed phrase, no token approvals, no gas.' },
  { icon: 'phone-check', title: 'Real installs only', text: 'Apps are verified with Google Play Integrity, and every event is checked on our servers.' },
  { icon: 'wifi', title: 'Works offline', text: "Games keep running without a connection. Earning just pauses until you're back." },
]
</script>

<template>
  <section class="container transparency">
    <SectionHeading
      eyebrow="Transparency"
      title="Nothing to take on faith."
      lede="Players see where every point stands. Brands see what they paid for. Anyone can check the chain."
    />
    <div class="transparency__states">
      <div v-for="s in states" :key="s.chip" class="state">
        <span class="chip" :class="s.tone">{{ s.chip }}</span>
        <h3 class="state__title">{{ s.title }}</h3>
        <p class="state__text">{{ s.text }}</p>
      </div>
    </div>
    <div class="transparency__features">
      <div v-for="f in features" :key="f.title" class="feature">
        <AppIcon :name="f.icon" :size="24" />
        <div>
          <h3 class="feature__title">{{ f.title }}</h3>
          <p class="feature__text">{{ f.text }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.transparency {
  padding-block: var(--section-y);
}

.transparency__states {
  margin-top: 56px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.state {
  flex: 1 1 300px;
  padding: 28px;
  border: 1px solid var(--line);
  border-radius: 24px;
}

.state .chip {
  min-height: 28px;
  padding: 0 12px;
  font-size: 12px;
}

.state__title {
  margin-top: 28px;
  font-size: 22px;
  letter-spacing: -0.02em;
  font-weight: 600;
}

.state__text {
  margin-top: 10px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}

.transparency__features {
  margin-top: 56px;
  padding-top: 40px;
  border-top: 1px solid var(--line);
  display: flex;
  flex-wrap: wrap;
  gap: 40px;
}

.feature {
  flex: 1 1 280px;
  display: flex;
  gap: 16px;
}

.feature__title {
  font-size: 17px;
  font-weight: 600;
}

.feature__text {
  margin-top: 6px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}
</style>
