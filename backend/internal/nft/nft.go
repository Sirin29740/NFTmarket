package nft

import "gorm.io/gorm"

type Metadata struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
type NFT struct {
	gorm.Model
	TokenID     uint64 `Gorm:"uniqueIndex" json:"tokenId"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	LocalPath   string `json:"localPath"`
	SyncStatus  int    `json:"syncStatus"`
	RetryCount  int    `json:"retryCount"`
}
