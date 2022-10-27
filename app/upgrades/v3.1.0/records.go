package v310

// curl -X GET "https://secret-4.api.trivium.network:1317/cosmos/staking/v1beta1/validators/secretvaloper1hscf4cjrhzsea5an5smt4z9aezhh4sf5jjrqka/delegations?pagination.limit=10500&pagination.count_total=true" -H  "x-cosmos-block-height: 5181125" > DAN.JSON
// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.23/365)*9+1) | floor) | tostring},)' DAN.JSON > to_mint.json

// Slash was 5%
// Lost APR is 23% for 9 days
// Chihuahua APR is 42% for x days

var recordsJSONString = `[
	{
		"address": "chihuahua12kupfxepn4lu2kjmed84rvhvjddy3wz6rj2llu",
		"amount": "1000000"
	},
	{
		"address": "chihuahua17nssn54y5nmjl874vcdlprhatstq3995p49hu4",
		"amount": "10000000"
	}
]`
