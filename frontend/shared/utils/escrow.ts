import { parseAbi } from 'viem'

/** The subset of contracts/src/CampaignEscrow.sol the dApp calls or reads. */
export const campaignEscrowAbi = parseAbi([
  'function fundOrder(bytes32 id, uint256 amount)',
  'function orders(bytes32 id) view returns (address payer, uint256 cap, uint256 funded, uint256 spent, bytes32 terms, uint64 releaseAfter)',
  'function remaining(bytes32 id) view returns (uint256)',
  'event Funded(bytes32 indexed orderId, address indexed payer, uint256 amount)',
])
