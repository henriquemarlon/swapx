// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { BaseAsyncSwap } from "OpenZeppelin/uniswap-hooks/base/BaseAsyncSwap.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";
import {CurrencySettler} from "OpenZeppelin/uniswap-hooks/utils/CurrencySettler.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {BeforeSwapDelta, BeforeSwapDeltaLibrary, toBeforeSwapDelta} from "v4-core/src/types/BeforeSwapDelta.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";

contract OrderBookHook is BaseAsyncSwap {
    using CurrencySettler for Currency;
    using SafeCast for uint256;


    constructor(IPoolManager _poolManager) BaseAsyncSwap(_poolManager) {}

    function _beforeSwap(address sender, PoolKey calldata key, IPoolManager.SwapParams calldata params, bytes calldata)
        internal
        virtual
        override
        returns (bytes4, BeforeSwapDelta, uint24)
    {
        // Async swaps are only possible on exact-input swaps, so exact-output swaps are executed by the `PoolManager` as normal
        if (params.amountSpecified < 0) {
            // Determine which currency is specified
            Currency specified = params.zeroForOne ? key.currency0 : key.currency1;

            // Get the positive specified amount
            uint256 specifiedAmount = uint256(-params.amountSpecified);

            // Mint ERC-6909 claim token for the specified currency and amount
            specified.take(poolManager, address(this), specifiedAmount, false);

            // Return delta that nets out specified amount to 0.
            return (this.beforeSwap.selector, toBeforeSwapDelta(specifiedAmount.toInt128(), 0), 0);
        } else {
            return (this.beforeSwap.selector, BeforeSwapDeltaLibrary.ZERO_DELTA, 0);
        }
    }

}
