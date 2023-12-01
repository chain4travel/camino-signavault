package service

import (
	"errors"
	"time"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/json"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"

	addr_utils "github.com/ava-labs/avalanchego/utils/formatting/address"
)

var _ DepositOfferService = (*depositOfferService)(nil)

type DepositOfferService interface {
	AddSignatures(args *dto.AddSignatureArgs) error
	GetSignatures(address, timestamp, signature string, multisig bool) (*[]model.DepositOfferSig, error)
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
	ErrNoAliasFound          = errors.New("no alias found for given address")
	ErrAddressesSigsMismatch = errors.New("number of addresses does not match number of signatures")
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

func (s *depositOfferService) AddSignatures(args *dto.AddSignatureArgs) error {
	if len(args.Addresses) != len(args.Signatures) {
		return ErrAddressesSigsMismatch
	}

	id, err := ids.FromString(args.DepositOfferID)
	if err != nil {
		return ErrParsingDepositOfferID
	}

	// if no timestamp is provided, use current time
	t := json.Uint64(time.Now().Unix())
	if args.Timestamp != 0 {
		t = json.Uint64(args.Timestamp)
	}

	reply, err := s.nodeService.GetAllDepositOffers(&platformvm.GetAllDepositOffersArgs{Timestamp: t})
	if err != nil {
		return err
	}
	var depositOffer *platformvm.APIDepositOffer
	for _, do := range reply.DepositOffers {
		if do.ID == id {
			depositOffer = do
			break
		}
	}
	if depositOffer == nil {
		return ErrDepositOfferNotFound
	}

	for i, a := range args.Addresses {
		addr, err := ids.ShortFromString(a)
		if err != nil {
			return ErrParsingAddress
		}
		signatureArgs := append(id[:], addr[:]...)
		signer, err := s.getAddressFromSignature(signatureArgs, args.Signatures[i])
		if err != nil {
			return ErrParsingSignature
		}
		if depositOffer.OwnerAddress != signer {
			return ErrInvalidSignature
		}

	}

	err = s.dao.AddSignatures(args.DepositOfferID, args.Addresses, args.Signatures)
	if err != nil {
		return err
	}
	return nil
}

func (s *depositOfferService) GetSignatures(address, timestamp, signature string, multisig bool) (*[]model.DepositOfferSig, error) {
	addr, err := ids.ShortFromString(address)
	if err != nil {
		return nil, ErrParsingAddress
	}
	signatureArgs := append(addr[:], []byte(timestamp)...)
	sigOwner, err := s.getAddressFromSignature(signatureArgs, signature)
	if err != nil {
		return nil, ErrParsingSignature
	}

	// if address is singlesig, check if it matches signature owner
	if !multisig && addr != sigOwner {
		return nil, ErrInvalidSignature
	} else if multisig {
		aliasInfo, err := s.nodeService.GetMultisigAlias(address)
		if err != nil {
			return nil, err
		} else if aliasInfo.Result.Addresses == nil {
			return nil, ErrNoAliasFound
		}

		signer, err := addr_utils.Format(util.PChainAlias, constants.GetHRP(s.config.NetworkId), sigOwner.Bytes())
		if err != nil {
			return nil, err
		}
		// if signer address is not part of multisig alias return error
		if !slices.Contains(aliasInfo.Result.Addresses, signer) {
			return nil, ErrInvalidSignature
		}
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
