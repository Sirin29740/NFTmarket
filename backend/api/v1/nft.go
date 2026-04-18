package v1

import (
	"NFTmarket/internal/nft"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MintReq struct {
	To       string `json:"to"`
	TokenURI string `json:"token_uri"`
}

func Mint(c *gin.Context) {
	var req MintReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	txHash, err := nft.MintNFT(req.To, req.TokenURI)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"tx_hash": txHash,
	})
}
func Getnftlist(c *gin.Context) {
	addr := c.Query("address")
	url := "https://eth-sepolia.g.alchemy.com/nft/v3/MJUJlpbdShFQfc2Dt7c4j/getNFTsForOwner?owner=" + addr

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(500, gin.H{"err": "url error"})
	}
	defer resp.Body.Close()
	var alchemyRes struct {
		OwnedNfts []struct {
			TokenID     uint64 `json:"tokenId"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Image       struct {
				CachedUrl   string `json:"cachedUrl"`
				OriginalUrl string `json:"originalUrl"`
			} `json:"image"`
		} `json:"ownedNfts"`
	}

	json.NewDecoder(resp.Body).Decode(&alchemyRes)
	var nftlist []nft.NFT
	for _, n := range alchemyRes.OwnedNfts {
		nftlist = append(nftlist, nft.NFT{
			TokenID:     n.TokenID,
			Name:        n.Name,
			Description: n.Description,
			Image:       n.Image.OriginalUrl,
		})
	}
	c.JSON(200, nftlist)
}
