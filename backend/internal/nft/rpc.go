package nft

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	EthClient *ethclient.Client
	once      sync.Once // 确保只初始化一次
)

func GetClient() (*ethclient.Client, error) {
	var err error
	once.Do(func() {
		// 这里把 https 改成 wss，为了支持我们刚才说的监听功能
		EthClient, err = ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/MJUJlpbdShFQfc2Dt7c4j")
	})
	return EthClient, err
}

func GetAuth() (*bind.TransactOpts, error) {
	privateKey, _ := crypto.HexToECDSA("980e805ac340b7550e552c767baa745a4d2bb0c8499f7d1201be128cbee557c8")
	chainID := big.NewInt(11155111) // Sepolia
	return bind.NewKeyedTransactorWithChainID(privateKey, chainID)
}
