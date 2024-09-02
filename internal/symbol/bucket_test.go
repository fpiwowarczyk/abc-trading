package symbol_test

import (
	"testing"

	"github.com/fpiwowarczyk/abc-trading/internal/symbol"
	"github.com/stretchr/testify/assert"
)

func Test_PointsFitIntoBucket(t *testing.T) {
	t.Parallel()

	t.Run("points fit into bucket", func(t *testing.T) {
		t.Parallel()
		b := &symbol.Bucket{
			Size:   10,
			Points: []float64{1, 2, 3, 4, 5},
		}

		canFit := b.CantFitIntoBucket(5)
		assert.False(t, canFit)
	})

	t.Run("points do not fit into bucket", func(t *testing.T) {
		t.Parallel()
		b := &symbol.Bucket{
			Size:   10,
			Points: []float64{1, 2, 3, 4, 5},
		}

		canFit := b.CantFitIntoBucket(6)
		assert.True(t, canFit)
	})
}
