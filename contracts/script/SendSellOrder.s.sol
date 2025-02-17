// NOTE: This is based on V4PreDeployed.s.sol
// You can make changes to base on V4Deployer.s.sol to deploy everything fresh as well

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "forge-std/Script.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";
import {PoolDonateTest} from "v4-core/src/test/PoolDonateTest.sol";
import {PoolTakeTest} from "v4-core/src/test/PoolTakeTest.sol";
import {PoolClaimsTest} from "v4-core/src/test/PoolClaimsTest.sol";
import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {HookMiner} from "./HookMiner.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {ISwapXTaskManager} from "src/interface/ISwapXHook.sol";
import {SwapXManagerMock} from "../test/mocks/SwapXManagerMock.sol";
import "forge-std/console.sol";

contract SendSellOrder is Script {
    PoolManager manager =
        PoolManager(0x1429859428C0aBc9C2C47C8Ee9FBaf82cFA0F20f);
    PoolSwapTest swapRouter =
        PoolSwapTest(0xB0D4afd8879eD9F52b28595d31B441D079B2Ca07);
    PoolModifyLiquidityTest modifyLiquidityRouter =
        PoolModifyLiquidityTest(0x162A433068F51e18b7d13932F27e66a3f99E6890);
    SwapXHook hook = SwapXHook(0x43B0c089e42cC9Dd636F730512F089602bF7e088);
    PoolKey key;

    Currency token0;
    Currency token1;

    function setUp() public {
        vm.startBroadcast();

        MockERC20 tokenA = MockERC20(
            0x0355B7B8cb128fA5692729Ab3AAa199C1753f726
        );
        MockERC20 tokenB = MockERC20(
            0x202CCe504e04bEd6fC0521238dDf04Bc9E8E15aB
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
            zeroForOne: false,
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
