package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
)

var (
	errTxNotExists      = errors.New("multisig transaction does not exist")
	errEmptySignature   = errors.New("signature is empty")
	errParsingSignature = errors.New("failed to retrieve address from signature")
	errAddressNotOwner  = errors.New("address is not an owner for this alias")
	errOwnerHasSigned   = errors.New("owner has already signed this alias")
	errThresholdParsing = errors.New("threshold is not a number")
	errPendingTx        = errors.New("there is already a pending tx for this alias")
)

const (
	defaultCacheSize = 256
)

type MultisigService interface {
	CreateMultisigTx(multisigTxArgs *dto.MultisigTxArgs) (*model.MultisigTx, error)
	GetAllMultisigTxForAlias(alias string, timestamp string, signature string) (*[]model.MultisigTx, error)
	GetMultisigTx(id string) (*model.MultisigTx, error)
	SignMultisigTx(id string, signer *dto.SignTxArgs) (*model.MultisigTx, error)
	IssueMultisigTx(issueTxArgs *dto.IssueTxArgs) (ids.ID, error)
}

type multisigService struct {
	config      *util.Config
	secpFactory crypto.FactorySECP256K1R
	dao         dao.MultisigTxDao
	nodeService NodeService
}

func NewMultisigService(config *util.Config, dao dao.MultisigTxDao, nodeService NodeService) MultisigService {
	return &multisigService{
		config: config,
		secpFactory: crypto.FactorySECP256K1R{
			Cache: cache.LRU{Size: defaultCacheSize},
		},
		dao:         dao,
		nodeService: nodeService,
	}
}

func (s *multisigService) CreateMultisigTx(multisigTxArgs *dto.MultisigTxArgs) (*model.MultisigTx, error) {
	var err error

	alias := multisigTxArgs.Alias

	exists, err := s.dao.PendingAliasExists(alias)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errPendingTx
	}

	aliasInfo, err := s.getAliasInfo(alias)
	if err != nil {
		return nil, err
	}

	signature := multisigTxArgs.Signature
	unsignedTx := multisigTxArgs.UnsignedTx
	outputOwners := multisigTxArgs.OutputOwners
	creator, err := s.getAddressFromSignature(unsignedTx, signature, true)
	if err != nil {
		return nil, errParsingSignature
	}
	threshold, err := strconv.Atoi(aliasInfo.Result.Threshold)
	if err != nil {
		return nil, errThresholdParsing
	}
	owners := aliasInfo.Result.Addresses

	if !s.isCreatorOwner(owners, creator) {
		return nil, errAddressNotOwner
	}
	// generate txId by hasing the unsignedTx
	id := fmt.Sprintf("%x", hashing.ComputeHash256([]byte(unsignedTx)))

	_, err = s.dao.CreateMultisigTx(id, alias, threshold, unsignedTx, creator, signature, outputOwners, owners)
	if err != nil {
		return nil, err
	}
	return s.GetMultisigTx(id)
}

func (s *multisigService) GetAllMultisigTxForAlias(alias string, timestamp string, signature string) (*[]model.MultisigTx, error) {
	signatureArgs := alias + timestamp
	owner, err := s.getAddressFromSignature(signatureArgs, signature, false)
	if err != nil {
		return nil, errParsingSignature
	}

	tx, err := s.dao.GetMultisigTx("", alias, owner)
	if err != nil {
		return nil, fmt.Errorf("couldn't get txs for alias %s: %w", alias, err)
	}
	if len(*tx) <= 0 {
		return &[]model.MultisigTx{}, nil
	}

	return tx, nil
}

func (s *multisigService) GetMultisigTx(id string) (*model.MultisigTx, error) {
	tx, err := s.dao.GetMultisigTx(id, "", "")
	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, errTxNotExists
	}

	return &(*tx)[0], nil
}

func (s *multisigService) SignMultisigTx(id string, signer *dto.SignTxArgs) (*model.MultisigTx, error) {
	multisigTx, err := s.GetMultisigTx(id)
	if err != nil {
		return nil, err
	}

	if signer.Signature == "" {
		return nil, errEmptySignature
	}

	signerAddr, err := s.getAddressFromSignature(multisigTx.UnsignedTx, signer.Signature, true)
	if err != nil {
		return nil, errParsingSignature
	}

	isOwner, isSigner := s.isOwner(multisigTx, signerAddr)
	if !isOwner {
		return nil, errAddressNotOwner
	}
	if isSigner {
		return nil, errOwnerHasSigned
	}

	_, err = s.dao.AddSigner(id, signer.Signature, signerAddr)
	if err != nil {
		return nil, err
	}

	return s.GetMultisigTx(id)
}

func (s *multisigService) IssueMultisigTx(sendTxArgs *dto.IssueTxArgs) (ids.ID, error) {
	tx, err := s.unmarshalTx(sendTxArgs.SignedTx)
	if err != nil {
		return ids.Empty, err
	}

	utxHash := hashing.ComputeHash256(tx.Unsigned.Bytes())
	utxHashStr := fmt.Sprintf("%x", utxHash)

	storedTx, err := s.GetMultisigTx(utxHashStr)
	if err != nil {
		return ids.Empty, err
	}

	signerAddr, err := s.getAddressFromSignature(storedTx.UnsignedTx, sendTxArgs.Signature, true)
	if err != nil {
		return ids.Empty, errParsingSignature
	}

	isOwner, _ := s.isOwner(storedTx, signerAddr)
	if !isOwner {
		return ids.Empty, errAddressNotOwner
	}

	txID, err := s.nodeService.IssueTx(tx.Bytes())
	_, err = s.dao.UpdateTransactionId(utxHashStr, txID.String())
	if err != nil {
		return ids.Empty, err
	}
	return txID, nil
}

func (s *multisigService) isOwner(multisigTx *model.MultisigTx, address string) (bool, bool) {
	for _, owner := range multisigTx.Owners {
		if owner.Address == address {
			return true, owner.Signature != ""
		}
	}
	return false, false
}

func (s *multisigService) isCreatorOwner(owners []string, address string) bool {
	for _, owner := range owners {
		if owner == address {
			return true
		}
	}
	return false
}

func (s *multisigService) getAliasInfo(alias string) (*model.AliasInfo, error) {
	aliasInfo, err := s.nodeService.GetMultisigAlias(alias)
	if err != nil {
		log.Printf("Getting info for alias %s failed: %v", alias, err)
		return nil, err
	}
	return aliasInfo, nil
}

func (s *multisigService) getAddressFromSignature(signatureArgs string, signature string, isHex bool) (string, error) {
	var signatureArgsBytes []byte
	var err error
	if isHex {
		signatureArgsBytes = common.Hex2Bytes(signatureArgs)
	} else {
		signatureArgsBytes = []byte(signatureArgs)
	}
	signatureArgsHash := hashing.ComputeHash256(signatureArgsBytes)
	signatureBytes := common.Hex2Bytes(signature)

	pub, err := s.secpFactory.RecoverHashPublicKey(signatureArgsHash, signatureBytes)
	if err != nil {
		return "", err
	}

	hrp := constants.NetworkIDToHRP[s.config.NetworkId]
	bech32Address, err := address.FormatBech32(hrp, pub.Address().Bytes())
	if err != nil {
		return "", err
	}

	return "P-" + bech32Address, nil
}

func (s *multisigService) unmarshalTx(txHexString string) (txs.Tx, error) {
	var tx txs.Tx
	txBytes := common.Hex2Bytes(txHexString)

	_, err := txs.Codec.Unmarshal(txBytes, &tx)
	if err != nil {
		return tx, err
	}

	return tx, nil
}
