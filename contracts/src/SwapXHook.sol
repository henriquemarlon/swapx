// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {BaseAsyncSwap} from "OpenZeppelin/uniswap-hooks/base/BaseAsyncSwap.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";
import {CurrencySettler} from "OpenZeppelin/uniswap-hooks/utils/CurrencySettler.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {BeforeSwapDelta, BeforeSwapDeltaLibrary, toBeforeSwapDelta} from "v4-core/src/types/BeforeSwapDelta.sol";
import {ISwapXHook, ISwapXTaskManager} from "./interface/ISwapXHook.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";

contract SwapXHook is ISwapXHook, BaseAsyncSwap {
    using SafeCast for uint256;
    using CurrencySettler for Currency;

    bytes constant ZERO_BYTES = new bytes(0);

    PoolKey public poolKey;
    Currency public currency0;
    Currency public currency1;
    ISwapXTaskManager public swapXTaskManager;

    struct Order {
        address account;
        uint256 sqrtPrice;
        uint256 amount;
        uint256 matchedAmount;
    }

    mapping(uint256 => bool) public buyOrderCancelled;
    mapping(uint256 => bool) public sellOrderCancelled;

    Order[] public buyOrders;
    Order[] public sellOrders;

    event OrderCreated(uint256 indexed orderId, address indexed account, uint256 sqrtPrice, uint256 amount, bool isBuy);
    event OrderFulfilled(
        uint256 indexed orderId, address indexed account, uint256 sqrtPrice, uint256 amount, bool isBuy
    );
    event OrderCancelled(
        uint256 indexed orderId, address indexed account, uint256 sqrtPrice, uint256 amount, bool isBuy
    );
    event OrderPartiallyFulfilled(
        uint256 indexed orderId, address indexed account, uint256 sqrtPrice, uint256 amount, bool isBuy
    );

    error OrderWasCancelled();
    error OrderDoesNotExist();
    error OrderAlreadyFulfilled();
    error OnlyOrderCreatorCanCancel();
    error OrderSqrtPricesDoNotMatch();

    constructor(IPoolManager _poolManager, ISwapXTaskManager _swapXTaskManager) BaseAsyncSwap(_poolManager) {
        swapXTaskManager = _swapXTaskManager;
    }

    function _beforeInitialize(address, PoolKey calldata key, uint160) internal virtual override returns (bytes4) {
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
        if (params.amountSpecified < 0 && hookData.length > 0) {
            Currency specified = params.zeroForOne ? key.currency0 : key.currency1;
            uint256 specifiedAmount = uint256(-params.amountSpecified);

            specified.take(poolManager, address(this), specifiedAmount, false);

            (uint256 sqrtPrice, address sender) = abi.decode(hookData, (uint256, address));

            Order memory order = Order(sender, sqrtPrice, specifiedAmount, 0);

            uint256 orderId = params.zeroForOne ? buyOrders.length : sellOrders.length;
            if (params.zeroForOne) {
                buyOrders.push(order);
            } else {
                sellOrders.push(order);
            }

            swapXTaskManager.createTask(abi.encode(orderId, sqrtPrice, specifiedAmount, params.zeroForOne ? 0 : 1));
            emit OrderCreated(orderId, sender, sqrtPrice, specifiedAmount, params.zeroForOne);

            return (this.beforeSwap.selector, toBeforeSwapDelta(specifiedAmount.toInt128(), 0), 0);
        }

        return (this.beforeSwap.selector, BeforeSwapDeltaLibrary.ZERO_DELTA, 0);
    }

    function getHookPermissions() public pure virtual override returns (Hooks.Permissions memory permissions) {
        return Hooks.Permissions({
            beforeInitialize: true,
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
        if (buyOrderId >= buyOrders.length || sellOrderId >= sellOrders.length) revert OrderDoesNotExist();
        if (buyOrderCancelled[buyOrderId] || sellOrderCancelled[sellOrderId]) revert OrderWasCancelled();

        Order storage buyOrder = buyOrders[buyOrderId];
        Order storage sellOrder = sellOrders[sellOrderId];

        if (buyOrder.matchedAmount == buyOrder.amount || sellOrder.matchedAmount == sellOrder.amount) {
            revert OrderAlreadyFulfilled();
        }

        if (buyOrder.sqrtPrice < sellOrder.sqrtPrice) revert OrderSqrtPricesDoNotMatch();

        uint256 tradeAmount = buyOrder.amount < sellOrder.amount ? buyOrder.amount : sellOrder.amount;

        currency1.transfer(buyOrder.account, tradeAmount);
        currency0.transfer(sellOrder.account, tradeAmount);

        buyOrder.matchedAmount += tradeAmount;
        sellOrder.matchedAmount += tradeAmount;

        if (buyOrder.matchedAmount == buyOrder.amount) {
            emit OrderFulfilled(buyOrderId, buyOrder.account, buyOrder.sqrtPrice, tradeAmount, true);
        } else {
            emit OrderPartiallyFulfilled(buyOrderId, buyOrder.account, buyOrder.sqrtPrice, tradeAmount, true);
        }

        if (sellOrder.matchedAmount == sellOrder.amount) {
            emit OrderFulfilled(sellOrderId, sellOrder.account, sellOrder.sqrtPrice, tradeAmount, false);
        } else {
            emit OrderPartiallyFulfilled(sellOrderId, sellOrder.account, sellOrder.sqrtPrice, tradeAmount, false);
        }
    }

    function cancelBuyOrder(uint256 orderId) public {
        if (orderId >= buyOrders.length) revert OrderDoesNotExist();
        if (buyOrderCancelled[orderId]) revert OrderWasCancelled();

        Order storage order = buyOrders[orderId];

        if (order.account != msg.sender) revert OnlyOrderCreatorCanCancel();
        if (order.matchedAmount == order.amount) revert OrderAlreadyFulfilled();

        buyOrderCancelled[orderId] = true;

        uint256 remainingAmount = order.amount - order.matchedAmount;
        if (remainingAmount > 0) {
            currency0.transfer(order.account, remainingAmount);
        }

        emit OrderCancelled(orderId, order.account, order.sqrtPrice, order.amount, true);
    }

    function cancelSellOrder(uint256 orderId) public {
        if (orderId >= sellOrders.length) revert OrderDoesNotExist();
        if (sellOrderCancelled[orderId]) revert OrderWasCancelled();

        Order storage order = sellOrders[orderId];

        if (order.account != msg.sender) revert OnlyOrderCreatorCanCancel();
        if (order.matchedAmount == order.amount) revert OrderAlreadyFulfilled();

        sellOrderCancelled[orderId] = true;

        uint256 remainingAmount = order.amount - order.matchedAmount;
        if (remainingAmount > 0) {
            currency1.transfer(order.account, remainingAmount);
        }

        emit OrderCancelled(orderId, order.account, order.sqrtPrice, order.amount, false);
    }
}
