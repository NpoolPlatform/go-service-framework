package wlog

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWlog(t *testing.T) {
	aErr := errors.New("ssss")
	bErr := fmt.Errorf("ssssz")
	cErr := fmt.Errorf("ssss")

	assert.Equal(t, false, Equal(aErr, bErr))
	assert.Equal(t, false, Equal(WrapError(aErr), bErr))
	assert.Equal(t, false, Equal(WrapError(aErr), WrapError(bErr)))
	assert.Equal(t, false, Equal(Errorf("%w", aErr), bErr))

	assert.Equal(t, true, Equal(aErr, aErr))
	assert.Equal(t, true, Equal(WrapError(aErr), aErr))
	assert.Equal(t, true, Equal(WrapError(aErr), WrapError(aErr)))
	assert.Equal(t, true, Equal(Errorf("%w", aErr), aErr))

	assert.Equal(t, true, Equal(aErr, cErr))
	assert.Equal(t, true, Equal(WrapError(aErr), cErr))
	assert.Equal(t, true, Equal(WrapError(aErr), WrapError(cErr)))
	assert.Equal(t, true, Equal(Errorf("%w", aErr), cErr))

	aaErr := WrapError(aErr)
	aaaErr := WrapError(aaErr)
	aaaaErr := WrapError(aaaErr)
	aaaaaErr := WrapError(aaaaErr)
	assert.Equal(t, true, Equal(aErr, aaErr))
	assert.Equal(t, true, Equal(aaErr, aaaErr))
	assert.Equal(t, true, Equal(aErr, aaaaErr))
	assert.Equal(t, true, Equal(aaErr, aaaaaErr))
}
