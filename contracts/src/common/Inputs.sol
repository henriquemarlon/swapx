// SPDX-License-Identifier: MIT

pragma solidity 0.8.28;

/// @title Inputs
/// @notice Defines the signatures of inputs.
interface Inputs {
    /// @notice An advance request from an EVM-compatible blockchain to a Cartesi Machine.
    /// @param chainId The chain ID
    /// @param appContract The application contract address
    /// @param msgSender The address of whoever sent the input
    /// @param blockHash The hash of the block in which the input was added
    /// @param blockNumber The number of the block in which the input was added
    /// @param blockTimestamp The timestamp of the block in which the input was added
    /// @param prevRandao The latest RANDAO mix of the post beacon state of the previous block
    /// @param payload The payload of the input
    /// @dev See EIP-4399 for safe usage of `prevRandao`.
    function EvmAdvance(
        uint256 chainId,
        address appContract,
        address msgSender,
        bytes32 blockHash,
        uint256 blockNumber,
        uint256 blockTimestamp,
        uint256 prevRandao,
        bytes memory payload
    ) external;
}