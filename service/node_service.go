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

	"github.com/ava-labs/avalanchego/vms/platformvm"

	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
)

var errAliasInfoNotFound = errors.New("could not find alias info from node - alias does not exist")

type NodeService interface {
	GetMultisigAlias(alias string) (*model.AliasInfo, error)
	IssueTx(txBytes []byte) (ids.ID, error)
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

	err = s.strictUnmarshal(resBody, &aliasInfo)
	if err != nil {
		return nil, errAliasInfoNotFound
	}

	return aliasInfo, nil
}

func (s *nodeService) IssueTx(txBytes []byte) (ids.ID, error) {
	return s.client.IssueTx(context.Background(), txBytes)
}

func (s *nodeService) strictUnmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
