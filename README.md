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
![image](https://github.com/user-attachments/assets/8974e20e-49c3-470b-921c-e5abf45234f3)

> 1 - The [`SwapXHook.sol`](https://github.com/henriquemarlon/swapx/blob/main/contracts/src/SwapXHook.sol) implementation is a Uniswap hook based on [`AsyncSwap`](https://docs.uniswap.org/contracts/v4/quickstart/hooks/async-swap). Instead of executing the swap at the market price, it implements a custom logic for [**limit orders**](https://www.investopedia.com/terms/l/limitorder.asp):
>
>   - The order value is transferred to the hook upon creation;
>   - The user can cancel the order and receive the funds back;
>   - When an order is created, a task is issued to the SwapX order book, which will efficiently and intelligently match orders, including aggregating multiple orders and ensuring that the best orders are matched with the incoming order.

> 2 - The assets in this case are token contracts that will be transacted in the swap between users through the pool and contracts that are part of the [**UniswapV4 SDK**](https://docs.uniswap.org/contracts/v4/overview).

> 3 - The **Operator**, which is part of the **Cartesi Coprocessor**, operates under the **crypto-economic security** of the [**EigenLayer restaking protocol**](https://docs.eigenlayer.xyz/eigenlayer/overview). This gives it the ability to perform operations with guarantees of the computation performed, while also having [**"skin in the game"**](https://docs.eigenlayer.xyz/eigenlayer/concepts/slashing/slashing-concept) through slashing penalties in case of malicious behavior.

> 4 - The **Cartesi Coprocessor** is an EigenLayer AVS with a network of operators, leveraging the runtime provided by the **Cartesi Machine**. To learn more, visit: https://docs.mugen.builders/cartesi-co-processor-tutorial/introduction.

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
   Enter Coprocessor address: <Devnet_task_issuer>
   Enter Machine Hash: <Machine Hash>
  ```

### Interacting

If the previous steps were followed precisely, specifically the one that sets up the local architecture, accessing [http://localhost:5100](http://localhost:5100) will present you with a block explorer where you can monitor the transactions occurring on the contract of interest. In this project, that contract is the one that implements the Uniswap hook via asyncSwap. After that, just search for the contract using that address on Otterscan. Then, follow the instructions below:

```bash
make demo
```

You can see even more details by accessing the logs tab of one of these transactions, and you'll come across something like this:

![logs](https://github.com/user-attachments/assets/d2161550-aa96-41b2-bb13-0e5ebe457ea3)
