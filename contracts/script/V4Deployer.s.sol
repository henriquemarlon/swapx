// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {Script} from "forge-std/Script.sol";
import {PoolManager} from "v4-core/src/PoolManager.sol";
import {PoolSwapTest} from "v4-core/src/test/PoolSwapTest.sol";
import {PoolModifyLiquidityTest} from "v4-core/src/test/PoolModifyLiquidityTest.sol";

import "forge-std/console.sol";

contract V4Deployer is Script {
    function run() public {
        vm.startBroadcast();

        address ownerManager = 0xa0Ee7A142d267C1f36714E4a8F75612F20a79720; // anvil last test account

        PoolManager manager = new PoolManager(ownerManager);
        console.log("Deployed PoolManager at", address(manager));
        PoolSwapTest swapRouter = new PoolSwapTest(manager);
        console.log("Deployed PoolSwapTest at", address(swapRouter));
        PoolModifyLiquidityTest modifyLiquidityRouter = new PoolModifyLiquidityTest(manager);
        console.log("Deployed PoolModifyLiquidityTest at", address(modifyLiquidityRouter));
        vm.stopBroadcast();
    }
}
