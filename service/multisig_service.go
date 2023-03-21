package service

import (
	"errors"
	"log"
	"strconv"

	"github.com/ava-labs/avalanchego/vms/platformvm/txs"

	"github.com/chain4travel/camino-signavault/dao"

	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/chain4travel/camino-signavault/util"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/hashing"

	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
)

var (
	errFailedToVerifyTX = errors.New("failed to verify transaction on chain")
	errTxNotVerified    = errors.New("multisig transaction is not verified on chain")
	errTxNotExists      = errors.New("multisig transaction does not exist")
	errEmptySignature   = errors.New("signature is empty")
	errEmptyTimestamp   = errors.New("timestamp is empty")
	errParsingSignature = errors.New("failed to retrieve address from signature")
	errAddressNotOwner  = errors.New("address is not an owner for this alias")
	errOwnerHasSigned   = errors.New("owner has already signed this alias")
	errThresholdParsing = errors.New("threshold is not a number")
)

const (
	defaultCacheSize = 256
)

type MultisigService interface {
	CreateMultisigTx(multisigTxArgs *dto.MultisigTxArgs) (*model.MultisigTx, error)
	GetAllMultisigTxForAlias(alias string, timestamp string, signature string) (*[]model.MultisigTx, error)
	GetMultisigTx(id int64) (*model.MultisigTx, error)
	SignMultisigTx(id int64, signer *dto.SignTxArgs) (*model.MultisigTx, error)
	CompleteMultisigTx(id int64, completeTx *dto.CompleteTxArgs) (bool, error)
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
	aliasInfo, err := s.getAliasInfo(alias)
	if err != nil {
		return nil, err
	}

	signature := multisigTxArgs.Signature
	unsignedTx := multisigTxArgs.UnsignedTx
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

	id, err := s.dao.CreateMultisigTx(alias, threshold, unsignedTx, creator, signature, owners)
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

	tx, err := s.dao.GetMultisigTx(-1, alias, owner)
	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return &[]model.MultisigTx{}, nil
	}

	return tx, nil
}

func (s *multisigService) GetMultisigTx(id int64) (*model.MultisigTx, error) {
	tx, err := s.dao.GetMultisigTx(id, "", "")
	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, errTxNotExists
	}

	return &(*tx)[0], nil
}

func (s *multisigService) SignMultisigTx(id int64, signer *dto.SignTxArgs) (*model.MultisigTx, error) {
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

func (s *multisigService) CompleteMultisigTx(id int64, completeTx *dto.CompleteTxArgs) (bool, error) {
	multisigTx, err := s.GetMultisigTx(id)
	if err != nil {
		return false, err
	}

	if completeTx.Signature == "" {
		return false, errEmptySignature
	}

	if completeTx.Timestamp == "" {
		return false, errEmptyTimestamp
	}

	signatureArgs := strconv.FormatInt(id, 10) + completeTx.Timestamp + completeTx.TransactionId
	signerAddr, err := s.getAddressFromSignature(signatureArgs, completeTx.Signature, false)
	if err != nil {
		return false, errParsingSignature
	}

	isOwner, _ := s.isOwner(multisigTx, signerAddr)
	if !isOwner {
		return false, errAddressNotOwner
	}

	isTxVerified, err := s.verifyTx(multisigTx, completeTx.TransactionId)
	if err != nil {
		log.Print(err)
		return false, errFailedToVerifyTX
	}
	if !isTxVerified {
		return false, errTxNotVerified
	}

	completed, err := s.dao.UpdateTransactionId(id, completeTx.TransactionId)
	if err != nil {
		return false, err
	}
	return completed, nil
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
	if isHex {
		signatureArgsBytes, _ = formatting.Decode(formatting.Hex, signatureArgs)
	} else {
		signatureArgsBytes = []byte(signatureArgs)
	}
	signatureArgsHash := hashing.ComputeHash256(signatureArgsBytes)
	signatureBytes, _ := formatting.Decode(formatting.Hex, signature)

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

func (s *multisigService) verifyTx(multisigTx *model.MultisigTx, txID string) (bool, error) {
	txRes, err := s.nodeService.GetTx(txID)
	if err != nil {
		return false, err
	}

	txBytes, err := formatting.Decode(formatting.Hex, txRes.Result.Tx)
	if err != nil {
		return false, err
	}

	var tx txs.Tx
	_, err = txs.Codec.Unmarshal(txBytes, &tx)
	if err != nil {
		return false, err
	}

	utxBytes := tx.Unsigned.Bytes()
	utxString, err := formatting.Encode(formatting.Hex, utxBytes)
	if err != nil {
		return false, err
	}

	return utxString == multisigTx.UnsignedTx, nil
}
