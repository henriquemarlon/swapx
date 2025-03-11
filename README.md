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
> This is an experimental project under continuous development and should be treated as such. **Its use in production/mainnet is not recommended.** Learn about the known limitations of the application by accessing [DISCLAIMER](./DISCLAIMER.md).

##  Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Running](#running)
- [Interacting](#interacting)

### Overview
WIP

[![Docs]][Link-docs] [![Deck]][Link-deck]
	
[Docs]: https://img.shields.io/badge/Documentation-6FE1E5?style=for-the-badge
[Link-docs]: https://cartesi.io

[Deck]: https://img.shields.io/badge/Deck-868380?style=for-the-badge
[Link-deck]: https://cartesi.io

### Architecture
WIP

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

1. Start the devnet coprocessor infrastructure:

```bash
make infra
```

2. Build and Publish the application:

```sh
cartesi-coprocessor publish --network devnet
```

3. Deploy [UniswapV4](https://docs.uniswap.org/contracts/v4/overview) contracts, `SwapXHook.sol` and `SwapXTaskManager.sol` contracts:

> [!WARNING]
> Before deploy, create a `.env` with the command bellow:
> ```bash
> make env
> ```
> This should look like:
> ```env
> RPC_URL=http://localhost:8545
> PRIVATE_KEY=0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba
> ```

3.1 Deploy UniswapV4 contracts:
   
```sh
make v4
```

3.2 Deploy `SwapXHook.sol` and `SwapXTaskManager.sol`:

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

<div align="justify">
If the previous steps were followed precisely, specifically the one that sets up the local architecture, accessing http://localhost:5100 will present you with a block explorer where you can monitor the transactions occurring on the contract of interest. In this project, that contract is the one that implements the Uniswap hook via asyncSwap. (Hint: To identify the address of the deployed SwapXHook.sol instance, simply access `broadcast/31337/HookDeployer.sol/run-latest.json` and look for the address.) After that, just search for the contract using that address on Otterscan. Then, follow the instructions below.
</div>
<br>

```bash
make demo
```

> [!IMPORTANT]
> You should observe, after a while, five calls targeting the signature method 0x7417ccfb, each covering one of the following scenarios:
> 
> |        | Scenario                                        | Description                                                                                                |
> |--------|-------------------------------------------------|------------------------------------------------------------------------------------------------------------|
> | 0      | Bid Fully Matched by Single Ask                 | A bid order is completely fulfilled by one single ask order.                                               |
> | 1      | Bid Fully Matched by Multiple Asks              | A bid order is completely matched by a combination of multiple ask orders.                                 |
> | 2      | Ask Fully Matched by Single Bid                 | An ask order is completely fulfilled by one single bid order.                                              |
> | 3      | Ask Fully Matched by Multiple Bids              | An ask order is completely matched by a combination of multiple bid orders.                                |
> | 4      | Ask Partially Matched but Bid Fully Fulfilled   | An ask order is only partially filled, while the bid order is fully satisfied.                             |
>

You can see even more details by accessing the logs tab of one of these transactions, and you'll come across something like this:

![logs](https://github.com/user-attachments/assets/d2161550-aa96-41b2-bb13-0e5ebe457ea3)
