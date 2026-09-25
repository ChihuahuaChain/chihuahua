package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// the historical governance proposals of chihuahua-1 must stay decodable
func TestHistoricalProposalTypesResolve(t *testing.T) {
	app := Setup(t)
	for _, typeURL := range []string{
		"/alliance.alliance.MsgUpdateAlliance",
		"/alliance.alliance.MsgCreateAllianceProposal",
		"/cosmos.distribution.v1beta1.CommunityPoolSpendProposal",
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
		"/cosmos.gov.v1.MsgUpdateParams",
		"/cosmos.gov.v1beta1.TextProposal",
		"/cosmos.mint.v1beta1.MsgUpdateParams",
		"/cosmos.params.v1beta1.ParameterChangeProposal",
		"/cosmos.staking.v1beta1.MsgUpdateParams",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal",
		"/ibc.core.client.v1.ClientUpdateProposal",
		"/ibc.core.client.v1.MsgRecoverClient",
		"/liquidity.v1beta1.MsgUpdateParams",
		"/osmosis.tokenfactory.v1beta1.MsgUpdateParams",
	} {
		_, err := app.InterfaceRegistry().Resolve(typeURL)
		require.NoError(t, err, typeURL)
	}
}
