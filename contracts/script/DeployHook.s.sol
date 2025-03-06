// NOTE: This is based on V4PreDeployed.s.sol
// You can make changes to base on V4Deployer.s.sol to deploy everything fresh as well

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "forge-std/Script.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {console} from "forge-std/console.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {HookMiner} from "../test/utils/HookMiner.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {SwapXTaskManager} from "src/SwapXTaskManager.sol";
import {ISwapXTaskManager} from "src/interface/ISwapXHook.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";

contract DeployHook is Script {
    PoolManager manager = PoolManager(0x8464135c8F25Da09e49BC8782676a84730C318bC);
    PoolSwapTest swapRouter = PoolSwapTest(0x71C95911E9a5D330f4D621842EC243EE1343292e);
    PoolModifyLiquidityTest modifyLiquidityRouter = PoolModifyLiquidityTest(0x948B3c65b89DF0B4894ABE91E6D02FE579834F8F);

    Currency token0;
    Currency token1;

    PoolKey key;

    function setUp() public {
        vm.startBroadcast();
        ISwapXTaskManager taskManager = deployTaskManager();
        console.log("Deployed task manager at", address(taskManager));

        MockERC20 tokenA = new MockERC20("Token0", "TK0", 18);
        MockERC20 tokenB = new MockERC20("Token1", "TK1", 18);

        console.log("tokenA", address(tokenA));
        console.log("tokenB", address(tokenB));

        if (address(tokenA) > address(tokenB)) {
            (token0, token1) = (Currency.wrap(address(tokenB)), Currency.wrap(address(tokenA)));
        } else {
            (token0, token1) = (Currency.wrap(address(tokenA)), Currency.wrap(address(tokenB)));
        }

        tokenA.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenB.approve(address(modifyLiquidityRouter), type(uint256).max);
        tokenA.approve(address(swapRouter), type(uint256).max);
        tokenB.approve(address(swapRouter), type(uint256).max);

        tokenA.mint(msg.sender, 100 * 10 ** 18);
        tokenB.mint(msg.sender, 100 * 10 ** 18);

        // Mine for hook address
        vm.stopBroadcast();

        uint160 flags =
            uint160(Hooks.BEFORE_SWAP_FLAG | Hooks.BEFORE_SWAP_RETURNS_DELTA_FLAG | Hooks.BEFORE_INITIALIZE_FLAG);

        address CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
        (address hookAddress, bytes32 salt) = HookMiner.find(
            CREATE2_DEPLOYER, flags, type(SwapXHook).creationCode, abi.encode(address(manager), address(taskManager))
        );

        vm.startBroadcast();
        SwapXHook hook = new SwapXHook{salt: salt}(manager, taskManager);
        require(address(hook) == hookAddress, "hook address mismatch");

        console.log("Deployed hook at", address(hook));

        key = PoolKey({currency0: token0, currency1: token1, fee: 3000, tickSpacing: 120, hooks: hook});

        manager.initialize(key, 79228162514264337593543950336);
        vm.stopBroadcast();
    }

    function deployTaskManager() internal returns (ISwapXTaskManager) {
        address taskIssuerAddress = vm.envAddress("TASK_ISSUER_ADDRESS");
        bytes32 machineHash = vm.envBytes32("MACHINE_HASH");
        SwapXTaskManager swapXTaskManager = new SwapXTaskManager(taskIssuerAddress, machineHash);
        return ISwapXTaskManager(address(swapXTaskManager));
    }

    function run() public {
        vm.startBroadcast();

        modifyLiquidityRouter.modifyLiquidity(
            key,
            IPoolManager.ModifyLiquidityParams({tickLower: -120, tickUpper: 120, liquidityDelta: 10e18, salt: 0}),
            new bytes(0)
        );
        vm.stopBroadcast();
    }
}
