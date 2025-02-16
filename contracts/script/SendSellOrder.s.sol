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
        PoolManager(0x5FbDB2315678afecb367f032d93F642f64180aa3);
    PoolSwapTest swapRouter =
        PoolSwapTest(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);
    PoolModifyLiquidityTest modifyLiquidityRouter =
        PoolModifyLiquidityTest(0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0);
    SwapXHook hook = SwapXHook(0x614439a0d066DE45FBD7190dac6De2C702D42088);

    PoolKey key;

    Currency token0;
    Currency token1;

    function setUp() public {
        vm.startBroadcast();

        MockERC20 tokenA = MockERC20(
            0x5f3f1dBD7B74C6B46e8c44f98792A1dAf8d69154
        );
        MockERC20 tokenB = MockERC20(
            0xb7278A61aa25c888815aFC32Ad3cC52fF24fE575
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
