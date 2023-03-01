package service

import (
	"context"
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/test"
	"github.com/chain4travel/camino-signavault/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"testing"
	"time"
)

var conn *sql.DB

func TestMain(m *testing.M) {
	var code = 1
	defer func() { os.Exit(code) }()

	ctx := context.Background()
	mysqlContainer, err := test.SetupMysql(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		_, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		conn, err = sql.Open("mysql", mysqlContainer.URI)

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
//		name: "AddMultisigTxSigner 1",
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
//			name: "AddMultisigTxSigner 2 identical to 1",
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
//			got, err := s.AddMultisigTxSigner(tt.args.id, tt.args.signer)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("AddMultisigTxSigner() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("AddMultisigTxSigner() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestMultisigService_CreateMultisigTx(t *testing.T) {
	type args struct {
		multisigTx *model.MultisigTx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateMultisigTx 1",
			args: args{
				multisigTx: &model.MultisigTx{
					Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
					Threshold:  2,
					UnsignedTx: "FFFFFFFC",
					Owners: []model.MultisigTxOwner{
						{
							Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2r",
						},
					},
					Signers: []model.MultisigTxSigner{
						{
							MultisigTxOwner: model.MultisigTxOwner{
								Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2r",
							},
							Signature: "FFFFFFFA",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CreateMultisigTx 2",
			args: args{
				multisigTx: &model.MultisigTx{
					Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
					Threshold:  3,
					UnsignedTx: "FFFFFFFD",
					Owners: []model.MultisigTxOwner{
						{
							Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
						},
					},
					Signers: []model.MultisigTxSigner{
						{
							MultisigTxOwner: model.MultisigTxOwner{
								Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
							},
							Signature: "FFFFFFFB",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CreateMultisigTx 3 with duplicate signers",
			args: args{
				multisigTx: &model.MultisigTx{
					Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
					Threshold:  3,
					UnsignedTx: "FFFFFFFD",
					Owners: []model.MultisigTxOwner{
						{
							Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
						},
					},
					Signers: []model.MultisigTxSigner{
						{
							MultisigTxOwner: model.MultisigTxOwner{
								Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
							},
							Signature: "FFFFFFFB",
						},
						{
							MultisigTxOwner: model.MultisigTxOwner{
								Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
							},
							Signature: "FFFFFFFB",
						},
						{
							MultisigTxOwner: model.MultisigTxOwner{
								Address: "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
							},
							Signature: "FFFFFFFB",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MultisigService{db.Db{DB: conn}}
			got, err := s.CreateMultisigTx(tt.args.multisigTx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isEqual(got, tt.args.multisigTx) && !tt.wantErr {
				t.Errorf("CreateMultisigTx() got = %v, want > %v", got, 0)
			}
		})
	}
}

func TestMultisigService_GetAllMultisigTx(t *testing.T) {
	s := &MultisigService{db.Db{DB: conn}}

	// array of mock multisig tx
	var mockMultisigTx []model.MultisigTx

	mock1 := model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
		Threshold:  2,
		UnsignedTx: "FFFFFFFQ",
		Owners: []model.MultisigTxOwner{
			{
				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
			},
		},
		Signers: []model.MultisigTxSigner{
			{
				MultisigTxOwner: model.MultisigTxOwner{
					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
				},
				Signature: "FFFFFFFA",
			},
		},
	}

	_, err := s.CreateMultisigTx(&mock1)
	if err != nil {
		t.Errorf("GetAllMultisigTx() error = %v", err)
		return
	}
	mock2 := model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxxx",
		Threshold:  3,
		UnsignedTx: "FFFFFFCC",
		Owners: []model.MultisigTxOwner{
			{
				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
			},
		},
		Signers: []model.MultisigTxSigner{
			{
				MultisigTxOwner: model.MultisigTxOwner{
					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
				},
				Signature: "FFFFFFFD",
			},
		},
	}

	_, err = s.CreateMultisigTx(&mock2)
	if err != nil {
		t.Errorf("GetAllMultisigTx() error = %v", err)
		return
	}
	mockMultisigTx = append(mockMultisigTx, mock1, mock2)

	tests := []struct {
		name    string
		want    *[]model.MultisigTx
		wantErr bool
	}{
		{
			name:    "GetAllMultisigTx",
			want:    &mockMultisigTx,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetAllMultisigTx()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return

			}
			var found bool
			for i := range *got {
				for k := range *tt.want {
					a := (*got)[i]
					b := (*tt.want)[k]
					if isEqual(&a, &b) {
						found = true
					}
				}
			}

			if !found && !tt.wantErr {
				t.Errorf("GetAllMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultisigService_GetAllMultisigTxForAlias(t *testing.T) {
	s := &MultisigService{db.Db{DB: conn}}

	// mock1
	var mockMultisigTx []model.MultisigTx
	mock1 := model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa",
		Threshold:  2,
		UnsignedTx: "FFFFFCCC",
		Owners: []model.MultisigTxOwner{
			{
				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
			},
		},
		Signers: []model.MultisigTxSigner{
			{
				MultisigTxOwner: model.MultisigTxOwner{
					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2q",
				},
				Signature: "FFFFFFFA",
			},
		},
	}
	_, err := s.CreateMultisigTx(&mock1)
	if err != nil {
		t.Errorf("GetAllMultisigTxForAlias() error = %v", err)
		return
	}

	// mock2
	mock2 := model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa",
		Threshold:  3,
		UnsignedTx: "FFFFCCCC",
		Owners: []model.MultisigTxOwner{
			{
				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
			},
		},
		Signers: []model.MultisigTxSigner{
			{
				MultisigTxOwner: model.MultisigTxOwner{
					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkiiii",
				},
				Signature: "FFFFFFFD",
			},
		},
	}
	_, err = s.CreateMultisigTx(&mock2)
	if err != nil {
		t.Errorf("GetAllMultisigTxForAlias() error = %v", err)
		return
	}
	// keep mocks in array
	mockMultisigTx = append(mockMultisigTx, mock1, mock2)

	type args struct {
		alias string
	}

	tests := []struct {
		name    string
		args    args
		want    *[]model.MultisigTx
		wantErr bool
	}{
		{
			name:    "GetAllMultisigTxForAlias - 1",
			args:    args{alias: "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcaaa"},
			want:    &mockMultisigTx,
			wantErr: false,
		},
		{
			name:    "GetAllMultisigTxForAlias - 2",
			args:    args{alias: "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcabc"},
			want:    &mockMultisigTx,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetAllMultisigTxForAlias(tt.args.alias)

			//if (err != nil) != tt.wantErr {
			if err != nil {
				t.Errorf("GetAllMultisigTxForAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// search for mock1 and mock2 in got array
			var found bool
			for i := range *got {
				for k := range *tt.want {
					a := (*got)[i]
					b := (*tt.want)[k]
					if isEqual(&a, &b) {
						found = true
					}
				}
			}

			if !found && !tt.wantErr {
				t.Errorf("GetAllMultisigTxForAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultisigService_GetMultisigTx(t *testing.T) {
	s := &MultisigService{db.Db{DB: conn}}
	_, err := s.CreateMultisigTx(&model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
		Threshold:  2,
		UnsignedTx: "FFFFFFFC",
		Owners: []model.MultisigTxOwner{
			{
				Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333aaaaaa",
			},
		},
		Signers: []model.MultisigTxSigner{
			{
				MultisigTxOwner: model.MultisigTxOwner{
					Address: "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333aaaaaa",
				},
				Signature: "FFFFFFFA",
			},
		},
	})
	if err != nil {
		return
	}

	type args struct {
		UnsignedTx string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.MultisigTx
		wantErr bool
	}{
		{
			name: "GetMultisigTx with correct transaction id",
			args: args{UnsignedTx: "FFFFFFFC"},
			want: &model.MultisigTx{
				Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
				Threshold:  2,
				UnsignedTx: "FFFFFFFC",
			},
			wantErr: false,
		},
		{
			name: "GetMultisigTx with wrong transaction id (nil result)",
			args: args{UnsignedTx: "FFFFFFFB"},
			want: &model.MultisigTx{
				Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzaaaaaa",
				Threshold:  4,
				UnsignedTx: "FFFFFCCC",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetMultisigTx(tt.args.UnsignedTx)

			if err != nil && !tt.wantErr {
				t.Errorf("GetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !isEqual(got, tt.want) && !tt.wantErr {
					t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func isEqual(a *model.MultisigTx, b *model.MultisigTx) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	// compare all fields of a and b excepts ids
	if a.Alias != b.Alias || a.Threshold != b.Threshold || a.UnsignedTx != b.UnsignedTx {
		if a.Owners != nil && b.Owners != nil {
			if len(a.Owners) != len(b.Owners) {
				return false
			}
			for i := range a.Owners {
				if a.Owners[i].Address != b.Owners[i].Address {
					return false
				}
			}
		}
		if a.Signers != nil && b.Signers != nil {
			if len(a.Signers) != len(b.Signers) {
				return false
			}
			for i := range a.Signers {
				if a.Signers[i].Address != b.Signers[i].Address || a.Signers[i].Signature != b.Signers[i].Signature {
					return false
				}
			}
		}
		return false
	}
	return true
}
