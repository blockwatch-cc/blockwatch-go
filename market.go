// Copyright (c) 2020 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package blockwatch

import (
	"time"
)

// Trade is a Go struct type that can hold raw tick data stored in market *:Trade
// tables.
type Trade struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"time"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	IsSell    bool      `json:"sell"`
}

// Ohlvc is a Go struct type that can hold data stored in market *:OHLCV time-series
type Ohlcv struct {
	Timestamp       time.Time `json:"time"`
	Open            float64   `json:"open"`
	Close           float64   `json:"close"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Vwap            float64   `json:"vwap"`
	Std             float64   `json:"stddev"`
	Mean            float64   `json:"mean"`
	TradeCount      int64     `json:"n_trades"`
	BuyCount        int64     `json:"n_buy"`
	SellCount       int64     `json:"n_sell"`
	BaseVolume      float64   `json:"vol_base"`
	QuoteVolume     float64   `json:"vol_quote"`
	BaseVolumeBuy   float64   `json:"vol_buy_base"`
	QuoteVolumeBuy  float64   `json:"vol_buy_quote"`
	BaseVolumeSell  float64   `json:"vol_sell_base"`
	QuoteVolumeSell float64   `json:"vol_sell_quote"`
}
