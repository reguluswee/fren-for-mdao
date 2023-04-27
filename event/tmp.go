package event

import (
	"time"

	"github.com/shopspring/decimal"
)

type MdaoData struct {
	Id      uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	AddTime time.Time

	Block   uint64
	Txhash  string
	Wallet  string
	Minter  string
	Round   uint64
	Term    uint64
	Rewards decimal.Decimal
	Loss    decimal.Decimal
	Ts      time.Time
}

type MdaoBlockError struct {
	Id      uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	AddTime time.Time
	Block   uint64
	Txhash  string
}

func (MdaoData) TableName() string {
	return "mdao_data"
}

func (MdaoBlockError) TableName() string {
	return "mdao_block_tx_error"
}
