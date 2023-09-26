/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestAddSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockDepositOfferDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockSig := dto.AddSignatureArgs{
		DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
		Address:        "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
		Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
	}

	offerID, err := ids.FromString("TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES")
	require.NoError(t, err)
	addr, err := ids.ShortFromString("7Sdex3LTEjsnswW38Eb48hQ9insctGrsN")
	require.NoError(t, err)
	addr2, err := ids.ShortFromString("6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV")
	require.NoError(t, err)
	offer := &platformvm.APIDepositOffer{
		ID:           offerID,
		OwnerAddress: addr,
	}
	offer2 := &platformvm.APIDepositOffer{
		ID:           offerID,
		OwnerAddress: addr2,
	}
	// first time return mock
	mockDao.EXPECT().AddSignature(mockSig.DepositOfferID, mockSig.Address, mockSig.Signature).Return(nil).Times(1)
	mockNodeService.EXPECT().GetAllDepositOffers(&platformvm.GetAllDepositOffersArgs{}).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{offer}}, nil).Times(1)
	mockNodeService.EXPECT().GetAllDepositOffers(&platformvm.GetAllDepositOffersArgs{}).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{offer2}}, nil).Times(1)
	mockNodeService.EXPECT().GetAllDepositOffers(&platformvm.GetAllDepositOffersArgs{}).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{}}, nil).AnyTimes()

	tests := []struct {
		name string
		args *dto.AddSignatureArgs
		err  error
	}{
		{
			name: "Add signature - success",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
				Address:        "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
			},
			err: nil,
		},
		{
			name: "Invalid signature",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
				Address:        "6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV",
				Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
			},
			err: ErrInvalidSignature,
		},
		{
			name: "Deposit offer not found",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
				Address:        "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
			},
			err: ErrDepositOfferNotFound,
		},
		{
			name: "Invalid deposit offer id",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "invalidid",
				Address:        "6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV",
				Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
			},
			err: ErrParsingDepositOfferID,
		},
		{
			name: "Invalid addr",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
				Address:        "s6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV",
				Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
			},
			err: ErrParsingAddress,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDepositOfferService(mockConfig, mockDao, mockNodeService)
			err := s.AddSignature(tt.args)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
func TestGetSignatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockDepositOfferDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockSig := model.DepositOfferSig{
		DepositOfferID: "TtF4d2QWbk5vzQGTEPrN48x6vwgAoAmKQ9cbp79inpQmcRKES",
		Address:        "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
		Signature:      "3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201",
	}

	// first time return mock
	mockDao.EXPECT().GetSignatures(mockSig.Address).Return(&[]model.DepositOfferSig{mockSig}, nil).Times(1)
	// second time return empty to simulate complete tx for address
	mockDao.EXPECT().GetSignatures(mockSig.Address).Return(&[]model.DepositOfferSig{}, nil).Times(1)

	type args struct {
		address   string
		timestamp string
		signature string
	}
	tests := []struct {
		name string
		args args
		want *[]model.DepositOfferSig
		err  error
	}{
		{
			name: "Get all by address - 1 match",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1695475705",
				signature: "f08074b1d7f3a379ffd6702f82e8d736a4fff13f162fca721aea3910e40fed2c2c974afd89cfd4c19d427a32e5a67ffe23901cdd9c30804ae8cb412e2be5599300",
			},
			want: &[]model.DepositOfferSig{mockSig},
			err:  nil,
		},
		{
			name: "Get all by address - 0 matches",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1695475705",
				signature: "f08074b1d7f3a379ffd6702f82e8d736a4fff13f162fca721aea3910e40fed2c2c974afd89cfd4c19d427a32e5a67ffe23901cdd9c30804ae8cb412e2be5599300",
			},
			want: &[]model.DepositOfferSig{},
			err:  nil,
		},
		{
			name: "Invalid address",
			args: args{
				address:   "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
				timestamp: "1678877386",
				signature: "47bf8e8601badef42a1157e07862157ded68fff927bc3809d5abb0d4a7c51cad3e53979193dc7069f73fe3f7b1b9e8a5946a1bd4782a565fe126a627634943dd01",
			},
			want: nil,
			err:  ErrParsingAddress,
		},
		{
			name: "Invalid signature format",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1678877386",
				signature: "invalidsignature",
			},
			want: nil,
			err:  ErrParsingSignature,
		},
		{
			name: "Invalid signature",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1678877386",
				signature: "f18074b1d7f3a379ffd6702f82e8d736a4fff13f162fca721aea3910e40fed2c2c974afd89cfd4c19d427a32e5a67ffe23901cdd9c30804ae8cb412e2be5599300",
			},
			want: nil,
			err:  ErrInvalidSignature,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDepositOfferService(mockConfig, mockDao, mockNodeService)
			got, err := s.GetSignatures(tt.args.address, tt.args.timestamp, tt.args.signature)
			require.ErrorIs(t, err, tt.err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignatures() got = %v, want %v", got, tt.want)
			}
		})
	}
}
