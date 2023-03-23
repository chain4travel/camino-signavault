package service

import (
	"fmt"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"log"
	"strconv"
	"strings"
	"testing"
)

const networkId = uint32(1002)

func TestCreateMultisigTx(t *testing.T) {
	// Create a tx
	//unsignedAddressStateTx := &txs.AddressStateTx{
	//	BaseTx: txs.BaseTx{BaseTx: avax.BaseTx{
	//		NetworkID:    networkId,
	//		BlockchainID: ids.ID{1},
	//		Outs:         []*avax.TransferableOutput{},
	//		Ins:          []*avax.TransferableInput{},
	//	}},
	//}
	//
	//signer := getSigner(3)
	//tx, err := signTX(unsignedAddressStateTx, signer)
	//require.NoError(t, err)
	//
	//// Get the unsigned tx bytes
	//utxBytes := tx.Unsigned.Bytes()
	//utxString, err := formatting.Encode(formatting.Hex, utxBytes)
	//require.NoError(t, err)
	//
	//signature, err := getSignatureFromTx(tx)
	//require.NoError(t, err)

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

//func TestGetMultisigTx(t *testing.T) {
//	msigAlias := "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy"
//	timestamp := "1678877386"
//
//	signer := getSigner(3)
//	signature, err := signPayload(signer, msigAlias, timestamp)
//	require.NoError(t, err)
//
//	type args struct {
//		msigAlias string
//		timestamp string
//		signature string
//	}
//	tests := []struct {
//		name string
//		args args
//		err  error
//	}{
//		{
//			name: "Happy path",
//			args: args{
//				msigAlias: msigAlias,
//				timestamp: timestamp,
//				signature: signature,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			config := &util.Config{}
//			s := NewMultisigService(config, dao.NewMultisigTxDao(&db.Db{DB: conn}), NewNodeService(config))
//
//			_, err := s.GetAllMultisigTxForAlias(tt.args.msigAlias, tt.args.timestamp, tt.args.signature)
//			if tt.err != nil {
//				require.Equal(t, tt.err, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}

func signTX(unsigned txs.UnsignedTx, signer *crypto.PrivateKeySECP256K1R) (*txs.Tx, error) {
	// Sign the tx
	tx, err := txs.NewSigned(unsigned, txs.Codec, [][]*crypto.PrivateKeySECP256K1R{{signer}})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func getSignatureFromTx(tx *txs.Tx) (string, error) {
	// Get the signature from the tx
	var sig [crypto.SECP256K1RSigLen]byte
	for _, v := range tx.Creds {
		if cred, ok := v.(*secp256k1fx.Credential); ok {
			sig = cred.Sigs[0]
			break
		}
	}
	signature, err := formatting.Encode(formatting.Hex, sig[:])
	if err != nil {
		return "", err
	}
	return signature, nil
}

func getSigner(prefundedKeyIndex int) *crypto.PrivateKeySECP256K1R {
	preFundedKeys := crypto.BuildTestKeys()
	addr := preFundedKeys[prefundedKeyIndex].PublicKey()
	log.Print(addr.Address().String())

	hrp := constants.NetworkIDToHRP[networkId]
	bech32Address, _ := address.FormatBech32(hrp, addr.Bytes())
	log.Print(bech32Address)
	signer := preFundedKeys[prefundedKeyIndex]
	return signer
}

func signPayload(signer *crypto.PrivateKeySECP256K1R, payload ...string) (string, error) {
	// Compute the hash of the payload
	hash := hashing.ComputeHash256([]byte(strings.Join(payload, "")))

	// Sign the hash
	sig, err := signer.SignHash(hash)
	if err != nil {
		return "", err
	}
	signature, err := formatting.Encode(formatting.Hex, sig[:])
	if err != nil {
		return "", err
	}
	return signature, nil
}

//func TestCompleteMultisigTx(t *testing.T) {
//	preFundedKeys := crypto.BuildTestKeys()
//	address := preFundedKeys[3].PublicKey().Address()
//	log.Print(address.String())
//
//	signer := preFundedKeys[3]
//	id := 1
//	txID := "2wKRJ8XKh8h7g2BKaDPXGjBwDJ3fMCuXrdCaeUgqpisJMwAui8"
//	timestamp := "1678877386"
//
//	// Compute the hash of the payload
//	hash := hashing.ComputeHash256([]byte(strconv.Itoa(id) + timestamp + txID))
//
//	// Sign the hash
//	sig, err := signer.SignHash(hash)
//	require.NoError(t, err)
//	signature, err := formatting.Encode(formatting.Hex, sig[:])
//	require.NoError(t, err)
//
//	type args struct {
//		id         int64
//		completeTx *dto.CompleteTxArgs
//	}
//	tests := []struct {
//		name string
//		args args
//		err  error
//	}{
//		{
//			name: "Happy path",
//			args: args{
//				id: 1,
//				completeTx: &dto.CompleteTxArgs{
//					TransactionId: txID,
//					Signature:     signature,
//					Timestamp:     timestamp,
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockConfig := &util.Config{
//				NetworkId: 1002,
//			}
//			mockNodeService := new(MockNodeService)
//			mockNodeService.On("GetTx", tt.args.completeTx.TransactionId).Return(&model.TxInfo{
//				Jsonrpc: "2.0",
//				Result: struct {
//					Tx       string `json:"tx"`
//					Encoding string `json:"encoding"`
//				}{
//					Tx:       "0x000000",
//					Encoding: "hex",
//				},
//				Id: 1,
//			}, nil)
//
//			s := NewMultisigService(mockConfig, dao.NewMultisigTxDao(&db.Db{DB: conn}), mockNodeService)
//			_, err := s.CompleteMultisigTx(tt.args.id, tt.args.completeTx)
//			if tt.err != nil {
//				require.Equal(t, tt.err, err)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}

//func TestMultisigService_AddMultisigTxSigner(t *testing.T) {
//
//	s := &MultisigService{db.Db{DB: conn}}
//	id, err := s.CreateMultisigTx(&model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
//		Threshold:  2,
//		UnsignedTx: "FFFFFFFC",
//	})
//	if err != nil {
//		return
//	}
//
//	type args struct {
//		id     int
//		signer *model.MultisigTxOwner
//	}
//	type result struct {
//		name    string
//		args    args
//		want    int64
//		wantErr bool
//	}
//	tests := []result{}
//	tests = append(tests, result{
//		name: "SignMultisigTx 1",
//		args: args{
//			id: int(id),
//			signer: &model.MultisigTxOwner{
//				Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2h",
//				Signature: "FFFFFFFA",
//			},
//		},
//		want:    1,
//		wantErr: false,
//	},
//		result{
//			name: "SignMultisigTx 2 identical to 1",
//			args: args{
//				id: int(id),
//				signer: &model.MultisigTxOwner{
//					Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2h",
//					Signature: "FFFFFFFA",
//				},
//			},
//			want:    0,
//			wantErr: true,
//		},
//	)
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MultisigService{db.Db{DB: conn}}
//			got, err := s.SignMultisigTx(tt.args.id, tt.args.signer)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("SignMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("SignMultisigTx() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMultisigService_GetAllMultisigTx(t *testing.T) {
//	s := &MultisigService{db.Db{DB: conn}}
//
//	// array of mock multisig tx
//	var mockMultisigTx []model.MultisigTx
//
//	mock1 := model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
//		Threshold:  2,
//		UnsignedTx: "FFFFFFFQ",
//		Owners: []model.MultisigTxOwner{
//			{
//				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
//			},
//		},
//		Signers: []model.MultisigTxSigner{
//			{
//				MultisigTxOwner: model.MultisigTxOwner{
//					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
//				},
//				Signature: "FFFFFFFA",
//			},
//		},
//	}
//
//	_, err := s.CreateMultisigTx(&mock1)
//	if err != nil {
//		t.Errorf("GetAllMultisigTx() error = %v", err)
//		return
//	}
//	mock2 := model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxxx",
//		Threshold:  3,
//		UnsignedTx: "FFFFFFCC",
//		Owners: []model.MultisigTxOwner{
//			{
//				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
//			},
//		},
//		Signers: []model.MultisigTxSigner{
//			{
//				MultisigTxOwner: model.MultisigTxOwner{
//					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
//				},
//				Signature: "FFFFFFFD",
//			},
//		},
//	}
//
//	_, err = s.CreateMultisigTx(&mock2)
//	if err != nil {
//		t.Errorf("GetAllMultisigTx() error = %v", err)
//		return
//	}
//	mockMultisigTx = append(mockMultisigTx, mock1, mock2)
//
//	tests := []struct {
//		name    string
//		want    *[]model.MultisigTx
//		wantErr bool
//	}{
//		{
//			name:    "GetAllMultisigTx",
//			want:    &mockMultisigTx,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := s.GetAllMultisigTx()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetAllMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
//				return
//
//			}
//			var found bool
//			for i := range *got {
//				for k := range *tt.want {
//					a := (*got)[i]
//					b := (*tt.want)[k]
//					if isEqual(&a, &b) {
//						found = true
//					}
//				}
//			}
//
//			if !found && !tt.wantErr {
//				t.Errorf("GetAllMultisigTx() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMultisigService_GetAllMultisigTxForAlias(t *testing.T) {
//	s := &MultisigService{db.Db{DB: conn}}
//
//	// mock1
//	var mockMultisigTx []model.MultisigTx
//	mock1 := model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa",
//		Threshold:  2,
//		UnsignedTx: "FFFFFCCC",
//		Owners: []model.MultisigTxOwner{
//			{
//				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
//			},
//		},
//		Signers: []model.MultisigTxSigner{
//			{
//				MultisigTxOwner: model.MultisigTxOwner{
//					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
//				},
//				Signature: "FFFFFFFA",
//			},
//		},
//	}
//	_, err := s.CreateMultisigTx(&mock1)
//	if err != nil {
//		t.Errorf("GetAllMultisigTxForAlias() error = %v", err)
//		return
//	}
//
//	// mock2
//	mock2 := model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa",
//		Threshold:  3,
//		UnsignedTx: "FFFFCCCC",
//		Owners: []model.MultisigTxOwner{
//			{
//				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
//			},
//		},
//		Signers: []model.MultisigTxSigner{
//			{
//				MultisigTxOwner: model.MultisigTxOwner{
//					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
//				},
//				Signature: "FFFFFFFD",
//			},
//		},
//	}
//	_, err = s.CreateMultisigTx(&mock2)
//	if err != nil {
//		t.Errorf("GetAllMultisigTxForAlias() error = %v", err)
//		return
//	}
//	// keep mocks in array
//	mockMultisigTx = append(mockMultisigTx, mock1, mock2)
//
//	type args struct {
//		alias string
//	}
//
//	tests := []struct {
//		name    string
//		args    args
//		want    *[]model.MultisigTx
//		wantErr bool
//	}{
//		{
//			name:    "GetAllMultisigTxForAlias - 1",
//			args:    args{alias: "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa"},
//			want:    &mockMultisigTx,
//			wantErr: false,
//		},
//		{
//			name:    "GetAllMultisigTxForAlias - 2",
//			args:    args{alias: "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcabc"},
//			want:    &mockMultisigTx,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := s.GetAllMultisigTxForAlias(tt.args.alias)
//
//			//if (err != nil) != tt.wantErr {
//			if err != nil {
//				t.Errorf("GetAllMultisigTxForAlias() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			// search for mock1 and mock2 in got array
//			var found bool
//			for i := range *got {
//				for k := range *tt.want {
//					a := (*got)[i]
//					b := (*tt.want)[k]
//					if isEqual(&a, &b) {
//						found = true
//					}
//				}
//			}
//
//			if !found && !tt.wantErr {
//				t.Errorf("GetAllMultisigTxForAlias() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMultisigService_GetMultisigTx(t *testing.T) {
//	s := &MultisigService{db.Db{DB: conn}}
//	_, err := s.CreateMultisigTx(&model.MultisigTx{
//		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
//		Threshold:  2,
//		UnsignedTx: "FFFFFFFC",
//		Owners: []model.MultisigTxOwner{
//			{
//				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333aaaaaa",
//			},
//		},
//		Signers: []model.MultisigTxSigner{
//			{
//				MultisigTxOwner: model.MultisigTxOwner{
//					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333aaaaaa",
//				},
//				Signature: "FFFFFFFA",
//			},
//		},
//	})
//	if err != nil {
//		return
//	}
//
//	type args struct {
//		UnsignedTx string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *model.MultisigTx
//		wantErr bool
//	}{
//		{
//			name: "GetMultisigTx with correct transaction id",
//			args: args{UnsignedTx: "FFFFFFFC"},
//			want: &model.MultisigTx{
//				Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
//				Threshold:  2,
//				UnsignedTx: "FFFFFFFC",
//			},
//			wantErr: false,
//		},
//		{
//			name: "GetMultisigTx with wrong transaction id (nil result)",
//			args: args{UnsignedTx: "FFFFFFFB"},
//			want: &model.MultisigTx{
//				Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
//				Threshold:  4,
//				UnsignedTx: "FFFFFCCC",
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := s.GetMultisigTx(tt.args.UnsignedTx)
//
//			if err != nil && !tt.wantErr {
//				t.Errorf("GetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != nil {
//				if !isEqual(got, tt.want) && !tt.wantErr {
//					t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
//				}
//			} else {
//				if !tt.wantErr {
//					t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
//				}
//			}
//		})
//	}
//}
//
//func isEqual(a *model.MultisigTx, b *model.MultisigTx) bool {
//	if (a == nil) != (b == nil) {
//		return false
//	}
//	// compare all fields of a and b excepts ids
//	if a.Alias != b.Alias || a.Threshold != b.Threshold || a.UnsignedTx != b.UnsignedTx {
//		if a.Owners != nil && b.Owners != nil {
//			if len(a.Owners) != len(b.Owners) {
//				return false
//			}
//			for i := range a.Owners {
//				if a.Owners[i].Address != b.Owners[i].Address {
//					return false
//				}
//			}
//		}
//		if a.Signers != nil && b.Signers != nil {
//			if len(a.Signers) != len(b.Signers) {
//				return false
//			}
//			for i := range a.Signers {
//				if a.Signers[i].Address != b.Signers[i].Address || a.Signers[i].Signature != b.Signers[i].Signature {
//					return false
//				}
//			}
//		}
//		return false
//	}
//	return true
//}
