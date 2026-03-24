package srv

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecorateError_ClassifiesScopedEOF(t *testing.T) {
	originalErr := errors.New("Get \"https://cluster.example/api/v1/nodes\": EOF")

	decoratedErr := DecorateError(originalErr)

	var scopedEOFErr ScopedEOFError
	require.ErrorAs(t, decoratedErr, &scopedEOFErr)
	require.ErrorIs(t, decoratedErr, originalErr)
	require.Contains(t, decoratedErr.Error(), originalErr.Error())
}

func TestDecorateError_DoesNotClassifyNonEOF(t *testing.T) {
	originalErr := errors.New("Get \"https://cluster.example/api/v1/nodes\": context deadline exceeded")

	decoratedErr := DecorateError(originalErr)

	var scopedEOFErr ScopedEOFError
	require.NotErrorAs(t, decoratedErr, &scopedEOFErr)
	require.Same(t, originalErr, decoratedErr)
}
