// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.26;

import "forge-std/Script.sol";
import {Owner} from "./Owner.sol";

contract Test is Script {
    function run() external {
        vm.startBroadcast();

        Owner owner = new Owner();

        bytes memory payload = owner.encode();

        console.logBytes(payload);

        bytes32 blockHash = owner.getCurrentBlockHash();
        console.logBytes32(blockHash);

        vm.stopBroadcast();
    }
}
