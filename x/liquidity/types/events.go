package types

// Event types for the liquidity module.
const (
	EventTypeCreatePool          = TypeMsgCreatePool
	EventTypeDepositWithinBatch  = TypeMsgDepositWithinBatch
	EventTypeWithdrawWithinBatch = TypeMsgWithdrawWithinBatch
	EventTypeSwapWithinBatch     = TypeMsgSwapWithinBatch
	EventTypeDirectSwap          = TypeMsgDirectSwap
	EventTypeDepositToPool       = "deposit_to_pool"
	EventTypeWithdrawFromPool    = "withdraw_from_pool"
	EventTypeSwapTransacted      = "swap_transacted"

	AttributeValuePoolId         = "pool_id"      //nolint:golint
	AttributeValuePoolTypeId     = "pool_type_id" //nolint:golint
	AttributeValuePoolName       = "pool_name"
	AttributeValueReserveAccount = "reserve_account"
	AttributeValuePoolCoinDenom  = "pool_coin_denom"
	AttributeValuePoolCoinAmount = "pool_coin_amount"
	AttributeValueBatchIndex     = "batch_index"
	AttributeValueMsgIndex       = "msg_index"

	AttributeValueDepositCoins = "deposit_coins"

	AttributeValueOfferCoinDenom         = "offer_coin_denom"
	AttributeValueOfferCoinAmount        = "offer_coin_amount"
	AttributeValueOfferCoinFeeAmount     = "offer_coin_fee_amount"
	AttributeValueExchangedCoinFeeAmount = "exchanged_coin_fee_amount"
	AttributeValueDemandCoinDenom        = "demand_coin_denom"
	AttributeValueOrderPrice             = "order_price"

	AttributeValueDepositor        = "depositor"
	AttributeValueRefundedCoins    = "refunded_coins"
	AttributeValueAcceptedCoins    = "accepted_coins"
	AttributeValueSuccess          = "success"
	AttributeValueNoMatch          = "no_match"
	AttributeValueWithdrawer       = "withdrawer"
	AttributeValueWithdrawCoins    = "withdraw_coins"
	AttributeValueWithdrawFeeCoins = "withdraw_fee_coins"
	AttributeValueSwapRequester    = "swap_requester"
	AttributeValueSwapTypeId       = "swap_type_id" //nolint:golint
	AttributeValueSwapPrice        = "swap_price"

	AttributeValueTransactedCoinAmount       = "transacted_coin_amount"
	AttributeValueRemainingOfferCoinAmount   = "remaining_offer_coin_amount"
	AttributeValueExchangedOfferCoinAmount   = "exchanged_offer_coin_amount"
	AttributeValueExchangedDemandCoinAmount  = "exchanged_demand_coin_amount"
	AttributeValueReservedOfferCoinFeeAmount = "reserved_offer_coin_fee_amount"
	AttributeValueOrderExpiryHeight          = "order_expiry_height"
	AttributeValueOrderExpired               = "order_has_expired"
	AttributeValueInvalidSwapErr             = "swap_err_message"

	AttributeValueCategory = ModuleName

	Success = "success"
	Failure = "failure"
)
