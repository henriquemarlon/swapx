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

contract Demo04 is Script {
    PoolKey key;

    MockERC20 tk0;
    MockERC20 tk1;

    SwapXHook hook;
    PoolSwapTest swapRouter;
    SwapXTaskManager swapXTaskManager;
    PoolModifyLiquidityTest modifyLiquidityRouter;

    Currency currency0;
    Currency currency1;

    function setUp() public {
        string memory root = vm.projectRoot();

        string memory v4Path = string.concat(root, "/broadcast/V4Deployer.s.sol/31337/run-latest.json");
        string memory v4Json = vm.readFile(v4Path);

        address swapRouterAddress = bytesToAddress(stdJson.parseRaw(v4Json, ".transactions[1].contractAddress"));
        swapRouter = PoolSwapTest(swapRouterAddress);
        console.log("Deployed PoolSwapTest at", address(swapRouter));

        address modifyLiquidityRouterAddress =
            bytesToAddress(stdJson.parseRaw(v4Json, ".transactions[2].contractAddress"));
        modifyLiquidityRouter = PoolModifyLiquidityTest(modifyLiquidityRouterAddress);
        console.log("Deployed PoolModifyLiquidityTest at", address(modifyLiquidityRouter));

        string memory hookPath = string.concat(root, "/broadcast/HookDeployer.s.sol/31337/run-latest.json");
        string memory hookJson = vm.readFile(hookPath);

        address tk0Address = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[1].contractAddress"));
        tk0 = MockERC20(tk0Address);
        console.log("Deployed TokenA at", address(tk0));

        address tk1Address = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[2].contractAddress"));
        tk1 = MockERC20(tk1Address);
        console.log("Deployed TokenB at", address(tk1));

        address swapXHookAddress = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[9].contractAddress"));
        hook = SwapXHook(swapXHookAddress);
        console.log("Deployed SwapXHook at", address(hook));

        vm.startBroadcast();
        if (address(tk0) > address(tk1)) {
            (currency0, currency1) = (Currency.wrap(address(tk1)), Currency.wrap(address(tk0)));
        } else {
            (currency0, currency1) = (Currency.wrap(address(tk0)), Currency.wrap(address(tk1)));
        }

        tk0.approve(address(modifyLiquidityRouter), type(uint256).max);
        tk1.approve(address(modifyLiquidityRouter), type(uint256).max);
        tk0.approve(address(swapRouter), type(uint256).max);
        tk1.approve(address(swapRouter), type(uint256).max);

        tk1.mint(msg.sender, 100 * 10 ** 18);
        tk0.mint(msg.sender, 100 * 10 ** 18);

        key = PoolKey({currency0: currency0, currency1: currency1, fee: 3000, tickSpacing: 120, hooks: hook});
        vm.stopBroadcast();
    }

    // Case: BuyOrderPartiallyFulfilledByMultipleSellOrders
    // Number of matchs: 3 ( because of the last swap )
    function run() public {
        PoolSwapTest.TestSettings memory testSettings =
            PoolSwapTest.TestSettings({takeClaims: false, settleUsingBurn: false});

        vm.startBroadcast();
        swapRouter.swap(
            key,
            IPoolManager.SwapParams({zeroForOne: false, amountSpecified: -20, sqrtPriceLimitX96: 100}),
            testSettings,
            abi.encode(100, msg.sender)
        );

        swapRouter.swap(
            key,
            IPoolManager.SwapParams({zeroForOne: false, amountSpecified: -30, sqrtPriceLimitX96: 120}),
            testSettings,
            abi.encode(120, msg.sender)
        );

        swapRouter.swap(
            key,
            IPoolManager.SwapParams({zeroForOne: true, amountSpecified: -60, sqrtPriceLimitX96: 120}),
            testSettings,
            abi.encode(120, msg.sender)
        );

        // match remaining amount
        swapRouter.swap(
            key,
            IPoolManager.SwapParams({zeroForOne: false, amountSpecified: -10, sqrtPriceLimitX96: 120}),
            testSettings,
            abi.encode(120, msg.sender)
        );
        vm.stopBroadcast();
    }

    function bytesToAddress(bytes memory bys) private pure returns (address addr) {
        assembly {
            addr := mload(add(bys, 32))
        }
    }
}
