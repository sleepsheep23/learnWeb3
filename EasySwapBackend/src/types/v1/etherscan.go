package types

import "time"

type EtherScanBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"` // 余额通常是字符串类型（单位 wei）
}

type EthereumBlockHeader struct {
	Number     uint64    `gorm:"primaryKey;column:number" json:"number"`
	Hash       string    `gorm:"column:hash" json:"hash"`
	ParentHash string    `gorm:"column:parent_hash" json:"parent_hash"`
	Timestamp  uint64    `gorm:"column:timestamp" json:"timestamp"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (EthereumBlockHeader) TableName() string {
	return "cg_eth_block_headers"
}

type WatchAddress struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Address   string    `gorm:"column:address" json:"address"`
	ChainID   int       `gorm:"column:chain_id" json:"chain_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (WatchAddress) TableName() string {
	return "cg_watch_addresses"
}
