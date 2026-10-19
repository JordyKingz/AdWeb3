// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;
interface IERC20 {
    function balanceOf(address) external view returns(uint256);
    function transfer(address,uint256) external returns(bool);
    function transferFrom(address,address,uint256) external returns(bool);
}
/// @notice Testnet escrow. Delivery acceptance is an explicitly trusted operator decision.
contract CampaignEscrow {
    struct Order { address payer; uint256 cap; uint256 funded; uint256 spent; bytes32 terms; uint64 releaseAfter; }
    IERC20 public immutable token;
    address public immutable operator;
    address public immutable operations;
    address public immutable reserve;
    address public immutable rewardPool;
    mapping(bytes32=>Order) public orders;
    mapping(bytes32=>bool) public processed;
    bool private entered;
    event Funded(bytes32 indexed orderId,address indexed payer,uint256 amount);
    event Released(bytes32 indexed settlementId,bytes32 indexed orderId,uint256 amount,uint256 players,uint256 ops,uint256 contingency);
    event Refunded(bytes32 indexed refundId,bytes32 indexed orderId,address indexed payer,uint256 amount);
    modifier onlyOperator(){ require(msg.sender==operator,'unauthorized'); _; }
    modifier nonReentrant(){ require(!entered,'reentrancy'); entered=true; _; entered=false; }
    constructor(address asset,address ops,address contingency,address pool) {
        require(asset.code.length>0 && ops!=address(0) && contingency!=address(0) && pool!=address(0),'invalid config');
        token=IERC20(asset); operator=msg.sender; operations=ops; reserve=contingency; rewardPool=pool;
    }
    function registerOrder(bytes32 id,address payer,uint256 cap,bytes32 terms,uint64 releaseAfter) external onlyOperator {
        require(id!=bytes32(0) && payer!=address(0) && cap>0 && terms!=bytes32(0) && releaseAfter>block.timestamp,'invalid order');
        require(orders[id].payer==address(0),'exists');
        orders[id]=Order(payer,cap,0,0,terms,releaseAfter);
    }
    function remaining(bytes32 id) public view returns(uint256){ Order storage o=orders[id]; return o.funded-o.spent; }
    function fundOrder(bytes32 id,uint256 amount) external nonReentrant {
        Order storage o=orders[id];
        require(msg.sender==o.payer && amount>0 && o.funded==0 && amount==o.cap,'invalid funding');
        uint256 beforeBalance=token.balanceOf(address(this));
        require(token.transferFrom(msg.sender,address(this),amount),'transfer failed');
        require(token.balanceOf(address(this))-beforeBalance==amount,'received mismatch');
        o.funded=amount;
        emit Funded(id,msg.sender,amount);
    }
    function releaseEarned(bytes32 settlementId,bytes32 orderId,uint256 amount,bytes32 allocationHash) external onlyOperator nonReentrant {
        Order storage o=orders[orderId];
        require(block.timestamp>=o.releaseAfter,'hold');
        require(allocationHash==o.terms && o.payer!=address(0),'terms');
        consume(settlementId,orderId,amount);
        uint256 players=amount*3000/10000;
        uint256 ops=amount*6000/10000;
        uint256 contingency=amount-players-ops;
        transferExact(rewardPool,players); transferExact(operations,ops); transferExact(reserve,contingency);
        emit Released(settlementId,orderId,amount,players,ops,contingency);
    }
    function refundOrder(bytes32 refundId,bytes32 orderId,address recipient,uint256 amount) external onlyOperator nonReentrant {
        require(recipient==orders[orderId].payer && recipient!=address(0),'payer');
        consume(refundId,orderId,amount); transferExact(recipient,amount);
        emit Refunded(refundId,orderId,recipient,amount);
    }
    function consume(bytes32 id,bytes32 orderId,uint256 amount) private {
        require(id!=bytes32(0) && !processed[id] && amount>0 && amount<=remaining(orderId),'invalid consumption');
        processed[id]=true; orders[orderId].spent+=amount;
    }
    function transferExact(address to,uint256 amount) private {
        if(amount==0) return;
        uint256 beforeBalance=token.balanceOf(to);
        require(token.transfer(to,amount),'transfer failed');
        require(token.balanceOf(to)-beforeBalance==amount,'received mismatch');
    }
}
