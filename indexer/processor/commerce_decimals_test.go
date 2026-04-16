package processor

import (
	"math/big"
	"testing"
)

func TestScaleTokenAmountToFloat64(t *testing.T) {
	t.Run("decimals6", func(t *testing.T) {
		amt := big.NewInt(123456789) // 123.456789
		got := scaleTokenAmountToFloat64(amt, 6)
		if got != 123.456789 {
			t.Fatalf("expected 123.456789, got %v", got)
		}
	})

	t.Run("decimals18_one", func(t *testing.T) {
		amt := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1e18 => 1
		got := scaleTokenAmountToFloat64(amt, 18)
		if got != 1 {
			t.Fatalf("expected 1, got %v", got)
		}
	})

	t.Run("roundsToSchemaScale", func(t *testing.T) {
		// 0.000000000000000001 (1 wei) should round to 0 with numeric(36,8)
		amt := big.NewInt(1)
		got := scaleTokenAmountToFloat64(amt, 18)
		if got != 0 {
			t.Fatalf("expected 0, got %v", got)
		}
	})
}

