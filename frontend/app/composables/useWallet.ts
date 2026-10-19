import { createWalletClient, custom, getAddress, numberToHex, type Address, type EIP1193Provider } from 'viem'

export interface WalletOption {
  id: string
  name: string
  icon: string
  rdns: string
  provider: EIP1193Provider
}

const STORAGE_KEY = 'adweb3:wallet'

// Browser-only state: providers are live objects and never exist on the server.
const options = shallowRef<WalletOption[]>([])
const listening = new WeakSet<EIP1193Provider>()
let discovered: Promise<void> | null = null
let active: WalletOption | null = null

/** Finds installed wallets via EIP-6963, falling back to a legacy injected provider. */
function discover(): Promise<void> {
  if (!import.meta.client) return Promise.resolve()
  discovered ??= new Promise((resolve) => {
    window.addEventListener('eip6963:announceProvider', (event: Event) => {
      const { info, provider } = (event as CustomEvent).detail ?? {}
      if (!info?.uuid || !provider || options.value.some((o) => o.rdns === info.rdns)) return
      options.value = [...options.value, { id: info.uuid, name: info.name, icon: info.icon, rdns: info.rdns, provider }]
    })
    window.dispatchEvent(new Event('eip6963:requestProvider'))
    setTimeout(() => {
      const injected = (window as unknown as { ethereum?: EIP1193Provider }).ethereum
      if (!options.value.length && injected) {
        options.value = [{ id: 'injected', name: 'Browser wallet', icon: '', rdns: 'injected', provider: injected }]
      }
      resolve()
    }, 350)
  })
  return discovered
}

export function useWallet() {
  const address = useState<Address | null>('wallet:address', () => null)
  const chainId = useState<number | null>('wallet:chainId', () => null)
  const walletName = useState<string | null>('wallet:name', () => null)
  const { chain } = useChain()

  const isConnected = computed(() => !!address.value)
  const wrongChain = computed(() => isConnected.value && chainId.value !== chain.id)

  function forget() {
    address.value = null
    chainId.value = null
    walletName.value = null
    active = null
    localStorage.removeItem(STORAGE_KEY)
  }

  function bind(option: WalletOption) {
    active = option
    walletName.value = option.name
    if (listening.has(option.provider)) return
    listening.add(option.provider)
    option.provider.on('accountsChanged', (accounts) => {
      if (active?.provider !== option.provider) return
      if (accounts[0]) address.value = getAddress(accounts[0])
      else forget()
    })
    option.provider.on('chainChanged', (id) => {
      if (active?.provider === option.provider) chainId.value = Number(id)
    })
  }

  async function connect(option: WalletOption): Promise<Address> {
    const accounts = await option.provider.request({ method: 'eth_requestAccounts' })
    if (!accounts[0]) throw new Error('The wallet did not share an account.')
    bind(option)
    address.value = getAddress(accounts[0])
    chainId.value = Number(await option.provider.request({ method: 'eth_chainId' }))
    localStorage.setItem(STORAGE_KEY, option.rdns)
    return address.value
  }

  /** Restores the last wallet without prompting, if it still authorizes this site. */
  async function reconnect() {
    const rdns = localStorage.getItem(STORAGE_KEY)
    if (!rdns || address.value) return
    await discover()
    const option = options.value.find((o) => o.rdns === rdns)
    if (!option) return
    const accounts = await option.provider.request({ method: 'eth_accounts' }).catch(() => [] as Address[])
    if (!accounts[0]) return
    bind(option)
    address.value = getAddress(accounts[0])
    chainId.value = Number(await option.provider.request({ method: 'eth_chainId' }))
  }

  function provider(): EIP1193Provider {
    if (!active) throw new Error('Connect a wallet first.')
    return active.provider
  }

  async function switchChain() {
    const hexId = numberToHex(chain.id)
    try {
      await provider().request({ method: 'wallet_switchEthereumChain', params: [{ chainId: hexId }] })
    } catch (err) {
      if ((err as { code?: number }).code !== 4902) throw err
      await provider().request({
        method: 'wallet_addEthereumChain',
        params: [
          {
            chainId: hexId,
            chainName: chain.name,
            nativeCurrency: chain.nativeCurrency,
            rpcUrls: [...chain.rpcUrls.default.http],
            blockExplorerUrls: chain.blockExplorers ? [chain.blockExplorers.default.url] : undefined,
          },
        ],
      })
    }
    chainId.value = chain.id
  }

  function walletClient() {
    if (!address.value) throw new Error('Connect a wallet first.')
    return createWalletClient({ account: address.value, chain, transport: custom(provider()) })
  }

  async function signMessage(message: string) {
    return walletClient().signMessage({ account: address.value!, message })
  }

  return {
    options: readonly(options),
    address: readonly(address),
    chainId: readonly(chainId),
    walletName: readonly(walletName),
    isConnected,
    wrongChain,
    discover,
    connect,
    reconnect,
    switchChain,
    walletClient,
    signMessage,
    disconnect: forget,
  }
}
