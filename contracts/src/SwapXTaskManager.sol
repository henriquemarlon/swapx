// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Inputs} from "./common/Inputs.sol";
import {CanonicalMachine} from "./common/CanonicalMachine.sol";
import {CoprocessorAdapter} from "../lib/coprocessor-base-contract/src/CoprocessorAdapter.sol";
import {SwapXHook} from "./SwapXHook.sol";

contract SwapXTaskManager is CoprocessorAdapter {
    error InputTooLarge(
        address appContract,
        uint256 inputLength,
        uint256 maxInputLength
    );

    constructor(
        address _taskIssuer,
        bytes32 _machineHash
    ) CoprocessorAdapter(_taskIssuer, _machineHash) {}

    function createTask(bytes memory payload) external payable {
        //TODO: define tokenomics model

        bytes memory input = abi.encodeCall(
            Inputs.EvmAdvance,
            (
                block.chainid,
                address(this),
                msg.sender,
                blockhash(block.number - 1),
                block.number,
                block.timestamp,
                block.prevrandao,
                payload
            )
        );

        if (input.length > CanonicalMachine.INPUT_MAX_SIZE) {
            revert InputTooLarge(
                address(this),
                input.length,
                CanonicalMachine.INPUT_MAX_SIZE
            );
        }

        callCoprocessor(input);
    }

    function handleNotice(bytes32 payloadHash, bytes memory notice) internal override {
        (uint256 buyOrderId, uint256 sellOrderId, address hookAddress) = abi.decode(notice, (uint256, uint256, address));
        SwapXHook hook = SwapXHook(hookAddress);
        hook.executeAsyncSwap(buyOrderId, sellOrderId);
    }
}