<script setup lang="ts">
import { erc20Abi } from 'viem'
import type { WalletOption } from '~/composables/useWallet'

const props = withDefaults(defineProps<{ showBalance?: boolean }>(), { showBalance: false })

const wallet = useWallet()
const { tokenAddress, tokenSymbol, tokenDecimals, publicClient } = useChain()

const dialog = ref<HTMLDialogElement | null>(null)
const menuOpen = ref(false)
const busy = ref<string | null>(null)
const error = ref<string | null>(null)
const balance = ref<string | null>(null)

async function loadBalance() {
  balance.value = null
  if (!props.showBalance || !tokenAddress || !wallet.address.value) return
  try {
    const units = await publicClient().readContract({
      address: tokenAddress,
      abi: erc20Abi,
      functionName: 'balanceOf',
      args: [wallet.address.value],
    })
    balance.value = formatToken(units, 2, tokenDecimals)
  } catch {
    balance.value = null
  }
}

watch(() => wallet.address.value, loadBalance, { immediate: true })
defineExpose({ refreshBalance: loadBalance })

function open() {
  error.value = null
  dialog.value?.showModal()
}

async function choose(option: WalletOption) {
  busy.value = option.id
  error.value = null
  try {
    await wallet.connect(option)
    dialog.value?.close()
  } catch (err) {
    error.value = readableError(err)
  } finally {
    busy.value = null
  }
}

async function copy() {
  if (wallet.address.value) await navigator.clipboard?.writeText(wallet.address.value)
  menuOpen.value = false
}

function disconnect() {
  wallet.disconnect()
  menuOpen.value = false
}
</script>

<template>
  <div class="wallet">
    <template v-if="wallet.address.value">
      <button
        type="button"
        class="wallet__chip"
        :aria-expanded="menuOpen"
        :aria-label="`Connected wallet ${wallet.address.value}`"
        @click="menuOpen = !menuOpen"
      >
        <span v-if="balance" class="mono wallet__balance">{{ balance }} {{ tokenSymbol }}</span>
        <span class="wallet__address mono"><AddressAvatar :size="20" />{{ shortAddress(wallet.address.value) }}</span>
      </button>
      <div v-if="menuOpen" class="wallet__menu" role="menu">
        <button type="button" role="menuitem" @click="copy"><AppIcon name="copy" :size="16" />Copy address</button>
        <button type="button" role="menuitem" @click="disconnect"><AppIcon name="logout" :size="16" />Disconnect</button>
      </div>
    </template>
    <button v-else type="button" class="btn btn--primary btn--sm" @click="open">Connect wallet</button>

    <dialog ref="dialog" class="wallet__dialog" aria-labelledby="wallet-dialog-title" @click.self="dialog?.close()">
      <div class="wallet__dialog-body">
        <div class="wallet__dialog-head">
          <h2 id="wallet-dialog-title">Connect a wallet</h2>
          <button type="button" class="wallet__close" aria-label="Close" @click="dialog?.close()">
            <AppIcon name="plus" :size="18" :stroke-width="2" style="transform: rotate(45deg)" />
          </button>
        </div>
        <p class="wallet__dialog-text">Connecting only shares your address. Nothing is signed or sent yet.</p>
        <WalletPicker :busy="busy" @select="choose" />
        <p v-if="error" class="field-error" role="alert">{{ error }}</p>
      </div>
    </dialog>
  </div>
</template>

<style scoped>
.wallet {
  position: relative;
}

.wallet__chip {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  padding: 0 6px 0 6px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #fff;
}

.wallet__balance {
  padding-left: 8px;
  font-size: 13px;
}

.wallet__address {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 32px;
  padding: 0 12px 0 6px;
  border-radius: 999px;
  background: var(--surface);
  font-size: 13px;
}

.wallet__menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  z-index: 20;
  min-width: 200px;
  padding: 6px;
  border-radius: 14px;
  border: 1px solid var(--line);
  background: #fff;
  box-shadow: var(--shadow-float);
}

.wallet__menu button {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-height: 44px;
  padding: 0 12px;
  border: 0;
  border-radius: 10px;
  background: none;
  font-size: 14px;
  text-align: left;
}

.wallet__menu button:hover {
  background: var(--surface);
}

.wallet__dialog {
  width: min(440px, calc(100vw - 32px));
  padding: 0;
  border: 1px solid var(--line);
  border-radius: 24px;
  box-shadow: 0 40px 80px -40px rgb(11 11 12 / 0.45);
}

.wallet__dialog::backdrop {
  background: rgb(11 11 12 / 0.35);
}

.wallet__dialog-body {
  padding: 28px;
}

.wallet__dialog-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.wallet__dialog-head h2 {
  font-size: 22px;
  letter-spacing: -0.02em;
  font-weight: 600;
}

.wallet__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  margin-right: -10px;
  border: 0;
  border-radius: 50%;
  background: none;
}

.wallet__close:hover {
  background: var(--surface);
}

.wallet__dialog-text {
  margin: 6px 0 20px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--muted);
}
</style>
