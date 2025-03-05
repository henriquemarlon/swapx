<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>An EVM Linux-powered coprocessor as an orderbook for UniswapV4 Hooks</i>
</div>
<div align="center">
<b>Cartesi Coprocessor orderbook powered by EigenLayer cryptoeconomic security</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/swapx?style=default&logo=opensourceinitiative&logoColor=white&color=79F7FA" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/swapx?style=default&logo=git&logoColor=white&color=868380" alt="last-commit">
</p>

##  Table of Contents

- [Prerequisites](#prerequisites)
- [Running](#running)
- [Interacting](#interacting)
- [Demo](#demo)

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
cd coprocessor
cartesi-coprocessor publish --network devnet
```

3. Deploy `SwapXHook.sol` and `SwapXTaskManager.sol` contracts:

> [!WARNING]
> 
> Before deploy the contract, create a `.env` file like this
> ```bash
> RPC_URL=http://localhost:8545
> PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
> MACHINE_HASH=""
> TASK_ISSUER_ADDRESS=""
> ```
> 
> - You can see the machine hash running `cartesi hash` in the folder `/coprocessor`;
> - You can see the task issuer address for the devnet enviroment running `cartesi-coprocessor address-book`;
   
```sh
make v4
make hook
```

### Interacting

```bash
make buy
```

```bash
make sell
```
