/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chain4travel/camino-signavault/caminogo/ids"
	"github.com/chain4travel/camino-signavault/util"
	"io"
	"net/http"
	"strings"

	"github.com/chain4travel/camino-signavault/caminogo/utils/formatting"
	"github.com/chain4travel/camino-signavault/model"
)

var errAliasInfoNotFound = errors.New("could not find alias info from node - alias does not exist")

type ID [32]byte
type NodeService interface {
	GetMultisigAlias(alias string) (*model.AliasInfo, error)
	IssueTx(txBytes []byte) (string, error)
}

type nodeService struct {
	config *util.Config
	//client platformvm.Client
}

func NewNodeService(config *util.Config) NodeService {
	return &nodeService{
		config: config,
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

//	func (s *nodeService) IssueTx(txBytes []byte) (ids.ID, error) {
//		return s.client.IssueTx(context.Background(), txBytes)
//	}
func (s *nodeService) IssueTx(txBytes []byte) (string, error) {
	requestURL := fmt.Sprintf("%s/ext/bc/P", s.config.CaminoNode)
	txStr, err := formatting.Encode(formatting.Hex, txBytes)
	bodyReader := strings.NewReader(`
			{
				"jsonrpc":"2.0",
				"id":1,
				"method":"platform.issueTx",
				"params":{
					"tx":" ` + txStr + `",
    				"encoding": "hex"
				}
			}`)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", errors.New("error creating request: " + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("client: error making http request: " + err.Error())
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("client: could not read response body: " + err.Error())
	}

	var txId *ids.ID

	err = s.unmarshal(resBody, &txId)
	if err != nil {
		return "", errAliasInfoNotFound
	}

	return txId.String(), nil
}

func (s *nodeService) unmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}

//func (id ID) idToString() string {
//	// We assume that the maximum size of a byte slice that
//	// can be stringified is at least the length of an ID
//	s, _ := cb58.Encode(id[:])
//	return s
//}
