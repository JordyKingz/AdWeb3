// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;
import {CampaignEscrow} from '../src/CampaignEscrow.sol';
import {TestToken} from './TestToken.sol';
interface Vm { function prank(address) external; function expectRevert() external; function warp(uint256) external; }
contract CampaignEscrowTest {
    Vm constant vm=Vm(address(uint160(uint256(keccak256('hevm cheat code')))));
    TestToken token;
    CampaignEscrow escrow;
    bytes32 constant ORDER=keccak256('order');
    bytes32 constant TERMS=keccak256('terms');
    address constant PAYER=address(0xA11CE);
    address constant OPS=address(0xB0B);
    address constant RESERVE=address(0xCAFE);
    address constant POOL=address(0xD00D);
    function setUp() public {
        token=new TestToken(); escrow=new CampaignEscrow(address(token),OPS,RESERVE,POOL);
        escrow.registerOrder(ORDER,PAYER,25_000_000,TERMS,uint64(block.timestamp+1 days));
        vm.warp(block.timestamp+1 days);
        token.mint(PAYER,25_000_000);
        vm.prank(PAYER); token.approve(address(escrow),25_000_000);
        vm.prank(PAYER); escrow.fundOrder(ORDER,25_000_000);
    }
    function testDepositProtectedAndRefundGoesToPayer() public {
        require(escrow.remaining(ORDER)==25_000_000,'deposit not protected');
        escrow.refundOrder(keccak256('refund'),ORDER,PAYER,15_000_000);
        require(token.balanceOf(PAYER)==15_000_000,'refund missing');
        require(escrow.remaining(ORDER)==10_000_000,'refund accounting');
    }
    function testReleaseAllocatesThirtySixtyTenOnce() public {
        bytes32 id=keccak256('settlement');
        escrow.releaseEarned(id,ORDER,10_000_000,TERMS);
        require(token.balanceOf(POOL)==3_000_000,'reward allocation');
        require(token.balanceOf(OPS)==6_000_000,'ops allocation');
        require(token.balanceOf(RESERVE)==1_000_000,'reserve allocation');
        vm.expectRevert(); escrow.releaseEarned(id,ORDER,10_000_000,TERMS);
    }
    function testRejectUnauthorizedAndChangedTermsOrDestination() public {
        vm.prank(PAYER); vm.expectRevert(); escrow.releaseEarned(keccak256('s'),ORDER,1,TERMS);
        vm.expectRevert(); escrow.releaseEarned(keccak256('s'),ORDER,1,bytes32(0));
        vm.expectRevert(); escrow.refundOrder(keccak256('r'),ORDER,OPS,1);
    }
    function testCannotOverdrawOrRefundTwice() public {
        escrow.refundOrder(keccak256('r'),ORDER,PAYER,25_000_000);
        vm.expectRevert(); escrow.refundOrder(keccak256('r'),ORDER,PAYER,1);
        vm.expectRevert(); escrow.releaseEarned(keccak256('s'),ORDER,1,TERMS);
    }
    function testFuzzConservation(uint64 accepted) public {
        uint256 amount=uint256(accepted)%25_000_001;
        if(amount==0) return;
        escrow.releaseEarned(keccak256('s'),ORDER,amount,TERMS);
        uint256 total=token.balanceOf(POOL)+token.balanceOf(OPS)+token.balanceOf(RESERVE)+escrow.remaining(ORDER);
        require(total==25_000_000,'conservation');
    }
}

contract EscrowBoundaryTest {
    Vm constant vm=Vm(address(uint160(uint256(keccak256('hevm cheat code')))));
    function testReleaseRequiresHoldAndOrderIsImmutable() public {
        TestToken token=new TestToken();CampaignEscrow e=new CampaignEscrow(address(token),address(1),address(2),address(3));bytes32 id=keccak256('o');bytes32 terms=keccak256('t');
        e.registerOrder(id,address(this),100,terms,uint64(block.timestamp+100));token.mint(address(this),100);token.approve(address(e),100);e.fundOrder(id,100);
        vm.expectRevert();e.registerOrder(id,address(this),100,terms,uint64(block.timestamp+100));
        vm.expectRevert();e.releaseEarned(keccak256('s'),id,50,terms);
        vm.warp(block.timestamp+100);e.releaseEarned(keccak256('s'),id,60,terms);
        vm.expectRevert();e.refundOrder(keccak256('r'),id,address(this),50);
        e.refundOrder(keccak256('r'),id,address(this),40);require(e.remaining(id)==0,'remaining');
    }
    function testRejectShortOrRepeatedFunding() public {
        TestToken token=new TestToken();CampaignEscrow e=new CampaignEscrow(address(token),address(1),address(2),address(3));bytes32 id=keccak256('o');
        e.registerOrder(id,address(this),100,keccak256('t'),uint64(block.timestamp+100));token.mint(address(this),200);token.approve(address(e),200);
        vm.expectRevert();e.fundOrder(id,99);e.fundOrder(id,100);vm.expectRevert();e.fundOrder(id,100);
    }
}
contract FeeToken {
 mapping(address=>uint256)public balanceOf;
 function approve(address,uint256)external pure returns(bool){return true;}
 function transferFrom(address,address to,uint256 amount)external returns(bool){balanceOf[to]+=amount-1;return true;}
 function transfer(address,uint256)external pure returns(bool){return true;}
}
contract EscrowTokenBehaviorTest {
 Vm constant vm=Vm(address(uint160(uint256(keccak256('hevm cheat code')))));
 function testRejectUnexpectedReceivedAmount()public{
  FeeToken token=new FeeToken();CampaignEscrow e=new CampaignEscrow(address(token),address(1),address(2),address(3));bytes32 id=keccak256('o');e.registerOrder(id,address(this),100,keccak256('t'),uint64(block.timestamp+100));vm.expectRevert();e.fundOrder(id,100);require(e.remaining(id)==0,'fee token credited');
 }
}
