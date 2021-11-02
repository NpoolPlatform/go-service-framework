package price

func VisualPriceToDBPrice(price float32) uint64 {
	return uint64(price * 1000000000000)
}

func DBPriceToVisualPrice(price uint64) float32 {
	return float32(price / 1000000000000.0)
}
