/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package handler

import (
	"bytes"
	"encoding/json"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestCreateMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMultisigService := service.NewMockMultisigService(ctrl)
	h := NewMultisigHandler(mockMultisigService)

	mock := &model.MultisigTx{
		Id:           "1",
		Alias:        "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		UnsignedTx:   "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Threshold:    2,
		OutputOwners: "outputOwners",
		Metadata:     "metadata",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "1",
				Address:      "address",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
		},
	}
	mockMultisigService.EXPECT().CreateMultisigTx(gomock.Any()).Return(mock, nil).AnyTimes()
	mockAsJson, _ := json.Marshal(mock)

	type args struct {
		Body string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody string
		isError  bool
	}{
		{
			name: "create multisig tx",
			args: args{
				Body: ` {
						"unsignedTx": "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
						"alias": "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
						"signature": "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
						"outputOwners": "OutputOwners",
						"metadata": "Metadata",
						"chainId": "11111111111111111111111111111111LpoYY"
						}`,
			},
			wantCode: http.StatusCreated,
			wantBody: string(mockAsJson),
			isError:  false,
		},
		{
			name: "create multisig tx with empty signature - should fail",
			args: args{
				Body: ` {
						"unsignedTx": "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
						"alias": "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
						"outputOwners": "OutputOwners",
						"metadata": "Metadata",
						"chainId": "11111111111111111111111111111111LpoYY"
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: "Error parsing multisig transaction from JSON",
			isError:  true,
		},
		{
			name: "create multisig tx with malformed json - should fail",
			args: args{
				Body: ` {
						"something": "wrong"
						}`,
			},
			wantCode: http.StatusBadRequest,
			wantBody: "Error parsing multisig transaction from JSON",
			isError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := &http.Request{
				Method: "POST",
				Header: make(http.Header),
				Body:   io.NopCloser(bytes.NewBuffer([]byte(tt.args.Body))),
			}
			req.Header.Add("Accept", "application/json")
			c.Request = req

			h.CreateMultisigTx(c)

			assert.Equal(t, tt.wantCode, w.Code)
			if !tt.isError {
				assert.Equal(t, tt.wantBody, w.Body.String())
			} else {
				// check if the error message is in the response
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestGetAllMultisigTxForAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMultisigService := service.NewMockMultisigService(ctrl)
	h := NewMultisigHandler(mockMultisigService)

	now := time.Now()
	mock := model.MultisigTx{
		Id:           "1",
		Alias:        "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		UnsignedTx:   "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Threshold:    2,
		OutputOwners: "outputOwners",
		Metadata:     "metadata",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "1",
				Address:      "address",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
		},
		Timestamp: &now,
	}
	mockResult := make([]model.MultisigTx, 0)
	mockResult = append(mockResult, mock)
	mockMultisigService.EXPECT().GetAllMultisigTxForAlias(mock.Alias, gomock.Any(), mock.Owners[0].Signature).Return(&mockResult, nil).Times(1)
	mockMultisigService.EXPECT().GetAllMultisigTxForAlias("P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saaza", gomock.Any(), gomock.Any()).Return(&[]model.MultisigTx{}, nil).Times(1)
	mockResultAsJson, _ := json.Marshal(mockResult)

	type args struct {
		Alias     string
		Signature string
		Timestamp string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody string
		isError  bool
	}{
		{
			name: "get all multisig tx for alias",
			args: args{
				Alias:     mock.Alias,
				Signature: mock.Owners[0].Signature,
				Timestamp: time.Now().String(),
			},
			wantCode: http.StatusOK,
			wantBody: string(mockResultAsJson),
			isError:  false,
		},
		{
			name: "get all multisig tx for wrong alias - should fail",
			args: args{
				Alias:     "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saaza",
				Signature: "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
				Timestamp: time.Now().String(),
			},
			wantCode: http.StatusOK,
			wantBody: "[]",
			isError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			urlQuery, _ := url.Parse("?signature=" + tt.args.Signature + "&timestamp=" + tt.args.Timestamp)
			req := &http.Request{
				Method: "GET",
				URL:    urlQuery,
				Header: make(http.Header),
				//Body:   io.NopCloser(bytes.NewBuffer([]byte(tt.args.Alias))),
			}
			req.Header.Add("Accept", "application/json")
			c.Request = req
			c.Params = gin.Params{
				{
					Key:   "alias",
					Value: tt.args.Alias,
				},
			}

			h.GetAllMultisigTxForAlias(c)

			assert.Equal(t, tt.wantCode, w.Code)
			if !tt.isError {
				assert.Equal(t, tt.wantBody, w.Body.String())
			} else {
				// check if the error message is in the response
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestIssueMultisigTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMultisigService := service.NewMockMultisigService(ctrl)
	h := NewMultisigHandler(mockMultisigService)

	txId, _ := ids.FromString("ccccc")
	mockResult := &dto.IssueTxResponse{
		TxID: txId.String(),
	}
	resultAsJson, _ := json.Marshal(mockResult)

	req := &dto.IssueTxArgs{
		SignedTx:  "aaaaa",
		Signature: "bbbbb",
	}
	reqAsJson, _ := json.Marshal(req)

	mockMultisigService.EXPECT().IssueMultisigTx(req).Return(txId, nil).AnyTimes()

	type args struct {
		Body string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody string
		isError  bool
	}{
		{
			name: "issue multisig tx",
			args: args{
				Body: string(reqAsJson),
			},
			wantCode: http.StatusOK,
			wantBody: string(resultAsJson),
			isError:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := &http.Request{
				Method: "POST",
				Header: make(http.Header),
				Body:   io.NopCloser(bytes.NewBuffer([]byte(tt.args.Body))),
			}
			req.Header.Add("Accept", "application/json")
			c.Request = req

			h.IssueMultisigTx(c)

			assert.Equal(t, tt.wantCode, w.Code)
			if !tt.isError {
				assert.Equal(t, tt.wantBody, w.Body.String())
			} else {
				// check if the error message is in the response
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}

		})
	}
}

func TestSignMultisigTx(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockMultisigService := service.NewMockMultisigService(ctrl)
	h := NewMultisigHandler(mockMultisigService)

	now := time.Now().UTC()
	mockResult := &model.MultisigTx{
		Id:           "1",
		UnsignedTx:   "000000002004000003ea010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Alias:        "P-kopernikus1k4przmfu79ypp4u7y98glmdpzwk0u3sc7saazy",
		Threshold:    2,
		OutputOwners: "outputOwners",
		Metadata:     "metadata",
		Owners: []model.MultisigTxOwner{
			{
				MultisigTxId: "1",
				Address:      "address",
				Signature:    "4d974561be4675853e0bc6062eac412228e94b16c6ba86dcfedccc1ef2b2a5156ab5aaddbd11f9d88786563fe9f3c17ca5e44a9936621b027b3179284dd86dc000",
			},
		},
		Timestamp: &now,
	}
	resultAsJson, _ := json.Marshal(mockResult)

	req := &dto.SignTxArgs{
		Signature: mockResult.Owners[0].Signature,
	}
	reqAsJson, _ := json.Marshal(req)

	mockMultisigService.EXPECT().SignMultisigTx(mockResult.Id, req).Return(mockResult, nil).Times(1)

	type args struct {
		id   string
		body string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody string
		isError  bool
	}{
		{
			name: "sign multisig tx",
			args: args{
				id:   mockResult.Id,
				body: string(reqAsJson),
			},
			wantCode: http.StatusOK,
			wantBody: string(resultAsJson),
			isError:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := &http.Request{
				Method: "PUT",
				Header: make(http.Header),
				Body:   io.NopCloser(bytes.NewBuffer([]byte(tt.args.body))),
			}
			req.Header.Add("Accept", "application/json")
			c.Request = req
			c.Params = gin.Params{
				{
					Key:   "id",
					Value: tt.args.id,
				},
			}

			h.SignMultisigTx(c)

			assert.Equal(t, tt.wantCode, w.Code)
			if !tt.isError {
				assert.Equal(t, tt.wantBody, w.Body.String())
			} else {
				// check if the error message is in the response
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestNewMultisigHandler(t *testing.T) {

	type args struct {
		multisigService service.MultisigService
	}

	tests := []struct {
		name string
		args args
		want MultisigHandler
	}{
		{
			name: "new multisig handler instance",
			args: args{
				multisigService: service.NewMultisigService(nil, nil, nil),
			},
			want: &multisigHandler{
				multisigService: service.NewMultisigService(nil, nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMultisigHandler(tt.args.multisigService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMultisigHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
