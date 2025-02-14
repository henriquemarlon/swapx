// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Inputs} from "./common/Inputs.sol";
import {CanonicalMachine} from "./common/CanonicalMachine.sol";
import {CoprocessorAdapter} from "../lib/coprocessor-base-contract/src/CoprocessorAdapter.sol";
import {SwapXHook} from "./SwapXHook.sol";

contract SwapX is CoprocessorAdapter {
    error InputTooLarge(
        address appContract,
        uint256 inputLength,
        uint256 maxInputLength
    );

    constructor(
        address _taskIssuer,
        bytes32 _machineHash
    ) CoprocessorAdapter(_taskIssuer, _machineHash) {}

    mapping(address => bytes[]) private _hookSwaps;

    function createTask(bytes memory input) external payable {
        //TODO: define tokenomics model

        bytes[] storage swaps = _hookSwaps[msg.sender];

        uint256 index = swaps.length;

        bytes memory payload = abi.encodeCall(
            Inputs.EvmAdvance,
            (
                block.chainid,
                address(this),
                msg.sender,
                blockhash(block.number),
                block.number,
                block.timestamp,
                block.prevrandao,
                index
            )
        );

        if (input.length > CanonicalMachine.INPUT_MAX_SIZE) {
            revert InputTooLarge(
                address(this),
                input.length,
                CanonicalMachine.INPUT_MAX_SIZE
            );
        }

        callCoprocessor(payload);

        // Store input for later retrieval through the GIO (GET_STORAGE_GIO 0x27)
        // What is the sslot and the element address given an hook address?
        swaps.push(input);
    }

    function handleNotice(bytes32 payloadHash, bytes memory notice) internal override {
        address destination;
        bytes memory decodedPayload;
        SwapXHook hook = SwapXHook(destination);
        hook.callHook();
    }
}