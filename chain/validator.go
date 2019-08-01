package chain

import (
	"github.com/drep-project/drep-chain/database"
	"github.com/drep-project/drep-chain/types"
	"math/big"
)

type BlockExecuteContext struct {
	Db      *database.Database
	Block   *types.Block
	Gp      *GasPool
	GasUsed *big.Int
	GasFee  *big.Int
}

func (blockExecuteContext *BlockExecuteContext) AddGasUsed(gas *big.Int) {
	blockExecuteContext.GasUsed = blockExecuteContext.GasUsed.Add(blockExecuteContext.GasUsed, gas)
}

func (blockExecuteContext *BlockExecuteContext) AddGasFee(fee *big.Int) {
	blockExecuteContext.GasFee = blockExecuteContext.GasFee.Add(blockExecuteContext.GasFee, fee)
}

type IBlockValidator interface {
	VerifyHeader(header, parent *types.BlockHeader) error

	VerifyBody(block *types.Block) error

	ExecuteBlock(context *BlockExecuteContext) error
}

type ITransactionValidator interface {
	ExecuteTransaction(db *database.Database, tx *types.Transaction, gp *GasPool, header *types.BlockHeader) (*types.Receipt, *big.Int, *big.Int, error)
}
