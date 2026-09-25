// Package alliance holds the protobuf types of the x/alliance module, removed
// in v10, so that the historical governance proposals that carry them can
// still be decoded and queried. The module logic is not included.
//
// The generated code comes from github.com/terra-money/alliance v0.4.3
// (Apache-2.0, see LICENSE).
package alliance

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const routerKey = "alliance"

// RegisterInterfaces registers the alliance messages and legacy proposal
// contents found in the historical governance proposals.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAlliance{},
		&MsgUpdateAlliance{},
		&MsgDeleteAlliance{},
	)
	registry.RegisterImplementations((*govv1beta1.Content)(nil),
		&MsgCreateAllianceProposal{},
		&MsgUpdateAllianceProposal{},
		&MsgDeleteAllianceProposal{},
	)
}

var (
	_ govv1beta1.Content = &MsgCreateAllianceProposal{}
	_ govv1beta1.Content = &MsgUpdateAllianceProposal{}
	_ govv1beta1.Content = &MsgDeleteAllianceProposal{}
)

func (m *MsgCreateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgCreateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgCreateAllianceProposal) ProposalRoute() string  { return routerKey }
func (m *MsgCreateAllianceProposal) ProposalType() string   { return "msg_create_alliance_proposal" }
func (m *MsgCreateAllianceProposal) ValidateBasic() error   { return nil }

func (m *MsgUpdateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgUpdateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgUpdateAllianceProposal) ProposalRoute() string  { return routerKey }
func (m *MsgUpdateAllianceProposal) ProposalType() string   { return "msg_update_alliance_proposal" }
func (m *MsgUpdateAllianceProposal) ValidateBasic() error   { return nil }

func (m *MsgDeleteAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgDeleteAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgDeleteAllianceProposal) ProposalRoute() string  { return routerKey }
func (m *MsgDeleteAllianceProposal) ProposalType() string   { return "msg_delete_alliance_proposal" }
func (m *MsgDeleteAllianceProposal) ValidateBasic() error   { return nil }
