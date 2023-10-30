package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestDepositOfferHandlerAddSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDepositOfferService := service.NewMockDepositOfferService(ctrl)

	mockDepositOfferService.EXPECT().AddSignatures(gomock.Any()).Return(nil).Times(1)
	mockDepositOfferService.EXPECT().AddSignatures(gomock.Any()).Return(fmt.Errorf("error")).Times(1)

	type fields struct {
		DepositOfferService service.DepositOfferService
	}
	type args struct {
		r       *httptest.ResponseRecorder
		ctx     func(r *httptest.ResponseRecorder) *gin.Context
		sigArgs *dto.AddSignatureArgs
	}
	tests := map[string]struct {
		fields fields
		args   args
		err    dto.SignavaultError
	}{
		"Missing required arg signature": {
			fields: fields{
				DepositOfferService: mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				sigArgs: &dto.AddSignatureArgs{
					DepositOfferID: "1",
					Addresses:      []string{"0x123"},
				},
			},
			err: dto.SignavaultError{
				Message: "Error parsing signature args",
			},
		},
		"Valid args - service success": {
			fields: fields{
				DepositOfferService: mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				sigArgs: &dto.AddSignatureArgs{
					DepositOfferID: "1",
					Addresses:      []string{"0x123"},
					Signatures:     []string{"0x123"},
				},
			},
		},
		"Valid args - service error": {
			fields: fields{
				DepositOfferService: mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				sigArgs: &dto.AddSignatureArgs{
					DepositOfferID: "1",
					Addresses:      []string{"0x123"},
					Signatures:     []string{"0x123"},
				},
			},
			err: dto.SignavaultError{
				Message: "Error inserting signature",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := &depositOfferHandler{
				DepositOfferService: tt.fields.DepositOfferService,
			}
			ctx := tt.args.ctx(tt.args.r)

			body, err := json.Marshal(tt.args.sigArgs)
			assert.Nil(t, err)
			ctx.Request = httptest.NewRequest("POST", "/deposit-offer", strings.NewReader(string(body)))
			h.AddSignature(ctx)

			if tt.err != (dto.SignavaultError{}) {
				var errorResponse dto.SignavaultError
				err = json.Unmarshal(tt.args.r.Body.Bytes(), &errorResponse)
				assert.Nil(t, err)
				assert.Equal(t, tt.err.Message, errorResponse.Message)
				assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
			} else {
				assert.Equal(t, http.StatusCreated, ctx.Writer.Status())
			}
		})
	}
}

func TestDepositOfferHandlerGetSignatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDepositOfferService := service.NewMockDepositOfferService(ctrl)

	sigs := &[]model.DepositOfferSig{
		{
			DepositOfferID: "1",
		},
	}
	// mocks
	addr := "0x123"
	signature := "0x123"
	timestamp := "0"
	multisig := true
	mockError := errors.New("error")
	mockDepositOfferService.EXPECT().GetSignatures(addr, timestamp, signature, multisig).Return(sigs, nil).Times(1)
	mockDepositOfferService.EXPECT().GetSignatures(addr, timestamp, signature, multisig).Return(nil, mockError).Times(1)

	type fields struct {
		DepositOfferService service.DepositOfferService
	}
	type args struct {
		r           *httptest.ResponseRecorder
		ctx         func(r *httptest.ResponseRecorder) *gin.Context
		pathParam   string
		queryParams map[string]string
	}
	tests := map[string]struct {
		name   string
		fields fields
		args   args
		err    dto.SignavaultError
	}{
		"missing signature": {
			fields: fields{
				mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				pathParam: addr,
				queryParams: map[string]string{
					"timestamp": timestamp,
					"multisig":  strconv.FormatBool(multisig),
				},
			},
			err: dto.SignavaultError{
				Message: "Missing query parameter 'signature'",
				Error:   "missing query parameter",
			},
		},
		"missing timestamp": {
			fields: fields{
				mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				pathParam: "0x123",
				queryParams: map[string]string{
					"signature": signature,
					"multisig":  strconv.FormatBool(multisig),
				},
			},
			err: dto.SignavaultError{
				Message: "Missing query parameter 'timestamp'",
				Error:   "missing query parameter",
			},
		},
		"missing multisig": {
			fields: fields{
				mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				pathParam: addr,
				queryParams: map[string]string{
					"signature": signature,
					"timestamp": timestamp,
				},
			},
			err: dto.SignavaultError{
				Message: "Missing query parameter 'multisig'",
				Error:   "missing query parameter",
			},
		},
		"valid args - service success": {
			fields: fields{
				mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				pathParam: addr,
				queryParams: map[string]string{
					"signature": signature,
					"timestamp": timestamp,
					"multisig":  strconv.FormatBool(multisig),
				},
			},
		},
		"valid args - service error": {
			fields: fields{
				mockDepositOfferService,
			},
			args: args{
				r: httptest.NewRecorder(),
				ctx: func(r *httptest.ResponseRecorder) *gin.Context {
					c, _ := gin.CreateTestContext(r)
					return c
				},
				pathParam: addr,
				queryParams: map[string]string{
					"signature": signature,
					"timestamp": timestamp,
					"multisig":  strconv.FormatBool(multisig),
				},
			},
			err: dto.SignavaultError{
				Message: fmt.Sprintf("Error getting all deposit offer signatures for address %s", addr),
				Error:   mockError.Error(),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := &depositOfferHandler{
				DepositOfferService: tt.fields.DepositOfferService,
			}
			ctx := tt.args.ctx(tt.args.r)

			urlQueryStr := "?"
			for k, v := range tt.args.queryParams {
				urlQueryStr += k + "=" + v + "&"
			}
			urlQuery, _ := url.Parse(urlQueryStr[:len(urlQueryStr)-1])

			ctx.Params = gin.Params{
				{
					Key:   "address",
					Value: tt.args.pathParam,
				},
			}
			ctx.Request = &http.Request{
				Method: "GET",
				URL:    urlQuery,
				Header: make(http.Header),
			}
			h.GetSignatures(ctx)

			if tt.err != (dto.SignavaultError{}) {
				var errorResponse dto.SignavaultError
				err := json.Unmarshal(tt.args.r.Body.Bytes(), &errorResponse)
				assert.Nil(t, err)
				assert.Equal(t, tt.err, errorResponse)
				assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
			} else {
				assert.Equal(t, http.StatusOK, ctx.Writer.Status())
			}

		})
	}
}
