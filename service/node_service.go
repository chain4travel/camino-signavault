package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
)

type NodeServiceInterface interface {
	GetMultisigAlias(alias string) (*model.AliasInfo, error)
	GetTx(txID string) (*model.TxInfo, error)
}

type NodeService struct {
	config *util.Config
}

func NewNodeService(config *util.Config) NodeServiceInterface {
	return &NodeService{
		config: config,
	}
}

func (s *NodeService) GetMultisigAlias(alias string) (*model.AliasInfo, error) {
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

	err = s.strictUnmarshal(resBody, &aliasInfo)
	if err != nil {
		return nil, errors.New("could not unmarshal alias info: " + err.Error())
	}

	return aliasInfo, nil
}

func (s *NodeService) GetTx(txID string) (*model.TxInfo, error) {
	config := util.GetInstance()
	requestURL := fmt.Sprintf("%s/ext/bc/P", config.CaminoNode)
	bodyReader := strings.NewReader(`
			{
				"jsonrpc":"2.0",
				"id":1,
				"method":"platform.getTx",
				"params":{
					"txID":"` + txID + `",
					"encoding": "hex"
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

	var txInfo *model.TxInfo

	err = s.strictUnmarshal(resBody, &txInfo)
	if err != nil {
		return nil, errors.New("could not unmarshal alias info: " + err.Error())
	}

	return txInfo, nil
}

func (s *NodeService) strictUnmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
