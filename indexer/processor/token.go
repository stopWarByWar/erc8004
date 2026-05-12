package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// geckoHTTPClient is the HTTP client used for all CoinGecko API calls.
// Uses no proxy to avoid inconsistent responses between CLI and Go.
var geckoHTTPClient = &http.Client{Timeout: 15 * time.Second}

// ─── Caches ───────────────────────────────────────────────────────────────────

// geckoIDCache caches CoinGecko coin ID per (platform, tokenAddr).
// Coin IDs rarely change, so this is a long-lived cache (no TTL).
var geckoIDCache   = map[string]string{}
var geckoIDCacheMu sync.RWMutex

func geckoIDKey(platform, tokenAddr string) string {
	return platform + ":" + strings.ToLower(tokenAddr)
}

func getCachedGeckoID(platform, tokenAddr string) (string, bool) {
	geckoIDCacheMu.RLock()
	id, ok := geckoIDCache[geckoIDKey(platform, tokenAddr)]
	geckoIDCacheMu.RUnlock()
	return id, ok
}

func setCachedGeckoID(platform, tokenAddr, id string) {
	geckoIDCacheMu.Lock()
	geckoIDCache[geckoIDKey(platform, tokenAddr)] = id
	geckoIDCacheMu.Unlock()
}

// Price cache: key = "platform:addr:date" → price.
// Different dates for the same token have different prices.
type tokenPriceEntry struct {
	price     float64
	fetchedAt time.Time
}

var globalPriceCache struct {
	mu sync.RWMutex
	m  map[string]tokenPriceEntry
}

const priceCacheTTL = 1 * time.Hour

func init() {
	globalPriceCache.m = make(map[string]tokenPriceEntry)
}

func priceCacheKey(platform, tokenAddr, date string) string {
	return platform + ":" + strings.ToLower(tokenAddr) + ":" + date
}

// ─── Token Info ────────────────────────────────────────────────────────────────

// TokenInfo holds symbol and decimals for an ERC-20 token.
type TokenInfo struct {
	Symbol   string
	Decimals uint8
}

// tokenInfoCache caches TokenInfo per token address.
type tokenInfoCache struct {
	mu     sync.RWMutex
	tokens map[string]TokenInfo
}

var globalTokenInfoCache tokenInfoCache

func init() {
	globalTokenInfoCache.tokens = make(map[string]TokenInfo)
}

// ResolveTokenInfo returns cached TokenInfo, fetching from chain if not cached.
func ResolveTokenInfo(ctx context.Context, client *ethclient.Client, addr string) (*TokenInfo, error) {
	globalTokenInfoCache.mu.RLock()
	if info, ok := globalTokenInfoCache.tokens[addr]; ok {
		globalTokenInfoCache.mu.RUnlock()
		return &info, nil
	}
	globalTokenInfoCache.mu.RUnlock()

	info, err := fetchTokenInfo(ctx, client, common.HexToAddress(addr))
	if err != nil {
		return &TokenInfo{Symbol: "???", Decimals: 18}, err
	}

	globalTokenInfoCache.mu.Lock()
	globalTokenInfoCache.tokens[addr] = *info
	globalTokenInfoCache.mu.Unlock()
	return info, nil
}

func fetchTokenInfo(ctx context.Context, client *ethclient.Client, addr common.Address) (*TokenInfo, error) {
	symbol, err := callSymbol(ctx, client, addr)
	if err != nil {
		return &TokenInfo{Symbol: "???", Decimals: 18}, err
	}
	decimals, err := callDecimals(ctx, client, addr)
	if err != nil {
		return &TokenInfo{Symbol: symbol, Decimals: 18}, err
	}
	return &TokenInfo{Symbol: symbol, Decimals: decimals}, nil
}

func callSymbol(ctx context.Context, client *ethclient.Client, addr common.Address) (string, error) {
	data := common.FromHex("0x95d89b41")
	out, err := client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: data}, nil)
	if err != nil {
		return "", err
	}
	return parseString32(out), nil
}

func callDecimals(ctx context.Context, client *ethclient.Client, addr common.Address) (uint8, error) {
	data := common.FromHex("0x313ce567")
	out, err := client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: data}, nil)
	if err != nil {
		return 18, err
	}
	if len(out) < 32 {
		return 18, nil
	}
	return uint8(out[31]), nil
}

func parseString32(b []byte) string {
	if len(b) < 32 {
		return ""
	}
	offset := int(new(big.Int).SetBytes(b[0:32]).Uint64())
	if offset > len(b) {
		return ""
	}
	length := int(new(big.Int).SetBytes(b[offset : offset+32]).Uint64())
	if offset+32+length > len(b) {
		return ""
	}
	return string(b[offset+32 : offset+32+length])
}

// ─── Price APIs ────────────────────────────────────────────────────────────────

// GetHistoricalUSDBudget returns amount * price at blockTimestamp (USD).
// Price is cached for 1 hour per (platform, tokenAddr, date).
func GetHistoricalUSDBudget(ctx context.Context, platform, tokenAddr string, amount *big.Int, decimals uint8, blockTimestamp uint64) (float64, error) {
	dateStr := time.Unix(int64(blockTimestamp), 0).UTC().Format("02-01-2006")
	cacheKey := priceCacheKey(platform, tokenAddr, dateStr)

	globalPriceCache.mu.RLock()
	cached, ok := globalPriceCache.m[cacheKey]
	globalPriceCache.mu.RUnlock()

	var price float64
	if ok && time.Since(cached.fetchedAt) < priceCacheTTL {
		price = cached.price
	} else {
		var err error
		price, err = fetchHistoricalPrice(ctx, platform, tokenAddr, dateStr)
		if err != nil {
			logrus.WithField("token", tokenAddr).Warnf("fetchHistoricalPrice failed, using default 1: %v", err)
			price = 1.0
		}
		globalPriceCache.mu.Lock()
		globalPriceCache.m[cacheKey] = tokenPriceEntry{price: price, fetchedAt: time.Now()}
		globalPriceCache.mu.Unlock()
	}

	if amount == nil || amount.Sign() == 0 {
		return 0.0, nil
	}
	f, _ := new(big.Float).SetInt(amount).Float64()
	return f * price / math.Pow10(int(decimals)), nil
}

// GetCurrentUSDBudget returns current USD value of amount.
func GetCurrentUSDBudget(ctx context.Context, platform, tokenAddr string, amount *big.Int, decimals uint8) (float64, error) {
	price, err := fetchPriceWithFallback(ctx, platform, tokenAddr)
	if err != nil {
		logrus.WithField("token", tokenAddr).Warnf("fetchCurrentPrice failed, using default 1: %v", err)
		price = 1.0
	}
	if amount == nil || amount.Sign() == 0 {
		return 0.0, nil
	}
	f, _ := new(big.Float).SetInt(amount).Float64()
	return f * price / math.Pow10(int(decimals)), nil
}

// ─── Job Budget Cache ─────────────────────────────────────────────────────────

var jobBudgetCache sync.Map

type jobBudgetEntry struct {
	PaymentToken    string
	PaymentDecimals uint8
	TokenSymbol      string
	Budget           *big.Int
	BudgetUSD        float64
}

func SetJobBudgetCache(jobID uint64, entry *jobBudgetEntry) {
	jobBudgetCache.Store(jobID, entry)
}

func GetJobBudgetFromCache(jobID uint64) (*jobBudgetEntry, bool) {
	val, ok := jobBudgetCache.Load(jobID)
	if !ok {
		return nil, false
	}
	return val.(*jobBudgetEntry), true
}

// ─── CoinGecko API ────────────────────────────────────────────────────────────

var coingeckoAPIKey string

func SetCoingeckoAPIKey(key string) { coingeckoAPIKey = key }

var chainIDToPlatform = map[string]string{
	"1":     "ethereum",
	"8453":  "ethereum",
	"84532": "ethereum",
	"56":    "binance-smart-chain",
	"42161": "arbitrum",
}

var staticTokenToGeckoID = map[string]string{
	"0xc02aa39b223fe8d0a0e5c4f27ead9083c756cc2": "weth",
	"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": "usd-coin",
	"0xdac17f958d2ee523a2206206994597c13d831ec7": "tether",
	"0xeeeee0eee0eeee0eeee0eeee0eeee0eeee0eeeee": "ethereum",
}

// resolveGeckoID returns CoinGecko ID: static map → cache → dynamic lookup.
func resolveGeckoID(ctx context.Context, platform, tokenAddr string) (string, error) {
	addr := strings.ToLower(tokenAddr)
	if id, ok := staticTokenToGeckoID[addr]; ok {
		return id, nil
	}
	if id, ok := getCachedGeckoID(platform, tokenAddr); ok {
		return id, nil
	}
	id, err := fetchGeckoIDByContract(ctx, platform, tokenAddr)
	if err != nil {
		return "", err
	}
	setCachedGeckoID(platform, tokenAddr, id)
	return id, nil
}

// fetchGeckoIDByContract tries free contract lookup → search fallback → Pro API.
func fetchGeckoIDByContract(ctx context.Context, platform, tokenAddr string) (string, error) {
	id, err := fetchGeckoIDByContractFreeAPI(ctx, platform, tokenAddr)
	if err == nil && id != "" {
		return id, nil
	}
	logrus.WithField("token", tokenAddr).Debugf("contract lookup failed, trying search: %v", err)

	id, err = searchGeckoIDByContract(ctx, platform, tokenAddr)
	if err == nil && id != "" {
		logrus.WithField("token", tokenAddr).Debugf("search fallback found: %s", id)
		return id, nil
	}

	if coingeckoAPIKey != "" {
		id, err := fetchGeckoIDByContractProAPI(ctx, platform, tokenAddr)
		if err == nil && id != "" {
			return id, nil
		}
		logrus.WithField("token", tokenAddr).Warnf("Pro contract lookup failed: %v", err)
	}
	return "", fmt.Errorf("coin ID not found for %s on %s", tokenAddr, platform)
}

// searchGeckoIDByContract uses CoinGecko search API to find coin ID by contract address.
func searchGeckoIDByContract(ctx context.Context, platform, tokenAddr string) (string, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/search?query=%s", strings.ToLower(tokenAddr))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search API returned status %d", resp.StatusCode)
	}

	var result struct {
		Coins []struct {
			ID        string            `json:"id"`
			Symbol    string            `json:"symbol"`
			Name      string            `json:"name"`
			Platforms map[string]string `json:"platforms"`
		} `json:"coins"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, coin := range result.Coins {
		if contract, ok := coin.Platforms[platform]; ok {
			if strings.ToLower(contract) == strings.ToLower(tokenAddr) {
				return coin.ID, nil
			}
		}
	}
	return "", fmt.Errorf("search: no coin matched contract %s on %s", tokenAddr, platform)
}

func fetchGeckoIDByContractProAPI(ctx context.Context, platform, tokenAddr string) (string, error) {
	url := fmt.Sprintf(
		"https://pro-api.coingecko.com/api/v3/coins/%s/contract/%s",
		platform, strings.ToLower(tokenAddr),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-cg-pro-api-key", coingeckoAPIKey)

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CoinGecko Pro API returned status %d", resp.StatusCode)
	}

	var result struct{ ID string `json:"id"` }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", fmt.Errorf("coin ID not found for %s", tokenAddr)
	}
	return result.ID, nil
}

func fetchGeckoIDByContractFreeAPI(ctx context.Context, platform, tokenAddr string) (string, error) {
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/%s/contract/%s",
		platform, strings.ToLower(tokenAddr),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 200))
		return "", fmt.Errorf("CoinGecko Free API returned status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct{ ID string `json:"id"` }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", fmt.Errorf("coin ID empty for %s", tokenAddr)
	}
	return result.ID, nil
}

// fetchPriceWithFallback tries Pro API → Free API → default 1.
func fetchPriceWithFallback(ctx context.Context, platform, tokenAddr string) (float64, error) {
	if coingeckoAPIKey != "" {
		price, err := fetchPriceFromProAPI(ctx, platform, tokenAddr)
		if err == nil {
			return price, nil
		}
		logrus.WithField("token", tokenAddr).Warnf("Pro API failed, fallback to free: %v", err)
	}

	price, err := fetchPriceFromFreeAPI(ctx, platform, tokenAddr)
	if err == nil {
		return price, nil
	}
	logrus.WithField("token", tokenAddr).Warnf("Free API failed, fallback to default 1: %v", err)
	return 1.0, nil
}

func fetchPriceFromProAPI(ctx context.Context, platform, tokenAddr string) (float64, error) {
	url := fmt.Sprintf(
		"https://pro-api.coingecko.com/api/v3/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd",
		platform, strings.ToLower(tokenAddr),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("x-cg-pro-api-key", coingeckoAPIKey)

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if price, ok := result[strings.ToLower(tokenAddr)]["usd"]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("price not found")
}

func fetchPriceFromFreeAPI(ctx context.Context, platform, tokenAddr string) (float64, error) {
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd",
		platform, strings.ToLower(tokenAddr),
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if price, ok := result[strings.ToLower(tokenAddr)]["usd"]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("price not found")
}

// fetchHistoricalPrice queries CoinGecko history endpoint for historical price.
func fetchHistoricalPrice(ctx context.Context, platform, tokenAddr, date string) (float64, error) {
	geckoID, err := resolveGeckoID(ctx, platform, tokenAddr)
	if err != nil {
		return 1.0, fmt.Errorf("resolveGeckoID failed: %w", err)
	}
	logrus.WithField("token", tokenAddr).Debugf("geckoID=%s, date=%s", geckoID, date)

	if coingeckoAPIKey != "" {
		price, err := fetchHistoricalFromProAPI(ctx, geckoID, date)
		if err == nil {
			return price, nil
		}
		logrus.WithField("token", tokenAddr).Warnf("Pro historical API failed, fallback to free: %v", err)
	}

	price, err := fetchHistoricalFromFreeAPI(ctx, geckoID, date)
	if err != nil {
		return 1.0, fmt.Errorf("fetchHistoricalFromFreeAPI failed: %w", err)
	}
	return price, nil
}

func fetchHistoricalFromProAPI(ctx context.Context, geckoID, date string) (float64, error) {
	url := fmt.Sprintf(
		"https://pro-api.coingecko.com/api/v3/coins/%s/history?date=%s&localization=false",
		geckoID, date,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("x-cg-pro-api-key", coingeckoAPIKey)

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Demo key rejection → let caller fall back to free API
	if resp.StatusCode == http.StatusUnauthorized {
		return 0, fmt.Errorf("pro-api rejected demo key (10011)")
	}

	var result struct {
		MarketData struct {
			CurrentPrice map[string]float64 `json:"current_price"`
		} `json:"market_data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if price, ok := result.MarketData.CurrentPrice["usd"]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("historical price not found")
}

func fetchHistoricalFromFreeAPI(ctx context.Context, geckoID, date string) (float64, error) {
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/%s/history?date=%s&localization=false",
		geckoID, date,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := geckoHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 200))
		return 0, fmt.Errorf("CoinGecko Free API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		MarketData struct {
			CurrentPrice map[string]float64 `json:"current_price"`
		} `json:"market_data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("JSON decode failed: %w", err)
	}
	if price, ok := result.MarketData.CurrentPrice["usd"]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("historical price not found")
}