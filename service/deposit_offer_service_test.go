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

func TestAddSignatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockDepositOfferDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
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
	sigs := []string{"3dcbee3d459c03c5741a7bbf434688af933e5d63c5ad38df16b08ed86f8e74012db9105d791fa408535b94a14660674938105f02ca22a69596a8f57ba1f448f201"}
	mockSig := dto.AddSignatureArgs{
		DepositOfferID: offerID.String(),
		Addresses:      []string{addr.String()},
		Signatures:     sigs,
	}
	mockMultipleSigs := dto.AddSignatureArgs{
		DepositOfferID: offerID.String(),
		Addresses:      []string{addr.String(), addr2.String()},
		Signatures:     append(sigs, "f2e5662693c3307f8ed970db60e95e45ca544ffed881fa3654a0f5ca508f248e0355f06d55c63289c3d387131976064cdcb3e1ee9e93e11607d5d352c75710fa00"),
	}
	// first time return mock
	mockDao.EXPECT().AddSignatures(mockSig.DepositOfferID, mockSig.Addresses, mockSig.Signatures).Return(nil).Times(1)
	mockDao.EXPECT().AddSignatures(mockMultipleSigs.DepositOfferID, mockMultipleSigs.Addresses, mockMultipleSigs.Signatures).Return(nil).Times(1)
	mockNodeService.EXPECT().GetAllDepositOffers(gomock.Any()).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{offer}}, nil).Times(3)
	mockNodeService.EXPECT().GetAllDepositOffers(gomock.Any()).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{offer2}}, nil).Times(2)
	mockNodeService.EXPECT().GetAllDepositOffers(gomock.Any()).
		Return(&platformvm.GetAllDepositOffersReply{DepositOffers: []*platformvm.APIDepositOffer{}}, nil).AnyTimes()

	tests := []struct {
		name string
		args *dto.AddSignatureArgs
		err  error
	}{
		{
			name: "Add signature - success",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{addr.String()},
				Signatures:     sigs,
			},
			err: nil,
		},
		{
			name: "Add multiple signatures at once - success",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      mockMultipleSigs.Addresses,
				Signatures:     mockMultipleSigs.Signatures,
			},
			err: nil,
		},
		{
			name: "Fail on verifying 1/2 of provided signatures",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{addr.String(), addr2.String()},
				Signatures:     append(sigs, "55f9ab9eb87ea1c6a2865781eaae6d2acbe1c7a497e9213c92b45c7de0bee3db0ef1ec3e6b401c52305c20568290db36af1f4204c7d8672b1c649bd18e0d5ad300"), //invalid 2nd sig
			},
			err: ErrInvalidSignature,
		},
		{
			name: "Mismatch addresses and signatures",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{addr.String(), addr2.String()},
				Signatures:     sigs,
			},
			err: ErrAddressesSigsMismatch,
		},
		{
			name: "Invalid signature",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{addr2.String()},
				Signatures:     sigs,
			},
			err: ErrInvalidSignature,
		},
		{
			name: "Invalid addr",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{"s6Y3kysjF9jnHnYkdS9yGAuoHyae2eNmeV"},
				Signatures:     sigs,
			},
			err: ErrParsingAddress,
		},
		{
			name: "Deposit offer not found",
			args: &dto.AddSignatureArgs{
				DepositOfferID: offerID.String(),
				Addresses:      []string{addr.String()},
				Signatures:     sigs,
			},
			err: ErrDepositOfferNotFound,
		},
		{
			name: "Invalid deposit offer id",
			args: &dto.AddSignatureArgs{
				DepositOfferID: "invalidid",
				Addresses:      []string{addr2.String()},
				Signatures:     sigs,
			},
			err: ErrParsingDepositOfferID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDepositOfferService(mockConfig, mockDao, mockNodeService)
			err := s.AddSignatures(tt.args)
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
	mockSigMultisig := model.DepositOfferSig{
		DepositOfferID: "WixPyBD9vuDNfbNhPcLbQoxYjDZ3bfNa51uLFqTKRVVB536Ng",
		Address:        "9UkT1iQk9kyQ3rvecXBeoknT687t2kApG", // P-kopernikus1t5qgr9hcmf2vxj7k0hz77kawf9yr389cxte5j0
		Signature:      "6198c95b5c9a586a7dd8970f601673b0c224a4fa3f10990e3bd5eddc83269f1f3de5c484c75591597ead7edd023b73ee3c502becfc00d1a3b245a8ee49a8cc2b00",
	}

	// first time return mock
	mockDao.EXPECT().GetSignatures(mockSig.Address).Return(&[]model.DepositOfferSig{mockSig}, nil).Times(1)
	// first time return mock
	mockDao.EXPECT().GetSignatures(mockSigMultisig.Address).Return(&[]model.DepositOfferSig{mockSigMultisig}, nil).Times(1)
	// second time return empty to simulate complete tx for address
	mockDao.EXPECT().GetSignatures(mockSig.Address).Return(&[]model.DepositOfferSig{}, nil).Times(1)

	// first 2 times valid multisig alias
	mockNodeService.EXPECT().GetMultisigAlias(mockSigMultisig.Address).Return(&model.AliasInfo{
		Result: model.Result{
			Addresses: []string{"P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
				"P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3"},
			Threshold: "2",
		},
	}, nil).Times(2)
	// third time invalid multisig alias
	mockNodeService.EXPECT().GetMultisigAlias(mockSig.Address).Return(&model.AliasInfo{
		Result: model.Result{
			Addresses: nil,
		},
	}, nil)

	type args struct {
		address   string
		timestamp string
		signature string
		multisig  bool
	}
	tests := []struct {
		name string
		args args
		want *[]model.DepositOfferSig
		err  error
	}{
		{
			name: "Singlesig - 1 match",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1695475705",
				signature: "f08074b1d7f3a379ffd6702f82e8d736a4fff13f162fca721aea3910e40fed2c2c974afd89cfd4c19d427a32e5a67ffe23901cdd9c30804ae8cb412e2be5599300",
				multisig:  false,
			},
			want: &[]model.DepositOfferSig{mockSig},
			err:  nil,
		},
		{
			name: "Multisig - 1 match",
			args: args{
				address:   "9UkT1iQk9kyQ3rvecXBeoknT687t2kApG",
				timestamp: "1696491289",
				signature: "5c484babcf10cfedf4d0e82a8e72acc60ff3a10aead755ddf09f9cb35e1c82b56168d637fb02ae9373823af732e1f19567d3f728d869a4d54d9bdc425663bff301",
				multisig:  true,
			},
			want: &[]model.DepositOfferSig{mockSigMultisig},
			err:  nil,
		},
		{
			name: "Get all by address - 0 matches",
			args: args{
				address:   "7Sdex3LTEjsnswW38Eb48hQ9insctGrsN",
				timestamp: "1695475705",
				signature: "f08074b1d7f3a379ffd6702f82e8d736a4fff13f162fca721aea3910e40fed2c2c974afd89cfd4c19d427a32e5a67ffe23901cdd9c30804ae8cb412e2be5599300",
				multisig:  false,
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
				multisig:  false,
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
				multisig:  false,
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
				multisig:  false,
			},
			want: nil,
			err:  ErrInvalidSignature,
		},
		{
			name: "Multisig - addr not part of alias",
			args: args{
				address:   mockSigMultisig.Address,
				timestamp: "1696491289",
				signature: "0fc8b75b58d73773e2731d4e6ec697fedd7dd2299f8515aae03ad4f95892498734ec7524ef80804379c751113d1d0f76d069f5eba0e13dc422314143044b10ed01",
				multisig:  true,
			},
			want: nil,
			err:  ErrInvalidSignature,
		},
		{
			name: "Multisig - no alias found",
			args: args{
				address:   mockSig.Address,
				timestamp: "1696491289",
				signature: "5c484babcf10cfedf4d0e82a8e72acc60ff3a10aead755ddf09f9cb35e1c82b56168d637fb02ae9373823af732e1f19567d3f728d869a4d54d9bdc425663bff301",
				multisig:  true,
			},
			want: nil,
			err:  ErrNoAliasFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewDepositOfferService(mockConfig, mockDao, mockNodeService)
			got, err := s.GetSignatures(tt.args.address, tt.args.timestamp, tt.args.signature, tt.args.multisig)
			require.ErrorIs(t, err, tt.err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignatures() got = %v, want %v", got, tt.want)
			}
		})
	}
}
