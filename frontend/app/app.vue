<script setup lang="ts">
const landingOnly = useRuntimeConfig().public.landingOnly
</script>

<template>
  <div v-if="landingOnly" class="pilot-notice" role="status">Pilot in development. Live campaigns and cash payouts are not available yet.</div>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<style>
.pilot-notice { padding: 10px 20px; text-align: center; font-size: 13px; background: #0B0B0C; color: #fff; }
</style>
