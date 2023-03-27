/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
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
	id := "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69"

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

	mockTx := model.MultisigTx{
		Id:            id,
		Alias:         "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		UnsignedTx:    unsignedTx,
		Threshold:     2,
		TransactionId: "",
		OutputOwners:  "OutputOwners",
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

	thresholdInt, _ := strconv.Atoi(mockAliasInfo.Result.Threshold)
	mockDao.EXPECT().CreateMultisigTx(mockTx.Id, mockTx.Alias, thresholdInt, mockTx.UnsignedTx, mockTx.Owners[0].Address, mockTx.Owners[0].Signature, mockTx.OutputOwners, mockAliasInfo.Result.Addresses).Return(mockTx.Id, nil)
	mockDao.EXPECT().GetMultisigTx(mockTx.Id, "", "").Return(&[]model.MultisigTx{mockTx}, nil).AnyTimes()
	mockDao.EXPECT().PendingAliasExists("P-kopernikus1fq0jc8svlyazhygkj0s36qnl6s0km0h3uuc99e").Return(true, nil)
	mockDao.EXPECT().PendingAliasExists(gomock.Any()).Return(false, nil).AnyTimes()
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
			err: errPendingTx,
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
	// todo: implement this test
}

func TestGetMultisigTx(t *testing.T) {
	// todo: implement this test
}

func TestSignMultisigTx(t *testing.T) {
	// todo: implement this test
}
