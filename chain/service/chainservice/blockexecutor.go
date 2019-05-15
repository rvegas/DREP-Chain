package service

import (
	"bytes"
	"fmt"
	"github.com/drep-project/drep-chain/common"
	"math/big"

	"github.com/drep-project/drep-chain/chain/params"
	chainTypes "github.com/drep-project/drep-chain/chain/types"
	"github.com/drep-project/drep-chain/crypto/secp256k1"
	"github.com/drep-project/drep-chain/crypto/secp256k1/schnorr"
	"github.com/drep-project/drep-chain/crypto/sha3"
	"github.com/drep-project/drep-chain/database"
)

type ChainBlockValidator struct {
	txValidator ITransactionValidator
	chain *ChainService
}

func NewChainBlockValidator(chainService *ChainService) *ChainBlockValidator {
	return &ChainBlockValidator{
		txValidator: chainService.TransactionValidator,
		chain: chainService,
	}
}

func (chainBlockValidator *ChainBlockValidator) VerifyHeader(header, parent *chainTypes.BlockHeader) error {
	// Verify chainId  matched
	if header.ChainId != chainBlockValidator.chain.ChainID() {
		return ErrChainId
	}
	// Verify version  mathch
	if header.Version != common.Version {
		return ErrVersion
	}
	//Verify header's previousHash is equal parent hash
	if header.PreviousHash !=  *parent.Hash() {
		return ErrPreHash
	}
	// Verify that the block number is parent's +1
	if header.Height-parent.Height != 1 {
		return ErrInvalidateBlockNumber
	}
	// pre block timestamp before this block time
	if header.Timestamp <= parent.Timestamp {
		return ErrInvalidateTimestamp
	}

	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit.Uint64() > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed.Uint64() > header.GasLimit.Uint64() {
		return fmt.Errorf("invalid gasUsed: have %v, gasLimit %v", header.GasUsed, header.GasLimit)
	}

	//TODO Verify that the gas limit remains within allowed bounds
	nextGasLimit := chainBlockValidator.chain.CalcGasLimit(parent, params.MinGasLimit, params.MaxGasLimit)
	if nextGasLimit.Cmp(&header.GasLimit) != 0 {
		return fmt.Errorf("invalid gas limit: have %v, want %v += %v", header.GasLimit, parent.GasLimit, nextGasLimit)
	}
	// check multisig
	// leader
	if !chainBlockValidator.isInLocalBp(&header.LeaderPubKey) {
		return ErrBpNotInList
	}
	// minor
	for _, minor := range header.MinorPubKeys {
		if !chainBlockValidator.isInLocalBp(&minor) {
			return ErrBpNotInList
		}
	}
	return nil
}

// isInLocalBp check the specific pubket  is a bp node
func (chainBlockValidator *ChainBlockValidator) isInLocalBp(key *secp256k1.PublicKey) bool {
	for _, bp := range chainBlockValidator.chain.Config.Producers {
		if bp.Pubkey.IsEqual(key) {
			return true
		}
	}
	return false
}

func (chainBlockValidator *ChainBlockValidator) VerifyBody(block *chainTypes.Block) error {

	if !(chainBlockValidator.VerifyMultiSig(block, chainBlockValidator.chain.Config.SkipCheckMutiSig || false)) {
		return ErrInvalidateBlockMultisig
	}

	// Header validity is known at this point, check the uncles and transactions
	header := block.Header
	if hash := chainBlockValidator.chain.DeriveMerkleRoot(block.Data.TxList); !bytes.Equal(hash, header.TxRoot) {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxRoot)
	}
	return nil
}

func (chainBlockValidator *ChainBlockValidator) VerifyMultiSig(b *chainTypes.Block, skipCheckSig bool) bool {
	if skipCheckSig { //just for solo
		return true
	}
	participators := []*secp256k1.PublicKey{}
	for index, val := range b.MultiSig.Bitmap {
		if val == 1 {
			producer := chainBlockValidator.chain.Config.Producers[index]
			participators = append(participators, producer.Pubkey)
		}
	}
	msg := b.ToMessage()
	sigmaPk := schnorr.CombinePubkeys(participators)
	return schnorr.Verify(sigmaPk, sha3.Keccak256(msg), b.MultiSig.Sig.R, b.MultiSig.Sig.S)
}

func (chainBlockValidator *ChainBlockValidator) ExecuteBlock(db *database.Database, block *chainTypes.Block, gp *GasPool) (*big.Int, *big.Int, error) {
	totalGasFee := big.NewInt(0)
	totalGasUsed := big.NewInt(0)
	if len(block.Data.TxList) < 0 {
		return totalGasUsed, totalGasFee, nil
	}
	for _, t := range block.Data.TxList {
		gasUsed, gasFee, err := chainBlockValidator.txValidator.ExecuteTransaction(db, t, gp, block.Header)
		if err != nil {
			return nil, nil, err
			//dlog.Debug("execute transaction fail", "txhash", t.Data, "reason", err.Error())
		}
		if gasFee != nil {
			totalGasFee.Add(totalGasFee, gasFee)
			totalGasUsed.Add(totalGasUsed, gasUsed)
		}
	}
	return totalGasUsed, totalGasFee, nil
}

type IBlockValidator interface {

	VerifyHeader(header, parent *chainTypes.BlockHeader) error

	VerifyBody(block *chainTypes.Block) error

	VerifyMultiSig(b *chainTypes.Block, skipCheckSig bool) bool

	ExecuteBlock(db *database.Database, block *chainTypes.Block, gp *GasPool) (*big.Int, *big.Int, error)
}