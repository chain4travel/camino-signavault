/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"github.com/ava-labs/avalanchego/ids"
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
	"time"
)

const networkId = uint32(1002)

func TestCreateMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)

	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	unsignedTx := "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	id := "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28"

	alias := "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy"
	mockAliasInfo := &model.AliasInfo{
		Jsonrpc: "2.0",
		Result: model.Result{
			Memo:      "0x",
			Addresses: []string{"P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68", "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3"},
			Threshold: "2",
		},
		Id: 1,
	}

	now := time.Now().UTC().Round(time.Second).Add(time.Hour * 24 * 14)
	mockTx := model.MultisigTx{
		Id:            id,
		UnsignedTx:    unsignedTx,
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		OutputOwners:  "OutputOwners",
		Metadata:      "",
		Expiration:    &now,
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: id,
				Address:      mockAliasInfo.Result.Addresses[0],
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
			{
				MultisigTxId: id,
				Address:      mockAliasInfo.Result.Addresses[1],
				Signature:    "",
			},
		},
	}

	mockDao.EXPECT().CreateMultisigTx(&mockTx).Return(mockTx.Id, nil)
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).AnyTimes()
	mockDao.EXPECT().PendingAliasExists("P-kopernikus1fq0jc8svlyazhygkj0s36qnl6s0km0h3uuc99e", "11111111111111111111111111111111LpoYY").Return(true, nil)
	mockDao.EXPECT().PendingAliasExists(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()
	mockNodeService.EXPECT().GetMultisigAlias(alias).Return(mockAliasInfo, nil)
	mockNodeService.EXPECT().GetMultisigAlias(gomock.Any()).Return(nil, errAliasInfoNotFound)

	type args struct {
		multisigTx *dto.MultisigTxArgs
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Alias with 2 owners",
			args: args{
				multisigTx: &dto.MultisigTxArgs{
					Alias:        alias,
					UnsignedTx:   mockTx.UnsignedTx,
					Signature:    mockTx.Owners[0].Signature,
					OutputOwners: mockTx.OutputOwners,
					Expiration:   mockTx.Expiration.Unix(),
				},
			},
			err: nil,
		},
		{
			name: "Wrong alias - no node info",
			args: args{
				multisigTx: &dto.MultisigTxArgs{
					Alias:        "P-kopernikus1fq0jc8svlyazhygkj0s36qnl6s0km0h3uuc99w",
					UnsignedTx:   mockTx.UnsignedTx,
					Signature:    mockTx.Owners[0].Signature,
					OutputOwners: mockTx.OutputOwners,
				},
			},
			err: errAliasInfoNotFound,
		},
		{
			name: "Duplicate alias",
			args: args{
				multisigTx: &dto.MultisigTxArgs{
					Alias:        "P-kopernikus1fq0jc8svlyazhygkj0s36qnl6s0km0h3uuc99e",
					UnsignedTx:   mockTx.UnsignedTx,
					Signature:    mockTx.Owners[0].Signature,
					OutputOwners: mockTx.OutputOwners,
				},
			},
			err: ErrPendingTx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)

			_, err := s.CreateMultisigTx(tt.args.multisigTx)

			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)

			}

		})
	}
}

func TestGetAllMultisigTxForAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockTx := model.MultisigTx{
		Id:            "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	// first time return mock
	mockDao.EXPECT().GetMultisigTx("", mockTx.Alias, mockTx.Owners[0].Address).Return(&[]model.MultisigTx{mockTx}, nil).Times(1)
	// second time return empty to simulate complete tx for alias
	mockDao.EXPECT().GetMultisigTx("", mockTx.Alias, mockTx.Owners[0].Address).Return(&[]model.MultisigTx{}, nil).Times(1)

	type args struct {
		alias     string
		timestamp string
		signature string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]model.MultisigTx
		wantErr bool
	}{
		{
			name: "Get all by alias",
			args: args{
				alias:     "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
				timestamp: "1678877386",
				signature: "47bf8e8601badef42a1157e07862157ded68fff927bc3809d5abb0d4a7c51cad3e53979193dc7069f73fe3f7b1b9e8a5946a1bd4782a565fe126a627634943dd01",
			},
			want:    &[]model.MultisigTx{mockTx},
			wantErr: false,
		},
		{
			name: "Get all with completed tx for alias",
			args: args{
				alias:     "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
				timestamp: "1678877386",
				signature: "47bf8e8601badef42a1157e07862157ded68fff927bc3809d5abb0d4a7c51cad3e53979193dc7069f73fe3f7b1b9e8a5946a1bd4782a565fe126a627634943dd01",
			},
			want:    &[]model.MultisigTx{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)
			got, err := s.GetAllMultisigTxForAlias(tt.args.alias, tt.args.timestamp, tt.args.signature)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMultisigTxForAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMultisigTxForAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockTx := model.MultisigTx{
		Id:            "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	// first time return mock
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).Times(1)
	mockDao.EXPECT().GetMultisigTx(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.MultisigTx
		wantErr bool
	}{
		{
			name: "Get multisig by id",
			args: args{
				id: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
			},
			want:    &mockTx,
			wantErr: false,
		},
		{
			name: "Get multisig by non existing id",
			args: args{
				id: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9999",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)
			got, err := s.GetMultisigTx(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockTx := model.MultisigTx{
		Id:            "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
			},
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	mockTxWithSigner := model.MultisigTx{
		Id:            "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d29",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d29",
				Address:      "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d29",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	// mock without signer
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).AnyTimes()
	mockDao.EXPECT().AddSigner(mockTx.Id, "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000", mockTx.Owners[0].Address).Return(true, nil).AnyTimes()
	// mock with existing signer
	mockDao.EXPECT().GetMultisigTx(mockTxWithSigner.Id, "", "").Return(&[]model.MultisigTx{mockTxWithSigner}, nil).AnyTimes()
	mockDao.EXPECT().AddSigner(mockTxWithSigner.Id, "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000", mockTx.Owners[0].Address).Return(false, nil).AnyTimes()

	type args struct {
		id       string
		signArgs *dto.SignTxArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *model.MultisigTx
		wantErr bool
	}{
		{
			name: "Sign multisig tx",
			args: args{
				id: mockTx.Id,
				signArgs: &dto.SignTxArgs{
					Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
				},
			},
			want:    &mockTx,
			wantErr: false,
		},
		{
			name: "Sign multisig tx with existing signature",
			args: args{
				id: mockTxWithSigner.Id,
				signArgs: &dto.SignTxArgs{
					Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)
			got, err := s.SignMultisigTx(tt.args.id, tt.args.signArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockTx := model.MultisigTx{
		Id:            "b62c43e3522eec9891723220785711274a979f295a11dd58156080ea462db5ac",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "b62c43e3522eec9891723220785711274a979f295a11dd58156080ea462db5ac",
				Address:      "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
			},
			{
				MultisigTxId: "b62c43e3522eec9891723220785711274a979f295a11dd58156080ea462db5ac",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	// mock without signer
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).AnyTimes()
	mockDao.EXPECT().UpdateTransactionId(mockTx.Id, gomock.Any()).Return(true, nil).AnyTimes()
	txId, _ := ids.FromString("3N3j8FpRtvx9UAJrsS6CTcsUQPCmRqf4Hjnfp81CuEJSMcqJ2")
	mockNodeService.EXPECT().IssueTx(gomock.Any()).Return(txId, nil).AnyTimes()

	type args struct {
		issueArgs *dto.IssueTxArgs
	}
	tests := []struct {
		name    string
		args    args
		want    ids.ID
		wantErr bool
	}{
		{
			name: "Issue multisig tx",
			args: args{
				issueArgs: &dto.IssueTxArgs{
					SignedTx:  "000000002007000003ea00000000000000000000000000000000000000000000000000000000000000000000000159eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a0000000700016bcc41d9bdc0000000000000000000000001000000015d008196f8da54c34bd67dc5ef5bae4948389cb8000000010903208c79e9d29ad5e5ea7caf771ecca4db7a218c44d7c3619deea62e6227640000000359eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a0000000500016bcc41e9000000000002000000000000000100000000000000000000000000000000000000000000000083b1ddd7b166dbe6305c22fed5f59065525c4e510000000a00000001000000005d008196f8da54c34bd67dc5ef5bae4948389cb8000000030000200c00000002dd3be02c98a8d121e6a0e3bb123117db44bfc0ec78cc73e5a0b87a92afccd6d71d4f952bba5ff34defd3626cd3b3c86816c384f4f2c5241a75393da4b77572b1006e19b48ad5ab9ed3e7d774bef7aae5d9047b773075c7372a3736022f7064e66a32a567eb5112d32061622a1cfd33ff4076579a0ab962fad9547816c095277d6e010000000200000000000000010000200c00000001a32fc319922bf20632f85f5c99c3ecdf88387cf28564452403e81d635c805c736d2217e83cc33dea85311039a2745fc4bcd28f4b822799b819dc858a6391dd810100000001000000000000200c00000002dd3be02c98a8d121e6a0e3bb123117db44bfc0ec78cc73e5a0b87a92afccd6d71d4f952bba5ff34defd3626cd3b3c86816c384f4f2c5241a75393da4b77572b1006e19b48ad5ab9ed3e7d774bef7aae5d9047b773075c7372a3736022f7064e66a32a567eb5112d32061622a1cfd33ff4076579a0ab962fad9547816c095277d6e01000000020000000000000001",
					Signature: "9b0d10e2b321b54edac30aae019bc0ceb639d3c1f312cd65d8dbafe735e14ccc39b11974f4efd29c11a9dccc140878ba689294a1c91d5d569a44b9665a0031fb01",
				},
			},
			want:    txId,
			wantErr: false,
		},
		{
			name: "Issue multisig tx - invalid signature",
			args: args{
				issueArgs: &dto.IssueTxArgs{
					SignedTx:  "000000002007000003ea00000000000000000000000000000000000000000000000000000000000000000000000159eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a0000000700016bcc41d9bdc0000000000000000000000001000000015d008196f8da54c34bd67dc5ef5bae4948389cb8000000010903208c79e9d29ad5e5ea7caf771ecca4db7a218c44d7c3619deea62e6227640000000359eb48b8b3a928ca9d6b90a0f3492ab47ebf06e9edc553cfb6bcd2d3f38e319a0000000500016bcc41e9000000000002000000000000000100000000000000000000000000000000000000000000000083b1ddd7b166dbe6305c22fed5f59065525c4e510000000a00000001000000005d008196f8da54c34bd67dc5ef5bae4948389cb8000000030000200c00000002dd3be02c98a8d121e6a0e3bb123117db44bfc0ec78cc73e5a0b87a92afccd6d71d4f952bba5ff34defd3626cd3b3c86816c384f4f2c5241a75393da4b77572b1006e19b48ad5ab9ed3e7d774bef7aae5d9047b773075c7372a3736022f7064e66a32a567eb5112d32061622a1cfd33ff4076579a0ab962fad9547816c095277d6e010000000200000000000000010000200c00000001a32fc319922bf20632f85f5c99c3ecdf88387cf28564452403e81d635c805c736d2217e83cc33dea85311039a2745fc4bcd28f4b822799b819dc858a6391dd810100000001000000000000200c00000002dd3be02c98a8d121e6a0e3bb123117db44bfc0ec78cc73e5a0b87a92afccd6d71d4f952bba5ff34defd3626cd3b3c86816c384f4f2c5241a75393da4b77572b1006e19b48ad5ab9ed3e7d774bef7aae5d9047b773075c7372a3736022f7064e66a32a567eb5112d32061622a1cfd33ff4076579a0ab962fad9547816c095277d6e01000000020000000000000001",
					Signature: "9b0d10e2b321b54edac30aae019bc0ceb639d3c1f312cd65d8dbafe735e14ccc39b11974f4efd29c11a9dccc140878ba689294a1c91d5d569a44b9665a0031fb02",
				},
			},
			want:    ids.Empty,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)
			got, err := s.IssueMultisigTx(tt.args.issueArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IssueMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCancelMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockNodeService := NewMockNodeService(ctrl)
	mockDao := dao.NewMockMultisigTxDao(ctrl)
	mockConfig := &util.Config{
		NetworkId: networkId,
	}

	mockTx := model.MultisigTx{
		Id:            "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
		UnsignedTx:    "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:     2,
		ChainId:       "11111111111111111111111111111111LpoYY",
		TransactionId: "",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus1yzq6k26nsyuzssj8j9k6x6x7fqgndtadk66948",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
			{
				MultisigTxId: "cec9762115a58339c0f5e9ae582c1879300c1ff7303f9b566a95cf5ebe2a9d28",
				Address:      "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
				Signature:    "",
			},
		},
	}

	// mock without signer
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).AnyTimes()
	mockDao.EXPECT().DeletePendingTx(mockTx.Id).Return(true, nil).AnyTimes()

	type args struct {
		cancelArgs *dto.CancelTxArgs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Cancel multisig tx",
			args: args{
				cancelArgs: &dto.CancelTxArgs{
					Id:        mockTx.Id,
					Timestamp: "1678877386",
					Signature: "47bf8e8601badef42a1157e07862157ded68fff927bc3809d5abb0d4a7c51cad3e53979193dc7069f73fe3f7b1b9e8a5946a1bd4782a565fe126a627634943dd01",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMultisigService(mockConfig, mockDao, mockNodeService)
			err := s.CancelMultisigTx(tt.args.cancelArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
