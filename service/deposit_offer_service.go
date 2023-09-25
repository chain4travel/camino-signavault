package service

import (
	"errors"
	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/ethereum/go-ethereum/common"
)

var _ DepositOfferService = (*depositOfferService)(nil)

type DepositOfferService interface {
	AddSignature(args *dto.AddSignatureArgs) error
	GetSignatures(address, timestamp, signature string) (*[]model.DepositOfferSig, error)
}

type depositOfferService struct {
	config      *util.Config
	secpFactory secp256k1.Factory
	dao         dao.DepositOfferDao
	nodeService NodeService
}

var (
	ErrParsingDepositOfferID = errors.New("error parsing deposit offer id")
	ErrParsingAddress        = errors.New("error parsing address")
	ErrDepositOfferNotFound  = errors.New("deposit offer not found")
	ErrInvalidSignature      = errors.New("invalid signature")
)

func NewDepositOfferService(config *util.Config, dao dao.DepositOfferDao, nodeService NodeService) DepositOfferService {
	return &depositOfferService{
		config: config,
		secpFactory: secp256k1.Factory{
			Cache: cache.LRU[ids.ID, *secp256k1.PublicKey]{Size: defaultCacheSize},
		},
		dao:         dao,
		nodeService: nodeService,
	}
}

func (s *depositOfferService) AddSignature(args *dto.AddSignatureArgs) error {
	id, err := ids.FromString(args.DepositOfferID)
	if err != nil {
		return ErrParsingDepositOfferID
	}
	addr, err := ids.ShortFromString(args.Address)
	if err != nil {
		return ErrParsingAddress
	}
	signatureArgs := append(id[:], addr[:]...)

	signer, err := s.getAddressFromSignature(signatureArgs, args.Signature)
	if err != nil {
		return ErrParsingSignature
	}

	reply, err := s.nodeService.GetAllDepositOffers(&platformvm.GetAllDepositOffersArgs{})
	reply.DepositOffers = append(reply.DepositOffers, &platformvm.APIDepositOffer{ID: id, OwnerAddress: signer})
	var depositOffer *platformvm.APIDepositOffer
	for _, do := range reply.DepositOffers {
		if do.ID == id {
			depositOffer = do
		}
	}
	if depositOffer == nil {
		return ErrDepositOfferNotFound
	}
	if depositOffer.OwnerAddress != signer {
		return ErrInvalidSignature
	}
	err = s.dao.AddSignature(args.DepositOfferID, args.Address, args.Signature)
	if err != nil {
		return err
	}
	return nil
}

func (s *depositOfferService) GetSignatures(address, timestamp, signature string) (*[]model.DepositOfferSig, error) {
	addr, err := ids.ShortFromString(address)
	if err != nil {
		return nil, ErrParsingAddress
	}
	signatureArgs := append(addr[:], []byte(timestamp)...)
	sigOwner, err := s.getAddressFromSignature(signatureArgs, signature)
	if err != nil {
		return nil, ErrParsingSignature
	}
	if addr != sigOwner {
		return nil, ErrInvalidSignature
	}
	return s.dao.GetSignatures(address)
}

func (s *depositOfferService) getAddressFromSignature(signatureArgs []byte, signature string) (ids.ShortID, error) {
	signatureArgsHash := hashing.ComputeHash256(signatureArgs)
	signatureBytes := common.FromHex(signature)

	pub, err := s.secpFactory.RecoverHashPublicKey(signatureArgsHash, signatureBytes)
	if err != nil {
		return ids.ShortEmpty, err
	}

	return pub.Address(), nil
}
