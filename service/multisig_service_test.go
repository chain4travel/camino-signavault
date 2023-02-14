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
	"reflect"
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

func TestMultisigService_AddMultisigTxSigner(t *testing.T) {

	s := &MultisigService{db.Db{DB: conn}}
	id, err := s.CreateMultisigTx(&model.MultisigTx{
		Alias:      "X-kopernikus1vscyf7czawylztn6ghhg0z27swwewxgzgpcxvv",
		Threshold:  2,
		UnsignedTx: "FFFFFFFC",
	})
	if err != nil {
		return
	}

	type args struct {
		id     int
		signer *model.MultisigTxSigner
	}
	type result struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}
	tests := []result{}
	tests = append(tests, result{
		name: "AddMultisigTxSigner 1",
		args: args{
			id: int(id),
			signer: &model.MultisigTxSigner{
				Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2h",
				Signature: "FFFFFFFA",
			},
		},
		want:    1,
		wantErr: false,
	},
		result{
			name: "AddMultisigTxSigner 2 identical to 1",
			args: args{
				id: int(id),
				signer: &model.MultisigTxSigner{
					Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2h",
					Signature: "FFFFFFFA",
				},
			},
			want:    0,
			wantErr: true,
		},
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MultisigService{db.Db{DB: conn}}
			got, err := s.AddMultisigTxSigner(tt.args.id, tt.args.signer)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMultisigTxSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddMultisigTxSigner() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					Signers: []model.MultisigTxSigner{
						{
							Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2r",
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
					Signers: []model.MultisigTxSigner{
						{
							Address:   "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
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
					Signers: []model.MultisigTxSigner{
						{
							Address:   "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
							Signature: "FFFFFFFB",
						},
						{
							Address:   "X-kopernikus10q3dc78tw70m89s8lgl9fgwe9tfu4a0sfr6yjr",
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
			if got <= 0 && !tt.wantErr {
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
		UnsignedTx: "FFFFFFFC",
		Signers: []model.MultisigTxSigner{
			{
				Address:   "X-kopernikus1vxmf8899y6x7dsam0xnr0hp6syzwz333dkvn2r",
				Signature: "FFFFFFFA",
			},
		},
	}

	id, err := s.CreateMultisigTx(&mock1)
	if err != nil {
		t.Errorf("GetAllMultisigTx() error = %v", err)
		return
	}
	mock1.Id = id
	mock1.Signers[0].Id = 1
	mock1.Signers[0].MultisigTxId = id

	mockMultisigTx = append(mockMultisigTx, mock1)

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
			//s := &MultisigService{}
			got, err := s.GetAllMultisigTx()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func TestMultisigService_GetAllMultisigTxForAlias(t *testing.T) {
//	type args struct {
//		alias string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *[]model.MultisigTx
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MultisigService{}
//			got, err := s.GetAllMultisigTxForAlias(tt.args.alias)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetAllMultisigTxForAlias() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetAllMultisigTxForAlias() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMultisigService_GetMultisigTx(t *testing.T) {
//	type args struct {
//		alias string
//		id    int
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *model.MultisigTx
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MultisigService{}
//			got, err := s.GetMultisigTx(tt.args.alias, tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetMultisigTx() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMultisigService_doGetMultisigTx(t *testing.T) {
//	type args struct {
//		alias string
//		id    int
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *[]model.MultisigTx
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &MultisigService{}
//			got, err := s.doGetMultisigTx(tt.args.alias, tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("doGetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("doGetMultisigTx() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestNewMultisigService(t *testing.T) {
//	tests := []struct {
//		name string
//		want *MultisigService
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMultisigService(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMultisigService() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
