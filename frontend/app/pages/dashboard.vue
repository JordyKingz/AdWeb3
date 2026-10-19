<script setup lang="ts">
import type { ActivityKind, PlayerSummary } from '#shared/types/api'
import type { Bar } from '~/components/BarChart.vue'
import type { IconName } from '~/components/AppIcon.vue'

definePageMeta({ layout: 'player' })
useSeoMeta({ title: 'Your rewards · AdWeb3' })

const { data: summary, error } = await useFetch<PlayerSummary>('/api/player/summary', { key: 'player-summary' })
const online = useOnline()
const net = useChain()
const wallet = useWallet()

const s = computed(() => summary.value!)
const finalizedPct = computed(() => Math.min(100, (Number(s.value.finalizedUnits) / Number(s.value.payoutThresholdUnits)) * 100))
const toGoUnits = computed(() => {
  const gap = BigInt(s.value.payoutThresholdUnits) - BigInt(s.value.finalizedUnits)
  return gap > 0n ? gap : 0n
})

const bars = computed<Bar[]>(() =>
  s.value.daily.map((d, i) => {
    const day = new Date(`${d.date}T00:00:00Z`).getUTCDate()
    return { label: i % 2 === 0 || d.today ? String(day) : '', value: d.points, highlight: d.today }
  }),
)
const chartLabel = computed(() => {
  const days = s.value.daily
  const offline = days.filter((d) => d.points === null).map((d) => shortDate(d.date))
  const best = days.reduce((a, b) => ((b.points ?? 0) > (a.points ?? 0) ? b : a))
  return `Validated points per day for the last ${days.length} days. Highest ${best.points} on ${shortDate(best.date)}. Offline with no points on ${offline.join(', ') || 'no days'}. Daily cap is ${s.value.dailyCap}.`
})

// Countdown to the next daily processing run.
const now = ref(Date.now())
let timer: ReturnType<typeof setInterval> | undefined
onMounted(() => {
  timer = setInterval(() => (now.value = Date.now()), 30_000)
})
onBeforeUnmount(() => clearInterval(timer))
const countdown = computed(() => {
  const ms = Math.max(0, Date.parse(s.value.nextProcessingAt) - now.value)
  const h = Math.floor(ms / 3_600_000)
  const m = Math.floor((ms % 3_600_000) / 60_000)
  return `${h}h ${String(m).padStart(2, '0')}m`
})

const ACTIVITY: Record<ActivityKind, { icon: IconName; tone: string; chip: string }> = {
  session: { icon: 'gamepad', tone: 'soft', chip: 'chip--soft' },
  offline: { icon: 'wifi-off', tone: 'dashed', chip: 'chip--outline-soft' },
  rejected: { icon: 'x-circle', tone: 'danger', chip: 'chip--danger' },
  finalized: { icon: 'check', tone: 'ink', chip: 'chip--ink' },
}

async function signOut() {
  await $fetch('/api/auth/logout', { method: 'POST' })
  wallet.disconnect()
  await refreshNuxtData('player-session')
  await navigateTo('/')
}
</script>

<template>
  <main class="container dash">
    <div v-if="error" class="gate">
      <span class="gate__icon"><AppIcon name="lock" :size="22" /></span>
      <h1 class="gate__title">Connect your wallet to see your rewards</h1>
      <p class="gate__text">Sign a short message to prove the wallet is yours. No gas, no approvals.</p>
      <NuxtLink to="/connect" class="btn btn--primary btn--lg">
        Connect wallet <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
      </NuxtLink>
    </div>

    <template v-else-if="summary">
      <div id="overview" class="dash__head">
        <div>
          <div class="dash__title-row">
            <h1 class="dash__title">Your rewards</h1>
            <span v-if="s.sample" class="chip chip--soft">Sample data</span>
          </div>
          <p class="dash__sub">{{ s.game.name }} · activity is processed daily at 00:00 UTC</p>
        </div>
        <div class="dash__clock"><AppIcon name="clock" :size="16" />Next processing in <span class="mono">{{ countdown }}</span></div>
      </div>

      <div v-if="!online" class="offline" role="status">
        <AppIcon name="wifi-off" :size="20" :stroke-width="2" />
        <p><strong>You're offline.</strong> Keep playing. Progress saves on this device. Earning is paused until you reconnect.</p>
      </div>

      <div class="balances">
        <section class="balance">
          <span class="chip chip--outline">Pending</span>
          <div class="balance__value mono">{{ formatInt(s.pendingPoints) }} <span>pts</span></div>
          <p class="balance__text">Validated play waiting for campaign revenue to clear. No fixed {{ net.tokenSymbol }} value yet.</p>
        </section>

        <section class="balance balance--main">
          <div class="balance__top">
            <span class="chip balance__chip-light">Finalized</span>
            <span class="balance__hint">backed &amp; locked</span>
          </div>
          <div class="balance__value balance__value--lg mono">{{ formatToken(s.finalizedUnits, 4) }} <span>{{ net.tokenSymbol }}</span></div>
          <div
            class="balance__meter"
            role="progressbar"
            :aria-valuenow="Math.round(finalizedPct)"
            aria-valuemin="0"
            aria-valuemax="100"
            aria-label="Progress to automatic payout"
          >
            <span :style="{ width: `${finalizedPct}%` }" />
          </div>
          <div class="balance__foot">
            <span>{{ formatToken(toGoUnits, 4) }} to go</span>
            <span>Auto-payout at <strong class="mono">{{ formatToken(s.payoutThresholdUnits) }} {{ net.tokenSymbol }}</strong></span>
          </div>
        </section>

        <section class="balance">
          <span class="chip chip--accent">Paid</span>
          <div class="balance__value mono">{{ formatToken(s.paidUnits, 4) }} <span>{{ net.tokenSymbol }}</span></div>
          <p class="balance__text">
            {{ s.payoutCount }} payouts to your wallet<template v-if="s.lastPayoutAt">. Last one {{ shortDate(s.lastPayoutAt) }}</template>.
          </p>
        </section>
      </div>

      <div class="row">
        <section class="card card--wide">
          <div class="card__head">
            <h2 class="card-title">Validated points, last {{ s.daily.length }} days</h2>
            <div class="legend">
              <span><i class="legend__swatch" />Validated</span>
              <span><i class="legend__swatch legend__swatch--dashed" />Offline, not earning</span>
            </div>
          </div>
          <BarChart
            class="card__chart"
            :bars="bars"
            :max="s.dailyCap"
            :ticks="[0, s.dailyCap / 3, (s.dailyCap / 3) * 2, s.dailyCap]"
            :height="200"
            cap-label="daily cap"
            :summary="chartLabel"
          />
        </section>

        <section class="card today">
          <h2 class="card-title">Today</h2>
          <div class="today__value mono">{{ s.today.points }} <span>/ {{ s.dailyCap }} pts</span></div>
          <div class="today__meter"><span :style="{ width: `${(s.today.points / s.dailyCap) * 100}%` }" /></div>
          <dl class="rows today__rows">
            <div><dt>Online play</dt><dd class="mono">{{ s.today.onlineMinutes }} min</dd></div>
            <div><dt>Sponsor cards seen</dt><dd class="mono">{{ s.today.sponsorCardsSeen }}</dd></div>
            <div><dt>Status</dt><dd>{{ online ? 'Earning' : 'Paused, offline' }}</dd></div>
          </dl>
          <p class="today__note">Points come from validated online play and sponsor cards you actually see. Minutes alone don't earn.</p>
        </section>
      </div>

      <div class="row">
        <section id="activity" class="card card--wide">
          <h2 class="card-title">Recent activity</h2>
          <ul class="activity">
            <li v-for="(item, i) in s.activity" :key="i">
              <span class="activity__icon" :class="`activity__icon--${ACTIVITY[item.kind].tone}`">
                <AppIcon :name="ACTIVITY[item.kind].icon" :size="18" />
              </span>
              <div class="activity__body">
                <div class="activity__title">{{ item.title }}</div>
                <div class="activity__detail">{{ item.detail }}</div>
              </div>
              <span class="mono activity__amount" :class="{ 'activity__amount--muted': item.amount === '0' }">{{ item.amount }}</span>
              <span class="chip chip--plain" :class="ACTIVITY[item.kind].chip">{{ item.status }}</span>
            </li>
          </ul>
        </section>

        <div id="wallet" class="side">
          <section class="card">
            <h2 class="card-title">Payout wallet</h2>
            <div class="wallet-row">
              <AddressAvatar :size="40" />
              <div>
                <div class="mono wallet-row__addr">{{ shortAddress(s.address) }}</div>
                <div class="wallet-row__sub">{{ net.chainName }} · verified {{ shortDate(s.linkedAt) }}</div>
              </div>
            </div>
            <NuxtLink to="/connect?change=1" class="btn btn--secondary btn--sm btn--block side__btn">Change wallet</NuxtLink>
            <p class="side__note">Changing needs a new signature. Rewards already finalized keep their original destination.</p>
            <button type="button" class="btn btn--ghost btn--sm btn--block" @click="signOut">
              <AppIcon name="logout" :size="16" />Sign out
            </button>
          </section>

          <section class="card">
            <h2 class="card-title">Linked games</h2>
            <div class="game">
              <span class="game__icon" aria-hidden="true"><i /><i /><i /><i /></span>
              <div class="game__body">
                <div class="game__name">{{ s.game.name }}</div>
                <div class="game__sub">{{ s.game.platform }} · this device</div>
              </div>
              <span class="chip chip--plain chip--accent">Linked</span>
            </div>
            <p class="game__more">More games are on the way.</p>
          </section>
        </div>
      </div>

      <section id="payouts" class="card payouts">
        <div class="card__head">
          <h2 class="card-title">Payouts</h2>
          <span class="card__hint">Sent automatically. Gas is covered by AdWeb3.</span>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">Payout ID</th>
                <th scope="col">Epoch</th>
                <th scope="col" class="num">Amount</th>
                <th scope="col" class="pad">Transaction</th>
                <th scope="col" class="num">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in s.payouts" :key="p.id">
                <td class="mono">{{ p.id }}</td>
                <td>{{ p.epoch }}</td>
                <td class="num mono">{{ formatToken(p.amountUnits, 4) }} {{ net.tokenSymbol }}</td>
                <td class="pad mono">
                  <span v-if="s.sample">{{ shortHash(p.txHash) }}</span>
                  <a v-else :href="net.txUrl(p.txHash)" target="_blank" rel="noopener">{{ shortHash(p.txHash) }}</a>
                </td>
                <td class="num">
                  <span class="chip chip--plain" :class="p.status === 'confirmed' ? 'chip--accent' : 'chip--soft'">
                    {{ p.status === 'confirmed' ? 'Confirmed' : 'Submitted' }}
                  </span>
                </td>
              </tr>
              <tr class="pending-row">
                <td class="mono">Next</td>
                <td>When finalized reaches {{ formatToken(s.payoutThresholdUnits) }}</td>
                <td class="num mono">{{ formatToken(s.finalizedUnits, 4) }} {{ net.tokenSymbol }}</td>
                <td class="pad">—</td>
                <td class="num"><span class="chip chip--plain chip--outline-soft">Accruing</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <p class="disclaimer">
        Rewards are a bonus funded by real ad revenue, not an income. Amounts depend on campaign funding and validated
        activity, and availability depends on your country.
      </p>
    </template>
  </main>
</template>

<style scoped>
.dash {
  padding-block: clamp(24px, 4vw, 48px) 56px;
}

.gate {
  max-width: 520px;
  margin: clamp(32px, 8vw, 96px) auto;
  padding: clamp(28px, 4vw, 40px);
  border: 1px solid var(--line);
  border-radius: 28px;
  background: #fff;
  text-align: center;
}

.gate__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--accent);
}

.gate__title {
  margin-top: 20px;
  font-size: 26px;
  letter-spacing: -0.03em;
  font-weight: 600;
}

.gate__text {
  margin: 10px 0 24px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}

.dash__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
  scroll-margin-top: 24px;
}

.dash__title-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.dash__title {
  font-size: clamp(30px, 3.4vw, 44px);
  letter-spacing: -0.04em;
  font-weight: 600;
}

.dash__sub {
  margin-top: 8px;
  font-size: 15px;
  color: var(--muted);
}

.dash__clock {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #fff;
  font-size: 13px;
  color: var(--ink-2);
}

.dash__clock .mono {
  font-weight: 500;
}

.offline {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 20px;
  padding: 16px 18px;
  border-radius: 16px;
  background: var(--ink);
  color: #c9c9c4;
  font-size: 14px;
  line-height: 1.45;
}

.offline :deep(.icon) {
  color: var(--accent);
}

.offline strong {
  color: #fff;
}

.balances,
.row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 12px;
}

.balances {
  margin-top: 24px;
}

.balance {
  flex: 1 1 280px;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 24px;
  border: 1px solid var(--line);
  border-radius: 22px;
  background: #fff;
}

.balance--main {
  flex: 1.3 1 320px;
  align-items: stretch;
  border-color: var(--ink);
  background: var(--ink);
  color: #fff;
}

.balance__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance__chip-light {
  background: #fff;
  color: var(--ink);
}

.balance__hint {
  font-size: 12px;
  color: var(--on-ink-muted);
}

.balance__value {
  margin-top: 24px;
  font-size: 40px;
  letter-spacing: -0.04em;
  line-height: 1;
}

.balance__value--lg {
  font-size: 48px;
}

.balance__value span {
  font-size: 15px;
  color: var(--muted);
}

.balance--main .balance__value span {
  font-size: 16px;
  color: var(--on-ink-muted);
}

.balance__text {
  margin-top: 14px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--muted);
}

.balance__meter {
  height: 8px;
  margin-top: 22px;
  border-radius: 999px;
  background: #2a2a2a;
  overflow: hidden;
}

.balance__meter span {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: var(--accent);
}

.balance__foot {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 6px;
  margin-top: 10px;
  font-size: 13px;
  color: var(--on-ink-muted);
}

.balance__foot strong {
  color: #fff;
  font-weight: 400;
}

.row > .card,
.side {
  flex: 1 1 300px;
  min-width: 0;
}

.row > .card--wide {
  flex: 2 1 520px;
}

.card__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
}

.card__hint {
  font-size: 13px;
  color: var(--muted);
}

.card__chart {
  margin-top: 24px;
}

.legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 16px;
  font-size: 12px;
  color: var(--muted);
}

.legend span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.legend__swatch {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  background: var(--ink);
}

.legend__swatch--dashed {
  background: none;
  border: 1.5px dashed var(--line-dashed);
}

.today {
  display: flex;
  flex-direction: column;
}

.today__value {
  margin-top: 18px;
  font-size: 34px;
  letter-spacing: -0.03em;
  line-height: 1;
}

.today__value span {
  font-size: 15px;
  color: var(--muted);
}

.today__meter {
  height: 8px;
  margin-top: 14px;
  border-radius: 999px;
  background: #f0f0ec;
  overflow: hidden;
}

.today__meter span {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: var(--ink);
}

.today__rows {
  margin-top: 18px;
}

.today__note {
  margin-top: auto;
  padding-top: 14px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--muted);
}

.activity {
  list-style: none;
  margin: 14px 0 0;
  padding: 0;
  scroll-margin-top: 24px;
}

.activity li {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 0;
  border-top: 1px solid var(--line-soft);
}

.activity__icon {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 12px;
}

.activity__icon--soft {
  background: var(--surface);
}

.activity__icon--dashed {
  border: 1.5px dashed var(--line-dashed);
  color: var(--muted);
}

.activity__icon--danger {
  background: var(--danger-bg);
  color: var(--danger);
}

.activity__icon--ink {
  background: var(--ink);
  color: var(--accent);
}

.activity__body {
  flex: 1;
  min-width: 0;
}

.activity__title {
  font-size: 15px;
  font-weight: 500;
}

.activity__detail {
  margin-top: 2px;
  font-size: 13px;
  color: var(--muted);
}

.activity__amount {
  font-size: 14px;
}

.activity__amount--muted {
  color: var(--muted);
}

.chip--outline-soft {
  border: 1px solid var(--line-strong);
  color: var(--ink-2);
}

.side {
  display: flex;
  flex-direction: column;
  gap: 12px;
  scroll-margin-top: 24px;
}

.wallet-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
}

.wallet-row__addr {
  font-size: 15px;
}

.wallet-row__sub {
  margin-top: 2px;
  font-size: 13px;
  color: var(--muted);
}

.side__btn {
  margin-top: 18px;
}

.side__note {
  margin: 12px 0 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--muted);
}

.game {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
}

.game__icon {
  flex: none;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 3px;
  width: 40px;
  height: 40px;
  padding: 8px;
  border-radius: 11px;
  background: #fff4e0;
}

.game__icon i {
  border-radius: 3px;
}

.game__icon i:nth-child(1) {
  background: #ff8fb1;
}

.game__icon i:nth-child(2) {
  background: #7aa7ff;
}

.game__icon i:nth-child(3) {
  background: #ffb547;
}

.game__icon i:nth-child(4) {
  background: #b9f23a;
}

.game__body {
  flex: 1;
  min-width: 0;
}

.game__name {
  font-size: 15px;
  font-weight: 500;
}

.game__sub {
  margin-top: 2px;
  font-size: 13px;
  color: var(--muted);
}

.game__more {
  margin-top: 12px;
  padding: 12px 14px;
  border: 1.5px dashed var(--line-strong);
  border-radius: 12px;
  font-size: 13px;
  color: var(--muted);
}

.payouts {
  margin-top: 12px;
  scroll-margin-top: 24px;
}

.table-wrap {
  margin-top: 14px;
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 620px;
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

.pad {
  padding-left: 24px;
}

.pending-row {
  color: var(--muted);
}

.disclaimer {
  max-width: 720px;
  margin-top: 28px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--muted);
}
</style>
