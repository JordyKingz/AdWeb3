<script setup lang="ts">
import { createSiweMessage } from 'viem/siwe'
import type { SessionInfo } from '#shared/types/api'
import type { WalletOption } from '~/composables/useWallet'

definePageMeta({ layout: false })
useSeoMeta({ title: 'Connect your wallet · AdWeb3' })

type Stage = 'choose' | 'connecting' | 'sign' | 'signing' | 'done'

const STATEMENT = 'Link this wallet to your AdWeb3 player account. This is not a transaction and costs no gas.'

const route = useRoute()
const wallet = useWallet()
const net = useChain()

const { data: existing } = await useFetch<SessionInfo>('/api/auth/session', { key: 'player-session' })

const stage = ref<Stage>(existing.value && !route.query.change ? 'done' : 'choose')
const busyId = ref<string | null>(null)
const error = ref<string | null>(null)
const message = ref('')
const expiresAt = ref('')
const linked = ref<SessionInfo | null>(existing.value ?? null)

const STEP_LABELS: Record<Stage, string> = {
  choose: 'Step 1 of 3 · Choose wallet',
  connecting: 'Step 1 of 3 · Connecting',
  sign: 'Step 2 of 3 · Sign message',
  signing: 'Step 2 of 3 · Signing',
  done: 'Step 3 of 3 · Done',
}
const progress = computed(() => (stage.value === 'done' ? 3 : stage.value === 'sign' || stage.value === 'signing' ? 2 : 1))

async function prepareMessage() {
  const { nonce, issuedAt, expirationTime } = await $fetch<{ nonce: string; issuedAt: string; expirationTime: string }>(
    '/api/auth/nonce',
  )
  message.value = createSiweMessage({
    address: wallet.address.value!,
    chainId: net.chain.id,
    domain: window.location.host,
    nonce,
    uri: `${window.location.origin}/connect`,
    version: '1',
    statement: STATEMENT,
    issuedAt: new Date(issuedAt),
    expirationTime: new Date(expirationTime),
  })
  expiresAt.value = expirationTime
  stage.value = 'sign'
}

async function choose(option: WalletOption) {
  error.value = null
  busyId.value = option.id
  stage.value = 'connecting'
  try {
    await wallet.connect(option)
    await prepareMessage()
  } catch (err) {
    stage.value = 'choose'
    error.value = readableError(err)
  } finally {
    busyId.value = null
  }
}

async function continueWithConnected() {
  error.value = null
  stage.value = 'connecting'
  try {
    await prepareMessage()
  } catch (err) {
    stage.value = 'choose'
    error.value = readableError(err)
  }
}

async function sign() {
  error.value = null
  stage.value = 'signing'
  try {
    const signature = await wallet.signMessage(message.value)
    linked.value = await $fetch<SessionInfo>('/api/auth/verify', {
      method: 'POST',
      body: { message: message.value, signature },
    })
    await refreshNuxtData('player-session')
    stage.value = 'done'
  } catch (err) {
    error.value = readableError(err)
    // A failed verification burns the nonce, so prepare a fresh message.
    stage.value = 'choose'
    if (wallet.isConnected.value) await prepareMessage().catch(() => {})
  }
}

function useDifferentWallet() {
  wallet.disconnect()
  error.value = null
  message.value = ''
  stage.value = 'choose'
}
</script>

<template>
  <div class="connect">
    <aside class="connect__aside">
      <AppLogo inverted />
      <div>
        <h1 class="connect__title">Connect once. Get paid automatically.</h1>
        <ul class="connect__points">
          <li>
            <span class="connect__icon"><AppIcon name="pen" :size="16" :stroke-width="2" /></span>
            <div>
              <div class="connect__point-title">A signature, not a transaction</div>
              <p class="connect__point-text">No gas, no token approvals, nothing leaves your wallet.</p>
            </div>
          </li>
          <li>
            <span class="connect__icon"><AppIcon name="lock" :size="16" :stroke-width="2" /></span>
            <div>
              <div class="connect__point-title">Your keys stay yours</div>
              <p class="connect__point-text">We will never ask for your seed phrase or private key.</p>
            </div>
          </li>
          <li>
            <span class="connect__icon"><AppIcon name="clock" :size="16" :stroke-width="2" /></span>
            <div>
              <div class="connect__point-title">Single use, expires fast</div>
              <p class="connect__point-text">Each link request works once and expires in five minutes.</p>
            </div>
          </li>
        </ul>
      </div>
      <p class="connect__foot">
        The sign-in request is tied to this browser session, so nobody else can complete it with a copied link.
      </p>
    </aside>

    <main class="connect__main">
      <div class="connect__card">
        <div class="connect__progress" aria-hidden="true">
          <span v-for="n in 3" :key="n" :class="{ on: n <= progress }" />
        </div>
        <p class="eyebrow connect__step">{{ STEP_LABELS[stage] }}</p>

        <ClientOnly>
          <!-- Choose -->
          <div v-if="stage === 'choose'">
            <h2 class="connect__h2">Connect a wallet</h2>
            <p class="connect__lede">Pick the wallet you want rewards paid to. You can switch later with a new signature.</p>
            <button v-if="wallet.address.value" type="button" class="connect__resume" @click="continueWithConnected">
              <AddressAvatar :size="28" />
              <span>
                <span class="connect__resume-title">Continue with {{ shortAddress(wallet.address.value) }}</span>
                <span class="connect__resume-sub">{{ wallet.walletName.value }}</span>
              </span>
              <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
            </button>
            <div class="connect__picker">
              <WalletPicker :busy="busyId" @select="choose" />
            </div>
            <p class="connect__context"><AppIcon name="phone" :size="16" />Linking to Candy Blocks on this device</p>
          </div>

          <!-- Connecting / signing -->
          <div v-else-if="stage === 'connecting' || stage === 'signing'" class="connect__waiting" role="status">
            <span class="spinner connect__spinner" aria-hidden="true" />
            <h2 class="connect__h2">{{ stage === 'connecting' ? 'Opening your wallet…' : 'Waiting for your signature…' }}</h2>
            <p class="connect__lede">
              {{ stage === 'connecting' ? 'Approve the connection request in your wallet.' : 'Check your wallet and confirm the message.' }}
            </p>
          </div>

          <!-- Sign -->
          <div v-else-if="stage === 'sign'">
            <h2 class="connect__h2">Sign to prove it's yours</h2>
            <div class="connect__who mono">
              <AddressAvatar />{{ shortAddress(wallet.address.value ?? '') }} · {{ wallet.walletName.value }}
            </div>
            <div class="connect__message">
              <div class="connect__message-head">Message you'll sign</div>
              <pre class="mono">{{ message }}</pre>
            </div>
            <ul class="connect__ticks">
              <li><AppIcon name="check" :size="14" :stroke-width="3" />No gas</li>
              <li><AppIcon name="check" :size="14" :stroke-width="3" />No token approval</li>
              <li><AppIcon name="check" :size="14" :stroke-width="3" />Expires {{ expiresAt ? utcTime(expiresAt) : '' }} UTC</li>
            </ul>
            <button type="button" class="btn btn--primary btn--lg btn--block connect__sign" @click="sign">Sign message</button>
            <button type="button" class="btn btn--ghost btn--block" @click="useDifferentWallet">Use a different wallet</button>
          </div>

          <!-- Done -->
          <div v-else class="connect__done" role="status">
            <CheckDot :size="56" />
            <h2 class="connect__h2">Wallet linked</h2>
            <p class="connect__lede">
              <span class="mono connect__addr">{{ shortAddress(linked?.address ?? '') }}</span> is your payout wallet. Finalized
              rewards are sent here automatically once you pass the payout threshold.
            </p>
            <dl class="rows connect__rows">
              <div><dt>Game</dt><dd>Candy Blocks · Android</dd></div>
              <div><dt>Network</dt><dd>{{ net.chainName }}</dd></div>
              <div><dt>Paid in</dt><dd>{{ net.tokenSymbol }}</dd></div>
            </dl>
            <NuxtLink to="/dashboard" class="btn btn--primary btn--lg btn--block">
              Open your dashboard <AppIcon name="arrow-right" :size="16" :stroke-width="2" />
            </NuxtLink>
            <button type="button" class="btn btn--ghost btn--block" @click="useDifferentWallet">Link a different wallet</button>
          </div>

          <p v-if="error" class="field-error connect__error" role="alert">{{ error }}</p>

          <template #fallback>
            <div class="connect__waiting"><span class="spinner connect__spinner" /></div>
          </template>
        </ClientOnly>
      </div>
    </main>
  </div>
</template>

<style scoped>
.connect {
  min-height: 100vh;
  display: flex;
  flex-wrap: wrap;
  background: var(--bg-app);
}

.connect__aside {
  flex: 1 1 420px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 48px;
  padding: clamp(28px, 5vw, 64px);
  background: var(--ink);
  color: #fff;
}

.connect__aside > :first-child {
  align-self: flex-start;
}

.connect__title {
  max-width: 460px;
  font-size: clamp(36px, 4vw, 56px);
  line-height: 1.02;
  letter-spacing: -0.045em;
  font-weight: 600;
}

.connect__points {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 22px;
  max-width: 420px;
  margin: 40px 0 0;
  padding: 0;
}

.connect__points li {
  display: flex;
  gap: 14px;
}

.connect__icon {
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent);
  color: var(--ink);
}

.connect__point-title {
  font-size: 16px;
  font-weight: 600;
}

.connect__point-text,
.connect__foot {
  margin-top: 4px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--on-ink-muted);
}

.connect__foot {
  max-width: 380px;
  font-size: 13px;
}

.connect__main {
  flex: 1.4 1 520px;
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: clamp(24px, 5vw, 64px) clamp(16px, 5vw, 64px);
}

.connect__card {
  width: 100%;
  max-width: 480px;
  padding: clamp(24px, 4vw, 36px);
  border: 1px solid var(--line);
  border-radius: 28px;
  background: #fff;
  box-shadow: 0 40px 80px -50px rgb(11 11 12 / 0.35);
}

.connect__progress {
  display: flex;
  gap: 6px;
}

.connect__progress span {
  flex: 1;
  height: 4px;
  border-radius: 999px;
  background: var(--line);
  transition: background 0.3s ease;
}

.connect__progress span.on {
  background: var(--ink);
}

.connect__step {
  margin-top: 12px;
}

.connect__h2 {
  margin-top: 10px;
  font-size: 28px;
  letter-spacing: -0.03em;
  font-weight: 600;
}

.connect__lede {
  margin-top: 8px;
  font-size: 15px;
  line-height: 1.55;
  color: var(--muted);
}

.connect__resume {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  min-height: 64px;
  margin-top: 24px;
  padding: 0 16px;
  border: 1.5px solid var(--ink);
  border-radius: 16px;
  background: #fff;
  text-align: left;
}

.connect__resume > span {
  flex: 1;
  min-width: 0;
}

.connect__resume-title {
  display: block;
  font-size: 15px;
  font-weight: 600;
}

.connect__resume-sub {
  display: block;
  font-size: 13px;
  color: var(--muted);
}

.connect__picker {
  margin-top: 24px;
}

.connect__resume + .connect__picker {
  margin-top: 10px;
}

.connect__context {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
  padding: 12px 14px;
  border-radius: 12px;
  background: var(--surface);
  font-size: 13px;
  color: var(--ink-2);
}

.connect__waiting {
  padding: 40px 0 16px;
  text-align: center;
}

.connect__spinner {
  width: 40px;
  height: 40px;
  border-width: 3px;
}

.connect__who {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 36px;
  margin-top: 14px;
  padding: 0 12px 0 6px;
  border-radius: 999px;
  background: var(--surface);
  font-size: 13px;
}

.connect__message {
  margin-top: 18px;
  border: 1px solid var(--line);
  border-radius: 16px;
  overflow: hidden;
}

.connect__message-head {
  padding: 10px 14px;
  border-bottom: 1px solid var(--line);
  background: var(--bg-app);
  font-size: 12px;
  color: var(--muted);
}

.connect__message pre {
  margin: 0;
  padding: 14px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--ink-2);
  white-space: pre-wrap;
  word-break: break-all;
}

.connect__ticks {
  list-style: none;
  display: flex;
  flex-wrap: wrap;
  gap: 8px 18px;
  margin: 16px 0 0;
  padding: 0;
  font-size: 13px;
  color: var(--ink-2);
}

.connect__ticks li {
  display: flex;
  align-items: center;
  gap: 6px;
}

.connect__sign {
  margin-top: 24px;
  margin-bottom: 8px;
}

.connect__done {
  padding-top: 12px;
}

.connect__addr {
  color: var(--ink);
}

.connect__rows {
  margin: 20px 0;
}

.connect__done .btn--ghost {
  margin-top: 8px;
}

.connect__error {
  margin-top: 16px;
}
</style>
