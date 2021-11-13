package price

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisualPriceToDBPrice(t *testing.T) {
	vp := 0.13
	dbp := VisualPriceToDBPrice(vp)
	assert.Equal(t, uint64(vp*priceScale), dbp)

	vp1 := DBPriceToVisualPrice(dbp)
	assert.Equal(t, vp, vp1)
}
