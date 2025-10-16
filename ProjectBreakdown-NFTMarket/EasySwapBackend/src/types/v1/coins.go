package types

import "time"

type Users struct {
	Id        int       `json:"id"`
	Address   string    `json:"address"`
	ChainId   int       `json:"chain_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u Users) TableName() string {
	return "cg_users"
}

type UserCoinsMapping struct {
	UserId    int       `json:"user_id"`
	CoinId    string    `json:"coin_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (u UserCoinsMapping) TableName() string {
	return "cg_user_coins_mapping"
}

type CoinsMetadata struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	ApiSymbol string    `json:"api_symbol"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (c CoinsMetadata) TableName() string {
	return "cg_coins_metadata"
}

type CoinsSearchApiResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	ApiSymbol     string `json:"api_symbol"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	Thumb         string `json:"thumb"`
	Large         string `json:"large"`
}

type CoinSimplePriceResponse struct {
	Usd           float64 `json:"usd"`
	UsdMarketCap  float64 `json:"usd_market_cap,omitempty"`
	Usd24HVol     float64 `json:"usd_24h_vol,omitempty"`
	Usd24HChange  float64 `json:"usd_24h_change,omitempty"`
	LastUpdatedAt int     `json:"last_updated_at,omitempty"`
}

type CoinWithLatestPrice struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type AddCoinFavoriteRequest struct {
	Id string `json:"id"`
}

type UserCoinAlert struct {
	ID            uint64    `json:"id" gorm:"primaryKey"`
	UserID        uint64    `json:"user_id"`
	CoinID        string    `json:"coin_id"`
	AlertType     string    `json:"alert_type"` // "above" 或 "below"
	TargetPrice   float64   `json:"target_price"`
	PriceRangeMin float64   `json:"price_range_min"`
	PriceRangeMax float64   `json:"price_range_max"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserCoinAlert) TableName() string {
	return "cg_user_coin_alerts"
}
