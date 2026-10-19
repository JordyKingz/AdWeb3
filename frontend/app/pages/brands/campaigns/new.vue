<script setup lang="ts">
import { parseUnits } from 'viem'
import {
  CPM_UNITS,
  CTA_OPTIONS,
  DURATIONS,
  FREQUENCY_CAPS,
  HEADLINE_MAX,
  MIN_BUDGET,
  PLACEMENTS,
  TOKEN_DECIMALS,
  isValidBudget,
  maxImpressions,
  type CtaOption,
  type Duration,
  type FrequencyCap,
  type PlacementId,
} from '#shared/utils/campaign'

definePageMeta({ layout: 'brand' })
useSeoMeta({ title: 'New campaign · AdWeb3' })

const DRAFT_KEY = 'adweb3:campaign-draft'
const DAY_MS = 86_400_000
const todayIso = () => new Date().toISOString().slice(0, 10)

const form = reactive({
  name: 'Spring launch',
  brand: 'Your Brand',
  url: 'https://yourbrand.com/spring',
  start: new Date(Date.now() + DAY_MS).toISOString().slice(0, 10),
  durationDays: 7 as Duration,
  placements: ['home', 'results'] as PlacementId[],
  incentivized: true,
  frequencyCap: 3 as FrequencyCap,
  headline: 'Your brand, seen by real players',
  cta: 'Learn more' as CtaOption,
  budget: '25',
})

const STEPS = ['Basics', 'Placements', 'Creative', 'Fund']
const step = ref(1)
const attempted = reactive<Record<number, boolean>>({})

const net = useChain()
const wallet = useWallet()
const funding = useCampaignFunding()
const locked = computed(() => funding.busy.value || funding.stage.value === 'funded')

// ---------- Validation ----------

function isHttpsUrl(value: string) {
  try {
    return new URL(value).protocol === 'https:'
  } catch {
    return false
  }
}

const errors = computed<Record<number, Record<string, string>>>(() => ({
  1: {
    name: form.name.trim() ? '' : 'Give the campaign a name.',
    brand: form.brand.trim() ? '' : 'Add your brand name.',
    url: isHttpsUrl(form.url) ? '' : 'Use a full https:// address.',
    start: form.start >= todayIso() ? '' : 'Pick today or a later date.',
  },
  2: { placements: form.placements.length ? '' : 'Choose at least one placement.' },
  3: { headline: form.headline.trim() ? '' : 'Add a headline.' },
  4: { budget: isValidBudget(form.budget) ? '' : `Enter ${MIN_BUDGET} to 100,000 USD₮, up to two decimals.` },
}))

const stepValid = (n: number) => Object.values(errors.value[n] ?? {}).every((e) => !e)
const fieldError = (n: number, field: string) => (attempted[n] ? errors.value[n]?.[field] || '' : '')

function goTo(target: number) {
  for (let i = 1; i < target; i++) {
    if (!stepValid(i)) {
      attempted[i] = true
      step.value = i
      return
    }
  }
  step.value = target
}

function next() {
  attempted[step.value] = true
  if (stepValid(step.value)) step.value = Math.min(STEPS.length, step.value + 1)
}

const stepDone = (n: number) => (n === 4 ? funding.stage.value === 'funded' : n < step.value && stepValid(n))

// ---------- Derived values ----------

const budgetUnits = computed(() => (isValidBudget(form.budget) ? parseUnits(form.budget, TOKEN_DECIMALS) : 0n))
const budgetLabel = computed(() => formatToken(budgetUnits.value))
const maxImps = computed(() => formatInt(maxImpressions(budgetUnits.value)))
const endDate = computed(() => {
  const start = Date.parse(`${form.start}T00:00:00Z`)
  return Number.isNaN(start) ? form.start : new Date(start + form.durationDays * DAY_MS).toISOString().slice(0, 10)
})
const placementsLabel = computed(() => {
  const names = PLACEMENTS.filter((p) => form.placements.includes(p.id)).map((p) => p.short)
  return names.length ? names.join(' + ') : 'None selected'
})

const checklist = computed(() => [
  { label: 'Basics', status: stepValid(1) ? 'Done' : 'Incomplete', done: stepValid(1) },
  {
    label: 'Placements',
    status: form.placements.length ? `${form.placements.length} selected` : 'Pick one',
    done: stepValid(2),
  },
  { label: 'Creative', status: stepValid(3) ? 'Ready for review' : 'Headline needed', done: stepValid(3) },
  {
    label: 'Funding',
    status:
      funding.stage.value === 'funded' ? 'In escrow' : funding.stage.value === 'confirming' ? 'Confirming' : 'Not funded',
    done: funding.stage.value === 'funded',
  },
])

// Any edit after an order exists means the next approval must create a fresh order.
watch(form, () => funding.invalidate(), { deep: true })

// ---------- Draft persistence ----------

const savedAt = ref<string | null>(null)

function saveDraft() {
  localStorage.setItem(DRAFT_KEY, JSON.stringify(form))
  savedAt.value = utcTime(new Date().toISOString())
}

onMounted(() => {
  try {
    const raw = localStorage.getItem(DRAFT_KEY)
    if (!raw) return
    const saved = JSON.parse(raw) as Partial<typeof form>
    for (const key of Object.keys(form) as (keyof typeof form)[]) {
      if (saved[key] !== undefined && typeof saved[key] === typeof form[key]) (form as Record<string, unknown>)[key] = saved[key]
    }
    if (form.start < todayIso()) form.start = todayIso()
  } catch {
    localStorage.removeItem(DRAFT_KEY)
  }
})

function approve() {
  attempted[4] = true
  if (![1, 2, 3, 4].every(stepValid)) return goTo(4)
  funding.approve({
    name: form.name.trim(),
    brand: form.brand.trim(),
    url: form.url.trim(),
    start: form.start,
    durationDays: form.durationDays,
    placements: [...form.placements],
    incentivized: form.incentivized,
    frequencyCap: form.frequencyCap,
    headline: form.headline.trim(),
    cta: form.cta,
    budget: form.budget,
  })
}

watch(
  () => funding.stage.value,
  (stage) => {
    if (stage === 'funded') localStorage.removeItem(DRAFT_KEY)
  },
)

const needsWallet = computed(() => net.fundingMode === 'onchain' && !wallet.isConnected.value)
</script>

<template>
  <div>
    <div class="page-head">
      <div>
        <p class="crumbs">Campaigns <span aria-hidden="true">/</span> New campaign</p>
        <h1 class="page-title">New campaign</h1>
      </div>
      <button type="button" class="btn btn--secondary btn--sm" :disabled="locked" @click="saveDraft">
        {{ savedAt ? `Saved ${savedAt} UTC` : 'Save draft' }}
      </button>
    </div>

    <BrandCampaignStepper :steps="STEPS" :current="step" :done="stepDone" @go="goTo" />

    <div class="layout">
      <section class="panel">
        <!-- Step 1 -->
        <div v-if="step === 1">
          <p class="eyebrow">Step 1 of 4</p>
          <h2 class="panel__title">Campaign basics</h2>
          <p class="panel__lede">Name it, point it somewhere, and choose when it runs. All times are UTC.</p>
          <div class="fields">
            <div class="field">
              <label for="c-name" class="field-label">Campaign name</label>
              <input id="c-name" v-model="form.name" type="text" class="input" maxlength="80" :disabled="locked" :aria-invalid="!!fieldError(1, 'name')" />
              <p v-if="fieldError(1, 'name')" class="field-error">{{ fieldError(1, 'name') }}</p>
            </div>
            <div class="field">
              <label for="c-brand" class="field-label">Brand name</label>
              <input id="c-brand" v-model="form.brand" type="text" class="input" maxlength="60" :disabled="locked" :aria-invalid="!!fieldError(1, 'brand')" />
              <p v-if="fieldError(1, 'brand')" class="field-error">{{ fieldError(1, 'brand') }}</p>
            </div>
            <div class="field field--full">
              <label for="c-url" class="field-label">Landing page</label>
              <input id="c-url" v-model="form.url" type="url" class="input" inputmode="url" :disabled="locked" :aria-invalid="!!fieldError(1, 'url')" />
              <p v-if="fieldError(1, 'url')" class="field-error">{{ fieldError(1, 'url') }}</p>
            </div>
            <div class="field">
              <label for="c-start" class="field-label">Start date (UTC)</label>
              <input id="c-start" v-model="form.start" type="date" class="input" :min="todayIso()" :disabled="locked" :aria-invalid="!!fieldError(1, 'start')" />
              <p v-if="fieldError(1, 'start')" class="field-error">{{ fieldError(1, 'start') }}</p>
            </div>
            <div class="field">
              <span id="dur-label" class="field-label">Serving window</span>
              <div class="segmented" role="group" aria-labelledby="dur-label">
                <button
                  v-for="d in DURATIONS"
                  :key="d"
                  type="button"
                  :aria-pressed="form.durationDays === d"
                  :disabled="locked"
                  @click="form.durationDays = d"
                >
                  {{ d }} days
                </button>
              </div>
            </div>
          </div>
          <p class="callout">
            <AppIcon name="clock" />
            Runs {{ shortDate(form.start) }} 00:00 → {{ shortDate(endDate) }} 00:00 UTC ({{ form.durationDays * 24 }} h). The end
            time is exclusive.
          </p>
        </div>

        <!-- Step 2 -->
        <div v-else-if="step === 2">
          <p class="eyebrow">Step 2 of 4</p>
          <h2 class="panel__title">Where it shows</h2>
          <p class="panel__lede">Contextual placements only. No personal tracking, no cross-app profiles.</p>
          <fieldset class="options">
            <legend class="sr-only">Placements</legend>
            <label
              v-for="p in PLACEMENTS"
              :key="p.id"
              class="option"
              :class="{ 'option--on': form.placements.includes(p.id) }"
            >
              <input v-model="form.placements" type="checkbox" :value="p.id" :disabled="locked" />
              <span class="option__body">
                <span class="option__head">
                  <span class="option__name">{{ p.name }}</span>
                  <span class="mono option__meta">{{ p.format.toUpperCase() }}</span>
                </span>
                <span class="option__text">{{ p.description }}</span>
              </span>
            </label>
          </fieldset>
          <p v-if="fieldError(2, 'placements')" class="field-error">{{ fieldError(2, 'placements') }}</p>

          <div class="toggle-row">
            <div>
              <div id="incent-label" class="toggle-row__title">Reward-enabled audience</div>
              <p class="toggle-row__text">
                Players here earn a 30% share of the revenue your campaign generates. It's disclosed on every report.
              </p>
            </div>
            <button
              type="button"
              role="switch"
              class="switch"
              :aria-checked="form.incentivized"
              aria-labelledby="incent-label"
              :disabled="locked"
              @click="form.incentivized = !form.incentivized"
            >
              <span class="switch__track"><span class="switch__knob" /></span>
            </button>
          </div>

          <div class="field field--narrow">
            <label for="p-cap" class="field-label">Frequency cap</label>
            <select id="p-cap" v-model.number="form.frequencyCap" class="input" :disabled="locked">
              <option v-for="cap in FREQUENCY_CAPS" :key="cap" :value="cap">{{ cap }} per player per day</option>
            </select>
          </div>
        </div>

        <!-- Step 3 -->
        <div v-else-if="step === 3">
          <p class="eyebrow">Step 3 of 4</p>
          <h2 class="panel__title">Your sponsor card</h2>
          <p class="panel__lede">One static card. Every creative is reviewed before it serves.</p>
          <div class="creative">
            <div class="creative__form">
              <div>
                <label for="cr-file" class="field-label">Card image</label>
                <div class="dropzone">
                  <AppIcon name="upload" :size="24" />
                  <div class="dropzone__title">Drop a PNG or JPG</div>
                  <div class="dropzone__meta">1200 × 628 · up to 1 MB</div>
                  <input id="cr-file" type="file" accept="image/png,image/jpeg" :disabled="locked" />
                </div>
              </div>
              <div>
                <label for="cr-head" class="field-label field-label--row">
                  Headline <span class="mono">{{ form.headline.length }}/{{ HEADLINE_MAX }}</span>
                </label>
                <input
                  id="cr-head"
                  v-model="form.headline"
                  type="text"
                  class="input"
                  :maxlength="HEADLINE_MAX"
                  :disabled="locked"
                  :aria-invalid="!!fieldError(3, 'headline')"
                />
                <p v-if="fieldError(3, 'headline')" class="field-error">{{ fieldError(3, 'headline') }}</p>
              </div>
              <div>
                <label for="cr-cta" class="field-label">Button label</label>
                <select id="cr-cta" v-model="form.cta" class="input" :disabled="locked">
                  <option v-for="c in CTA_OPTIONS" :key="c" :value="c">{{ c }}</option>
                </select>
              </div>
            </div>
            <div class="creative__preview">
              <div class="field-label">Preview</div>
              <BrandSponsorCardPreview :brand="form.brand" :headline="form.headline" :cta="form.cta" />
              <p class="hint">Shown on home and results screens only. Paused or unapproved creatives never bill.</p>
            </div>
          </div>
        </div>

        <!-- Step 4 -->
        <div v-else>
          <p class="eyebrow">Step 4 of 4</p>
          <h2 class="panel__title">Budget &amp; funding</h2>
          <p class="panel__lede">
            Your budget goes into escrow. You're billed only for validated, viewable impressions, and anything unspent comes back
            to you.
          </p>

          <div class="budget">
            <label for="f-budget" class="field-label">Budget</label>
            <div class="budget__input" :class="{ 'budget__input--invalid': !!fieldError(4, 'budget') }">
              <input
                id="f-budget"
                v-model="form.budget"
                type="text"
                inputmode="decimal"
                class="mono"
                :disabled="locked"
                :aria-invalid="!!fieldError(4, 'budget')"
                @blur="attempted[4] = true"
              />
              <span>USD₮</span>
            </div>
            <div class="quick">
              <button
                v-for="v in [25, 50, 100, 250]"
                :key="v"
                type="button"
                class="mono"
                :aria-pressed="form.budget === String(v)"
                :disabled="locked"
                @click="form.budget = String(v)"
              >
                {{ v }} USD₮
              </button>
            </div>
            <p v-if="fieldError(4, 'budget')" class="field-error" role="alert">{{ fieldError(4, 'budget') }}</p>
          </div>

          <div class="tiles">
            <div class="tile">
              <div class="tile__label">Rate</div>
              <div class="tile__value mono">{{ formatToken(CPM_UNITS) }} USD₮</div>
              <div class="tile__note">per 1,000 validated impressions</div>
            </div>
            <div class="tile">
              <div class="tile__label">Max billed impressions</div>
              <div class="tile__value mono">{{ maxImps }}</div>
              <div class="tile__note">billing stops at your budget</div>
            </div>
            <div class="tile">
              <div class="tile__label">Delivery</div>
              <div class="tile__value">Not guaranteed</div>
              <div class="tile__note">we only sell traffic we've measured</div>
            </div>
          </div>

          <div class="split-box">
            <div class="split-box__head">
              <span>Earned revenue split</span><span class="mono">LOCKED TO THIS CAMPAIGN · V1</span>
            </div>
            <div class="split-box__bar" aria-hidden="true">
              <span style="flex: 30" class="seg seg--accent" /><span style="flex: 60" class="seg seg--ink" /><span
                style="flex: 10"
                class="seg seg--grey"
              />
            </div>
            <div class="split-box__legend"><span>30% players</span><span>60% operations</span><span>10% contingency</span></div>
          </div>

          <dl class="rows contract">
            <div>
              <dt>Network</dt>
              <dd class="mono">{{ net.chainName }} · {{ net.chain.id }}</dd>
            </div>
            <template v-if="net.fundingMode === 'onchain'">
              <div>
                <dt>Token</dt>
                <dd class="mono">
                  <a :href="net.addressUrl(net.tokenAddress!)" target="_blank" rel="noopener">
                    {{ net.tokenSymbol }} · {{ shortAddress(net.tokenAddress!) }}
                  </a>
                </dd>
              </div>
              <div>
                <dt>Escrow contract</dt>
                <dd class="mono">
                  <a :href="net.addressUrl(net.escrowAddress!)" target="_blank" rel="noopener">{{ shortAddress(net.escrowAddress!) }}</a>
                </dd>
              </div>
            </template>
            <div v-else>
              <dt>Mode</dt>
              <dd><span class="chip chip--soft">Demo · no transactions</span></dd>
            </div>
          </dl>

          <div class="fund">
            <ClientOnly>
              <div v-if="needsWallet" class="fund__row">
                <div>
                  <div class="fund__title">Connect a wallet to fund</div>
                  <p class="fund__text">The wallet that pays is the one refunds go back to.</p>
                </div>
                <WalletChip />
              </div>

              <div v-else-if="net.fundingMode === 'onchain' && wallet.wrongChain.value" class="fund__row">
                <div>
                  <div class="fund__title">Wrong network</div>
                  <p class="fund__text">Switch your wallet to {{ net.chainName }} to continue.</p>
                </div>
                <button type="button" class="btn btn--primary" @click="wallet.switchChain().catch(() => {})">
                  Switch to {{ net.chainName }}
                </button>
              </div>

              <div v-else-if="['idle', 'creating', 'approving'].includes(funding.stage.value)" class="fund__row">
                <div>
                  <div class="fund__title">1 of 2 · Approve USD₮</div>
                  <p class="fund__text">Lets the escrow contract pull exactly {{ budgetLabel }} USD₮. Nothing more.</p>
                </div>
                <button type="button" class="btn btn--primary" :disabled="funding.busy.value" @click="approve">
                  <span v-if="funding.busy.value" class="spinner spinner--light" aria-hidden="true" />
                  {{ funding.stage.value === 'idle' ? `Approve ${budgetLabel} USD₮` : 'Check your wallet…' }}
                </button>
              </div>

              <div v-else-if="['approved', 'depositing'].includes(funding.stage.value)" class="fund__row">
                <div>
                  <div class="fund__done"><CheckDot :size="18" />Approved {{ budgetLabel }} USD₮</div>
                  <div class="fund__title">2 of 2 · Deposit into escrow</div>
                </div>
                <button type="button" class="btn btn--primary" :disabled="funding.busy.value" @click="funding.deposit()">
                  <span v-if="funding.busy.value" class="spinner spinner--light" aria-hidden="true" />
                  {{ funding.busy.value ? 'Check your wallet…' : `Deposit ${budgetLabel} USD₮` }}
                </button>
              </div>

              <div v-else-if="funding.stage.value === 'confirming'" class="fund__row" role="status">
                <div class="fund__pending">
                  <span class="spinner" aria-hidden="true" />
                  <div>
                    <div class="fund__title">Confirming on {{ net.chainName }}…</div>
                    <p class="fund__text mono">
                      <a v-if="funding.txHash.value" :href="net.txUrl(funding.txHash.value)" target="_blank" rel="noopener">
                        tx {{ shortHash(funding.txHash.value) }}
                      </a>
                      <template v-else>Demo · waiting for finality</template>
                    </p>
                  </div>
                </div>
              </div>

              <div v-else class="fund__row" role="status">
                <div class="fund__success">
                  <CheckDot :size="36" />
                  <div>
                    <div class="fund__title">Campaign funded</div>
                    <p class="fund__text">{{ budgetLabel }} USD₮ is in escrow. {{ form.name }} goes live after creative review.</p>
                    <p v-if="funding.txHash.value" class="fund__text mono">
                      <a :href="net.txUrl(funding.txHash.value)" target="_blank" rel="noopener">tx {{ shortHash(funding.txHash.value) }}</a>
                      · confirmed
                    </p>
                  </div>
                </div>
                <NuxtLink :to="`/brands/campaigns/${funding.order.value?.id}`" class="btn btn--primary">
                  View campaign <AppIcon name="arrow-right" :size="15" :stroke-width="2" />
                </NuxtLink>
              </div>

              <p v-if="funding.error.value" class="field-error" role="alert">{{ funding.error.value }}</p>
              <template #fallback><div class="fund__row"><span class="spinner" /></div></template>
            </ClientOnly>
          </div>
        </div>

        <div class="panel__nav">
          <button v-if="step > 1" type="button" class="btn btn--secondary" @click="step--">Back</button>
          <span v-else />
          <button v-if="step < STEPS.length" type="button" class="btn btn--primary" @click="next">
            {{ step === 3 ? 'Continue to funding' : 'Continue' }}
            <AppIcon name="arrow-right" :size="15" :stroke-width="2" />
          </button>
        </div>
      </section>

      <aside class="summary" aria-label="Campaign summary">
        <div class="summary__label">Summary</div>
        <div class="summary__name">{{ form.name || 'Untitled campaign' }}</div>
        <dl class="rows">
          <div>
            <dt>Window</dt>
            <dd>{{ shortDate(form.start) }} → {{ shortDate(endDate, true) }}</dd>
          </div>
          <div>
            <dt>Placements</dt>
            <dd>{{ placementsLabel }}</dd>
          </div>
          <div>
            <dt>Audience</dt>
            <dd>{{ form.incentivized ? 'Reward-enabled' : 'Standard' }}</dd>
          </div>
          <div>
            <dt>Frequency cap</dt>
            <dd>{{ form.frequencyCap }} / player / day</dd>
          </div>
          <div>
            <dt>Rate</dt>
            <dd class="mono">{{ formatToken(CPM_UNITS) }} CPM</dd>
          </div>
          <div>
            <dt>Max impressions</dt>
            <dd class="mono">{{ maxImps }}</dd>
          </div>
        </dl>
        <div class="summary__due">
          <div class="summary__due-label">Into escrow</div>
          <div class="summary__due-value mono">{{ budgetLabel }} <span>USD₮</span></div>
          <p class="summary__due-note">Billed on validated delivery. Unspent budget is refunded after the adjustment window.</p>
        </div>
        <ul class="checklist">
          <li v-for="c in checklist" :key="c.label">
            <CheckDot v-if="c.done" />
            <span v-else class="checklist__todo" aria-hidden="true" />
            <span class="checklist__label">{{ c.label }}</span>
            <span class="checklist__status">{{ c.status }}</span>
          </li>
        </ul>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
}

.crumbs {
  font-size: 14px;
  color: var(--muted);
}

.page-title {
  margin-top: 8px;
  font-size: clamp(28px, 3vw, 38px);
  letter-spacing: -0.035em;
  font-weight: 600;
}

.layout {
  margin-top: 24px;
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 20px;
}

.panel {
  flex: 3 1 520px;
  min-width: 0;
  padding: clamp(20px, 3vw, 32px);
  border: 1px solid var(--line);
  border-radius: 22px;
  background: #fff;
}

.panel__title {
  margin-top: 8px;
  font-size: 26px;
  letter-spacing: -0.03em;
  font-weight: 600;
}

.panel__lede {
  margin-top: 8px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--muted);
}

.panel__nav {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid var(--line-soft);
}

.fields {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-top: 28px;
}

.field {
  flex: 1 1 260px;
  min-width: 0;
}

.field--full {
  flex-basis: 100%;
}

.field--narrow {
  max-width: 320px;
  margin-top: 20px;
}

.field-label--row {
  display: flex;
  justify-content: space-between;
}

.field-label--row .mono {
  font-weight: 400;
  color: var(--muted);
}

.segmented {
  display: flex;
  gap: 8px;
}

.segmented button {
  flex: 1;
  min-height: 48px;
  border-radius: 12px;
  border: 1px solid var(--line-strong);
  background: #fff;
  font-size: 14px;
  font-weight: 500;
}

.segmented button[aria-pressed='true'] {
  border-color: var(--ink);
  background: var(--ink);
  color: #fff;
}

.callout {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
  padding: 14px 16px;
  border-radius: 12px;
  background: var(--surface);
  font-size: 14px;
  line-height: 1.45;
  color: var(--ink-2);
}

.options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 28px 0 0;
  padding: 0;
  border: 0;
}

.option {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 18px;
  border: 1.5px solid var(--line);
  border-radius: 16px;
  cursor: pointer;
}

.option--on {
  border-color: var(--ink);
}

.option input {
  flex: none;
  width: 20px;
  height: 20px;
  margin: 2px 0 0;
  accent-color: var(--ink);
}

.option__body {
  flex: 1;
  min-width: 0;
}

.option__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 8px 12px;
}

.option__name {
  font-size: 16px;
  font-weight: 600;
}

.option__meta {
  font-size: 12px;
  color: var(--muted);
}

.option__text {
  display: block;
  margin-top: 6px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--muted);
}

.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 16px;
  padding: 18px;
  border-radius: 16px;
  background: var(--surface);
}

.toggle-row__title {
  font-size: 16px;
  font-weight: 600;
}

.toggle-row__text {
  margin-top: 6px;
  max-width: 520px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--muted);
}

.switch {
  flex: none;
  display: flex;
  align-items: center;
  width: 60px;
  height: 44px;
  padding: 0;
  border: 0;
  background: none;
}

.switch__track {
  display: flex;
  justify-content: flex-start;
  width: 52px;
  height: 30px;
  padding: 3px;
  border-radius: 999px;
  background: #cfcfca;
  transition: background 0.2s ease;
}

.switch__knob {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 3px rgb(11 11 12 / 0.3);
  transition: transform 0.2s ease;
}

.switch[aria-checked='true'] .switch__track {
  background: var(--ink);
}

.switch[aria-checked='true'] .switch__knob {
  transform: translateX(22px);
}

.creative {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-top: 28px;
}

.creative__form {
  flex: 1 1 300px;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.creative__preview {
  flex: 1 1 280px;
  min-width: 0;
}

.dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
  border: 1.5px dashed #cfcfca;
  border-radius: 16px;
  background: var(--bg-app);
  text-align: center;
}

.dropzone__title {
  margin-top: 10px;
  font-size: 15px;
  font-weight: 500;
}

.dropzone__meta {
  margin-top: 4px;
  font-size: 13px;
  color: var(--muted);
}

.dropzone input {
  margin-top: 14px;
  max-width: 100%;
  font-size: 13px;
}

.hint {
  margin-top: 12px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--muted);
}

.budget {
  margin-top: 28px;
}

.budget__input {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 68px;
  padding: 0 18px;
  border: 1px solid var(--line-strong);
  border-radius: 14px;
  background: #fff;
}

.budget__input--invalid {
  border-color: var(--danger);
}

.budget__input input {
  flex: 1;
  min-width: 0;
  padding: 4px 0;
  border: 0;
  background: transparent;
  font-size: 30px;
  letter-spacing: -0.02em;
}

.budget__input span {
  font-size: 15px;
  font-weight: 500;
  color: var(--muted);
}

.quick {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.quick button {
  min-height: 40px;
  padding: 0 14px;
  border-radius: 999px;
  border: 1px solid var(--line-strong);
  background: #fff;
  font-size: 13px;
}

.quick button[aria-pressed='true'] {
  border-color: var(--ink);
  background: var(--ink);
  color: #fff;
}

.tiles {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 24px;
}

.tile {
  flex: 1 1 180px;
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 14px;
}

.tile__label,
.tile__note {
  font-size: 12px;
  color: var(--muted);
}

.tile__value {
  margin-top: 6px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.tile__value.mono {
  font-weight: 400;
}

.tile__note {
  margin-top: 4px;
}

.split-box {
  margin-top: 24px;
  padding: 18px;
  border-radius: 16px;
  background: var(--surface);
}

.split-box__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
}

.split-box__head .mono {
  font-size: 12px;
  font-weight: 400;
  color: var(--muted);
}

.split-box__bar {
  display: flex;
  gap: 3px;
  height: 10px;
  margin-top: 14px;
}

.seg {
  border-radius: 999px;
}

.seg--accent {
  background: var(--accent);
  box-shadow: inset 0 0 0 1px rgb(11 11 12 / 0.2);
}

.seg--ink {
  background: var(--ink);
}

.seg--grey {
  background: var(--line-dashed);
}

.split-box__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 20px;
  margin-top: 12px;
  font-size: 13px;
  color: var(--ink-2);
}

.contract {
  margin-top: 20px;
}

.contract a {
  text-decoration: none;
}

.contract a:hover {
  text-decoration: underline;
}

.fund {
  margin-top: 16px;
  padding: 20px;
  border: 1.5px solid var(--ink);
  border-radius: 18px;
}

.fund__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  min-height: 50px;
}

.fund__title {
  font-size: 15px;
  font-weight: 600;
}

.fund__text {
  margin-top: 4px;
  font-size: 13px;
  line-height: 1.45;
  color: var(--muted);
}

.fund__done {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-size: 13px;
  color: var(--ink-2);
}

.fund__pending,
.fund__success {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.fund__success {
  align-items: flex-start;
}

.spinner--light {
  width: 16px;
  height: 16px;
  border-width: 2px;
  border-color: rgb(255 255 255 / 0.3);
  border-top-color: #fff;
}

.summary {
  flex: 1 1 300px;
  min-width: 0;
  padding: 24px;
  border: 1px solid var(--line);
  border-radius: 22px;
  background: #fff;
}

.summary__label {
  font-size: 13px;
  font-weight: 500;
  color: var(--muted);
}

.summary__name {
  margin: 6px 0 18px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.02em;
  overflow-wrap: anywhere;
}

.summary__due {
  margin-top: 8px;
  padding: 18px;
  border-radius: 16px;
  background: var(--ink);
  color: #fff;
}

.summary__due-label,
.summary__due-note,
.summary__due-value span {
  font-size: 12px;
  color: var(--on-ink-muted);
}

.summary__due-value {
  margin-top: 6px;
  font-size: 30px;
  letter-spacing: -0.03em;
}

.summary__due-value span {
  font-size: 14px;
}

.summary__due-note {
  margin-top: 8px;
  line-height: 1.45;
}

.checklist {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 20px 0 0;
  padding: 0;
}

.checklist li {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
}

.checklist__todo {
  flex: none;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 1.5px dashed var(--line-dashed);
}

.checklist__label {
  flex: 1;
}

.checklist__status {
  font-size: 12px;
  color: var(--muted);
}
</style>
