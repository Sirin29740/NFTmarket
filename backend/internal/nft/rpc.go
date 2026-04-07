package nft

import (
	"NFTmarket/internal/nft/contract"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	EthClient   *ethclient.Client
	NFTContract *contract.MyNft
	Auth        *bind.TransactOpts
)

func InitRPC() (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/MJUJlpbdShFQfc2Dt7c4j")
	if err != nil {
		return nil, nil, err
	}

	privateKey, _ := crypto.HexToECDSA("980e805ac340b7550e552c767baa745a4d2bb0c8499f7d1201be128cbee557c8")
	chainID := big.NewInt(11155111)

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.GasLimit = 300000

	return client, auth, nil
}
