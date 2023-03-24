package service

import (
	"fmt"
	"github.com/ava-labs/avalanchego/utils/hashing"
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

	unsignedTx := "0x00000000200400003039010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4a36162"
	id := fmt.Sprintf("%x", hashing.ComputeHash256([]byte(unsignedTx)))

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
				Signature:    "0x83a657db18ff50438d418db9bde239a47bca2d910114aa0cc68ac736053c96c46b300f2c28d0ff6c8587ae916b66b5d575a731d8ecc37abee3c96bdc23dcd88b007c40d266",
			},
			{
				MultisigTxId: id,
				Address:      mockAliasInfo.Result.Addresses[1],
				Signature:    "",
			},
		},
	}

	thresholdInt, _ := strconv.Atoi(mockAliasInfo.Result.Threshold)
	// id string, alias string, threshold int, unsignedTx string, creator string, signature string, outputOwners string, owners []string
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
