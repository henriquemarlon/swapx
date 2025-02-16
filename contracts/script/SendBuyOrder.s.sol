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
import {HookMiner} from "./HookMiner.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {ISwapXTaskManager} from "src/interface/ISwapXHook.sol";
import {SwapXManagerMock} from "../test/mocks/SwapXManagerMock.sol";
import "forge-std/console.sol";

contract SendBuyOrder is Script {
    PoolManager manager =
        PoolManager(0x68B1D87F95878fE05B998F19b66F4baba5De1aed);
    PoolSwapTest swapRouter =
        PoolSwapTest(0x3Aa5ebB10DC797CAC828524e59A333d0A371443c);
    PoolModifyLiquidityTest modifyLiquidityRouter =
        PoolModifyLiquidityTest(0xc6e7DF5E7b4f2A278906862b61205850344D4e7d);
    SwapXHook hook = SwapXHook(0x0e75A2f72c53548E5b45E8a03179C69D6C0Ce088);

    PoolKey key;

    Currency token0;
    Currency token1;

    function setUp() public {
        vm.startBroadcast();

        MockERC20 tokenA = MockERC20(
            0x4ed7c70F96B99c776995fB64377f0d4aB3B0e1C1
        );
        MockERC20 tokenB = MockERC20(
            0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
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
