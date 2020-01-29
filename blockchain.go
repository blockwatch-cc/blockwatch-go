// Copyright (c) 2020 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package blockwatch

import (
	"time"
)

// Block is a Go struct type that can hold raw data about a blockchain block as
// stored in blockchain *:BLOCK tables.
type Block struct {
	RowID                uint64    `json:"row_id"`
	ParentID             uint64    `json:"parent_id"`
	Orphan               bool      `json:"is_orphan"`
	Hash                 string    `json:"hash"`
	Timestamp            time.Time `json:"time"`
	MedianTime           time.Time `json:"mediantime"`
	Height               uint64    `json:"height"`
	Version              int64     `json:"version"`
	Size                 uint64    `json:"size"`
	Weight               uint64    `json:"weight"`
	Bits                 uint64    `json:"bits"`
	ChainWork            float64   `json:"chainwork"`
	Difficulty           float64   `json:"difficulty"`
	Coinbase             []byte    `json:"coinbase"`
	AddressesSeen        uint64    `json:"n_addr"`
	AddressesCreated     uint64    `json:"n_new_addr"`
	AddressesEmptied     uint64    `json:"n_empty_addr"`
	AddressesFunded      uint64    `json:"n_funded_addr"`
	TransactionCount     uint64    `json:"n_tx"`
	UtxoConsumed         uint64    `json:"n_vin"`
	UtxoCreated          uint64    `json:"n_vout"`
	SpendableUtxoCreated uint64    `json:"n_vout_spendable"`
	TransactionVolume    float64   `json:"volume"`
	MiningReward         float64   `json:"reward"`
	TransactionFees      float64   `json:"fee"`
	BurnedCoins          float64   `json:"burned"`
	DaysDestroyed        float64   `json:"days_destroyed"`
	Solvetime            uint64    `json:"solvetime"`
}

// Tx is a Go struct type that can hold raw data about a blockchain transaction as
// stored in blockchain *:TX tables.
type Tx struct {
	RowID          uint64    `json:"row_id"`
	Timestamp      time.Time `json:"time"`
	Height         uint64    `json:"height"`
	Position       uint64    `json:"tx_n"`
	TransactionID  string    `json:"tx_id"`
	Locktime       int64     `json:"locktime"`
	Size           int64     `json:"size"`
	VirtualSize    int64     `json:"vsize"`
	Version        int64     `json:"version"`
	SpentInputs    int64     `json:"n_in"`
	CreatedOutputs int64     `json:"n_out"`
	Type           string    `json:"type"`
	HasData        bool      `json:"has_data"`
	Volume         float64   `json:"volume"`
	Fee            float64   `json:"fee"`
	DaysDestroyed  float64   `json:"days_destroyed"`
}

// Chain is a Go struct type that can hold running blockchain totals as
// stored in blockchain *:CHAIN tables.
type Chain struct {
	Height            uint64    `json:"height"`
	Timestamp         time.Time `json:"time"`
	Difficulty        float64   `json:"difficulty"`
	AvgHashrate3h     float64   `json:"hashrate_3h"`
	AvgHashrate12h    float64   `json:"hashrate_12h"`
	TotalWork         float64   `json:"total_work"`
	TotalSize         uint64    `json:"total_size"`
	TotalTransactions uint64    `json:"total_tx"`
	TotalUtxos        uint64    `json:"total_utxo"`
	TotalAddresses    uint64    `json:"total_addr"`
	FundedAddresses   uint64    `json:"funded_addr"`
	TotalSupply       float64   `json:"total_supply"`
	MintedSupply      float64   `json:"minted_supply"`
	MinedSupply       float64   `json:"mined_supply"`
	CurrentSupply     float64   `json:"current_supply"`
	LockedSupply      float64   `json:"locked_supply"`
	BurnedSupply      float64   `json:"burned_supply"`
}

// Flow is a Go struct type that can hold blockchain flow data as
// stored in blockchain *:FLOW tables.
type Flow struct {
	RowID             uint64    `json:"row_id"`
	FundingTime       time.Time `json:"fund_time"`
	FundingHeight     uint64    `json:"fund_height"`
	FundingPosition   uint64    `json:"fund_txpos"`
	FundingOutput     uint64    `json:"fund_vout"`
	FundingTxID       string    `json:"fund_txid"`
	Volume            float64   `json:"volume"`
	CoinGenerationMin uint64    `json:"coin_gen_min"`
	CoinGenerationMax uint64    `json:"coin_gen_max"`
	AddressCount      uint64    `json:"n_addr"`
	SignatureCount    uint64    `json:"n_req_sig"`
	AddressType       string    `json:"addr_type"`
	Address           string    `json:"addr"`
	Data              []byte    `json:"data"`
	IsBurned          bool      `json:"is_burned"`
	IsSpendable       bool      `json:"is_spendable"`
	IsSpent           bool      `json:"is_spent"`
	SpendingTime      time.Time `json:"spend_time"`
	SpendingHeight    uint64    `json:"spend_height"`
	SpendingPosition  uint64    `json:"spend_txpos"`
	SpendingInput     uint64    `json:"spend_vin"`
	SpendingTxID      string    `json:"spend_txid"`
}

// AddressAgeStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:AGE time series.
type AddressAgeStats struct {
	Timestamp       time.Time `json:"time"`
	Y1DormantAddr   int64     `json:"y1_addr"`
	Y1DormantFunds  float64   `json:"y1_funds"`
	Y2DormantAddr   int64     `json:"y2_addr"`
	Y2DormantFunds  float64   `json:"y2_funds"`
	Y3DormantAddr   int64     `json:"y3_addr"`
	Y3DormantFunds  float64   `json:"y3_funds"`
	Y4DormantAddr   int64     `json:"y4_addr"`
	Y4DormantFunds  float64   `json:"y4_funds"`
	Y5DormantAddr   int64     `json:"y5_addr"`
	Y5DormantFunds  float64   `json:"y5_funds"`
	Y6DormantAddr   int64     `json:"y6_addr"`
	Y6DormantFunds  float64   `json:"y6_funds"`
	Y7DormantAddr   int64     `json:"y7_addr"`
	Y7DormantFunds  float64   `json:"y7_funds"`
	Y8DormantAddr   int64     `json:"y8_addr"`
	Y8DormantFunds  float64   `json:"y8_funds"`
	Y9DormantAddr   int64     `json:"y9_addr"`
	Y9DormantFunds  float64   `json:"y9_funds"`
	Y10DormantAddr  int64     `json:"y10_addr"`
	Y10DormantFunds float64   `json:"y10_funds"`
	Y11DormantAddr  int64     `json:"y11_addr"`
	Y11DormantFunds float64   `json:"y11_funds"`
	Y12DormantAddr  int64     `json:"y12_addr"`
	Y12DormantFunds float64   `json:"y12_funds"`
	Y13DormantAddr  int64     `json:"y13_addr"`
	Y13DormantFunds float64   `json:"y13_funds"`
	Y14DormantAddr  int64     `json:"y14_addr"`
	Y14DormantFunds float64   `json:"y14_funds"`
	Y15DormantAddr  int64     `json:"y15_addr"`
	Y15DormantFunds float64   `json:"y15_funds"`
	Y16DormantAddr  int64     `json:"y16_addr"`
	Y16DormantFunds float64   `json:"y16_funds"`
	Y17DormantAddr  int64     `json:"y17_addr"`
	Y17DormantFunds float64   `json:"y17_funds"`
	Y18DormantAddr  int64     `json:"y18_addr"`
	Y18DormantFunds float64   `json:"y18_funds"`
	Y19DormantAddr  int64     `json:"y19_addr"`
	Y19DormantFunds float64   `json:"y19_funds"`
	Y20DormantAddr  int64     `json:"y20_addr"`
	Y20DormantFunds float64   `json:"y20_funds"`
}

// AddressActivityStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:ACTIVITY time series.
type AddressActivityStats struct {
	Timestamp             time.Time `json:"time"`
	BlockCount            int64     `json:"n_blocks"`
	AverageBlockTime      float64   `json:"avg_solvetime"`
	BlockchainGrowth      int64     `json:"size_growth"`
	SumRewards            float64   `json:"sum_rewards"`
	SumBurned             float64   `json:"sum_burned"`
	SumCoinDaysDestroyed  float64   `json:"sum_cdd"`
	TransactionVolume     float64   `json:"sum_vol"`
	TransactionFees       float64   `json:"sum_fee"`
	TransactionCount      int64     `json:"n_tx"`
	UtxoCreated           int64     `json:"n_vout"`
	SpendableUtxoCreated  int64     `json:"n_svout"`
	UnspendableTxoCreated int64     `json:"n_uvout"`
	UtxoConsumed          int64     `json:"n_vin"`
	AddressesSeen         int64     `json:"n_addr_active"`
	AddressesCreated      int64     `json:"n_addr_new"`
	AddressesFunded       int64     `json:"n_addr_funded"`
	AddressesEmpty        int64     `json:"n_addr_empty"`
	AddressReusePercent   float64   `json:"pct_addr_reuse"`
	Top1ByVolume          float64   `json:"vol_top1"`
	Top10ByVolume         float64   `json:"vol_top10"`
	Top100ByVolume        float64   `json:"vol_top100"`
	Top1kByVolume         float64   `json:"vol_top1k"`
	Top10kByVolume        float64   `json:"vol_top10k"`
	Top100kByVolume       float64   `json:"vol_top100k"`
	Top1ByTx              int64     `json:"tx_top1"`
	Top10ByTx             int64     `json:"tx_top10"`
	Top100ByTx            int64     `json:"tx_top100"`
	Top1kByTx             int64     `json:"tx_top1k"`
	Top10kByTx            int64     `json:"tx_top10k"`
	Top100kByTx           int64     `json:"tx_top100k"`
}

// AddressBalanceStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:BALANCE time series.
type AddressBalanceStats struct {
	Timestamp                         time.Time `json:"time"`
	Top1Richest                       float64   `json:"rich_top1"`
	Top10Richest                      float64   `json:"rich_top10"`
	Top100Richest                     float64   `json:"rich_top100"`
	Top1kRichest                      float64   `json:"rich_top1k"`
	Top10kRichest                     float64   `json:"rich_top10k"`
	Top100kRichest                    float64   `json:"rich_top100k"`
	SumFundsAtAddressesAbove1atom     float64   `json:"funds_e0"`
	SumFundsAtAddressesAbove10atoms   float64   `json:"funds_e1"`
	SumFundsAtAddressesAbove100atoms  float64   `json:"funds_e2"`
	SumFundsAtAddressesAbove1katoms   float64   `json:"funds_e3"`
	SumFundsAtAddressesAbove10katoms  float64   `json:"funds_e4"`
	SumFundsAtAddressesAbove100katoms float64   `json:"funds_e5"`
	SumFundsAtAddressesAbove1Matoms   float64   `json:"funds_e6"`
	SumFundsAtAddressesAbove10Matoms  float64   `json:"funds_e7"`
	SumFundsAtAddressesAbove1coin     float64   `json:"funds_e8"`
	SumFundsAtAddressesAbove10coins   float64   `json:"funds_e9"`
	SumFundsAtAddressesAbove100coins  float64   `json:"funds_e10"`
	SumFundsAtAddressesAbove1kcoins   float64   `json:"funds_e11"`
	SumFundsAtAddressesAbove10kcoins  float64   `json:"funds_e12"`
	SumFundsAtAddressesAbove100kcoins float64   `json:"funds_e13"`
	SumFundsAtAddressesAbove1Mcoins   float64   `json:"funds_e14"`
	SumFundsAtAddressesAbove10Mcoins  float64   `json:"funds_e15"`
	SumFundsAtAddressesAbove100Mcoins float64   `json:"funds_e16"`
	AddresssCountAbove1atom           int64     `json:"addrs_e0"`
	AddresssCountAbove10atoms         int64     `json:"addrs_e1"`
	AddresssCountAbove100atoms        int64     `json:"addrs_e2"`
	AddresssCountAbove1katoms         int64     `json:"addrs_e3"`
	AddresssCountAbove10katoms        int64     `json:"addrs_e4"`
	AddresssCountAbove100katoms       int64     `json:"addrs_e5"`
	AddresssCountAbove1Matoms         int64     `json:"addrs_e6"`
	AddresssCountAbove10Matoms        int64     `json:"addrs_e7"`
	AddresssCountAbove1coin           int64     `json:"addrs_e8"`
	AddresssCountAbove10coins         int64     `json:"addrs_e9"`
	AddresssCountAbove100coins        int64     `json:"addrs_e10"`
	AddresssCountAbove1kcoins         int64     `json:"addrs_e11"`
	AddresssCountAbove10kcoins        int64     `json:"addrs_e12"`
	AddresssCountAbove100kcoins       int64     `json:"addrs_e13"`
	AddresssCountAbove1Mcoins         int64     `json:"addrs_e14"`
	AddresssCountAbove10Mcoins        int64     `json:"addrs_e15"`
	AddresssCountAbove100Mcoins       int64     `json:"addrs_e16"`
}

// SupplyStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:SUPPY time series.
type SupplyStats struct {
	Timestamp               time.Time `json:"time"`
	TotalSupply             float64   `json:"total"`
	CurrentSupply           float64   `json:"current"`
	CirculatingSupply       float64   `json:"circulating"`
	MinedSupply             float64   `json:"mined"`
	LockedSupply            float64   `json:"locked"`
	BurnedSupply            float64   `json:"burned"`
	UntouchedSupply         float64   `json:"untouched"`
	HodlSupply3M            float64   `json:"hodl_3m"`
	TransactingSupply3M     float64   `json:"tx_3m"`
	DaysDestroyed3M         float64   `json:"cdd_3m"`
	InflationLast24h        float64   `json:"inflation"`
	AnnualizedInflationRate float64   `json:"inflation_rate"`
}

// TxStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:TX time series.
type TxStats struct {
	Timestamp              time.Time `json:"time"`
	Type                   string    `json:"type"`
	Count                  int64     `json:"n_tx"`
	MinFee                 float64   `json:"min_fee"`
	MaxFee                 float64   `json:"max_fee"`
	MeanFee                float64   `json:"mean_fee"`
	MedianFee              float64   `json:"median_fee"`
	SumFees                float64   `json:"sum_fee"`
	MinFeeRate             float64   `json:"min_fee_rate"`
	MaxFeeRate             float64   `json:"max_fee_rate"`
	MeanFeeRate            float64   `json:"mean_fee_rate"`
	MedianFeeRate          float64   `json:"median_fee_rate"`
	SumFeeRates            float64   `json:"sum_fee_rate"`
	MinSize                int64     `json:"min_size"`
	MaxSize                int64     `json:"max_size"`
	MeanSize               float64   `json:"mean_size"`
	MedianSize             float64   `json:"median_size"`
	SumSizes               int64     `json:"sum_size"`
	MinInputs              int64     `json:"min_n_vin"`
	MaxInputs              int64     `json:"max_n_vin"`
	MeanInputs             float64   `json:"mean_n_vin"`
	MedianInputs           float64   `json:"median_n_vin"`
	SumInputs              int64     `json:"sum_n_vin"`
	MinOutputs             int64     `json:"min_n_vout"`
	MaxOutputs             int64     `json:"max_n_vout"`
	MeanOutputs            float64   `json:"mean_n_vout"`
	MedianOutputs          float64   `json:"median_n_vout"`
	SumOutputs             int64     `json:"sum_n_vout"`
	MinVolume              float64   `json:"min_vol"`
	MaxVolume              float64   `json:"max_vol"`
	MeanVolume             float64   `json:"mean_vol"`
	MedianVolume           float64   `json:"median_vol"`
	SumVolume              float64   `json:"sum_vol"`
	MinDaysDestroyed       float64   `json:"min_cdd"`
	MaxDaysDestroyed       float64   `json:"max_cdd"`
	MeanDaysDestroyed      float64   `json:"mean_cdd"`
	MedianDaysDestroyed    float64   `json:"median_cdd"`
	SumDaysDestroyed       float64   `json:"sum_cdd"`
	MinAvgDaysDestroyed    float64   `json:"min_add"`
	MaxAvgDaysDestroyed    float64   `json:"max_add"`
	MeanAvgDaysDestroyed   float64   `json:"mean_add"`
	MedianAvgDaysDestroyed float64   `json:"median_add"`
	SumAvgDaysDestroyed    float64   `json:"sum_add"`
}

// UtxoStats is a Go struct type that can hold blockchain statistics data as
// stored in blockchain *-EOD:UTXO time series.
type UtxoStats struct {
	Timestamp time.Time `json:"time"`
	Type      string    `json:"type"`
	Count     int64     `json:"n_out"`
	Volume    float64   `json:"vol"`
}
