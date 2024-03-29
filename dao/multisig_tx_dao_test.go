package dao

import (
	"context"
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/test"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

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
		conn, err = sql.Open("mysql", mysqlContainer.URI+"?multiStatements=true&parseTime=true")
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

		// load test data
		path = filepath.Join(util.GetRootPath(), "test", "data", "test_data.sql")
		c, ioErr := os.ReadFile(path)
		if ioErr != nil {
			log.Fatal(err)
		}
		script := string(c)
		_, err = conn.Exec(script)
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

func TestNewMultisigTxDao(t *testing.T) {
	type args struct {
		db *db.Db
	}
	tests := []struct {
		name string
		args args
		want MultisigTxDao
	}{
		{
			name: "Happy path",
			args: args{
				db: &db.Db{DB: conn},
			},
			want: &multisigTxDao{
				db: &db.Db{DB: conn},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMultisigTxDao(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMultisigTxDao() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMultisigTx(t *testing.T) {
	type fields struct {
		db *db.Db
	}
	type args struct {
		multisigTx *model.MultisigTx
	}
	exp := time.Now().Add(time.Hour * 24 * 7)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Create multisig tx",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				multisigTx: &model.MultisigTx{
					Id:           "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69",
					UnsignedTx:   "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					Alias:        "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
					Threshold:    2,
					ChainId:      "11111111111111111111111111111111LpoYY",
					OutputOwners: "OutputOwners",
					Metadata:     "metadata",
					Expiration:   &exp,
					Owners: []model.MultisigTxOwner{
						{
							Address:   "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
							Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
						},
						{
							Address:   "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
							Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
						},
					},
				},
			},
			want:    "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69",
			wantErr: false,
		},
		{
			name: "Create duplicate id tx - should fail",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				multisigTx: &model.MultisigTx{
					Id:           "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69",
					UnsignedTx:   "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					Alias:        "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
					Threshold:    2,
					ChainId:      "11111111111111111111111111111111LpoYY",
					OutputOwners: "OutputOwners",
					Metadata:     "metadata",
					Owners: []model.MultisigTxOwner{
						{
							Address:   "P-kopernikus18jma8ppw3nhx5r4ap8clazz0dps7rv5uuvjh68",
							Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
						},
						{
							Address:   "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
							Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &multisigTxDao{
				db: tt.fields.db,
			}
			got, err := d.CreateMultisigTx(tt.args.multisigTx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateMultisigTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddSigner(t *testing.T) {
	type fields struct {
		db *db.Db
	}
	type args struct {
		id            string
		signature     string
		signerAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Add signer with signature",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				id:            "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69",
				signature:     "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
				signerAddress: "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Add signer without signature",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				id:            "bc6246f58b5aba675f4071bd1a13d7a774384e42f301208d1c2b0f22ee602e69",
				signerAddress: "P-kopernikus1g65uqn6t77p656w64023nh8nd9updzmxh8ttv3",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &multisigTxDao{
				db: tt.fields.db,
			}
			got, err := d.AddSigner(tt.args.id, tt.args.signature, tt.args.signerAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddSigner() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMultisigTx(t *testing.T) {
	type fields struct {
		db *db.Db
	}
	type args struct {
		id    string
		alias string
		owner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]model.MultisigTx
		wantErr bool
	}{
		{
			name: "Get multisig tx by existing id",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				id: "1",
			},
			want: &[]model.MultisigTx{
				{
					Id:           "1",
					Alias:        "alias",
					Threshold:    2,
					UnsignedTx:   "unsigned_tx",
					OutputOwners: "output_owners",
					Metadata:     "metadata",
					Owners: []model.MultisigTxOwner{
						{
							MultisigTxId: "1",
							Address:      "address",
							Signature:    "signature",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Get multisig tx by non existing id",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				id: "99",
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &multisigTxDao{
				db: tt.fields.db,
			}
			got, err := d.GetMultisigTx(tt.args.id, tt.args.alias, tt.args.owner, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMultisigTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}

			// we check fields one by one because of the timestamp field which is generated by the database
			// and do not want to check it right now. Just that it is not empty.
			assert.Equal(t, (*got)[0].Id, (*tt.want)[0].Id)
			assert.Equal(t, (*got)[0].UnsignedTx, (*tt.want)[0].UnsignedTx)
			assert.Equal(t, (*got)[0].Alias, (*tt.want)[0].Alias)
			assert.Equal(t, (*got)[0].Threshold, (*tt.want)[0].Threshold)
			assert.Equal(t, (*got)[0].TransactionId, (*tt.want)[0].TransactionId)
			assert.Equal(t, (*got)[0].OutputOwners, (*tt.want)[0].OutputOwners)
			assert.Equal(t, (*got)[0].Metadata, (*tt.want)[0].Metadata)
			assert.Equal(t, (*got)[0].Owners, (*tt.want)[0].Owners)
			assert.NotEmpty(t, (*got)[0].Timestamp)
		})
	}
}

func TestPendingAliasExists(t *testing.T) {
	type fields struct {
		db *db.Db
	}
	type args struct {
		alias   string
		chainId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test pending tx for existing msig alias",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				alias:   "alias_3",
				chainId: "11111111111111111111111111111111LpoYY",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test pending tx for existing msig alias on different chain",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				alias:   "alias_3",
				chainId: "jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test complete tx for existing msig alias on different chain",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				alias:   "alias_6",
				chainId: "jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test pending tx for existing complete msig alias",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				alias:   "alias_2",
				chainId: "11111111111111111111111111111111LpoYY",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test pending tx for non existing msig alias",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				alias:   "test",
				chainId: "11111111111111111111111111111111LpoYY",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &multisigTxDao{
				db: tt.fields.db,
			}
			got, err := d.PendingAliasExists(tt.args.alias, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("PendingAliasExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PendingAliasExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTransactionId(t *testing.T) {
	type fields struct {
		db *db.Db
	}
	type args struct {
		id            string
		transactionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Update transaction id for existing multisig tx",
			fields: fields{
				db: &db.Db{DB: conn},
			},
			args: args{
				id:            "3",
				transactionId: "transaction_id_3",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &multisigTxDao{
				db: tt.fields.db,
			}
			got, err := d.UpdateTransactionId(tt.args.id, tt.args.transactionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateTransactionId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
