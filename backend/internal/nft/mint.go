package nft

import (
	"NFTmarket/internal/nft/contract"

	"github.com/ethereum/go-ethereum/common"
)

const MyNftAddress = "0x5368c8a780B05f4A5aa267c54503b84Ed688EbE9" // anvil 部署出来的

func MintNFT(to string, tokenURI string) (string, error) {
	client, auth, err := InitRPC()
	if err != nil {
		return "", err
	}

	nftAddr := common.HexToAddress(MyNftAddress)
	instance, err := contract.NewMyNft(nftAddr, client)
	if err != nil {
		return "", err
	}

	tx, err := instance.Mint(auth, common.HexToAddress(to), tokenURI)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
