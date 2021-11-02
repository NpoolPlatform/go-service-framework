package price

import (
	"strconv"

	"github.com/shopspring/decimal"
)

const priceScale = 1000000000000

func VisualPriceToDBPrice(price float64) uint64 {
	return uint64(price * priceScale)
}

func DBPriceToVisualPrice(price uint64) float64 {
	myPrice := decimal.NewFromInt(int64(price)).Div(decimal.NewFromFloat(priceScale))
	fPrice, err := strconv.ParseFloat(myPrice.String(), 64)
	if err != nil {
		return float64(price)
	}
	return fPrice
}
