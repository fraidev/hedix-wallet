package models

// Asset represents the type of crypto asset
type Asset string

const (
	BTC Asset = "BTC"
	ETH Asset = "ETH"
	USD Asset = "USD"
)

// Decimal places for each asset
// BTC: 8 decimals (satoshis)
// ETH: 18 decimals (wei)
// USD: 2 decimals (cents)
const (
	BTCDecimals = 8
	ETHDecimals = 18
	USDDecimals = 2
)

// GetDecimals returns the number of decimal places for an asset
func (a Asset) GetDecimals() int {
	switch a {
	case BTC:
		return BTCDecimals
	case ETH:
		return ETHDecimals
	case USD:
		return USDDecimals
	default:
		return 0
	}
}
