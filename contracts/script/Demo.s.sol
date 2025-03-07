// NOTE: This is based on V4PreDeployed.s.sol
// You can make changes to base on V4Deployer.s.sol to deploy everything fresh as well

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
    PoolKey key;

    MockERC20 tokenA;
    MockERC20 tokenB;

    SwapXHook hook;
    PoolSwapTest swapRouter;
    SwapXTaskManager swapXTaskManager;
    PoolModifyLiquidityTest modifyLiquidityRouter;

    Currency currency0;
    Currency currency1;

    function setUp() public {
        vm.startBroadcast();

        string memory root = vm.projectRoot();

        ////////

        string memory v4Path = string.concat(root, "/broadcast/V4Deployer.s.sol/31337/run-latest.json");
        string memory v4Json = vm.readFile(v4Path);

        address swapRouterAddress = bytesToAddress(stdJson.parseRaw(v4Json, ".transactions[1].contractAddress"));
        swapRouter = PoolSwapTest(swapRouterAddress);
        console.log("Deployed PoolSwapTest at", address(swapRouter));

        address modifyLiquidityRouterAddress =
            bytesToAddress(stdJson.parseRaw(v4Json, ".transactions[2].contractAddress"));
        modifyLiquidityRouter = PoolModifyLiquidityTest(modifyLiquidityRouterAddress);
        console.log("Deployed PoolModifyLiquidityTest at", address(modifyLiquidityRouter));

        //////

        string memory hookPath = string.concat(root, "/broadcast/HookDeployer.s.sol/31337/run-latest.json");
        string memory hookJson = vm.readFile(hookPath);

        address tokenAAddress = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[1].contractAddress"));
        tokenA = MockERC20(tk0Address);
        console.log("Deployed TokenA at", address(tokenA));

        address tokenBAddress = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[2].contractAddress"));
        tokenB = MockERC20(tk1Address);
        console.log("Deployed TokenB at", address(tokenB));

        address swapXHookAddress = bytesToAddress(stdJson.parseRaw(hookJson, ".transactions[9].contractAddress"));
        hook = SwapXHook(swapXHookAddress);
        console.log("Deployed SwapXHook at", address(hook));

        if (address(tokenA) > address(tokenB)) {
            (currency0, currency1) = (Currency.wrap(address(tokenB)), Currency.wrap(address(tokenA)));
        } else {
            (currency0, currency1) = (Currency.wrap(address(tokenA)), Currency.wrap(address(tokenB)));
        }

        tokenA.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenB.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenA.approve(address(swapRouter), type(uint256).max);
        tokenB.approve(address(swapRouter), type(uint256).max);

        tokenB.mint(msg.sender, 100 * 10 ** 18);
        tokenA.mint(msg.sender, 100 * 10 ** 18);

        key = PoolKey({currency0: currency0, currency1: currency1, fee: 3000, tickSpacing: 120, hooks: hook});
        vm.stopBroadcast();
    }

    function run() public {
        vm.startBroadcast();

        IPoolManager.SwapParams memory buySwapParams = IPoolManager.SwapParams({
            zeroForOne: true,
            amountSpecified: -100,
            sqrtPriceLimitX96: 79228162514264337593543950336
        });

        IPoolManager.SwapParams memory sellSwapParams = IPoolManager.SwapParams({
            zeroForOne: false,
            amountSpecified: -100,
            sqrtPriceLimitX96: 79228162514264337593543950336
        });

        PoolSwapTest.TestSettings memory testSettings =
            PoolSwapTest.TestSettings({takeClaims: false, settleUsingBurn: false});

        swapRouter.swap(key, buySwapParams, testSettings, abi.encode(1000000000000000000, msg.sender));

        swapRouter.swap(key, sellSwapParams, testSettings, abi.encode(1000000000000000000, msg.sender));
        vm.stopBroadcast();
    }

    function bytesToAddress(bytes memory bys) private pure returns (address addr) {
        assembly {
            addr := mload(add(bys, 32))
        }
    }
}