package price

func VisualPriceToDBPrice(price float64) uint64 {
	return uint64(price * 1000000000000)
}

func DBPriceToVisualPrice(price uint64) float64 {
	return float64(price / 1000000000000.0)
}
