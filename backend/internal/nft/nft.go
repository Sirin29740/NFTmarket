package nft

type Metadata struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
type NFT struct {
	TokenID     string `json:"tokenId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
