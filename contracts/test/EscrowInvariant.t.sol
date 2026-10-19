// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;
import {CampaignEscrow} from '../src/CampaignEscrow.sol';
import {TestToken} from './TestToken.sol';
interface ClockVm { function warp(uint256) external; }
contract EscrowHandler {
    TestToken public token;
    CampaignEscrow public escrow;
    bytes32 public constant ORDER=keccak256('invariant-order');
    bytes32 public constant TERMS=keccak256('invariant-terms');
    uint256 private sequence;
    constructor(){
        token=new TestToken();escrow=new CampaignEscrow(address(token),address(0x1111),address(0x2222),address(0x3333));
        escrow.registerOrder(ORDER,address(this),25_000_000,TERMS,uint64(block.timestamp+1));
        token.mint(address(this),25_000_000);token.approve(address(escrow),25_000_000);escrow.fundOrder(ORDER,25_000_000);
        ClockVm(address(uint160(uint256(keccak256('hevm cheat code'))))).warp(block.timestamp+1);
    }
    function release(uint64 candidate)external{uint256 left=escrow.remaining(ORDER);if(left==0)return;uint256 amount=uint256(candidate)%left+1;escrow.releaseEarned(bytes32(++sequence),ORDER,amount,TERMS);}
    function refund(uint64 candidate)external{uint256 left=escrow.remaining(ORDER);if(left==0)return;uint256 amount=uint256(candidate)%left+1;escrow.refundOrder(bytes32(++sequence),ORDER,address(this),amount);}
}
contract EscrowInvariantTest {
    EscrowHandler public handler;
    function setUp()public{handler=new EscrowHandler();}
    function targetContracts()external view returns(address[] memory list){list=new address[](1);list[0]=address(handler);}
    struct FuzzSelector { address addr; bytes4[] selectors; }
    struct FuzzArtifactSelector { string artifact; bytes4[] selectors; }
    struct FuzzInterface { address addr; string[] artifacts; }
    function targetArtifactSelectors()external pure returns(FuzzArtifactSelector[] memory){return new FuzzArtifactSelector[](0);}
    function targetArtifacts()external pure returns(string[] memory){return new string[](0);}
    function excludeArtifacts()external pure returns(string[] memory){return new string[](0);}
    function targetSenders()external pure returns(address[] memory){return new address[](0);}
    function excludeSenders()external pure returns(address[] memory){return new address[](0);}
    function excludeContracts()external pure returns(address[] memory){return new address[](0);}
    function targetInterfaces()external pure returns(FuzzInterface[] memory){return new FuzzInterface[](0);}
    function targetSelectors()external pure returns(FuzzSelector[] memory){return new FuzzSelector[](0);}
    function excludeSelectors()external pure returns(FuzzSelector[] memory){return new FuzzSelector[](0);}
    function invariantConservesEveryDepositedUnit()public view{
        TestToken token=handler.token();CampaignEscrow escrow=handler.escrow();
        uint256 total=token.balanceOf(address(handler))+token.balanceOf(address(escrow))+token.balanceOf(address(0x1111))+token.balanceOf(address(0x2222))+token.balanceOf(address(0x3333));
        require(total==25_000_000,'conservation');require(token.balanceOf(address(escrow))==escrow.remaining(handler.ORDER()),'remaining backing');
    }
}
