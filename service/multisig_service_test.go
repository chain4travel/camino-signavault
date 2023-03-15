package service

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/utils/hashing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"

	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/test"
	"github.com/chain4travel/camino-signavault/util"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	code := 1
	defer func() { os.Exit(code) }()

	ctx := context.Background()
	mysqlContainer, err := test.SetupMysql(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		conn, err = sql.Open("mysql", mysqlContainer.URI)
		if err != nil {
			log.Fatal(err)
		}

		// run migration
		path := "file://" + util.GetRootPath() + "/db/migrations"
		m, err := migrate.New(path, "mysql://"+mysqlContainer.URI)
		if err != nil {
			log.Fatal(err)
		}
		err = m.Up()
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = conn.Close(); err != nil {
				panic(err)
			}
		}()
	}
	code = m.Run()
}

func TestCreateMultisigTx(t *testing.T) {
	preFundedKeys := crypto.BuildTestKeys()
	address := preFundedKeys[3].PublicKey().Address()
	log.Print(address.String())

	signers := [][]*crypto.PrivateKeySECP256K1R{
		{preFundedKeys[3]},
	}

	// Create a tx
	unsignedAddressStateTx := &txs.AddressStateTx{
		BaseTx: txs.BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    uint32(12345),
			BlockchainID: ids.ID{1},
			Outs:         []*avax.TransferableOutput{},
			Ins:          []*avax.TransferableInput{},
		}},
	}

	// Sign the tx
	tx, err := txs.NewSigned(unsignedAddressStateTx, txs.Codec, signers)
	require.NoError(t, err)

	// Get the unsigned tx bytes
	utxBytes := tx.Unsigned.Bytes()
	utxString, err := formatting.Encode(formatting.Hex, utxBytes)
	require.NoError(t, err)

	// Get the signature from the tx
	var sig [crypto.SECP256K1RSigLen]byte
	for _, v := range tx.Creds {
		if cred, ok := v.(*secp256k1fx.Credential); ok {
			sig = cred.Sigs[0]
			break
		}
	}
	signature, err := formatting.Encode(formatting.Hex, sig[:])
	require.NoError(t, err)

	type args struct {
		multisigTx *dto.MultisigTxArgs
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Happy path",
			args: args{
				multisigTx: &dto.MultisigTxArgs{
					Alias:      "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
					UnsignedTx: utxString,
					Signature:  signature,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MultisigService{
				db: db.Db{DB: conn},
				SECPFactory: crypto.FactorySECP256K1R{
					Cache: cache.LRU{Size: defaultCacheSize},
				},
			}

			_, err := s.CreateMultisigTx(tt.args.multisigTx)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetMultisigTx(t *testing.T) {

	preFundedKeys := crypto.BuildTestKeys()
	address := preFundedKeys[3].PublicKey().Address()
	log.Print(address.String())

	signer := preFundedKeys[3]

	msigAlias := "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy"
	timestamp := "1678877386"

	// Compute the hash of the payload
	hash := hashing.ComputeHash256([]byte(msigAlias + timestamp))

	// Sign the hash
	sig, err := signer.SignHash(hash)
	require.NoError(t, err)
	signature, err := formatting.Encode(formatting.Hex, sig[:])
	require.NoError(t, err)

	type args struct {
		msigAlias string
		timestamp string
		signature string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Happy path",
			args: args{
				msigAlias: msigAlias,
				timestamp: timestamp,
				signature: signature,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MultisigService{
				db: db.Db{DB: conn},
				SECPFactory: crypto.FactorySECP256K1R{
					Cache: cache.LRU{Size: defaultCacheSize},
				},
			}

			_, err := s.GetAllMultisigTxForAlias(tt.args.msigAlias, tt.args.timestamp, tt.args.signature)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

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
