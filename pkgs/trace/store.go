package trace

import (
	chainTypes "github.com/drep-project/drep-chain/types"
	"github.com/drep-project/drep-chain/crypto"
)

// IStore determine the interface to be implementation for storage
type IStore interface {
	ExistRecord(block *chainTypes.Block) (bool, error)

	InsertRecord(block *chainTypes.Block)

	DelRecord(block *chainTypes.Block)

	GetRawTransaction(txHash *crypto.Hash) ([]byte, error)

	GetTransaction(txHash *crypto.Hash) (*chainTypes.RpcTransaction, error)

	GetSendTransactionsByAddr(addr *crypto.CommonAddress, pageIndex, pageSize int) []*chainTypes.RpcTransaction

	GetReceiveTransactionsByAddr(addr *crypto.CommonAddress, pageIndex, pageSize int) []*chainTypes.RpcTransaction

	Close()
}
