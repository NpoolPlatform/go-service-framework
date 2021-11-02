package price

func VisualPriceToDBPrice(price float32) int64 {
	return int64(price * 1000000000000)
}

func DBPriceToVisualPrice(price int64) float32 {
	return float32(price / 1000000000000.0)
}
