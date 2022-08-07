package price

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisualPriceToDBPrice(t *testing.T) {
	vp := 0.13
	dbp := VisualPriceToDBPrice(vp)
	assert.Equal(t, uint64(vp*priceScale12), dbp)

	vp1 := DBPriceToVisualPrice(dbp)
	assert.Equal(t, vp, vp1)

	vp = -0.13
	dbp1 := VisualPriceToDBSignPrice(vp)
	assert.Equal(t, int64(vp*priceScale12), dbp1)
}
