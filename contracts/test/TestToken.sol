// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;
contract TestToken {
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    function mint(address to, uint256 amount) external { balanceOf[to] += amount; }
    function approve(address to,uint256 amount) external returns(bool) { allowance[msg.sender][to]=amount; return true; }
    function transfer(address to,uint256 amount) external returns(bool) { require(balanceOf[msg.sender]>=amount); balanceOf[msg.sender]-=amount; balanceOf[to]+=amount; return true; }
    function transferFrom(address from,address to,uint256 amount) external returns(bool) { require(allowance[from][msg.sender]>=amount); allowance[from][msg.sender]-=amount; require(balanceOf[from]>=amount); balanceOf[from]-=amount; balanceOf[to]+=amount; return true; }
}
