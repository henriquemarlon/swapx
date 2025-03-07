// NOTE: This is based on V4PreDeployed.s.sol
// You can make changes to base on V4Deployer.s.sol to deploy everything fresh as well

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "forge-std/Script.sol";
import {SwapXHook} from "src/SwapXHook.sol";
import {console} from "forge-std/console.sol";
import {stdJson} from "forge-std/StdJson.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {HookMiner} from "./utils/HookMiner.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {SwapXTaskManager} from "src/SwapXTaskManager.sol";
import {ISwapXTaskManager} from "src/interface/ISwapXHook.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";

contract HookDeployer is Script {
    PoolKey key;

    ISwapXTaskManager taskManager;
    SwapXTaskManager swapXTaskManager;

    PoolManager manager;
    PoolSwapTest swapRouter;
    PoolModifyLiquidityTest modifyLiquidityRouter;

    Currency token0;
    Currency token1;

    function setUp() public {
        vm.startBroadcast();

        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/broadcast/V4Deployer.s.sol/31337/run-latest.json");
        string memory json = vm.readFile(path);

        address managerAddress = bytesToAddress(stdJson.parseRaw(json, ".transactions[0].contractAddress"));
        manager = PoolManager(managerAddress);
        console.log("Deployed PoolManager at", address(manager));

        address swapRouterAddress = bytesToAddress(stdJson.parseRaw(json, ".transactions[1].contractAddress"));
        swapRouter = PoolSwapTest(swapRouterAddress);
        console.log("Deployed PoolSwapTest at", address(swapRouter));

        address modifyLiquidityRouterAddress =
            bytesToAddress(stdJson.parseRaw(json, ".transactions[2].contractAddress"));
        modifyLiquidityRouter = PoolModifyLiquidityTest(modifyLiquidityRouterAddress);
        console.log("Deployed PoolModifyLiquidityTest at", address(modifyLiquidityRouter));

        address taskIssuerAddress = vm.parseAddress(vm.prompt("Enter Coprocessor address"));
        bytes32 machineHash = vm.parseBytes32(vm.prompt("Enter machine hash"));
        
        swapXTaskManager = new SwapXTaskManager(taskIssuerAddress, machineHash);
        console.log("Deployed SwapXTaskManager at", address(swapXTaskManager));
        
        taskManager = ISwapXTaskManager(address(swapXTaskManager));

        MockERC20 tokenA = new MockERC20("Token0", "TK0", 18);
        MockERC20 tokenB = new MockERC20("Token1", "TK1", 18);

        console.log("Deployed TokenA at", address(tokenA));
        console.log("Deployed TokenB at", address(tokenB));

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

        uint160 flags =
            uint160(Hooks.BEFORE_SWAP_FLAG | Hooks.BEFORE_SWAP_RETURNS_DELTA_FLAG | Hooks.BEFORE_INITIALIZE_FLAG);

        address CREATE2_DEPLOYER = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
        (address hookAddress, bytes32 salt) = HookMiner.find(
            CREATE2_DEPLOYER, flags, type(SwapXHook).creationCode, abi.encode(address(manager), address(taskManager))
        );

        SwapXHook hook = new SwapXHook{salt: salt}(manager, taskManager);
        require(address(hook) == hookAddress, "hook address mismatch");

        console.log("Deployed SwapXHook at", address(hook));

        key = PoolKey({currency0: token0, currency1: token1, fee: 3000, tickSpacing: 120, hooks: hook});

        manager.initialize(key, 79228162514264337593543950336);
        vm.stopBroadcast();
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

    function bytesToAddress(bytes memory bys) private pure returns (address addr) {
        assembly {
            addr := mload(add(bys, 32))
        }
    }
}
