// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "../lib/forge-std/src/Script.sol";
import {SwapXTaskManager} from "../src/SwapXTaskManager.sol";

contract DeployTaskManager is Script {
    function run() external returns (SwapXTaskManager) {
        address taskIssuerAddress = vm.envAddress("TASK_ISSUER_ADDRESS");
        bytes32 machineHash = vm.envBytes32("MACHINE_HASH");

        vm.startBroadcast();
        SwapXTaskManager swapXTaskManager = new SwapXTaskManager(taskIssuerAddress, machineHash);
        vm.stopBroadcast();

        return swapXTaskManager;
    }
}