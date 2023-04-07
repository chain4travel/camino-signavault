/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

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
	ErrTxNotExists      = errors.New("multisig transaction does not exist")
	ErrEmptySignature   = errors.New("signature is empty")
	ErrParsingSignature = errors.New("failed to retrieve address from signature")
	ErrAddressNotOwner  = errors.New("address is not an owner for this alias")
	ErrOwnerHasSigned   = errors.New("owner has already signed this alias")
	ErrThresholdParsing = errors.New("threshold is not a number")
	ErrParsingTx        = errors.New("error parsing signed tx")
	ErrPendingTx        = errors.New("there is already a pending tx for this alias")
	ErrExpired          = errors.New("expiration date has passed")
)

const (
	defaultCacheSize      = 256
	defaultExpirationDays = 14
)

// Wraps the UnsignedTx to force marshalling typeID
type codecWrapper = struct {
	txs.UnsignedTx `serialize:"true"`
}

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
		return nil, ErrPendingTx
	}

	aliasInfo, err := s.getAliasInfo(alias)
	if err != nil {
		return nil, err
	}

	// expiration date
	var expiresAt *time.Time
	var t time.Time
	exp := multisigTxArgs.Expiration
	now := time.Now().UTC()
	expirationDays := s.config.TxExpiration
	// if the value is 0, use the default expiration
	if expirationDays <= 0 {
		expirationDays = defaultExpirationDays
	}
	if exp == 0 {
		t = now.Add(time.Hour * 24 * time.Duration(expirationDays))
	} else {
		t = time.Unix(exp, 0).UTC()

	}
	expiresAt = &t
	if expiresAt != nil && expiresAt.Before(now) {
		return nil, ErrExpired
	}

	signature := multisigTxArgs.Signature
	unsignedTx := multisigTxArgs.UnsignedTx
	outputOwners := multisigTxArgs.OutputOwners
	metadata := multisigTxArgs.Metadata
	creator, err := s.getAddressFromSignature(unsignedTx, signature, true)
	if err != nil {
		return nil, ErrParsingSignature
	}
	threshold, err := strconv.Atoi(aliasInfo.Result.Threshold)
	if err != nil {
		return nil, ErrThresholdParsing
	}
	owners := aliasInfo.Result.Addresses
	chainId := multisigTxArgs.ChainId

	if !s.isCreatorOwner(owners, creator) {
		return nil, ErrAddressNotOwner
	}
	// generate txId by hasing the unsignedTx
	id, err := s.generateId(unsignedTx)
	if err != nil {
		return nil, err
	}

	_, err = s.dao.CreateMultisigTx(id, alias, threshold, chainId, unsignedTx, creator, signature, outputOwners, metadata, owners, expiresAt)
	if err != nil {
		return nil, err
	}
	return s.GetMultisigTx(id)
}

func (s *multisigService) GetAllMultisigTxForAlias(alias string, timestamp string, signature string) (*[]model.MultisigTx, error) {
	signatureArgs := alias + timestamp
	owner, err := s.getAddressFromSignature(signatureArgs, signature, false)
	if err != nil {
		return nil, ErrParsingSignature
	}

	tx, err := s.dao.GetMultisigTx("", alias, owner)
	if err != nil {
		return nil, fmt.Errorf("couldn't get txs for alias %s: %w", alias, err)
	}

	if tx == nil || len(*tx) <= 0 {
		return &[]model.MultisigTx{}, nil
	}

	return tx, nil
}

func (s *multisigService) GetMultisigTx(id string) (*model.MultisigTx, error) {
	tx, err := s.dao.GetMultisigTx(id, "", "")
	if err != nil {
		return nil, err
	}
	if tx == nil || len(*tx) <= 0 {
		return nil, ErrTxNotExists
	}

	return &(*tx)[0], nil
}

func (s *multisigService) SignMultisigTx(id string, signer *dto.SignTxArgs) (*model.MultisigTx, error) {
	multisigTx, err := s.GetMultisigTx(id)
	if err != nil {
		return nil, err
	}

	if signer.Signature == "" {
		return nil, ErrEmptySignature
	}

	signerAddr, err := s.getAddressFromSignature(multisigTx.UnsignedTx, signer.Signature, true)
	if err != nil {
		return nil, ErrParsingSignature
	}

	isOwner, isSigner := s.isOwner(multisigTx, signerAddr)
	if !isOwner {
		return nil, ErrAddressNotOwner
	}
	if isSigner {
		return nil, ErrOwnerHasSigned
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

	utxBytes, _ := txs.Codec.Marshal(txs.Version, codecWrapper{tx.Unsigned})
	utxHash := hashing.ComputeHash256(utxBytes)
	utxHashStr := fmt.Sprintf("%x", utxHash)

	storedTx, err := s.GetMultisigTx(utxHashStr)
	if err != nil {
		return ids.Empty, err
	}

	signerAddr, err := s.getAddressFromSignature(sendTxArgs.SignedTx, sendTxArgs.Signature, true)
	if err != nil {
		return ids.Empty, ErrParsingSignature
	}

	isOwner, _ := s.isOwner(storedTx, signerAddr)
	if !isOwner {
		return ids.Empty, ErrAddressNotOwner
	}

	signedBytes, err := txs.Codec.Marshal(txs.Version, tx)
	if err != nil {
		return ids.Empty, ErrParsingTx
	}

	txID, err := s.nodeService.IssueTx(signedBytes)
	if err != nil {
		return ids.Empty, err
	}
	_, _ = s.dao.UpdateTransactionId(utxHashStr, txID.String())
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
		signatureArgsBytes = common.FromHex(signatureArgs)
	} else {
		signatureArgsBytes = []byte(signatureArgs)
	}

	signatureArgsHash := hashing.ComputeHash256(signatureArgsBytes)
	signatureBytes := common.FromHex(signature)

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
	txBytes := common.FromHex(txHexString)

	_, err := txs.Codec.Unmarshal(txBytes, &tx)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (s *multisigService) generateId(unsignedTx string) (string, error) {
	txBytes := common.FromHex(unsignedTx)
	return fmt.Sprintf("%x", hashing.ComputeHash256(txBytes)), nil
}
