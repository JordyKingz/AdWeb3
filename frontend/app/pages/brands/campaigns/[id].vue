<script setup lang="ts">
import type { CampaignReport } from '#shared/types/api'
import type { Bar } from '~/components/BarChart.vue'

definePageMeta({ layout: 'brand' })

const route = useRoute()
const id = computed(() => String(route.params.id))
const { data: report, error } = await useFetch<CampaignReport>(() => `/api/campaigns/${id.value}`, { key: `campaign-${id.value}` })

useSeoMeta({ title: () => `${report.value?.name ?? 'Campaign'} · AdWeb3` })

const net = useChain()
const toggling = ref(false)
const actionError = ref<string | null>(null)

const STATUS = {
  live: { label: 'Live', chip: 'chip--accent' },
  paused: { label: 'Paused', chip: 'chip--soft' },
  in_review: { label: 'In review', chip: 'chip--outline' },
  awaiting_funding: { label: 'Awaiting funding', chip: 'chip--outline' },
} as const

const r = computed(() => report.value!)
const isRunning = computed(() => r.value.status === 'live' || r.value.status === 'paused')
const invalidTotal = computed(() => r.value.invalid.reduce((sum, i) => sum + i.count, 0))
const invalidMax = computed(() => Math.max(1, ...r.value.invalid.map((i) => i.count)))

const chartBars = computed<Bar[]>(() =>
  r.value.delivery.daily.map((d) => ({ label: String(new Date(`${d.date}T00:00:00Z`).getUTCDate()), value: d.value, highlight: d.today && d.value !== null })),
)
const chartMax = computed(() => {
  const peak = Math.max(0, ...r.value.delivery.daily.map((d) => d.value ?? 0))
  return Math.max(200, Math.ceil(peak / 200) * 200)
})
const chartTicks = computed(() => Array.from({ length: chartMax.value / 200 + 1 }, (_, i) => i * 200))
const chartLabel = computed(() => {
  const served = r.value.delivery.daily.filter((d) => d.value !== null)
  if (!served.length) return 'No validated impressions yet.'
  return `Validated impressions per day: ${served.map((d) => `${shortDate(d.date)}, ${d.value}`).join('. ')}.`
})

const funnel = computed(() => {
  const d = r.value.delivery
  return [
    { label: 'Served', value: d.served, tone: '#cfcfca' },
    { label: 'Rendered', value: d.rendered, tone: '#a8a8a2' },
    { label: 'Viewable', value: d.viewable, tone: '#6b6b66' },
    { label: 'Validated · billed', value: d.validated, tone: 'var(--ink)', strong: true },
  ]
})

const totalClicks = computed(() => r.value.placements.reduce((sum, p) => sum + p.clicks, 0))
const deliveredPct = computed(() =>
  r.value.delivery.maxImpressions ? Math.round((r.value.delivery.validated / r.value.delivery.maxImpressions) * 100) : 0,
)
const billedPct = computed(() => {
  const budget = Number(r.value.budgetUnits)
  return budget ? Math.round((Number(r.value.billedUnits) / budget) * 100) : 0
})

async function togglePause() {
  toggling.value = true
  actionError.value = null
  try {
    report.value = await $fetch<CampaignReport>(`/api/campaigns/${id.value}/status`, {
      method: 'POST',
      body: { paused: r.value.status === 'live' },
    })
    await refreshNuxtData('brand-overview')
  } catch (err) {
    actionError.value = readableError(err)
  } finally {
    toggling.value = false
  }
}

function exportCsv() {
  const rows = [
    ['date', 'validated_impressions'],
    ...r.value.delivery.daily.filter((d) => d.value !== null).map((d) => [d.date, String(d.value)]),
  ]
  const blob = new Blob([rows.map((row) => row.join(',')).join('\n')], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = Object.assign(document.createElement('a'), { href: url, download: `${r.value.id}-delivery.csv` })
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div v-if="error" class="empty">
    <h1 class="page-title">Campaign not found</h1>
    <p class="empty__text">It may have been created in another session. Campaigns in the mock API live in memory.</p>
    <NuxtLink to="/brands/campaigns/new" class="btn btn--primary">Create a campaign</NuxtLink>
  </div>

  <div v-else-if="report">
    <p class="crumbs">Campaigns <span aria-hidden="true">/</span> {{ r.name }}</p>
    <div class="head">
      <div class="head__main">
        <div class="head__title-row">
          <h1 class="page-title">{{ r.name }}</h1>
          <span class="chip" :class="STATUS[r.status].chip"><span class="dot" />{{ STATUS[r.status].label }}</span>
          <span v-if="r.sample" class="chip chip--soft">Sample data</span>
        </div>
        <p class="head__meta">
          Candy Blocks · {{ r.placementsLabel }} · {{ shortDate(r.start) }} → {{ shortDate(r.end, true) }} UTC
          <template v-if="isRunning"> · Day {{ r.dayOfWindow }} of {{ r.durationDays }}</template>
        </p>
      </div>
      <div class="head__actions">
        <button v-if="isRunning" type="button" class="btn btn--secondary btn--sm" :disabled="toggling" @click="togglePause">
          <AppIcon :name="r.status === 'paused' ? 'play' : 'pause'" :size="14" />
          {{ r.status === 'paused' ? 'Resume' : 'Pause' }}
        </button>
        <button type="button" class="btn btn--secondary btn--sm" :disabled="!r.delivery.validated" @click="exportCsv">
          <AppIcon name="download" :size="15" :stroke-width="2" />Export CSV
        </button>
        <NuxtLink to="/brands/campaigns/new" class="btn btn--primary btn--sm">
          <AppIcon name="plus" :size="15" :stroke-width="2" />New campaign
        </NuxtLink>
      </div>
    </div>
    <p v-if="actionError" class="field-error" role="alert">{{ actionError }}</p>

    <div v-if="r.status === 'paused'" class="banner" role="status">
      <AppIcon name="pause" :size="16" />
      Paused. No new ad decisions are issued. Validated delivery stays billed, and the rest stays refundable.
    </div>
    <div v-else-if="r.status === 'in_review'" class="banner banner--light" role="status">
      <AppIcon name="eye" :size="18" />
      Funded and in creative review. Delivery starts once your sponsor card is approved; nothing is billed before then.
    </div>

    <div class="kpis">
      <div class="kpi">
        <div class="kpi__label">Validated impressions</div>
        <div class="kpi__value mono">{{ formatInt(r.delivery.validated) }}</div>
        <div class="meter"><span :style="{ width: `${deliveredPct}%` }" /></div>
        <div class="kpi__note">{{ deliveredPct }}% of {{ formatInt(r.delivery.maxImpressions) }} max</div>
      </div>
      <div class="kpi">
        <div class="kpi__label">Billed so far</div>
        <div class="kpi__value mono">{{ formatToken(r.billedUnits) }} <span>USD₮</span></div>
        <div class="meter"><span :style="{ width: `${billedPct}%` }" /></div>
        <div class="kpi__note">of a {{ formatToken(r.budgetUnits) }} USD₮ budget</div>
      </div>
      <div class="kpi">
        <div class="kpi__label">Refundable if it ended now</div>
        <div class="kpi__value mono">{{ formatToken(r.refundableUnits) }} <span>USD₮</span></div>
        <div class="kpi__note kpi__note--gap">held in escrow, never paid out</div>
      </div>
      <div class="kpi">
        <div class="kpi__label">Invalid traffic filtered</div>
        <div class="kpi__value mono">{{ formatInt(invalidTotal) }}</div>
        <div class="kpi__note kpi__note--gap">{{ formatPercent(invalidTotal, r.delivery.viewable) }} of viewable · not billed</div>
      </div>
    </div>

    <div class="grid">
      <section class="card card--wide">
        <div class="card__head">
          <h2 class="card-title">Validated impressions per day</h2>
          <span class="card__hint">UTC days · today is partial</span>
        </div>
        <BarChart class="card__chart" :bars="chartBars" :max="chartMax" :ticks="chartTicks" :summary="chartLabel" value-labels />
      </section>

      <section class="card">
        <h2 class="card-title">Delivery funnel</h2>
        <p class="card__hint">Only the last step is billed.</p>
        <div class="funnel">
          <div v-for="(step, i) in funnel" :key="step.label" class="funnel__step" :class="{ 'funnel__step--strong': step.strong }">
            <div class="funnel__row">
              <span>{{ step.label }}</span>
              <span class="mono">
                {{ formatInt(step.value) }}<template v-if="i > 0"> · {{ formatPercent(step.value, funnel[0]!.value) }}</template>
              </span>
            </div>
            <div class="funnel__track">
              <span :style="{ width: funnel[0]!.value ? `${(step.value / funnel[0]!.value) * 100}%` : '0%', background: step.tone }" />
            </div>
          </div>
        </div>
      </section>
    </div>

    <div class="grid">
      <section class="card">
        <h2 class="card-title">Why {{ formatInt(invalidTotal) }} impressions weren't billed</h2>
        <ul v-if="r.invalid.length" class="reasons">
          <li v-for="reason in r.invalid" :key="reason.reason">
            <span class="reasons__label">{{ reason.reason }}</span>
            <span class="reasons__track" aria-hidden="true"><span :style="{ width: `${(reason.count / invalidMax) * 100}%` }" /></span>
            <span class="mono reasons__count">{{ reason.count }}</span>
          </li>
        </ul>
        <p v-else class="card__hint card__hint--gap">Nothing filtered yet.</p>
      </section>

      <section class="card card--wide">
        <h2 class="card-title">By placement</h2>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">Placement</th>
                <th scope="col" class="num">Validated</th>
                <th scope="col" class="num">Clicks</th>
                <th scope="col" class="num">CTR</th>
                <th scope="col" class="num">Billed</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in r.placements" :key="p.name">
                <td>
                  <div class="cell-title">{{ p.name }}</div>
                  <div class="cell-sub">Candy Blocks · Android</div>
                </td>
                <td class="num mono">{{ formatInt(p.validated) }}</td>
                <td class="num mono">{{ formatInt(p.clicks) }}</td>
                <td class="num mono">{{ formatPercent(p.clicks, p.validated) }}</td>
                <td class="num mono">{{ formatToken(p.billedUnits) }}</td>
              </tr>
              <tr class="total">
                <td>Total</td>
                <td class="num mono">{{ formatInt(r.delivery.validated) }}</td>
                <td class="num mono">{{ formatInt(totalClicks) }}</td>
                <td class="num mono">{{ formatPercent(totalClicks, r.delivery.validated) }}</td>
                <td class="num mono">{{ formatToken(r.billedUnits) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <div class="grid">
      <section class="card card--ink card--wide">
        <div class="card__head">
          <h2 class="card-title">Escrow &amp; settlement</h2>
          <span class="mono card__tag">SPLIT LOCKED · V1</span>
        </div>
        <div class="escrow">
          <div>
            <div class="escrow__label">Deposited</div>
            <div class="escrow__value mono">{{ formatToken(r.budgetUnits) }}</div>
          </div>
          <div>
            <div class="escrow__label">Earned</div>
            <div class="escrow__value mono">{{ formatToken(r.billedUnits) }}</div>
          </div>
          <div>
            <div class="escrow__label">Still refundable</div>
            <div class="escrow__value escrow__value--accent mono">{{ formatToken(r.refundableUnits) }}</div>
          </div>
        </div>
        <div class="escrow__bar" aria-hidden="true">
          <span style="flex: 30; background: var(--accent)" /><span style="flex: 60; background: #fff" /><span
            style="flex: 10; background: #6b6b66"
          />
        </div>
        <div class="escrow__legend">
          <span><span class="mono">{{ formatToken(r.split.players, 3) }}</span> to players</span>
          <span><span class="mono">{{ formatToken(r.split.operations, 3) }}</span> to operations</span>
          <span><span class="mono">{{ formatToken(r.split.contingency, 3) }}</span> to contingency</span>
        </div>
        <div class="escrow__foot">
          <span v-if="r.funding.txHash">
            Deposit <a :href="net.txUrl(r.funding.txHash)" target="_blank" rel="noopener" class="mono">{{ shortHash(r.funding.txHash) }}</a>
            · {{ r.funding.fundedAt ? shortDate(r.funding.fundedAt) : '' }}
          </span>
          <span v-else-if="r.funding.mode === 'demo'">Demo funding · no on-chain transaction</span>
          <span class="mono">Order {{ shortHash(r.funding.orderId) }}</span>
        </div>
      </section>

      <section class="card">
        <h2 class="card-title">Purchased terms</h2>
        <dl class="rows terms">
          <div>
            <dt>Rate</dt>
            <dd class="mono">{{ formatToken(r.terms.cpmUnits) }} CPM</dd>
          </div>
          <div>
            <dt>Window</dt>
            <dd class="mono">{{ r.durationDays * 24 }} h</dd>
          </div>
          <div>
            <dt>Frequency cap</dt>
            <dd>{{ r.terms.frequencyCap }} / player / day</dd>
          </div>
          <div>
            <dt>Audience</dt>
            <dd>{{ r.terms.incentivized ? 'Reward-enabled' : 'Standard' }}</dd>
          </div>
          <div>
            <dt>Creative</dt>
            <dd>
              <span class="chip chip--plain" :class="r.terms.creativeStatus === 'approved' ? 'chip--accent' : 'chip--soft'">
                {{ r.terms.creativeStatus === 'approved' ? 'Approved' : 'In review' }}
              </span>
            </dd>
          </div>
        </dl>
        <p class="card__hint card__hint--gap">Unspent budget is refunded after the adjustment window closes.</p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.crumbs {
  font-size: 14px;
  color: var(--muted);
}

.page-title {
  font-size: clamp(28px, 3vw, 38px);
  letter-spacing: -0.035em;
  font-weight: 600;
}

.empty {
  max-width: 480px;
  padding: 48px 0;
}

.empty__text {
  margin: 12px 0 24px;
  color: var(--muted);
  line-height: 1.5;
}

.head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
  margin-top: 8px;
}

.head__main {
  min-width: 0;
}

.head__title-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.head__title-row .chip {
  min-height: 28px;
  font-size: 12px;
}

.dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
}

.head__meta {
  margin-top: 8px;
  font-size: 14px;
  color: var(--muted);
}

.head__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.banner {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding: 14px 16px;
  border-radius: 14px;
  background: var(--ink);
  color: #fff;
  font-size: 14px;
  line-height: 1.45;
}

.banner :deep(.icon) {
  color: var(--accent);
}

.banner--light {
  background: #fff;
  border: 1px solid var(--line);
  color: var(--ink-2);
}

.banner--light :deep(.icon) {
  color: var(--ink);
}

.kpis {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 24px;
}

.kpi {
  flex: 1 1 220px;
  padding: 20px;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #fff;
}

.kpi__label,
.kpi__note {
  font-size: 13px;
  color: var(--muted);
}

.kpi__note {
  margin-top: 8px;
  font-size: 12px;
}

.kpi__note--gap {
  margin-top: 20px;
}

.kpi__value {
  margin-top: 10px;
  font-size: 32px;
  letter-spacing: -0.03em;
  line-height: 1;
}

.kpi__value span {
  font-size: 14px;
  color: var(--muted);
}

.meter {
  height: 6px;
  margin-top: 14px;
  border-radius: 999px;
  background: var(--line-soft);
  overflow: hidden;
}

.meter span {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: var(--ink);
}

.grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 12px;
}

.grid > .card {
  flex: 1 1 300px;
  min-width: 0;
  border-radius: 18px;
}

.grid > .card--wide {
  flex: 2 1 480px;
}

.card__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
}

.card__hint {
  font-size: 12px;
  color: var(--muted);
}

.card__hint--gap {
  margin-top: 14px;
}

.card__chart {
  margin-top: 24px;
}

.card__tag {
  font-size: 11px;
  letter-spacing: 0.06em;
  color: var(--on-ink-muted);
}

.funnel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 20px;
}

.funnel__row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.funnel__step--strong .funnel__row {
  font-weight: 600;
}

.funnel__track {
  height: 10px;
  margin-top: 6px;
  border-radius: 999px;
  background: #f0f0ec;
}

.funnel__track span {
  display: block;
  height: 100%;
  border-radius: 999px;
}

.reasons {
  list-style: none;
  margin: 18px 0 0;
  padding: 0;
}

.reasons li {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-top: 1px solid var(--line-soft);
  font-size: 14px;
}

.reasons__label {
  flex: 1;
}

.reasons__track {
  width: 90px;
  height: 6px;
  border-radius: 999px;
  background: #f0f0ec;
}

.reasons__track span {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: var(--ink);
}

.reasons__count {
  width: 32px;
  text-align: right;
}

.table-wrap {
  margin-top: 14px;
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 480px;
  border-collapse: collapse;
  font-size: 14px;
}

th {
  padding: 10px 0;
  font-size: 12px;
  font-weight: 500;
  color: var(--muted);
  text-align: left;
}

td {
  padding: 14px 0;
  border-top: 1px solid var(--line-soft);
}

.num {
  text-align: right;
}

.cell-title {
  font-weight: 500;
}

.cell-sub {
  font-size: 12px;
  color: var(--muted);
}

.total td {
  border-top: 1.5px solid var(--ink);
  font-weight: 600;
}

.escrow {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-top: 20px;
}

.escrow > div {
  flex: 1 1 120px;
}

.escrow__label {
  font-size: 12px;
  color: var(--on-ink-muted);
}

.escrow__value {
  margin-top: 6px;
  font-size: 24px;
  letter-spacing: -0.02em;
}

.escrow__value--accent {
  color: var(--accent);
}

.escrow__bar {
  display: flex;
  gap: 3px;
  height: 12px;
  margin-top: 24px;
}

.escrow__bar span {
  border-radius: 999px;
}

.escrow__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 28px;
  margin-top: 14px;
  font-size: 13px;
  color: var(--on-ink-muted);
}

.escrow__legend .mono {
  color: #fff;
}

.escrow__foot {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
  margin-top: 24px;
  padding-top: 18px;
  border-top: 1px solid #2a2a2a;
  font-size: 13px;
  color: var(--on-ink-muted);
}

.escrow__foot a {
  color: #fff;
}

.terms {
  margin-top: 14px;
}
</style>
