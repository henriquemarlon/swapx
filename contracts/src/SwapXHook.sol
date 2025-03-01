// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {BaseAsyncSwap} from "OpenZeppelin/uniswap-hooks/base/BaseAsyncSwap.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";
import {CurrencySettler} from "OpenZeppelin/uniswap-hooks/utils/CurrencySettler.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {BeforeSwapDelta, BeforeSwapDeltaLibrary, toBeforeSwapDelta} from "v4-core/src/types/BeforeSwapDelta.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";
import {ISwapXHook, ISwapXTaskManager} from "./interface/ISwapXHook.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";

contract SwapXHook is ISwapXHook, BaseAsyncSwap {
    using CurrencySettler for Currency;
    using SafeCast for uint256;

    bytes constant ZERO_BYTES = new bytes(0);

    Currency public currency0;
    Currency public currency1;

    PoolKey public poolKey;

    ISwapXTaskManager public swapXTaskManager;

    struct Order {
        address account;
        uint256 sqrtPrice;
        uint256 amount;
    }

    // Buy order -> zeroForOne == true (deposits currency0, withdraws currency1)
    // Sell order -> zeroForOne == false (deposits currency1, withdraws currency0)

    mapping(uint256 => bool) public buyOrderFulfilledOrCancelled; // buyOrderId => isFulfilled (buy = zeroForOne == true)
    mapping(uint256 => bool) public sellOrderFulfilledOrCancelled; // sellOrderId => isFulfilled (sell = zeroForOne == false)

    Order[] public buyOrders; // (buy = zeroForOne == true)
    Order[] public sellOrders; // (sell = zeroForOne == false)

    //events

    event BuyOrderCreated(
        uint256 indexed buyOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );
    event SellOrderCreated(
        uint256 indexed sellOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );

    event BuyOrderFulfilled(
        uint256 indexed buyOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );
    event SellOrderFulfilled(
        uint256 indexed sellOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );

    event BuyOrderPartiallyFulfilled(
        uint256 indexed buyOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );
    event SellOrderPartiallyFulfilled(
        uint256 indexed sellOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );

    event BuyOrderCancelled(
        uint256 indexed buyOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );
    event SellOrderCancelled(
        uint256 indexed sellOrderId,
        address indexed account,
        uint256 sqrtPrice,
        uint256 amount
    );

    // errors

    error OrderDoesNotExist();
    error OrderAlreadyFulfilledOrCancelled();
    error OnlyOrderCreatorCanCancel();
    error OrderAmountsDoNotMatch();
    error OrderSqrtPricesDoNotMatch();

    constructor(
        IPoolManager _poolManager,
        ISwapXTaskManager _swapXTaskManager
    ) BaseAsyncSwap(_poolManager) {
        swapXTaskManager = _swapXTaskManager;
    }

    function _beforeInitialize(
        address,
        PoolKey calldata key,
        uint160
    ) internal virtual override returns (bytes4) {
        currency0 = key.currency0;
        currency1 = key.currency1;
        return this.beforeInitialize.selector;
    }

    function _beforeSwap(
        address,
        PoolKey calldata key,
        IPoolManager.SwapParams calldata params,
        bytes calldata hookData
    ) internal virtual override returns (bytes4, BeforeSwapDelta, uint24) {
        // Async swaps are only possible on exact-input swaps, so exact-output swaps are executed by the `PoolManager` as normal
        if (params.amountSpecified < 0 && hookData.length > 0) {
            // Determine which currency is specified
            Currency specified = params.zeroForOne
                ? key.currency0
                : key.currency1;

            // Get the positive specified amount
            uint256 specifiedAmount = uint256(-params.amountSpecified);

            // Mint ERC-6909 claim token for the specified currency and amount
            specified.take(poolManager, address(this), specifiedAmount, false);

            (uint256 sqrtPrice, address sender) = abi.decode(
                hookData,
                (uint256, address)
            );

            Order memory order = Order({
                account: sender,
                sqrtPrice: sqrtPrice,
                amount: specifiedAmount
            });

            if (params.zeroForOne) {
                emit BuyOrderCreated(
                    buyOrders.length,
                    sender,
                    sqrtPrice,
                    specifiedAmount
                );
                buyOrders.push(order);
                swapXTaskManager.createTask(
                    abi.encode(
                        buyOrders.length,
                        order.sqrtPrice,
                        order.amount,
                        uint256(0)
                    )
                );
            } else {
                emit SellOrderCreated(
                    sellOrders.length,
                    sender,
                    sqrtPrice,
                    specifiedAmount
                );
                sellOrders.push(order);
                swapXTaskManager.createTask(
                    abi.encode(
                        buyOrders.length,
                        order.sqrtPrice,
                        order.amount,
                        uint256(1)
                    )
                );
            }

            // Return delta that nets out specified amount to 0.
            return (
                this.beforeSwap.selector,
                toBeforeSwapDelta(specifiedAmount.toInt128(), 0),
                0
            );
        } else {
            return (
                this.beforeSwap.selector,
                BeforeSwapDeltaLibrary.ZERO_DELTA,
                0
            );
        }
    }

    function getHookPermissions()
        public
        pure
        virtual
        override
        returns (Hooks.Permissions memory permissions)
    {
        return
            Hooks.Permissions({
                beforeInitialize: true, // adding beforeInitialize to the hook permissions
                afterInitialize: false,
                beforeAddLiquidity: false,
                beforeRemoveLiquidity: false,
                afterAddLiquidity: false,
                afterRemoveLiquidity: false,
                beforeSwap: true,
                afterSwap: false,
                beforeDonate: false,
                afterDonate: false,
                beforeSwapReturnDelta: true,
                afterSwapReturnDelta: false,
                afterAddLiquidityReturnDelta: false,
                afterRemoveLiquidityReturnDelta: false
            });
    }

    function executeAsyncSwap(uint256 buyOrderId, uint256 sellOrderId) public {
        if (
            buyOrderId >= buyOrders.length || sellOrderId >= sellOrders.length
        ) {
            revert OrderDoesNotExist();
        }

        if (
            buyOrderFulfilledOrCancelled[buyOrderId] ||
            sellOrderFulfilledOrCancelled[sellOrderId]
        ) {
            revert OrderAlreadyFulfilledOrCancelled();
        }

        Order storage buyOrder = buyOrders[buyOrderId];
        Order storage sellOrder = sellOrders[sellOrderId];

        if (buyOrder.sqrtPrice < sellOrder.sqrtPrice) {
            revert OrderSqrtPricesDoNotMatch();
        }

        uint256 tradeAmount = buyOrder.amount < sellOrder.amount
            ? buyOrder.amount
            : sellOrder.amount;

        currency1.transfer(buyOrder.account, tradeAmount);
        currency0.transfer(sellOrder.account, tradeAmount);

        if (tradeAmount == buyOrder.amount) {
            buyOrderFulfilledOrCancelled[buyOrderId] = true;
            emit BuyOrderFulfilled(
                buyOrderId,
                buyOrder.account,
                buyOrder.sqrtPrice,
                tradeAmount
            );
        } else {
            buyOrder.amount -= tradeAmount;
            emit BuyOrderPartiallyFulfilled(
                buyOrderId,
                buyOrder.account,
                buyOrder.sqrtPrice,
                tradeAmount
            );
        }

        if (tradeAmount == sellOrder.amount) {
            sellOrderFulfilledOrCancelled[sellOrderId] = true;
            emit SellOrderFulfilled(
                sellOrderId,
                sellOrder.account,
                sellOrder.sqrtPrice,
                tradeAmount
            );
        } else {
            sellOrder.amount -= tradeAmount;
            emit SellOrderPartiallyFulfilled(
                sellOrderId,
                sellOrder.account,
                sellOrder.sqrtPrice,
                tradeAmount
            );
        }
    }

    function cancelBuyOrder(uint256 orderId) public {
        if (orderId >= buyOrders.length) {
            revert OrderDoesNotExist();
        }

        Order memory buyOrder = buyOrders[orderId];

        if (buyOrder.account != msg.sender) {
            revert OnlyOrderCreatorCanCancel();
        }

        if (buyOrderFulfilledOrCancelled[orderId]) {
            revert OrderAlreadyFulfilledOrCancelled();
        }

        buyOrderFulfilledOrCancelled[orderId] = true;

        currency0.transfer(buyOrder.account, buyOrder.amount);

        emit BuyOrderCancelled(
            orderId,
            buyOrder.account,
            buyOrder.sqrtPrice,
            buyOrder.amount
        );
    }

    function cancelSellOrder(uint256 orderId) public {
        if (orderId >= sellOrders.length) {
            revert OrderDoesNotExist();
        }

        Order memory sellOrder = sellOrders[orderId];

        if (sellOrder.account != msg.sender) {
            revert OnlyOrderCreatorCanCancel();
        }

        if (sellOrderFulfilledOrCancelled[orderId]) {
            revert OrderAlreadyFulfilledOrCancelled();
        }

        sellOrderFulfilledOrCancelled[orderId] = true;

        currency1.transfer(sellOrder.account, sellOrder.amount);

        emit SellOrderCancelled(
            orderId,
            sellOrder.account,
            sellOrder.sqrtPrice,
            sellOrder.amount
        );
    }
}
