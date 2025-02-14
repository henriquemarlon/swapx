// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

interface ISwapXHook {
    //TODO: define the function signature
    function executeAsyncSwap(uint256 buyOrderId, uint256 sellOrderId) external;
}