package price

import (
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	priceScale12 = 1000000000000
	Precision12  = 12
)

func VisualPriceToDBPrice(price float64) uint64 {
	myPrice := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(priceScale12))
	iPrice, err := strconv.ParseUint(myPrice.String(), 10, 64)
	if err != nil {
		return uint64(price * priceScale12)
	}
	return iPrice
}

func DBPriceToVisualPrice(price uint64) float64 {
	myPrice := decimal.NewFromInt(int64(price)).Div(decimal.NewFromInt(priceScale12))
	fPrice, err := strconv.ParseFloat(myPrice.String(), 64)
	if err != nil {
		return float64(price)
	}
	return fPrice
}

func VisualPriceToDBSignPrice(price float64) int64 {
	myPrice := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(priceScale12))
	iPrice, err := strconv.ParseInt(myPrice.String(), 10, 64)
	if err != nil {
		return int64(price * priceScale12)
	}
	return iPrice
}

const (
	priceScale6 = 1000000
	Precision6  = 6
)

func VisualPriceToDBPrice6(price float64) uint64 {
	myPrice := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(priceScale6))
	iPrice, err := strconv.ParseUint(myPrice.String(), 10, 64)
	if err != nil {
		return uint64(price * priceScale6)
	}
	return iPrice
}

func DBPriceToVisualPrice6(price uint64) float64 {
	myPrice := decimal.NewFromInt(int64(price)).Div(decimal.NewFromInt(priceScale6))
	fPrice, err := strconv.ParseFloat(myPrice.String(), 64)
	if err != nil {
		return float64(price)
	}
	return fPrice
}

func VisualPriceToDBSignPrice6(price float64) int64 {
	myPrice := decimal.NewFromFloat(price).Mul(decimal.NewFromInt(priceScale6))
	iPrice, err := strconv.ParseInt(myPrice.String(), 10, 64)
	if err != nil {
		return int64(price * priceScale6)
	}
	return iPrice
}
