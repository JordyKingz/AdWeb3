<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
</script>

<template>
  <section id="pilot" class="container cta-wrap">
    <div class="cta">
      <div>
        <h2 class="cta__title">Put your budget where the players are.</h2>
        <p v-if="landingOnly" class="cta__pilot-note">We are building the Android pilot. Campaign funding and player rewards open after testing is complete.</p>
      </div>
      <div v-if="!landingOnly" class="cta__actions">
        <NuxtLink to="/brands/campaigns/new" class="btn btn--primary btn--lg">
          Start a campaign <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
        </NuxtLink>
        <NuxtLink to="/connect" class="btn btn--outline btn--lg">Play &amp; earn</NuxtLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.cta-wrap {
  padding-bottom: var(--section-y);
}

.cta {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  justify-content: space-between;
  gap: 32px;
  padding: clamp(36px, 6vw, 80px);
  border-radius: 32px;
  background: var(--accent);
}

.cta__title {
  flex: 1 1 520px;
  font-size: clamp(38px, 5vw, 70px);
  line-height: 1;
  letter-spacing: -0.045em;
  font-weight: 600;
}

.cta__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.cta__pilot-note { margin-top: 18px; max-width: 560px; line-height: 1.6; }
</style>
