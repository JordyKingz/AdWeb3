<script setup lang="ts">
const faqs = [
  {
    q: 'How much can I earn?',
    a: "Honestly, a small bonus, not an income. Your share depends on how much brands spend and how much validated time you play. We'll never promise a rate we haven't funded.",
  },
  {
    q: 'Do I need crypto to play?',
    a: 'No. The games are free and work without a wallet. Connect one only when you want your share paid out.',
  },
  {
    q: 'Is connecting my wallet safe?',
    a: 'You sign a plain-text message that proves the wallet is yours. It works once and expires in minutes. We never ask for your seed phrase, your private key or a token approval.',
  },
  {
    q: "What happens when I'm offline?",
    a: "Keep playing. Progress saves on your device and syncs later. Earning pauses until you're back online, so every reward stays backed by validated activity.",
  },
  {
    q: 'When do I get paid?',
    a: "Activity is processed daily. Once a campaign's revenue clears its adjustment window, your share is finalized and paid automatically when your balance passes the payout threshold.",
  },
  {
    q: 'What does a brand actually pay for?',
    a: 'Validated, viewable impressions, nothing else. Duplicates, background time and invalid traffic are filtered out and itemized in your report.',
  },
]
</script>

<template>
  <section id="faq" class="faq">
    <div class="container faq__inner">
      <div class="faq__intro">
        <p class="eyebrow">FAQ</p>
        <h2 class="display-2 faq__title">Straight answers.</h2>
        <p class="faq__contact">Building the Android pilot. <a href="https://github.com/JordyKingz/AdWeb3">Follow development on GitHub</a>.</p>
      </div>
      <div class="faq__list">
        <details v-for="item in faqs" :key="item.q" class="faq__item">
          <summary>
            {{ item.q }}
            <span class="faq__plus"><AppIcon name="plus" :size="14" :stroke-width="2.2" /></span>
          </summary>
          <p class="faq__answer">{{ item.a }}</p>
        </details>
      </div>
    </div>
  </section>
</template>

<style scoped>
.faq {
  border-top: 1px solid var(--line);
  scroll-margin-top: 64px;
}

.faq__inner {
  display: flex;
  flex-wrap: wrap;
  gap: 48px;
  padding-block: var(--section-y);
}

.faq__intro {
  flex: 1 1 300px;
}

.faq__title {
  margin-top: 16px;
}

.faq__contact {
  margin-top: 20px;
  font-size: 16px;
  line-height: 1.55;
  color: var(--muted);
}

.faq__contact a {
  color: var(--ink);
  font-weight: 500;
}

.faq__list {
  flex: 2 1 560px;
  min-width: 0;
  border-top: 1px solid var(--line);
}

.faq__item {
  border-bottom: 1px solid var(--line);
}

.faq__item summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  min-height: 76px;
  font-size: 19px;
  font-weight: 500;
  letter-spacing: -0.01em;
  list-style: none;
  cursor: pointer;
}

.faq__item summary::-webkit-details-marker {
  display: none;
}

.faq__plus {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--surface);
  transition: transform 0.2s ease;
}

.faq__item[open] .faq__plus {
  transform: rotate(45deg);
}

.faq__answer {
  max-width: 640px;
  padding-bottom: 24px;
  font-size: 16px;
  line-height: 1.6;
  color: var(--muted);
}
</style>
