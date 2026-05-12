package processor

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestGetHistoricalUSDBudget(t *testing.T) {
	tests := []struct {
		platform  string // CoinGecko platform name: "ethereum", "binance-smart-chain"
		tokenAddr string // token contract address
		amount    int64  // token amount in smallest unit
		decimals  uint8
		timestamp uint64 // Unix timestamp
		wantMin   float64 // minimum expected price (sanity check)
	}{
		// ETH @ ~$2507, 2025-05-12 (within CoinGecko free 365-day window)
		{
			platform:  "ethereum",
			tokenAddr: "0xeeeee0eee0eeee0eeee0eeee0eeee0eeee0eeeee", // ETH (matches staticTokenToGeckoID)
			amount:    1_000000000000000000,                         // 1 ETH
			decimals:  18,
			timestamp: 1747064855,                                  // 2025-05-12 00:00:00 UTC
			wantMin:   2000,                                         // ETH was ~$2500 in May 2025
		},
		// USDC @ $1.00, 2025-05-12
		{
			platform:  "ethereum",
			tokenAddr: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
			amount:    1_000000, // 1 USDC
			decimals:  6,
			timestamp: 1747064855,
			wantMin:   0.9,
		},
		// BSC USDT @ $1.00, 2025-05-12
		{
			platform:  "binance-smart-chain",
			tokenAddr: "0x55d398326f99059fF775485246999027B3197955", // USDT on BSC
			amount:    1_000000000000000000,                         // 1 USDT
			decimals:  18,
			timestamp: 1747064855,
			wantMin:   0.9,
		},
		// BSC token 0x40b8129B786D766267A7a118cF8C07E31CDB6Fde, 2026-05-12
		{
			platform:  "binance-smart-chain",
			tokenAddr: "0x40b8129B786D766267A7a118cF8C07E31CDB6Fde",
			amount:    1_000000000000000000,
			decimals:  18,
			timestamp: 1778515200, // 2026-05-12 00:00:00 UTC
			wantMin:   0.05,      // Unibase ~$0.11-$0.15 in May 2026
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%s_%s_%d", tt.platform, tt.tokenAddr, tt.timestamp)
		t.Run(name, func(t *testing.T) {
			if i > 0 {
				time.Sleep(5 * time.Second) // avoid CoinGecko rate limit
			}
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			amount := big.NewInt(tt.amount)
			usd, err := GetHistoricalUSDBudget(ctx, tt.platform, tt.tokenAddr, amount, tt.decimals, tt.timestamp)
			if err != nil {
				t.Fatalf("GetHistoricalUSDBudget error: %v", err)
			}

			t.Logf("platform=%s token=%s timestamp=%d → USD=%.4f", tt.platform, tt.tokenAddr, tt.timestamp, usd)

			if usd < tt.wantMin {
				t.Errorf("USD price too low: got %.4f, want >= %.4f", usd, tt.wantMin)
			}
		})
	}
}