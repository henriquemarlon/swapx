<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>EVM Linux Coprocessor as an orderbook for UniswapV4 Hooks</i>
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

> [!NOTE]
> WIP

### Demo
WIP



| **Scenario** | **Description** | **Output** |
|-------------|---------------|----------|
| **4. Ask fully matched by a single Bid** | An ask order is fully matched by a single bid at the same price. | `sellToBuy` populated, `buyToSell = nil` |
| **5. Ask fully matched by multiple Bids (prioritizing larger orders)** | An ask order is matched by multiple bids at the same price, prioritizing the larger quantity first. | `sellToBuy` populated, `buyToSell = nil` |
| **6. Ask partially matched (remaining amount persists)** | An ask order is partially matched, but some amount remains in the order book. | `sellToBuy` populated, `buyToSell = nil` |
| **8. Ask arrives but finds no matching Bid** | An ask order is placed, but no bids are available to match. | `nil, nil` |
| **10. Ask finds multiple Bids at different prices (prioritizing higher prices first)** | An ask is matched with bids at different price levels, prioritizing the highest price first. | `sellToBuy` populated, `buyToSell = nil` |
| **11. Larger order included in match even if only partially filled** | A bid or ask with a large amount is partially matched, but it is still included in the match list. | `buyToSell` or `sellToBuy` populated, with remaining amount in the order book. |

map[string][]string -> map[uint64][]uint64