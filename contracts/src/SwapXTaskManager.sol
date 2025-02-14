// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Inputs} from "./common/Inputs.sol";
import {CanonicalMachine} from "./common/CanonicalMachine.sol";
import {CoprocessorAdapter} from "../lib/coprocessor-base-contract/src/CoprocessorAdapter.sol";
import {SwapXHook} from "./SwapXHook.sol";

contract SwapX is CoprocessorAdapter {

    event TaskCreated(address indexed taskIssuer, bytes payload);


    error InputTooLarge(
        address appContract,
        uint256 inputLength,
        uint256 maxInputLength
    );

    constructor(
        address _taskIssuer,
        bytes32 _machineHash
    ) CoprocessorAdapter(_taskIssuer, _machineHash) {}

    function createTask(bytes memory input) external payable {
        //TODO: define tokenomics model

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
                input
            )
        );

        if (input.length > CanonicalMachine.INPUT_MAX_SIZE) {
            revert InputTooLarge(
                address(this),
                input.length,
                CanonicalMachine.INPUT_MAX_SIZE
            );
        }

        emit TaskCreated(msg.sender, input);

        callCoprocessor(payload);

    }

    function handleNotice(bytes32 payloadHash, bytes memory notice) internal override {

        (uint256 buyOrderId, uint256 sellOrderId, address hookAddress) = abi.decode(notice, (uint256, uint256, address));
        SwapXHook hook = SwapXHook(hookAddress);

        hook.executeAsyncSwap(buyOrderId, sellOrderId);
    }
}