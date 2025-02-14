// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Hooks} from "v4-core/src/libraries/Hooks.sol";

interface ISwapXHook {
    //TODO: define the function signature
    function executeAsyncSwap(uint256 buyOrderId, uint256 sellOrderId) external;
    function cancelBuyOrder(uint256 orderId) external;
    function cancelSellOrder(uint256 orderId) external;
}