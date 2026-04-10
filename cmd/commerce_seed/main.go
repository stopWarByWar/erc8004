// commerce_seed 在链上按顺序调用 AgenticCommerce，生成 Indexer 可消费的 Job 生命周期事件。
// 典型用法（BSC Testnet chain_id=97）见 docs/resources/bsc-testnet-commerce-seed.md。
package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	cabi "agent_identity/abi"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/yaml.v2"
)

const erc20ApproveABI = `[{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
const erc20BalanceAllowanceABI = `[
  {"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
  {"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}
]`

type seedConfig struct {
	RPCURL   string `yaml:"rpc_url"`
	ChainID  int64  `yaml:"chain_id"`
	Commerce string `yaml:"commerce"`

	Client struct {
		PrivateKey string `yaml:"private_key"`
		Address    string `yaml:"address"`
	} `yaml:"client"`
	Provider struct {
		PrivateKey string `yaml:"private_key"`
		Address    string `yaml:"address"`
	} `yaml:"provider"`
	Evaluator struct {
		PrivateKey string `yaml:"private_key"`
		Address    string `yaml:"address"`
	} `yaml:"evaluator"`

	CreateJob struct {
		Provider    string `yaml:"provider"`
		Evaluator   string `yaml:"evaluator"`
		ExpiredAt   int64  `yaml:"expired_at"`
		Description string `yaml:"description"`
		Hook        string `yaml:"hook"`
	} `yaml:"create_job"`

	Budget      string `yaml:"budget"`
	OptParams   string `yaml:"opt_params"`
	Deliverable string `yaml:"deliverable"`
	Reason      string `yaml:"reason"`
	GasLimit    uint64 `yaml:"gas_limit"`

	Command string `yaml:"command"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "commerce_seed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("commerce_seed", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	cfgPath := fs.String("f", "./config/commerce_seed.local.yaml", "seed config YAML path")
	command := fs.String("command", "", "override command (otherwise uses config.command)")

	// init-config: 生成三钱包并写入配置文件（同时不会从环境变量读取任何内容）
	initConfig := fs.Bool("init-config", false, "generate 3 wallets and write config file, then exit")
	outPath := fs.String("out", "", "output config path for -init-config (default: -f)")
	force := fs.Bool("force", false, "overwrite existing config file when -init-config is set")

	// 对于需要 job-id 的分步命令，允许 CLI 传入（不强制写入配置）
	jobIDStr := fs.String("job-id", "", "decimal job id for set-budget / approve-fund / fund-only / submit-only / complete-only")
	paymentToken := fs.String("payment-token", "0x0000000000000000000000000000000000000000", "ERC-20 token address for payment (default: ETH)")
	providerAgentID := fs.String("provider-agent-id", "0", "provider agent ID for createJob (default: 0)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *initConfig {
		out := *outPath
		if strings.TrimSpace(out) == "" {
			out = *cfgPath
		}
		return writeInitConfig(out, *force)
	}

	cfg, err := readSeedConfig(*cfgPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(*command) != "" {
		cfg.Command = *command
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = 97
	}
	if strings.TrimSpace(cfg.RPCURL) == "" || strings.TrimSpace(cfg.Commerce) == "" {
		return errors.New("config requires rpc_url and commerce")
	}

	commerceAddr := common.HexToAddress(cfg.Commerce)
	if commerceAddr == (common.Address{}) {
		return errors.New("invalid config.commerce address")
	}
	chainIDBig := big.NewInt(cfg.ChainID)

	ctx := context.Background()
	ec, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return fmt.Errorf("dial rpc: %w", err)
	}
	defer ec.Close()

	contract, err := cabi.NewAgenticCommerce(commerceAddr, ec)
	if err != nil {
		return fmt.Errorf("bind commerce: %w", err)
	}

	switch strings.TrimSpace(cfg.Command) {
	case "check":
		return runCheck(ctx, ec, contract, commerceAddr, chainIDBig, cfg, *jobIDStr, *paymentToken)
	case "happy-path":
		return runHappyPath(ctx, ec, contract, commerceAddr, chainIDBig,
			cfg.Client.PrivateKey, cfg.Provider.PrivateKey, cfg.Evaluator.PrivateKey,
			cfg.CreateJob.Provider, cfg.CreateJob.Evaluator, cfg.Budget, cfg.CreateJob.ExpiredAt,
			cfg.CreateJob.Description, cfg.CreateJob.Hook, cfg.OptParams, cfg.Deliverable, cfg.Reason, cfg.GasLimit,
			*paymentToken, *providerAgentID)
	case "create-job":
		return runCreateJob(ctx, ec, contract, commerceAddr, chainIDBig,
			cfg.Client.PrivateKey, cfg.CreateJob.Provider, cfg.CreateJob.Evaluator, cfg.CreateJob.ExpiredAt,
			cfg.CreateJob.Description, cfg.CreateJob.Hook, cfg.GasLimit, *paymentToken, *providerAgentID)
	case "set-budget":
		return runSetBudget(ctx, ec, contract, commerceAddr, chainIDBig, cfg.Client.PrivateKey, *jobIDStr, cfg.Budget, cfg.OptParams, cfg.GasLimit, *paymentToken)
	case "approve-fund":
		return runApproveFund(ctx, ec, contract, commerceAddr, chainIDBig, cfg.Client.PrivateKey, *jobIDStr, cfg.Budget, cfg.OptParams, cfg.GasLimit, *paymentToken)
	case "fund-only":
		return runFundOnly(ctx, ec, contract, chainIDBig, cfg.Client.PrivateKey, *jobIDStr, cfg.Budget, cfg.OptParams, cfg.GasLimit, *paymentToken)
	case "submit-only":
		return runSubmitOnly(ctx, ec, contract, chainIDBig, cfg.Provider.PrivateKey, *jobIDStr, cfg.CreateJob.Description, cfg.Deliverable, cfg.OptParams, cfg.GasLimit)
	case "complete-only":
		return runCompleteOnly(ctx, ec, contract, chainIDBig, cfg.Evaluator.PrivateKey, cfg.Provider.PrivateKey, *jobIDStr, cfg.Reason, cfg.OptParams, cfg.GasLimit)
	default:
		return fmt.Errorf("unknown command %q (supported: check, happy-path, create-job, set-budget, approve-fund, fund-only, submit-only, complete-only)", cfg.Command)
	}
}

func readSeedConfig(path string) (*seedConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	cfg := &seedConfig{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("parse yaml %s: %w", path, err)
	}
	if strings.TrimSpace(cfg.Command) == "" {
		cfg.Command = "happy-path"
	}
	if strings.TrimSpace(cfg.CreateJob.Hook) == "" {
		cfg.CreateJob.Hook = "0x0000000000000000000000000000000000000000"
	}
	if strings.TrimSpace(cfg.Budget) == "" {
		cfg.Budget = "1000000"
	}
	if strings.TrimSpace(cfg.OptParams) == "" {
		cfg.OptParams = "0x"
	}
	if strings.TrimSpace(cfg.Reason) == "" {
		cfg.Reason = "0x0000000000000000000000000000000000000000000000000000000000000000"
	}
	return cfg, nil
}

func writeInitConfig(outPath string, force bool) error {
	outPath = strings.TrimSpace(outPath)
	if outPath == "" {
		return errors.New("empty -out path")
	}
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("config already exists: %s (use -force to overwrite)", outPath)
		}
	}

	gen := func() (privHex string, addr string, err error) {
		k, err := crypto.GenerateKey()
		if err != nil {
			return "", "", err
		}
		privHex = "0x" + hex.EncodeToString(crypto.FromECDSA(k))
		addr = crypto.PubkeyToAddress(k.PublicKey).Hex()
		return privHex, addr, nil
	}

	cPriv, cAddr, err := gen()
	if err != nil {
		return err
	}
	pPriv, pAddr, err := gen()
	if err != nil {
		return err
	}
	ePriv, eAddr, err := gen()
	if err != nil {
		return err
	}

	cfg := seedConfig{
		RPCURL:    "https://bsc-testnet-rpc.example",
		ChainID:   97,
		Commerce:  "0xYourAgenticCommerce",
		Budget:    "1000000",
		OptParams: "0x",
		Reason:    "0x0000000000000000000000000000000000000000000000000000000000000000",
		GasLimit:  0,
		Command:   "happy-path",
	}
	cfg.Client.PrivateKey = cPriv
	cfg.Client.Address = cAddr
	cfg.Provider.PrivateKey = pPriv
	cfg.Provider.Address = pAddr
	cfg.Evaluator.PrivateKey = ePriv
	cfg.Evaluator.Address = eAddr

	cfg.CreateJob.Provider = pAddr
	cfg.CreateJob.Evaluator = eAddr
	cfg.CreateJob.ExpiredAt = 0
	cfg.CreateJob.Description = "commerce_seed test job"
	cfg.CreateJob.Hook = "0x0000000000000000000000000000000000000000"

	y, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(outPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(outPath, y, 0o600); err != nil {
		return fmt.Errorf("write config %s: %w", outPath, err)
	}

	fmt.Printf("wrote config: %s\n", outPath)
	fmt.Printf("client.address: %s\nprovider.address: %s\nevaluator.address: %s\n", cAddr, pAddr, eAddr)
	fmt.Printf("next: edit rpc_url / commerce, fund wallets with tBNB + paymentToken, then run: go run ./cmd/commerce_seed/ -f %s\n", outPath)
	return nil
}

func parsePrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	hexKey = strings.TrimSpace(hexKey)
	if hexKey == "" {
		return nil, errors.New("empty private key")
	}
	hexKey = strings.TrimPrefix(hexKey, "0x")
	if len(hexKey) != 64 {
		return nil, fmt.Errorf("private key must be 32 bytes hex, got len %d", len(hexKey))
	}
	raw, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	return crypto.ToECDSA(raw)
}

func transactor(key *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	auth.Context = context.Background()
	return auth, nil
}

func applyGasLimit(auth *bind.TransactOpts, limit uint64) {
	if limit > 0 {
		auth.GasLimit = limit
	}
}

func waitMine(ctx context.Context, ec *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, ec, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("tx reverted: %s status=%d", tx.Hash().Hex(), receipt.Status)
	}
	return receipt, nil
}

func parseJobIDFromReceipt(contract *cabi.AgenticCommerce, receipt *types.Receipt) (*big.Int, error) {
	for i := range receipt.Logs {
		if receipt.Logs[i] == nil {
			continue
		}
		ev, err := contract.ParseJobCreated(*receipt.Logs[i])
		if err != nil {
			continue
		}
		return ev.JobId, nil
	}
	return nil, errors.New("no JobCreated log in receipt")
}

func findJobIDByFilterLogs(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, receipt *types.Receipt) (*big.Int, error) {
	parsed, err := cabi.AgenticCommerceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	ev := parsed.Events["JobCreated"]
	q := ethereum.FilterQuery{
		FromBlock: receipt.BlockNumber,
		ToBlock:   receipt.BlockNumber,
		Addresses: []common.Address{commerceAddr},
		Topics:    [][]common.Hash{{ev.ID}},
	}
	logs, err := ec.FilterLogs(ctx, q)
	if err != nil {
		return nil, err
	}
	for i := range logs {
		if logs[i].TxHash != receipt.TxHash {
			continue
		}
		parsedEv, err := contract.ParseJobCreated(logs[i])
		if err != nil {
			continue
		}
		return parsedEv.JobId, nil
	}
	return nil, errors.New("JobCreated not found via FilterLogs in tx")
}

func printReceiptLogSummary(receipt *types.Receipt) {
	if receipt == nil {
		return
	}
	fmt.Printf("\n[receipt]\n")
	fmt.Printf("txHash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("blockNumber: %s\n", receipt.BlockNumber.String())
	fmt.Printf("logs: %d\n", len(receipt.Logs))
	for i := range receipt.Logs {
		l := receipt.Logs[i]
		if l == nil {
			continue
		}
		t0 := common.Hash{}
		if len(l.Topics) > 0 {
			t0 = l.Topics[0]
		}
		fmt.Printf("log[%d]: address=%s topics0=%s topics=%d\n", i, l.Address.Hex(), t0.Hex(), len(l.Topics))
	}
}

func findJobIDAfterCreate(ctx context.Context, contract *cabi.AgenticCommerce, expectedClient common.Address, expectedProvider common.Address, beforeCounter *big.Int) (*big.Int, error) {
	afterCounter, err := contract.JobCounter(nil)
	if err != nil {
		return nil, fmt.Errorf("jobCounter(after): %w", err)
	}
	if beforeCounter != nil && afterCounter.Cmp(beforeCounter) <= 0 {
		return nil, fmt.Errorf("jobCounter did not increase (before=%s after=%s)", beforeCounter.String(), afterCounter.String())
	}

	// candidate job ids: afterCounter (1-based) and afterCounter-1 (0-based)
	candidates := []*big.Int{new(big.Int).Set(afterCounter)}
	if afterCounter.Sign() > 0 {
		candidates = append(candidates, new(big.Int).Sub(afterCounter, big.NewInt(1)))
	}
	for _, id := range candidates {
		if id.Sign() < 0 {
			continue
		}
		j, err := contract.GetJob(nil, id)
		if err != nil {
			continue
		}
		if j.Client == expectedClient && j.Provider == expectedProvider {
			return id, nil
		}
	}
	return nil, fmt.Errorf("cannot determine jobId from jobCounter=%s (candidates=%s,%s)", afterCounter.String(), candidates[0].String(), candidates[len(candidates)-1].String())
}

func findJobIDFromReceiptHeuristic(receipt *types.Receipt, commerceAddr common.Address, expectedClient common.Address, expectedProvider common.Address) (*big.Int, error) {
	if receipt == nil {
		return nil, errors.New("nil receipt")
	}
	for i := range receipt.Logs {
		l := receipt.Logs[i]
		if l == nil {
			continue
		}
		if l.Address != commerceAddr {
			continue
		}
		// Try to detect which topics correspond to client/provider (addresses).
		var t1Addr, t2Addr common.Address
		if len(l.Topics) > 1 {
			t1Addr = common.BytesToAddress(l.Topics[1].Bytes()[12:])
		}
		if len(l.Topics) > 2 {
			t2Addr = common.BytesToAddress(l.Topics[2].Bytes()[12:])
		}

		// Case A: one topic matches client, the other matches provider => jobId likely in data (first 32 bytes)
		if (t1Addr == expectedClient && t2Addr == expectedProvider) || (t2Addr == expectedClient && t1Addr == expectedProvider) {
			if len(l.Data) >= 32 {
				return new(big.Int).SetBytes(l.Data[:32]), nil
			}
		}

		// Case B: one topic matches client, the other is probably jobId (uint256 indexed)
		if t1Addr == expectedClient {
			return new(big.Int).SetBytes(l.Topics[2].Bytes()), nil
		}
		if t2Addr == expectedClient {
			return new(big.Int).SetBytes(l.Topics[1].Bytes()), nil
		}

		// Fallback: assume topics[1] is jobId (common pattern jobId indexed)
		if len(l.Topics) > 1 {
			return new(big.Int).SetBytes(l.Topics[1].Bytes()), nil
		}
	}
	return nil, errors.New("cannot derive jobId from receipt logs")
}

func parseOptParams(hexStr string) ([]byte, error) {
	hexStr = strings.TrimSpace(hexStr)
	if hexStr == "" || hexStr == "0x" {
		return nil, nil
	}
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("opt-params: %w", err)
	}
	return b, nil
}

func parseBytes32FromHex(s string, def *[32]byte) ([32]byte, error) {
	t := strings.TrimSpace(s)
	if t == "" || t == "0x" {
		if def != nil {
			return *def, nil
		}
		return [32]byte{}, errors.New("empty hex for 32 bytes")
	}
	t = strings.TrimPrefix(t, "0x")
	if len(t) != 64 {
		return [32]byte{}, fmt.Errorf("expected 64 hex chars (32 bytes), got %d", len(t))
	}
	raw, err := hex.DecodeString(t)
	if err != nil {
		return [32]byte{}, err
	}
	var out [32]byte
	copy(out[:], raw)
	return out, nil
}

func parseBudget(s string) (*big.Int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty budget")
	}
	b, ok := new(big.Int).SetString(s, 10)
	if !ok || b.Sign() <= 0 {
		return nil, fmt.Errorf("invalid budget %q", s)
	}
	return b, nil
}

func parseJobIDDecimal(s string) (*big.Int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("job-id is required")
	}
	b, ok := new(big.Int).SetString(s, 10)
	if !ok || b.Sign() < 0 {
		return nil, fmt.Errorf("invalid job-id %q", s)
	}
	return b, nil
}

func erc20Approve(ctx context.Context, ec *ethclient.Client, token common.Address, ownerKey *ecdsa.PrivateKey, spender common.Address, amount *big.Int, chainID *big.Int, gasLimit uint64) (*types.Transaction, error) {
	parsed, err := abi.JSON(strings.NewReader(erc20ApproveABI))
	if err != nil {
		return nil, err
	}
	data, err := parsed.Pack("approve", spender, amount)
	if err != nil {
		return nil, err
	}
	ownerAddr := crypto.PubkeyToAddress(ownerKey.PublicKey)
	nonce, err := ec.PendingNonceAt(ctx, ownerAddr)
	if err != nil {
		return nil, err
	}
	gasPrice, err := ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: ownerAddr, To: &token, Gas: 0, GasPrice: gasPrice, Value: big.NewInt(0), Data: data}
	gl := gasLimit
	if gl == 0 {
		est, err := ec.EstimateGas(ctx, msg)
		if err != nil {
			gl = 100000
		} else {
			gl = est
		}
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gl,
		To:       &token,
		Value:    big.NewInt(0),
		Data:     data,
	})
	signer := types.LatestSignerForChainID(chainID)
	signed, err := types.SignTx(tx, signer, ownerKey)
	if err != nil {
		return nil, err
	}
	if err := ec.SendTransaction(ctx, signed); err != nil {
		return nil, err
	}
	return signed, nil
}

func printIndexerHints(ctx context.Context, ec *ethclient.Client, commerceAddr common.Address, chainID int64) error {
	head, err := ec.BlockNumber(ctx)
	if err != nil {
		return err
	}
	const buf = uint64(50)
	var start uint64
	if head > buf {
		start = head - buf
	}
	fmt.Printf("\n--- indexer / DB hints ---\n")
	fmt.Printf("chain_id: %d\n", chainID)
	fmt.Printf("commerce.addr: %s\n", commerceAddr.Hex())
	fmt.Printf("suggested commerce.start_block: %d (head=%d, buffer=%d)\n", start, head, buf)
	fmt.Printf("paymentToken: call contract paymentToken() or check deployment; client must hold balance and approve commerce.\n")
	return nil
}

func expiryUnix(expiredAt int64) *big.Int {
	if expiredAt > 0 {
		return big.NewInt(expiredAt)
	}
	return big.NewInt(time.Now().Add(7 * 24 * time.Hour).Unix())
}

func deliverableBytes32(desc string, deliverableHex string) ([32]byte, error) {
	if strings.TrimSpace(deliverableHex) != "" {
		return parseBytes32FromHex(deliverableHex, nil)
	}
	h := crypto.Keccak256([]byte(desc))
	var out [32]byte
	copy(out[:], h)
	return out, nil
}

func evaluatorKeyOrProvider(evalHex, provHex string) (string, error) {
	if strings.TrimSpace(evalHex) != "" {
		return evalHex, nil
	}
	if strings.TrimSpace(provHex) == "" {
		return "", errors.New("evaluator-key or provider-key is required")
	}
	return provHex, nil
}

func addrFromPrivOrConfig(privHex, addrHex string) (common.Address, error) {
	if strings.TrimSpace(addrHex) != "" {
		a := common.HexToAddress(addrHex)
		if a == (common.Address{}) {
			return common.Address{}, fmt.Errorf("invalid address %q", addrHex)
		}
		return a, nil
	}
	k, err := parsePrivateKey(privHex)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(k.PublicKey), nil
}

func erc20CallBigInt(ctx context.Context, ec *ethclient.Client, token common.Address, method string, args ...any) (*big.Int, error) {
	parsed, err := abi.JSON(strings.NewReader(erc20BalanceAllowanceABI))
	if err != nil {
		return nil, err
	}
	data, err := parsed.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := ec.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	res, err := parsed.Unpack(method, out)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("unexpected %s outputs: %d", method, len(res))
	}
	v, ok := res[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected %s output type %T", method, res[0])
	}
	return v, nil
}

func runCheck(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, chainID *big.Int, cfg *seedConfig, jobIDStr string, paymentToken string) error {
	// 解析地址（优先用配置中的 address；为空则由私钥推导）
	clientAddr, err := addrFromPrivOrConfig(cfg.Client.PrivateKey, cfg.Client.Address)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}
	providerAddr, err := addrFromPrivOrConfig(cfg.Provider.PrivateKey, cfg.Provider.Address)
	if err != nil {
		return fmt.Errorf("provider: %w", err)
	}
	evaluatorKeyHex, err := evaluatorKeyOrProvider(cfg.Evaluator.PrivateKey, cfg.Provider.PrivateKey)
	if err != nil {
		return err
	}
	evaluatorAddr, err := addrFromPrivOrConfig(evaluatorKeyHex, cfg.Evaluator.Address)
	if err != nil {
		return fmt.Errorf("evaluator: %w", err)
	}

	budget, err := parseBudget(cfg.Budget)
	if err != nil {
		return err
	}

	// 1) Native 余额（tBNB）
	cBal, err := ec.BalanceAt(ctx, clientAddr, nil)
	if err != nil {
		return err
	}
	pBal, err := ec.BalanceAt(ctx, providerAddr, nil)
	if err != nil {
		return err
	}
	eBal, err := ec.BalanceAt(ctx, evaluatorAddr, nil)
	if err != nil {
		return err
	}

	// 2) paymentToken 余额与 allowance
	tokenAddr := common.HexToAddress(paymentToken)
	cTokenBal, err := erc20CallBigInt(ctx, ec, tokenAddr, "balanceOf", clientAddr)
	if err != nil {
		return fmt.Errorf("erc20 balanceOf(client): %w", err)
	}
	allow, err := erc20CallBigInt(ctx, ec, tokenAddr, "allowance", clientAddr, commerceAddr)
	if err != nil {
		return fmt.Errorf("erc20 allowance(client->commerce): %w", err)
	}

	fmt.Printf("--- preflight check ---\n")
	fmt.Printf("chain_id: %s\n", chainID.String())
	fmt.Printf("commerce: %s\n", commerceAddr.Hex())
	fmt.Printf("paymentToken: %s\n", tokenAddr.Hex())
	fmt.Printf("\n[accounts]\n")
	fmt.Printf("client:    %s\n", clientAddr.Hex())
	fmt.Printf("provider:  %s\n", providerAddr.Hex())
	fmt.Printf("evaluator: %s\n", evaluatorAddr.Hex())
	fmt.Printf("\n[native balances]\n")
	fmt.Printf("client.tBNB:    %s wei\n", cBal.String())
	fmt.Printf("provider.tBNB:  %s wei\n", pBal.String())
	fmt.Printf("evaluator.tBNB: %s wei\n", eBal.String())
	fmt.Printf("\n[token]\n")
	fmt.Printf("budget(required): %s\n", budget.String())
	fmt.Printf("client.tokenBalance: %s\n", cTokenBal.String())
	fmt.Printf("client.allowance->commerce: %s\n", allow.String())

	if cTokenBal.Cmp(budget) < 0 {
		fmt.Printf("\nRESULT: FAIL — client paymentToken balance < budget\n")
	} else if allow.Cmp(budget) < 0 {
		fmt.Printf("\nRESULT: WARN — allowance < budget (happy-path/approve-fund will send approve)\n")
	} else {
		fmt.Printf("\nRESULT: OK — token balance and allowance look sufficient\n")
	}

	// 可选：如果提供 job-id，则顺便打印链上 job 的预算，提前发现 BudgetMismatch
	if strings.TrimSpace(jobIDStr) != "" {
		jobID, err := parseJobIDDecimal(jobIDStr)
		if err != nil {
			return err
		}
		j, err := contract.GetJob(nil, jobID)
		if err != nil {
			return fmt.Errorf("getJob(%s): %w", jobID.String(), err)
		}
		fmt.Printf("\n[job]\n")
		fmt.Printf("jobId: %s\n", jobID.String())
		fmt.Printf("job.budget(on-chain): %s\n", j.Budget.String())
		if j.Budget.Cmp(budget) != 0 {
			fmt.Printf("job.budget != config.budget — fund may revert with BudgetMismatch (consider set-budget)\n")
		}
	}

	return nil
}

func runHappyPath(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, chainID *big.Int,
	clientKeyHex, providerKeyHex, evaluatorKeyHex, providerStr, evaluatorStr, budgetStr string, expiredAt int64, desc, hookStr, optParamsHex, deliverableHex, reasonHex string, gasLimit uint64,
	paymentToken string, providerAgentID string,
) error {
	clientKey, err := parsePrivateKey(clientKeyHex)
	if err != nil {
		return fmt.Errorf("client-key: %w", err)
	}
	providerKey, err := parsePrivateKey(providerKeyHex)
	if err != nil {
		return fmt.Errorf("provider-key: %w", err)
	}
	evalKeyHex, err := evaluatorKeyOrProvider(evaluatorKeyHex, providerKeyHex)
	if err != nil {
		return err
	}
	evalKey, err := parsePrivateKey(evalKeyHex)
	if err != nil {
		return fmt.Errorf("evaluator-key: %w", err)
	}

	if strings.TrimSpace(providerStr) == "" || strings.TrimSpace(evaluatorStr) == "" {
		return errors.New("-provider and -evaluator addresses are required")
	}
	provider := common.HexToAddress(providerStr)
	evaluator := common.HexToAddress(evaluatorStr)
	hook := common.HexToAddress(hookStr)

	budget, err := parseBudget(budgetStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}
	reason, err := parseBytes32FromHex(reasonHex, nil)
	if err != nil {
		return err
	}
	deliv, err := deliverableBytes32(desc, deliverableHex)
	if err != nil {
		return err
	}

	tokenAddr := common.HexToAddress(paymentToken)
	fmt.Printf("paymentToken: %s\n", tokenAddr.Hex())

	// 1) createJob
	cAuth, err := transactor(clientKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(cAuth, gasLimit)
	providerAgentIDBig, _ := new(big.Int).SetString(providerAgentID, 10)
	tx1, err := contract.CreateJob(cAuth, provider, evaluator, expiryUnix(expiredAt), desc, hook, providerAgentIDBig)
	if err != nil {
		return fmt.Errorf("createJob: %w", err)
	}
	fmt.Printf("createJob tx: %s\n", tx1.Hash().Hex())
	r1, err := waitMine(ctx, ec, tx1)
	if err != nil {
		return err
	}
	jobID, err := parseJobIDFromReceipt(contract, r1)
	if err != nil {
		// 合约 event 版本可能与本地 ABI 不一致（例如 indexed 数变化/参数变化）
		// 回退：用 jobCounter/getJob 推导本次新 jobId
		printReceiptLogSummary(r1)
		jobID, err = findJobIDFromReceiptHeuristic(r1, commerceAddr, cAuth.From, provider)
		if err != nil {
			// 再尝试通过 FilterLogs（仅当 topic 匹配本地 ABI）
			jobID, _ = findJobIDByFilterLogs(ctx, ec, contract, commerceAddr, r1)
			if jobID == nil {
				return err
			}
		}
	}
	fmt.Printf("jobId: %s\n", jobID.String())

	// 2) setBudget (链上 budget 须与 fund 的 expectedBudget 一致)
	tx2, err := contract.SetBudget(cAuth, jobID, tokenAddr, budget, optParams)
	if err != nil {
		return fmt.Errorf("setBudget: %w", err)
	}
	fmt.Printf("setBudget tx: %s\n", tx2.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx2); err != nil {
		return err
	}

	// 3) approve + fund
	tx3, err := erc20Approve(ctx, ec, tokenAddr, clientKey, commerceAddr, budget, chainID, 0)
	if err != nil {
		return fmt.Errorf("approve: %w", err)
	}
	fmt.Printf("ERC20 approve tx: %s\n", tx3.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx3); err != nil {
		return err
	}

	tx4, err := contract.Fund(cAuth, jobID, budget, optParams)
	if err != nil {
		return fmt.Errorf("fund: %w", err)
	}
	fmt.Printf("fund tx: %s\n", tx4.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx4); err != nil {
		return err
	}

	// 4) submit
	pAuth, err := transactor(providerKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(pAuth, gasLimit)
	tx5, err := contract.Submit(pAuth, jobID, deliv, optParams)
	if err != nil {
		return fmt.Errorf("submit: %w", err)
	}
	fmt.Printf("submit tx: %s\n", tx5.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx5); err != nil {
		return err
	}

	// 5) complete
	eAuth, err := transactor(evalKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(eAuth, gasLimit)
	tx6, err := contract.Complete(eAuth, jobID, reason, optParams)
	if err != nil {
		return fmt.Errorf("complete: %w", err)
	}
	fmt.Printf("complete tx: %s\n", tx6.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx6); err != nil {
		return err
	}

	_ = printIndexerHints(ctx, ec, commerceAddr, chainID.Int64())
	return nil
}

func runCreateJob(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, chainID *big.Int,
	clientKeyHex, providerStr, evaluatorStr string, expiredAt int64, desc, hookStr string, gasLimit uint64,
	paymentToken string, providerAgentID string,
) error {
	clientKey, err := parsePrivateKey(clientKeyHex)
	if err != nil {
		return fmt.Errorf("client-key: %w", err)
	}
	if strings.TrimSpace(providerStr) == "" || strings.TrimSpace(evaluatorStr) == "" {
		return errors.New("-provider and -evaluator addresses are required")
	}
	provider := common.HexToAddress(providerStr)
	evaluator := common.HexToAddress(evaluatorStr)
	hook := common.HexToAddress(hookStr)

	tokenAddr := common.HexToAddress(paymentToken)
	fmt.Printf("paymentToken: %s\n", tokenAddr.Hex())

	beforeCounter, err := contract.JobCounter(nil)
	if err != nil {
		return fmt.Errorf("jobCounter(before): %w", err)
	}

	cAuth, err := transactor(clientKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(cAuth, gasLimit)
	providerAgentIDBig, _ := new(big.Int).SetString(providerAgentID, 10)
	tx, err := contract.CreateJob(cAuth, provider, evaluator, expiryUnix(expiredAt), desc, hook, providerAgentIDBig)
	if err != nil {
		return fmt.Errorf("createJob: %w", err)
	}
	fmt.Printf("createJob tx: %s\n", tx.Hash().Hex())
	r, err := waitMine(ctx, ec, tx)
	if err != nil {
		return err
	}
	jobID, err := parseJobIDFromReceipt(contract, r)
	if err != nil {
		printReceiptLogSummary(r)
		jobID, err = findJobIDAfterCreate(ctx, contract, cAuth.From, provider, beforeCounter)
		if err != nil {
			return err
		}
	}
	fmt.Printf("jobId: %s\n", jobID.String())
	fmt.Printf("next: -command set-budget then approve-fund (or use happy-path)\n")
	return printIndexerHints(ctx, ec, commerceAddr, chainID.Int64())
}

func runSetBudget(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, chainID *big.Int,
	clientKeyHex, jobIDStr, budgetStr, optParamsHex string, gasLimit uint64, paymentToken string,
) error {
	clientKey, err := parsePrivateKey(clientKeyHex)
	if err != nil {
		return fmt.Errorf("client-key: %w", err)
	}
	jobID, err := parseJobIDDecimal(jobIDStr)
	if err != nil {
		return err
	}
	budget, err := parseBudget(budgetStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}
	cAuth, err := transactor(clientKey, chainID)
	if err != nil {
		return err
	}
	tokenAddr := common.HexToAddress(paymentToken)
	applyGasLimit(cAuth, gasLimit)
	tx, err := contract.SetBudget(cAuth, jobID, tokenAddr, budget, optParams)
	if err != nil {
		return fmt.Errorf("setBudget: %w", err)
	}
	fmt.Printf("setBudget tx: %s\n", tx.Hash().Hex())
	if _, err := waitMine(ctx, ec, tx); err != nil {
		return err
	}
	return printIndexerHints(ctx, ec, commerceAddr, chainID.Int64())
}

func runApproveFund(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, commerceAddr common.Address, chainID *big.Int,
	clientKeyHex, jobIDStr, budgetStr, optParamsHex string, gasLimit uint64, paymentToken string,
) error {
	clientKey, err := parsePrivateKey(clientKeyHex)
	if err != nil {
		return fmt.Errorf("client-key: %w", err)
	}
	jobID, err := parseJobIDDecimal(jobIDStr)
	if err != nil {
		return err
	}
	budget, err := parseBudget(budgetStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}

	tokenAddr := common.HexToAddress(paymentToken)
	j, err := contract.GetJob(nil, jobID)
	if err != nil {
		return fmt.Errorf("getJob: %w", err)
	}
	if j.Budget.Cmp(budget) != 0 {
		fmt.Fprintf(os.Stderr, "warning: -budget %s != on-chain job.budget %s; fund may revert with BudgetMismatch\n", budget.String(), j.Budget.String())
	}

	cAuth, err := transactor(clientKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(cAuth, gasLimit)

	txA, err := erc20Approve(ctx, ec, tokenAddr, clientKey, commerceAddr, budget, chainID, 0)
	if err != nil {
		return fmt.Errorf("approve: %w", err)
	}
	fmt.Printf("ERC20 approve tx: %s\n", txA.Hash().Hex())
	if _, err := waitMine(ctx, ec, txA); err != nil {
		return err
	}

	txF, err := contract.Fund(cAuth, jobID, budget, optParams)
	if err != nil {
		return fmt.Errorf("fund: %w", err)
	}
	fmt.Printf("fund tx: %s\n", txF.Hash().Hex())
	if _, err := waitMine(ctx, ec, txF); err != nil {
		return err
	}
	return printIndexerHints(ctx, ec, commerceAddr, chainID.Int64())
}

func runFundOnly(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, chainID *big.Int,
	clientKeyHex, jobIDStr, budgetStr, optParamsHex string, gasLimit uint64, paymentToken string,
) error {
	clientKey, err := parsePrivateKey(clientKeyHex)
	if err != nil {
		return fmt.Errorf("client-key: %w", err)
	}
	jobID, err := parseJobIDDecimal(jobIDStr)
	if err != nil {
		return err
	}
	budget, err := parseBudget(budgetStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}
	cAuth, err := transactor(clientKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(cAuth, gasLimit)
	tx, err := contract.Fund(cAuth, jobID, budget, optParams)
	if err != nil {
		return fmt.Errorf("fund: %w", err)
	}
	fmt.Printf("fund tx: %s\n", tx.Hash().Hex())
	_, err = waitMine(ctx, ec, tx)
	return err
}

func runSubmitOnly(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, chainID *big.Int,
	providerKeyHex, jobIDStr, desc, deliverableHex, optParamsHex string, gasLimit uint64,
) error {
	providerKey, err := parsePrivateKey(providerKeyHex)
	if err != nil {
		return fmt.Errorf("provider-key: %w", err)
	}
	jobID, err := parseJobIDDecimal(jobIDStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}
	deliv, err := deliverableBytes32(desc, deliverableHex)
	if err != nil {
		return err
	}
	pAuth, err := transactor(providerKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(pAuth, gasLimit)
	tx, err := contract.Submit(pAuth, jobID, deliv, optParams)
	if err != nil {
		return fmt.Errorf("submit: %w", err)
	}
	fmt.Printf("submit tx: %s\n", tx.Hash().Hex())
	_, err = waitMine(ctx, ec, tx)
	return err
}

func runCompleteOnly(ctx context.Context, ec *ethclient.Client, contract *cabi.AgenticCommerce, chainID *big.Int,
	evaluatorKeyHex, providerKeyHex, jobIDStr, reasonHex, optParamsHex string, gasLimit uint64,
) error {
	evalKeyHex, err := evaluatorKeyOrProvider(evaluatorKeyHex, providerKeyHex)
	if err != nil {
		return err
	}
	evalKey, err := parsePrivateKey(evalKeyHex)
	if err != nil {
		return fmt.Errorf("evaluator-key: %w", err)
	}
	jobID, err := parseJobIDDecimal(jobIDStr)
	if err != nil {
		return err
	}
	optParams, err := parseOptParams(optParamsHex)
	if err != nil {
		return err
	}
	reason, err := parseBytes32FromHex(reasonHex, nil)
	if err != nil {
		return err
	}
	eAuth, err := transactor(evalKey, chainID)
	if err != nil {
		return err
	}
	applyGasLimit(eAuth, gasLimit)
	tx, err := contract.Complete(eAuth, jobID, reason, optParams)
	if err != nil {
		return fmt.Errorf("complete: %w", err)
	}
	fmt.Printf("complete tx: %s\n", tx.Hash().Hex())
	_, err = waitMine(ctx, ec, tx)
	return err
}
