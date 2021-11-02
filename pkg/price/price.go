package price

import (
	"strconv"

	"github.com/shopspring/decimal"
)

const priceScale = 1000000000000

func VisualPriceToDBPrice(price float64) uint64 {
	myPrice := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(priceScale))
	iPrice, err := strconv.ParseUint(myPrice.String(), 10, 64)
	if err != nil {
		return uint64(price * priceScale)
	}
	return iPrice
}

func DBPriceToVisualPrice(price uint64) float64 {
	myPrice := decimal.NewFromInt(int64(price)).Div(decimal.NewFromInt(priceScale))
	fPrice, err := strconv.ParseFloat(myPrice.String(), 64)
	if err != nil {
		return float64(price)
	}
	return fPrice
}
