// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "forge-std/console.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {Script} from "forge-std/Script.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {SwapXTaskManager} from "src/SwapXTaskManager.sol";
import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";

contract Demo is Script {
    function run() public {
        vm.startBroadcast();

        // Test 0: Corresponds to TestBidFullyMatchedBySingleAsk
        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -50,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(50, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -50,
            sqrtPriceLimitX96: 90
        }), testSettings, abi.encode(50, msg.sender));

        // Test 2: Corresponds to TestBidFullyMatchedByMultipleAsks
        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -100,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(100, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -40,
            sqrtPriceLimitX96: 90
        }), testSettings, abi.encode(40, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -60,
            sqrtPriceLimitX96: 85
        }), testSettings, abi.encode(60, msg.sender));

        // Test 4: Corresponds to TestBidPartiallyMatched
        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -80,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(80, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -50,
            sqrtPriceLimitX96: 90
        }), testSettings, abi.encode(50, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -40,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(40, msg.sender));

        // Test 6: Corresponds to TestAskFullyMatchedBySingleBid
        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -50,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(50, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -50,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(50, msg.sender));

        // Test 8: Corresponds to TestAskFullyMatchedByMultipleBids
        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -50,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(50, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -40,
            sqrtPriceLimitX96: 100
        }), testSettings, abi.encode(40, msg.sender));

        swapRouter.swap(key, IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -100,
            sqrtPriceLimitX96: 90
        }), testSettings, abi.encode(100, msg.sender));

        vm.stopBroadcast();
    }
}