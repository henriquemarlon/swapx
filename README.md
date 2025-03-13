<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>An EVM Linux-powered coprocessor as an orderbook for UniswapV4</i>
</div>
<div align="center">
<b>Cartesi Coprocessor orderbook powered by EigenLayer cryptoeconomic security</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/swapx?style=default&logo=opensourceinitiative&logoColor=white&color=79F7FA" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/swapx?style=default&logo=git&logoColor=white&color=868380" alt="last-commit">
</p>

> [!CAUTION]
> This is an experimental project under continuous development and should be treated as such. **Its use in production/mainnet is not recommended.** Learn about the known limitations of the application by accessing: [DISCLAIMER](./DISCLAIMER.md).

##  Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Running](#running)
- [Interacting](#interacting)

### Overview
SwapX integrates a decentralized orderbook with Uniswap v4 hooks, replacing the traditional AMM logic with an asynchronous swap system and limit orders. Leveraging EigenLayer's cryptoeconomic security through the Cartesi Coprocessor, it enables swaps to be scheduled, optimized, and auditable, ensuring greater liquidity efficiency and reducing trader uncertainty. This approach eliminates the need for immediate execution, allowing for more sophisticated and flexible strategies for market makers, liquidity protocols, and derivatives.

[![Docs]][Link-docs] [![Deck]][Link-deck]
	
[Docs]: https://img.shields.io/badge/Documentation-6FE1E5?style=for-the-badge
[Link-docs]: https://cartesi.io

[Deck]: https://img.shields.io/badge/Deck-868380?style=for-the-badge
[Link-deck]: https://cartesi.io

### Architecture
![image](https://github.com/user-attachments/assets/f12834c7-769f-4f60-b714-06690cd74f62)

> 1 - The [`SwapXHook.sol`](https://github.com/henriquemarlon/swapx/blob/main/contracts/src/SwapXHook.sol) implementation is a Uniswap hook based on AsyncSwap[^1]. Instead of executing the swap at the market price, it implements a custom logic for [**limit orders**](https://www.investopedia.com/terms/l/limitorder.asp):
>
>   - The order value is transferred to the hook upon creation;
>   - The user can cancel the order and receive the funds back;
>   - When an order is created, a task is issued to the SwapX order book, which will efficiently and intelligently match orders, including aggregating multiple orders and ensuring that the best orders are matched with the incoming order.

> 2 - The **Assets** in this case are token contracts that will be transacted in the swap between users through the pool and contracts that are part of the [**UniswapV4 SDK**](https://docs.uniswap.org/contracts/v4/overview).

> 3 - The **Operator**, which is part of the **Cartesi Coprocessor**, operates under the **crypto-economic security** of the [**EigenLayer restaking protocol**](https://docs.eigenlayer.xyz/eigenlayer/overview). This gives it the ability to perform operations with guarantees of the computation performed, while also having [**"skin in the game"**](https://docs.eigenlayer.xyz/eigenlayer/concepts/slashing/slashing-concept) through slashing penalties in case of malicious behavior.

> 4 - The **Cartesi Coprocessor** is an EigenLayer AVS that operates through a network of operators, leveraging the runtime provided by the **Cartesi Machine**. It is triggered when a new TaskIssued(bytes32, bytes, address) event is emitted. To learn more, visit: https://docs.mugen.builders/cartesi-co-processor-tutorial/introduction.

> 5 - The **Base Layer Access** is enabled through domain `0x27`, present in the GIO, which allows calls via[`eth_getStorageAt`](https://www.quicknode.com/docs/ethereum/eth_getStorageAt). Based on this, the application can access the base layer (read-only), and this is how previously created orders (up to the previous block) are loaded into the application to be processed within the order book. The request specification is as follows:  
>   ```json
>   {"domain": "0x27", "id": "0x<block_hash:32_bytes><address:20_bytes><storage_slot:32_bytes>"}
>   ```

###  Prerequisites

1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

    To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
    docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download)

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps. To install it, run:

   ```shell
   npm i -g @cartesi/cli
   ```

4. [Install the Cartesi Coprocessor CLI](https://docs.mugen.builders/cartesi-co-processor-tutorial/installation)

###  Running
> [!WARNING]
> Before proceeding, make sure that the variable `ETHEREUM_ENDPOINT: http://anvil:8545` is set in the environment of the **operator service** inside the [**./third_party/cartesi-coprocessor/docker-compose-devnet.yaml**](./third_party/cartesi-coprocessor/docker-compose-devnet.yaml) file.

1. Start the devnet coprocessor infrastructure:

   ```bash
   make infra
   ```

2. Build and Publish the application:

   ```sh
   cartesi-coprocessor publish --network devnet
   ```
   
> [!WARNING]
> Before the next step, create a `.env` with the command bellow:
> ```bash
> make env
> ```
> This should look like:
> ```env
> RPC_URL=http://localhost:8545
> PRIVATE_KEY=0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba
> ``

3. Deploy [UniswapV4](https://docs.uniswap.org/contracts/v4/overview) contracts, `SwapXHook.sol` and `SwapXTaskManager.sol` contracts:`

- 3.1 Deploy UniswapV4 contracts:
   
   ```sh
   make v4
   ```

> [!NOTE]
> The following step requires some extra information provided by the command bellow:
> ```bash
> cartesi-coprocessor address-book
> ```
> Output sample:
> ```bash
> Machine Hash         0xdb1d7833f57f79c379e01b97ac5a398da31df195b1901746523be0bc348ccc88
> Devnet_task_issuer   0x95401dc811bb5740090279Ba06cfA8fcF6113778
> Testnet_task_issuer  0xff35E413F5e22A9e1Cc02F92dcb78a5076c1aaf3
> payment_token        0xc5a5C42992dECbae36851359345FE25997F5C42d
> ```

- 3.2 Deploy `SwapXHook.sol` and `SwapXTaskManager.sol`:

   ```bash
   make hook
   ```

   Output sample:

  ```bash
   [⠊] Compiling...
   No files changed, compilation skipped
   Enter Coprocessor address: <devnet_task_issuer>
   Enter Machine Hash: <machine_hash>
  ```

### Interacting

> [!IMPORTANT] 
> If the previous steps were followed precisely, specifically the one that sets up the local infrastructure, accessing [http://localhost:5100](http://localhost:5100) will present you with a block explorer where you can monitor the transactions occurring on the contract of interest. In this project, that contract is the one that implements the Uniswap hook via AsyncSwap[^1]. After that, just search for the contract using that address on Otterscan.

```bash
make demo
```


> [!NOTE]
> You should observe, after a while, four calls targeting the signature method 0x7417ccfb, each covering one of the following scenarios:
> 
> |        | Scenario                                       | Description                                                                                                |
> |--------|----------------------------------------------|------------------------------------------------------------------------------------------------------------|
> | 0      | Buy Order Fulfilled by One Sell Order       | A buy order is completely fulfilled by a single sell order.                                                |
> | 1      | Buy Order Fulfilled by Multiple Sell Orders | A buy order is completely matched by a combination of multiple sell orders.                                |
> | 2      | Sell Order Fulfilled by One Buy Order       | A sell order is completely fulfilled by a single buy order.                                                |
> | 3      | Sell Order Fulfilled by Multiple Buy Orders | A sell order is completely matched by a combination of multiple buy orders.                                |
> | 4      | Buy Order Partially Fulfilled by Multiple Sell Orders | A buy order is only partially matched, while multiple sell orders fulfill part of it.                     |

You can see even more details by accessing the logs tab of one of these transactions, and you'll come across something like this:

![logs](https://github.com/user-attachments/assets/d2161550-aa96-41b2-bb13-0e5ebe457ea3)

[^1]: You can see [here](https://docs.uniswap.org/contracts/v4/quickstart/hooks/async-swap#Configure-a-AsyncSwap-Hook) the reference for enabling the AsyncSwap in a UniswapV4 hook, and [here](https://github.com/henriquemarlon/swapx/blob/demo/contracts/src/SwapXHook.sol#L109) is where it was defined within the application.