// NOTE: This is based on V4PreDeployed.s.sol
// You can make changes to base on V4Deployer.s.sol to deploy everything fresh as well

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "forge-std/Script.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";
import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {ISwapXTaskManager} from "src/interface/ISwapXHook.sol";
import {SwapXManagerMock} from "../test/mocks/SwapXManagerMock.sol";
import "forge-std/console.sol";

contract SendBuyOrder is Script {
    PoolManager manager =
        PoolManager(0x8464135c8F25Da09e49BC8782676a84730C318bC);
    PoolSwapTest swapRouter =
        PoolSwapTest(0x71C95911E9a5D330f4D621842EC243EE1343292e);
    PoolModifyLiquidityTest modifyLiquidityRouter =
        PoolModifyLiquidityTest(0x948B3c65b89DF0B4894ABE91E6D02FE579834F8F);
    SwapXHook hook = SwapXHook(0x033D39E55607694e824828C897047ff6059DA088);

    PoolKey key;

    Currency token0;
    Currency token1;

    function setUp() public {
        vm.startBroadcast();

        MockERC20 tokenA = MockERC20(
            0x2572e04Caf46ba8692Bd6B4CBDc46DAA3cA9647E
        );
        MockERC20 tokenB = MockERC20(
            0x72F375F23BCDA00078Ac12e7e9E7f6a8CA523e7D
        );

        if (address(tokenA) > address(tokenB)) {
            (token0, token1) = (
                Currency.wrap(address(tokenB)),
                Currency.wrap(address(tokenA))
            );
        } else {
            (token0, token1) = (
                Currency.wrap(address(tokenA)),
                Currency.wrap(address(tokenB))
            );
        }

        tokenA.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenB.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenA.approve(address(swapRouter), type(uint256).max);
        tokenB.approve(address(swapRouter), type(uint256).max);

        tokenA.mint(msg.sender, 100 * 10 ** 18);
        tokenB.mint(msg.sender, 100 * 10 ** 18);

        key = PoolKey({
            currency0: token0,
            currency1: token1,
            fee: 3000,
            tickSpacing: 120,
            hooks: hook
        });

        vm.stopBroadcast();
    }

    function run() public {
        vm.startBroadcast();

        // buy params
        IPoolManager.SwapParams memory swapParams = IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -100,
            sqrtPriceLimitX96: 56022770974786139918731938227
        });

        PoolSwapTest.TestSettings memory testSettings = PoolSwapTest
            .TestSettings({takeClaims: false, settleUsingBurn: false});

        uint256 sqrtPrice = 1000000000000000000;

        swapRouter.swap(
            key,
            swapParams,
            testSettings,
            abi.encode(sqrtPrice, msg.sender)
        );

        vm.stopBroadcast();
    }
}
