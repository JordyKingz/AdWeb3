import type { BrandOverview } from '#shared/types/api'

export default defineEventHandler(async (): Promise<BrandOverview> => {
  const reports = (await listCampaigns()).map(buildReport)
  const inEscrow = reports.reduce((sum, r) => sum + BigInt(r.refundableUnits), 0n)
  return {
    inEscrowUnits: inEscrow.toString(),
    campaigns: reports.map((r) => ({ id: r.id, name: r.name, status: r.status })),
  }
})
