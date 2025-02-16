// SPDX-License-Identifier: MIT

pragma solidity 0.8.26;

interface Inputs {
    function EvmAdvance(
        uint256 chainId,
        address taskManager,
        address msgSender,
        bytes32 blockHash,
        uint256 blockNumber,
        uint256 blockTimestamp,
        uint256 prevRandao,
        bytes memory payload
    ) external;
}
