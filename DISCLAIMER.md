# DISCLAIMER

Below is a list of concerns that should be taken into consideration in the current version.

1 - Lack of State Root Verification in GIO Responses:

Currently, the application relies on an external API ( GIO - domain 0x27 ) to fetch Ethereum contract storage data, but this API does not yet implement state root verification. This introduces a potential trust issue, as we cannot independently verify the integrity and authenticity of the returned data.

Without cross-checking the retrieved storage values against the stateRoot of the corresponding block, we are inherently trusting the API to provide correct and untampered data. This leaves the system vulnerable to:

- Malicious or compromised API nodes providing falsified storage values.
- Outdated or inconsistent data if the API does not correctly synchronize with the blockchain.
- Lack of cryptographic guarantees regarding the correctness of the returned state.

To mitigate these risks, we should consider implementing verification using eth_getProof, which provides Merkle proofs for account storage. By validating these proofs against the stateRoot of the target block, we can ensure that the fetched storage values were indeed part of Ethereum’s canonical state at that block height.