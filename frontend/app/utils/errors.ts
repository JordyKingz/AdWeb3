/** Turns wallet, viem and $fetch errors into one readable sentence. */
export function readableError(err: unknown): string {
  const e = err as {
    code?: number
    cause?: { code?: number }
    shortMessage?: string
    statusMessage?: string
    data?: { statusMessage?: string; message?: string }
    message?: string
  }
  if (e?.code === 4001 || e?.cause?.code === 4001 || /user rejected|denied/i.test(e?.shortMessage ?? e?.message ?? '')) {
    return 'You declined the request in your wallet.'
  }
  return e?.data?.statusMessage ?? e?.statusMessage ?? e?.shortMessage ?? e?.message ?? 'Something went wrong. Please try again.'
}

export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
