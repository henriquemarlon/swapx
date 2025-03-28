## <div align="center">Disclaimer</div>

Below is a list of concerns that should be taken into consideration in the current version.

1 - Lack of State Root Verification in GIO Responses:

Currently, the application relies on an external API ( GIO - domain 0x27 ) to fetch Ethereum contract storage data, but this API does not yet implement state root verification. This introduces a potential trust issue, as we cannot independently verify the integrity and authenticity of the returned data.

Without cross-checking the retrieved storage values against the stateRoot of the corresponding block, we are inherently trusting the API to provide correct and untampered data. This leaves the system vulnerable to:

- Malicious or compromised API nodes providing falsified storage values.
- Outdated or inconsistent data if the API does not correctly synchronize with the blockchain.
- Lack of cryptographic guarantees regarding the correctness of the returned state.

To mitigate these risks, we should consider implementing verification using eth_getProof, which provides Merkle proofs for account storage. By validating these proofs against the stateRoot of the target block, we can ensure that the fetched storage values were indeed part of Ethereum’s canonical state at that block height.

Despite that, this is not a critical factor because we assume the operator is honest and, therefore, connected to a reliable node.

2 - Non-blocking Execution:

Currently, we do not have a queue system for processing outputs in the coprocessor, and we do not block the processing of a new input if there are not yet **N** confirmations that the output generated by the previous input has been executed. This is a problem in the following analogous scenario:  

We have a contract that stores swap orders, and our backend within the coprocessor reads the state of this array via **GIO**, for example. Each order has a status that gets updated upon execution. Our backend logic works as follows:  
- If the status is **A**, with **A** being arbitrary, we use that order for matching in the order book.  
- If the status is **B**, with **B** being arbitrary, we do not use it because it has already been executed.  

Now, considering the inputs/tasks in the system:  
- **Input A**: Signals that a new order has been submitted and checks if there are orders in the previous block (there aren’t).  
- **Input B**: Signals a new order, reads the order from the previous block, and tries to match them (eventually, a match occurs, and a notice is issued for the actual swap execution).  
- **Input C**: Also introduces a new order.  

This concern is based on behavior we observed on Otterscan. This is the order of transactions:  

1st Swap order **A** is created.  
2nd Swap order **B** is created.  
3rd Swap order **C** is created.  
4th The swap is executed between **A** and **B** (orders only change status when executed).  

Thus, there is a possibility that swap **C** might be reading the outdated status of orders **A** and **B**.

3 - Non-Fault-Tolerant Output Execution:

The current design of the coprocessor envisions batch execution of outputs derived from the same input. Taking this to the extreme, if one of the multiple outputs of an application fails and reverts during execution, all other outputs will also be reverted, which is highly undesirable. In our application, we have a similar potential issue:

One of the order executions defined after the order book matching process may revert due to an arbitrary error, effectively **blocking the execution of other swaps** since the entire transaction in which the batch of orders was executed will revert due to a single failing order.

4 - Undefined Economic Model of the Coprocessor:

The coprocessor infrastructure inherently has costs, whether from the execution of outputs (**gas fees**) or the infrastructure maintained by the operators. However, the current setup has not yet defined how coprocessor calls ( issueTask ) will be charged to cover these costs and ensure financial incentives for maintaining a network with multiple operators.