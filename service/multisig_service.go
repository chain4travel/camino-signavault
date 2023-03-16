package service

import (
	"errors"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/util"
	"log"
	"strconv"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/hashing"

	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
)

const (
	defaultCacheSize = 256
)

type MultisigService struct {
	db          db.Db
	SECPFactory crypto.FactorySECP256K1R
	dao         *dao.MultisigTxDao
}

func NewMultisigService(db db.Db) *MultisigService {
	txDao := dao.NewMultisigTxDao(db)

	return &MultisigService{
		db: db,
		SECPFactory: crypto.FactorySECP256K1R{
			Cache: cache.LRU{Size: defaultCacheSize},
		},
		dao: txDao,
	}
}

func (s *MultisigService) CreateMultisigTx(multisigTxArgs *dto.MultisigTxArgs) (*model.MultisigTx, error) {
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
		return nil, errors.New("failed to retrieve address from signature")
	}
	threshold, err := strconv.Atoi(aliasInfo.Result.Threshold)
	if err != nil {
		return nil, errors.New("threshold is not a number")
	}
	owners := aliasInfo.Result.Addresses

	if !s.isCreatorOwner(owners, creator) {
		return nil, errors.New("creator of multisig transaction is not an owner")
	}

	id, err := s.dao.CreateMultisigTx(alias, threshold, unsignedTx, creator, signature, owners)
	if err != nil {
		return nil, err
	}
	return s.GetMultisigTx(id)
}

func (s *MultisigService) GetAllMultisigTxForAlias(alias string, timestamp string, signature string) (*[]model.MultisigTx, error) {
	signatureArgs := alias + timestamp
	owner, err := s.getAddressFromSignature(signatureArgs, signature, false)
	if err != nil {
		return nil, errors.New("failed to retrieve address from signature")
	}

	tx, err := s.dao.GetMultisigTx(-1, alias, owner)
	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return &[]model.MultisigTx{}, nil
	}

	//ts, err := strconv.ParseInt(timestamp, 10, 64)
	//if err != nil {
	//	return nil, errors.New("error parsing timestamp")
	//}

	return tx, nil
}

func (s *MultisigService) CompleteMultisigTx(id int64, completeTx *dto.CompleteTxArgs) (bool, error) {
	multisigTx, err := s.GetMultisigTx(id)
	if err != nil {
		return false, err
	}

	if completeTx.Signature == "" {
		return false, errors.New("signature is empty")
	}

	if completeTx.Timestamp == "" {
		return false, errors.New("timestamp is empty")
	}

	signatureArgs := multisigTx.Alias + completeTx.Timestamp
	signerAddr, err := s.getAddressFromSignature(signatureArgs, completeTx.Signature, false)
	if err != nil {
		return false, errors.New("failed to retrieve address from signature")
	}

	isOwner, _ := s.isOwner(multisigTx, signerAddr)
	if !isOwner {
		return false, errors.New("address is not an owner address for this alias")
	}

	completed, err := s.dao.UpdateTransactionId(id, completeTx.TransactionId)
	if err != nil {
		return false, err
	}
	return completed, nil
}

func (s *MultisigService) GetMultisigTx(id int64) (*model.MultisigTx, error) {
	tx, err := s.dao.GetMultisigTx(id, "", "")
	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, nil
	}

	return &(*tx)[0], nil
}

func (s *MultisigService) SignMultisigTx(id int64, signer *dto.SignTxArgs) (*model.MultisigTx, error) {
	multisigTx, err := s.GetMultisigTx(id)
	if err != nil {
		return nil, err
	}

	if signer.Signature == "" {
		return nil, errors.New("signature is empty")
	}

	if signer.Timestamp == "" {
		return nil, errors.New("timestamp is empty")
	}

	signatureArgs := multisigTx.Alias + signer.Timestamp
	signerAddr, err := s.getAddressFromSignature(signatureArgs, signer.Signature, false)
	if err != nil {
		return nil, errors.New("failed to retrieve address from signature")
	}

	isOwner, isSigner := s.isOwner(multisigTx, signerAddr)
	if !isOwner {
		return nil, errors.New("address is not an owner for this alias")
	}
	if isSigner {
		return nil, errors.New("owner has already signed this alias")
	}

	_, err = s.dao.AddSigner(id, signerAddr, signer.Signature)
	if err != nil {
		return nil, err
	}

	return s.GetMultisigTx(id)
}

func (s *MultisigService) isOwner(multisigTx *model.MultisigTx, address string) (bool, bool) {
	for _, owner := range multisigTx.Owners {
		if owner.Address == address {
			return true, owner.Signature != ""
		}
	}
	return false, false
}

func (s *MultisigService) isCreatorOwner(owners []string, address string) bool {
	for _, owner := range owners {
		if owner == address {
			return true
		}
	}
	return false
}

func (s *MultisigService) getAliasInfo(alias string) (*model.AliasInfo, error) {
	nodeService := NewNodeService()
	aliasInfo, err := nodeService.GetMultisigAlias(alias)
	if err != nil {
		log.Printf("Getting info for alias %s failed: %v", alias, err)
		return nil, err
	}
	return aliasInfo, nil
}

func (s *MultisigService) getAddressFromSignature(signatureArgs string, signature string, isHex bool) (string, error) {
	var signatureArgsBytes []byte
	if isHex {
		signatureArgsBytes, _ = formatting.Decode(formatting.Hex, signatureArgs)
	} else {
		signatureArgsBytes = []byte(signatureArgs)
	}
	signatureArgsHash := hashing.ComputeHash256(signatureArgsBytes)
	signatureBytes, _ := formatting.Decode(formatting.Hex, signature)

	pub, err := s.SECPFactory.RecoverHashPublicKey(signatureArgsHash, signatureBytes)
	if err != nil {
		return "", err
	}

	config := util.GetInstance()
	hrp := constants.NetworkIDToHRP[config.NetworkId]
	bech32Address, err := address.FormatBech32(hrp, pub.Address().Bytes())
	if err != nil {
		return "", err
	}

	return "P-" + bech32Address, nil
}
