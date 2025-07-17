/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	utilsjson "github.com/ava-labs/avalanchego/utils/json"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/avalanchego/vms/platformvm/deposit"

	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
)

var errAliasInfoNotFound = errors.New("could not find address info from node - address does not exist")

type NodeService interface {
	GetMultisigAlias(alias string) (*model.AliasInfo, error)
	IssueTx(txBytes []byte) (ids.ID, error)
	GetAllDepositOffers(args *platformvm.GetAllDepositOffersArgs) ([]*deposit.Offer, error)
	GetDepositOffer(offerID ids.ID, timestamp int64) (*deposit.Offer, error)
}

type nodeService struct {
	config *util.Config
	client platformvm.Client
}

func NewNodeService(config *util.Config) NodeService {
	return &nodeService{
		config: config,
		client: platformvm.NewClient(config.CaminoNode),
	}
}

func (s *nodeService) GetMultisigAlias(alias string) (*model.AliasInfo, error) {
	requestURL := fmt.Sprintf("%s/ext/bc/P", s.config.CaminoNode)
	bodyReader := strings.NewReader(`
			{
				"jsonrpc":"2.0",
				"id":1,
				"method":"platform.getMultisigAlias",
				"params":{
					"Address":"` + alias + `"
				}
			}`)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, errors.New("error creating request: " + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("client: error making http request: " + err.Error())
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("client: could not read response body: " + err.Error())
	}

	var aliasInfo *model.AliasInfo

	err = s.unmarshal(resBody, &aliasInfo)
	if err != nil {
		return nil, errAliasInfoNotFound
	}

	return aliasInfo, nil
}

func (s *nodeService) IssueTx(txBytes []byte) (ids.ID, error) {
	return s.client.IssueTx(context.Background(), txBytes)
}

func (s *nodeService) GetAllDepositOffers(args *platformvm.GetAllDepositOffersArgs) ([]*deposit.Offer, error) {
	return s.client.GetAllDepositOffers(context.Background(), args)
}

func (s *nodeService) GetDepositOffer(offerID ids.ID, timestamp int64) (*deposit.Offer, error) {
	offers, err := s.client.GetAllDepositOffers(context.Background(), &platformvm.GetAllDepositOffersArgs{Timestamp: utilsjson.Uint64(timestamp)})
	if err != nil {
		return nil, err
	}
	for _, offer := range offers {
		if offer.ID == offerID {
			return offer, nil
		}
	}
	return nil, ErrDepositOfferNotFound
}

func (s *nodeService) unmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}
