// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

interface Inputs {
    function EvmAdvance(
        uint256 chainId,
        address taskManager,
        address msgSender,
        bytes32 blockHash,
        uint256 blockNumber,
        uint256 blockTimestamp,
        bytes memory payload
    ) external;
}

/**
 * @title Owner
 * @dev Set & change owner
 */
contract Owner {
    function encode() public view returns(bytes memory){
        bytes memory payload = abi.encodeCall(
            Inputs.EvmAdvance,
            (
                block.chainid,
                address(this),
                msg.sender,
                blockhash(block.number-1),
                block.number,
                block.timestamp,
                abi.encode(uint256(12), uint256(18))
            )
        );
        return payload;
    }

    function getCurrentBlockHash() public view returns (bytes32) {
        return blockhash(block.number);
    }
} 
