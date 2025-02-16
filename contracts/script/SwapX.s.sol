// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "../lib/forge-std/src/Script.sol";
import {SwapXTaskManager} from "../src/SwapXTaskManager.sol";
import {SwapXHook} from "../src/SwapXHook.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {ISwapXTaskManager} from "../src/interface/ISwapXHook.sol";


contract DeploySwapX is Script {
    function run() external returns (SwapXTaskManager, SwapXHook) {
        // These values should be replaced with your actual values
        address taskIssuerAddress = vm.envAddress("TASK_ISSUER_ADDRESS");
        bytes32 machineHash = vm.envBytes32("MACHINE_HASH");

        vm.startBroadcast();
        SwapXTaskManager swapXTaskManager = new SwapXTaskManager(taskIssuerAddress, machineHash);
        SwapXHook swapXHook = new SwapXHook(IPoolManager(address(new PoolManager(address(0)))), ISwapXTaskManager(address(swapXTaskManager)));
        vm.stopBroadcast();
        
        return (swapXTaskManager, swapXHook);
    }
}