// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AgenticCommerceUpgradeableJob is an auto generated low-level Go binding around an user-defined struct.
type AgenticCommerceUpgradeableJob struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}

// AgenticCommerceMetaData contains all meta data concerning the AgenticCommerce contract.
var AgenticCommerceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"trustedForwarder_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BudgetMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpiryTooShort\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeesTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HookCallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HookNotWhitelisted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidJob\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProviderNotSet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroBudget\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BudgetSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EvaluatorFeePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"HookWhitelistUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reason\",\"type\":\"bytes32\"}],\"name\":\"JobCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiredAt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"}],\"name\":\"JobCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"}],\"name\":\"JobExpired\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"JobFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rejector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"reason\",\"type\":\"bytes32\"}],\"name\":\"JobRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"deliverable\",\"type\":\"bytes32\"}],\"name\":\"JobSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"}],\"name\":\"ProviderSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"role\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int8\",\"name\":\"signal\",\"type\":\"int8\"}],\"name\":\"ReputationSignal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HOOK_GAS_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"}],\"name\":\"claimRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"reason\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"complete\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiredAt\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"}],\"name\":\"createJob\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"evaluatorFeeBP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedBudget\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"fund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedBudget\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"fundWithPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"}],\"name\":\"getJob\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"budget\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiredAt\",\"type\":\"uint256\"},{\"internalType\":\"enumAgenticCommerceUpgradeable.JobStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"}],\"internalType\":\"structAgenticCommerceUpgradeable.Job\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"paymentToken_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"treasury_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isTrustedForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"jobCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"}],\"name\":\"jobHasBudget\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"hasBudget\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"jobs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"evaluator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"budget\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiredAt\",\"type\":\"uint256\"},{\"internalType\":\"enumAgenticCommerceUpgradeable.JobStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paymentToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformFeeBP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformTreasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"reason\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"reject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"setBudget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeBP_\",\"type\":\"uint256\"}],\"name\":\"setEvaluatorFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"hook\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setHookWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"feeBP_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"treasury_\",\"type\":\"address\"}],\"name\":\"setPlatformFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"provider_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"setProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"deliverable\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"optParams\",\"type\":\"bytes\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedForwarder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"whitelistedHooks\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AgenticCommerceABI is the input ABI used to generate the binding from.
// Deprecated: Use AgenticCommerceMetaData.ABI instead.
var AgenticCommerceABI = AgenticCommerceMetaData.ABI

// AgenticCommerce is an auto generated Go binding around an Ethereum contract.
type AgenticCommerce struct {
	AgenticCommerceCaller     // Read-only binding to the contract
	AgenticCommerceTransactor // Write-only binding to the contract
	AgenticCommerceFilterer   // Log filterer for contract events
}

// AgenticCommerceCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgenticCommerceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgenticCommerceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgenticCommerceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgenticCommerceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgenticCommerceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgenticCommerceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgenticCommerceSession struct {
	Contract     *AgenticCommerce  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgenticCommerceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgenticCommerceCallerSession struct {
	Contract *AgenticCommerceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AgenticCommerceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgenticCommerceTransactorSession struct {
	Contract     *AgenticCommerceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AgenticCommerceRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgenticCommerceRaw struct {
	Contract *AgenticCommerce // Generic contract binding to access the raw methods on
}

// AgenticCommerceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgenticCommerceCallerRaw struct {
	Contract *AgenticCommerceCaller // Generic read-only contract binding to access the raw methods on
}

// AgenticCommerceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgenticCommerceTransactorRaw struct {
	Contract *AgenticCommerceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgenticCommerce creates a new instance of AgenticCommerce, bound to a specific deployed contract.
func NewAgenticCommerce(address common.Address, backend bind.ContractBackend) (*AgenticCommerce, error) {
	contract, err := bindAgenticCommerce(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerce{AgenticCommerceCaller: AgenticCommerceCaller{contract: contract}, AgenticCommerceTransactor: AgenticCommerceTransactor{contract: contract}, AgenticCommerceFilterer: AgenticCommerceFilterer{contract: contract}}, nil
}

// NewAgenticCommerceCaller creates a new read-only instance of AgenticCommerce, bound to a specific deployed contract.
func NewAgenticCommerceCaller(address common.Address, caller bind.ContractCaller) (*AgenticCommerceCaller, error) {
	contract, err := bindAgenticCommerce(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceCaller{contract: contract}, nil
}

// NewAgenticCommerceTransactor creates a new write-only instance of AgenticCommerce, bound to a specific deployed contract.
func NewAgenticCommerceTransactor(address common.Address, transactor bind.ContractTransactor) (*AgenticCommerceTransactor, error) {
	contract, err := bindAgenticCommerce(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceTransactor{contract: contract}, nil
}

// NewAgenticCommerceFilterer creates a new log filterer instance of AgenticCommerce, bound to a specific deployed contract.
func NewAgenticCommerceFilterer(address common.Address, filterer bind.ContractFilterer) (*AgenticCommerceFilterer, error) {
	contract, err := bindAgenticCommerce(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceFilterer{contract: contract}, nil
}

// bindAgenticCommerce binds a generic wrapper to an already deployed contract.
func bindAgenticCommerce(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgenticCommerceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgenticCommerce *AgenticCommerceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgenticCommerce.Contract.AgenticCommerceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgenticCommerce *AgenticCommerceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.AgenticCommerceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgenticCommerce *AgenticCommerceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.AgenticCommerceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgenticCommerce *AgenticCommerceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgenticCommerce.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgenticCommerce *AgenticCommerceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgenticCommerce *AgenticCommerceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceSession) ADMINROLE() ([32]byte, error) {
	return _AgenticCommerce.Contract.ADMINROLE(&_AgenticCommerce.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCallerSession) ADMINROLE() ([32]byte, error) {
	return _AgenticCommerce.Contract.ADMINROLE(&_AgenticCommerce.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgenticCommerce.Contract.DEFAULTADMINROLE(&_AgenticCommerce.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgenticCommerce.Contract.DEFAULTADMINROLE(&_AgenticCommerce.CallOpts)
}

// HOOKGASLIMIT is a free data retrieval call binding the contract method 0xff54740f.
//
// Solidity: function HOOK_GAS_LIMIT() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCaller) HOOKGASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "HOOK_GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HOOKGASLIMIT is a free data retrieval call binding the contract method 0xff54740f.
//
// Solidity: function HOOK_GAS_LIMIT() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceSession) HOOKGASLIMIT() (*big.Int, error) {
	return _AgenticCommerce.Contract.HOOKGASLIMIT(&_AgenticCommerce.CallOpts)
}

// HOOKGASLIMIT is a free data retrieval call binding the contract method 0xff54740f.
//
// Solidity: function HOOK_GAS_LIMIT() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCallerSession) HOOKGASLIMIT() (*big.Int, error) {
	return _AgenticCommerce.Contract.HOOKGASLIMIT(&_AgenticCommerce.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AgenticCommerce *AgenticCommerceCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AgenticCommerce *AgenticCommerceSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AgenticCommerce.Contract.UPGRADEINTERFACEVERSION(&_AgenticCommerce.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AgenticCommerce *AgenticCommerceCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AgenticCommerce.Contract.UPGRADEINTERFACEVERSION(&_AgenticCommerce.CallOpts)
}

// EvaluatorFeeBP is a free data retrieval call binding the contract method 0x2f0e31f4.
//
// Solidity: function evaluatorFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCaller) EvaluatorFeeBP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "evaluatorFeeBP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EvaluatorFeeBP is a free data retrieval call binding the contract method 0x2f0e31f4.
//
// Solidity: function evaluatorFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceSession) EvaluatorFeeBP() (*big.Int, error) {
	return _AgenticCommerce.Contract.EvaluatorFeeBP(&_AgenticCommerce.CallOpts)
}

// EvaluatorFeeBP is a free data retrieval call binding the contract method 0x2f0e31f4.
//
// Solidity: function evaluatorFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCallerSession) EvaluatorFeeBP() (*big.Int, error) {
	return _AgenticCommerce.Contract.EvaluatorFeeBP(&_AgenticCommerce.CallOpts)
}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns((uint256,address,address,address,string,uint256,uint256,uint8,address))
func (_AgenticCommerce *AgenticCommerceCaller) GetJob(opts *bind.CallOpts, jobId *big.Int) (AgenticCommerceUpgradeableJob, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "getJob", jobId)

	if err != nil {
		return *new(AgenticCommerceUpgradeableJob), err
	}

	out0 := *abi.ConvertType(out[0], new(AgenticCommerceUpgradeableJob)).(*AgenticCommerceUpgradeableJob)

	return out0, err

}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns((uint256,address,address,address,string,uint256,uint256,uint8,address))
func (_AgenticCommerce *AgenticCommerceSession) GetJob(jobId *big.Int) (AgenticCommerceUpgradeableJob, error) {
	return _AgenticCommerce.Contract.GetJob(&_AgenticCommerce.CallOpts, jobId)
}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns((uint256,address,address,address,string,uint256,uint256,uint8,address))
func (_AgenticCommerce *AgenticCommerceCallerSession) GetJob(jobId *big.Int) (AgenticCommerceUpgradeableJob, error) {
	return _AgenticCommerce.Contract.GetJob(&_AgenticCommerce.CallOpts, jobId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgenticCommerce.Contract.GetRoleAdmin(&_AgenticCommerce.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgenticCommerce.Contract.GetRoleAdmin(&_AgenticCommerce.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgenticCommerce *AgenticCommerceSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgenticCommerce.Contract.HasRole(&_AgenticCommerce.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgenticCommerce.Contract.HasRole(&_AgenticCommerce.CallOpts, role, account)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCaller) IsTrustedForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "isTrustedForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgenticCommerce *AgenticCommerceSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _AgenticCommerce.Contract.IsTrustedForwarder(&_AgenticCommerce.CallOpts, forwarder)
}

// IsTrustedForwarder is a free data retrieval call binding the contract method 0x572b6c05.
//
// Solidity: function isTrustedForwarder(address forwarder) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCallerSession) IsTrustedForwarder(forwarder common.Address) (bool, error) {
	return _AgenticCommerce.Contract.IsTrustedForwarder(&_AgenticCommerce.CallOpts, forwarder)
}

// JobCounter is a free data retrieval call binding the contract method 0x50355d76.
//
// Solidity: function jobCounter() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCaller) JobCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "jobCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// JobCounter is a free data retrieval call binding the contract method 0x50355d76.
//
// Solidity: function jobCounter() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceSession) JobCounter() (*big.Int, error) {
	return _AgenticCommerce.Contract.JobCounter(&_AgenticCommerce.CallOpts)
}

// JobCounter is a free data retrieval call binding the contract method 0x50355d76.
//
// Solidity: function jobCounter() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCallerSession) JobCounter() (*big.Int, error) {
	return _AgenticCommerce.Contract.JobCounter(&_AgenticCommerce.CallOpts)
}

// JobHasBudget is a free data retrieval call binding the contract method 0xfabc3329.
//
// Solidity: function jobHasBudget(uint256 jobId) view returns(bool hasBudget)
func (_AgenticCommerce *AgenticCommerceCaller) JobHasBudget(opts *bind.CallOpts, jobId *big.Int) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "jobHasBudget", jobId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// JobHasBudget is a free data retrieval call binding the contract method 0xfabc3329.
//
// Solidity: function jobHasBudget(uint256 jobId) view returns(bool hasBudget)
func (_AgenticCommerce *AgenticCommerceSession) JobHasBudget(jobId *big.Int) (bool, error) {
	return _AgenticCommerce.Contract.JobHasBudget(&_AgenticCommerce.CallOpts, jobId)
}

// JobHasBudget is a free data retrieval call binding the contract method 0xfabc3329.
//
// Solidity: function jobHasBudget(uint256 jobId) view returns(bool hasBudget)
func (_AgenticCommerce *AgenticCommerceCallerSession) JobHasBudget(jobId *big.Int) (bool, error) {
	return _AgenticCommerce.Contract.JobHasBudget(&_AgenticCommerce.CallOpts, jobId)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(uint256 id, address client, address provider, address evaluator, string description, uint256 budget, uint256 expiredAt, uint8 status, address hook)
func (_AgenticCommerce *AgenticCommerceCaller) Jobs(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "jobs", arg0)

	outstruct := new(struct {
		Id          *big.Int
		Client      common.Address
		Provider    common.Address
		Evaluator   common.Address
		Description string
		Budget      *big.Int
		ExpiredAt   *big.Int
		Status      uint8
		Hook        common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Client = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Provider = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Evaluator = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Description = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.Budget = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ExpiredAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.Hook = *abi.ConvertType(out[8], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(uint256 id, address client, address provider, address evaluator, string description, uint256 budget, uint256 expiredAt, uint8 status, address hook)
func (_AgenticCommerce *AgenticCommerceSession) Jobs(arg0 *big.Int) (struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}, error) {
	return _AgenticCommerce.Contract.Jobs(&_AgenticCommerce.CallOpts, arg0)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(uint256 id, address client, address provider, address evaluator, string description, uint256 budget, uint256 expiredAt, uint8 status, address hook)
func (_AgenticCommerce *AgenticCommerceCallerSession) Jobs(arg0 *big.Int) (struct {
	Id          *big.Int
	Client      common.Address
	Provider    common.Address
	Evaluator   common.Address
	Description string
	Budget      *big.Int
	ExpiredAt   *big.Int
	Status      uint8
	Hook        common.Address
}, error) {
	return _AgenticCommerce.Contract.Jobs(&_AgenticCommerce.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgenticCommerce *AgenticCommerceCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgenticCommerce *AgenticCommerceSession) Paused() (bool, error) {
	return _AgenticCommerce.Contract.Paused(&_AgenticCommerce.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgenticCommerce *AgenticCommerceCallerSession) Paused() (bool, error) {
	return _AgenticCommerce.Contract.Paused(&_AgenticCommerce.CallOpts)
}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_AgenticCommerce *AgenticCommerceCaller) PaymentToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "paymentToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_AgenticCommerce *AgenticCommerceSession) PaymentToken() (common.Address, error) {
	return _AgenticCommerce.Contract.PaymentToken(&_AgenticCommerce.CallOpts)
}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_AgenticCommerce *AgenticCommerceCallerSession) PaymentToken() (common.Address, error) {
	return _AgenticCommerce.Contract.PaymentToken(&_AgenticCommerce.CallOpts)
}

// PlatformFeeBP is a free data retrieval call binding the contract method 0xff96092a.
//
// Solidity: function platformFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCaller) PlatformFeeBP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "platformFeeBP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PlatformFeeBP is a free data retrieval call binding the contract method 0xff96092a.
//
// Solidity: function platformFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceSession) PlatformFeeBP() (*big.Int, error) {
	return _AgenticCommerce.Contract.PlatformFeeBP(&_AgenticCommerce.CallOpts)
}

// PlatformFeeBP is a free data retrieval call binding the contract method 0xff96092a.
//
// Solidity: function platformFeeBP() view returns(uint256)
func (_AgenticCommerce *AgenticCommerceCallerSession) PlatformFeeBP() (*big.Int, error) {
	return _AgenticCommerce.Contract.PlatformFeeBP(&_AgenticCommerce.CallOpts)
}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_AgenticCommerce *AgenticCommerceCaller) PlatformTreasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "platformTreasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_AgenticCommerce *AgenticCommerceSession) PlatformTreasury() (common.Address, error) {
	return _AgenticCommerce.Contract.PlatformTreasury(&_AgenticCommerce.CallOpts)
}

// PlatformTreasury is a free data retrieval call binding the contract method 0xe138818c.
//
// Solidity: function platformTreasury() view returns(address)
func (_AgenticCommerce *AgenticCommerceCallerSession) PlatformTreasury() (common.Address, error) {
	return _AgenticCommerce.Contract.PlatformTreasury(&_AgenticCommerce.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceSession) ProxiableUUID() ([32]byte, error) {
	return _AgenticCommerce.Contract.ProxiableUUID(&_AgenticCommerce.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AgenticCommerce *AgenticCommerceCallerSession) ProxiableUUID() ([32]byte, error) {
	return _AgenticCommerce.Contract.ProxiableUUID(&_AgenticCommerce.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgenticCommerce *AgenticCommerceSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgenticCommerce.Contract.SupportsInterface(&_AgenticCommerce.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgenticCommerce.Contract.SupportsInterface(&_AgenticCommerce.CallOpts, interfaceId)
}

// TrustedForwarder is a free data retrieval call binding the contract method 0x7da0a877.
//
// Solidity: function trustedForwarder() view returns(address)
func (_AgenticCommerce *AgenticCommerceCaller) TrustedForwarder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "trustedForwarder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedForwarder is a free data retrieval call binding the contract method 0x7da0a877.
//
// Solidity: function trustedForwarder() view returns(address)
func (_AgenticCommerce *AgenticCommerceSession) TrustedForwarder() (common.Address, error) {
	return _AgenticCommerce.Contract.TrustedForwarder(&_AgenticCommerce.CallOpts)
}

// TrustedForwarder is a free data retrieval call binding the contract method 0x7da0a877.
//
// Solidity: function trustedForwarder() view returns(address)
func (_AgenticCommerce *AgenticCommerceCallerSession) TrustedForwarder() (common.Address, error) {
	return _AgenticCommerce.Contract.TrustedForwarder(&_AgenticCommerce.CallOpts)
}

// WhitelistedHooks is a free data retrieval call binding the contract method 0x6d3b96c3.
//
// Solidity: function whitelistedHooks(address ) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCaller) WhitelistedHooks(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AgenticCommerce.contract.Call(opts, &out, "whitelistedHooks", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WhitelistedHooks is a free data retrieval call binding the contract method 0x6d3b96c3.
//
// Solidity: function whitelistedHooks(address ) view returns(bool)
func (_AgenticCommerce *AgenticCommerceSession) WhitelistedHooks(arg0 common.Address) (bool, error) {
	return _AgenticCommerce.Contract.WhitelistedHooks(&_AgenticCommerce.CallOpts, arg0)
}

// WhitelistedHooks is a free data retrieval call binding the contract method 0x6d3b96c3.
//
// Solidity: function whitelistedHooks(address ) view returns(bool)
func (_AgenticCommerce *AgenticCommerceCallerSession) WhitelistedHooks(arg0 common.Address) (bool, error) {
	return _AgenticCommerce.Contract.WhitelistedHooks(&_AgenticCommerce.CallOpts, arg0)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0x5b7baf64.
//
// Solidity: function claimRefund(uint256 jobId) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) ClaimRefund(opts *bind.TransactOpts, jobId *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "claimRefund", jobId)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0x5b7baf64.
//
// Solidity: function claimRefund(uint256 jobId) returns()
func (_AgenticCommerce *AgenticCommerceSession) ClaimRefund(jobId *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.ClaimRefund(&_AgenticCommerce.TransactOpts, jobId)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0x5b7baf64.
//
// Solidity: function claimRefund(uint256 jobId) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) ClaimRefund(jobId *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.ClaimRefund(&_AgenticCommerce.TransactOpts, jobId)
}

// Complete is a paid mutator transaction binding the contract method 0xd75bbdf3.
//
// Solidity: function complete(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Complete(opts *bind.TransactOpts, jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "complete", jobId, reason, optParams)
}

// Complete is a paid mutator transaction binding the contract method 0xd75bbdf3.
//
// Solidity: function complete(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) Complete(jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Complete(&_AgenticCommerce.TransactOpts, jobId, reason, optParams)
}

// Complete is a paid mutator transaction binding the contract method 0xd75bbdf3.
//
// Solidity: function complete(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Complete(jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Complete(&_AgenticCommerce.TransactOpts, jobId, reason, optParams)
}

// CreateJob is a paid mutator transaction binding the contract method 0x41528812.
//
// Solidity: function createJob(address provider, address evaluator, uint256 expiredAt, string description, address hook) returns(uint256)
func (_AgenticCommerce *AgenticCommerceTransactor) CreateJob(opts *bind.TransactOpts, provider common.Address, evaluator common.Address, expiredAt *big.Int, description string, hook common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "createJob", provider, evaluator, expiredAt, description, hook)
}

// CreateJob is a paid mutator transaction binding the contract method 0x41528812.
//
// Solidity: function createJob(address provider, address evaluator, uint256 expiredAt, string description, address hook) returns(uint256)
func (_AgenticCommerce *AgenticCommerceSession) CreateJob(provider common.Address, evaluator common.Address, expiredAt *big.Int, description string, hook common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.CreateJob(&_AgenticCommerce.TransactOpts, provider, evaluator, expiredAt, description, hook)
}

// CreateJob is a paid mutator transaction binding the contract method 0x41528812.
//
// Solidity: function createJob(address provider, address evaluator, uint256 expiredAt, string description, address hook) returns(uint256)
func (_AgenticCommerce *AgenticCommerceTransactorSession) CreateJob(provider common.Address, evaluator common.Address, expiredAt *big.Int, description string, hook common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.CreateJob(&_AgenticCommerce.TransactOpts, provider, evaluator, expiredAt, description, hook)
}

// Fund is a paid mutator transaction binding the contract method 0xd2e13f50.
//
// Solidity: function fund(uint256 jobId, uint256 expectedBudget, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Fund(opts *bind.TransactOpts, jobId *big.Int, expectedBudget *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "fund", jobId, expectedBudget, optParams)
}

// Fund is a paid mutator transaction binding the contract method 0xd2e13f50.
//
// Solidity: function fund(uint256 jobId, uint256 expectedBudget, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) Fund(jobId *big.Int, expectedBudget *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Fund(&_AgenticCommerce.TransactOpts, jobId, expectedBudget, optParams)
}

// Fund is a paid mutator transaction binding the contract method 0xd2e13f50.
//
// Solidity: function fund(uint256 jobId, uint256 expectedBudget, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Fund(jobId *big.Int, expectedBudget *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Fund(&_AgenticCommerce.TransactOpts, jobId, expectedBudget, optParams)
}

// FundWithPermit is a paid mutator transaction binding the contract method 0x71c720aa.
//
// Solidity: function fundWithPermit(uint256 jobId, uint256 expectedBudget, bytes optParams, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) FundWithPermit(opts *bind.TransactOpts, jobId *big.Int, expectedBudget *big.Int, optParams []byte, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "fundWithPermit", jobId, expectedBudget, optParams, deadline, v, r, s)
}

// FundWithPermit is a paid mutator transaction binding the contract method 0x71c720aa.
//
// Solidity: function fundWithPermit(uint256 jobId, uint256 expectedBudget, bytes optParams, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AgenticCommerce *AgenticCommerceSession) FundWithPermit(jobId *big.Int, expectedBudget *big.Int, optParams []byte, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.FundWithPermit(&_AgenticCommerce.TransactOpts, jobId, expectedBudget, optParams, deadline, v, r, s)
}

// FundWithPermit is a paid mutator transaction binding the contract method 0x71c720aa.
//
// Solidity: function fundWithPermit(uint256 jobId, uint256 expectedBudget, bytes optParams, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) FundWithPermit(jobId *big.Int, expectedBudget *big.Int, optParams []byte, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.FundWithPermit(&_AgenticCommerce.TransactOpts, jobId, expectedBudget, optParams, deadline, v, r, s)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.GrantRole(&_AgenticCommerce.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.GrantRole(&_AgenticCommerce.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address paymentToken_, address treasury_, address admin_) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Initialize(opts *bind.TransactOpts, paymentToken_ common.Address, treasury_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "initialize", paymentToken_, treasury_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address paymentToken_, address treasury_, address admin_) returns()
func (_AgenticCommerce *AgenticCommerceSession) Initialize(paymentToken_ common.Address, treasury_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Initialize(&_AgenticCommerce.TransactOpts, paymentToken_, treasury_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address paymentToken_, address treasury_, address admin_) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Initialize(paymentToken_ common.Address, treasury_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Initialize(&_AgenticCommerce.TransactOpts, paymentToken_, treasury_, admin_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgenticCommerce *AgenticCommerceSession) Pause() (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Pause(&_AgenticCommerce.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Pause() (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Pause(&_AgenticCommerce.TransactOpts)
}

// Reject is a paid mutator transaction binding the contract method 0x41dd26f5.
//
// Solidity: function reject(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Reject(opts *bind.TransactOpts, jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "reject", jobId, reason, optParams)
}

// Reject is a paid mutator transaction binding the contract method 0x41dd26f5.
//
// Solidity: function reject(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) Reject(jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Reject(&_AgenticCommerce.TransactOpts, jobId, reason, optParams)
}

// Reject is a paid mutator transaction binding the contract method 0x41dd26f5.
//
// Solidity: function reject(uint256 jobId, bytes32 reason, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Reject(jobId *big.Int, reason [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Reject(&_AgenticCommerce.TransactOpts, jobId, reason, optParams)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgenticCommerce *AgenticCommerceSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.RenounceRole(&_AgenticCommerce.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.RenounceRole(&_AgenticCommerce.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.RevokeRole(&_AgenticCommerce.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.RevokeRole(&_AgenticCommerce.TransactOpts, role, account)
}

// SetBudget is a paid mutator transaction binding the contract method 0xdd4ae9d4.
//
// Solidity: function setBudget(uint256 jobId, uint256 amount, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) SetBudget(opts *bind.TransactOpts, jobId *big.Int, amount *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "setBudget", jobId, amount, optParams)
}

// SetBudget is a paid mutator transaction binding the contract method 0xdd4ae9d4.
//
// Solidity: function setBudget(uint256 jobId, uint256 amount, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) SetBudget(jobId *big.Int, amount *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetBudget(&_AgenticCommerce.TransactOpts, jobId, amount, optParams)
}

// SetBudget is a paid mutator transaction binding the contract method 0xdd4ae9d4.
//
// Solidity: function setBudget(uint256 jobId, uint256 amount, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) SetBudget(jobId *big.Int, amount *big.Int, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetBudget(&_AgenticCommerce.TransactOpts, jobId, amount, optParams)
}

// SetEvaluatorFee is a paid mutator transaction binding the contract method 0x84f15090.
//
// Solidity: function setEvaluatorFee(uint256 feeBP_) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) SetEvaluatorFee(opts *bind.TransactOpts, feeBP_ *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "setEvaluatorFee", feeBP_)
}

// SetEvaluatorFee is a paid mutator transaction binding the contract method 0x84f15090.
//
// Solidity: function setEvaluatorFee(uint256 feeBP_) returns()
func (_AgenticCommerce *AgenticCommerceSession) SetEvaluatorFee(feeBP_ *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetEvaluatorFee(&_AgenticCommerce.TransactOpts, feeBP_)
}

// SetEvaluatorFee is a paid mutator transaction binding the contract method 0x84f15090.
//
// Solidity: function setEvaluatorFee(uint256 feeBP_) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) SetEvaluatorFee(feeBP_ *big.Int) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetEvaluatorFee(&_AgenticCommerce.TransactOpts, feeBP_)
}

// SetHookWhitelist is a paid mutator transaction binding the contract method 0xce79eb60.
//
// Solidity: function setHookWhitelist(address hook, bool status) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) SetHookWhitelist(opts *bind.TransactOpts, hook common.Address, status bool) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "setHookWhitelist", hook, status)
}

// SetHookWhitelist is a paid mutator transaction binding the contract method 0xce79eb60.
//
// Solidity: function setHookWhitelist(address hook, bool status) returns()
func (_AgenticCommerce *AgenticCommerceSession) SetHookWhitelist(hook common.Address, status bool) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetHookWhitelist(&_AgenticCommerce.TransactOpts, hook, status)
}

// SetHookWhitelist is a paid mutator transaction binding the contract method 0xce79eb60.
//
// Solidity: function setHookWhitelist(address hook, bool status) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) SetHookWhitelist(hook common.Address, status bool) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetHookWhitelist(&_AgenticCommerce.TransactOpts, hook, status)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0xb4d884f6.
//
// Solidity: function setPlatformFee(uint256 feeBP_, address treasury_) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) SetPlatformFee(opts *bind.TransactOpts, feeBP_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "setPlatformFee", feeBP_, treasury_)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0xb4d884f6.
//
// Solidity: function setPlatformFee(uint256 feeBP_, address treasury_) returns()
func (_AgenticCommerce *AgenticCommerceSession) SetPlatformFee(feeBP_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetPlatformFee(&_AgenticCommerce.TransactOpts, feeBP_, treasury_)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0xb4d884f6.
//
// Solidity: function setPlatformFee(uint256 feeBP_, address treasury_) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) SetPlatformFee(feeBP_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetPlatformFee(&_AgenticCommerce.TransactOpts, feeBP_, treasury_)
}

// SetProvider is a paid mutator transaction binding the contract method 0xc9a84bb9.
//
// Solidity: function setProvider(uint256 jobId, address provider_, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) SetProvider(opts *bind.TransactOpts, jobId *big.Int, provider_ common.Address, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "setProvider", jobId, provider_, optParams)
}

// SetProvider is a paid mutator transaction binding the contract method 0xc9a84bb9.
//
// Solidity: function setProvider(uint256 jobId, address provider_, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) SetProvider(jobId *big.Int, provider_ common.Address, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetProvider(&_AgenticCommerce.TransactOpts, jobId, provider_, optParams)
}

// SetProvider is a paid mutator transaction binding the contract method 0xc9a84bb9.
//
// Solidity: function setProvider(uint256 jobId, address provider_, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) SetProvider(jobId *big.Int, provider_ common.Address, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.SetProvider(&_AgenticCommerce.TransactOpts, jobId, provider_, optParams)
}

// Submit is a paid mutator transaction binding the contract method 0x9e63798d.
//
// Solidity: function submit(uint256 jobId, bytes32 deliverable, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Submit(opts *bind.TransactOpts, jobId *big.Int, deliverable [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "submit", jobId, deliverable, optParams)
}

// Submit is a paid mutator transaction binding the contract method 0x9e63798d.
//
// Solidity: function submit(uint256 jobId, bytes32 deliverable, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceSession) Submit(jobId *big.Int, deliverable [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Submit(&_AgenticCommerce.TransactOpts, jobId, deliverable, optParams)
}

// Submit is a paid mutator transaction binding the contract method 0x9e63798d.
//
// Solidity: function submit(uint256 jobId, bytes32 deliverable, bytes optParams) returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Submit(jobId *big.Int, deliverable [32]byte, optParams []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Submit(&_AgenticCommerce.TransactOpts, jobId, deliverable, optParams)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgenticCommerce *AgenticCommerceTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgenticCommerce *AgenticCommerceSession) Unpause() (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Unpause(&_AgenticCommerce.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) Unpause() (*types.Transaction, error) {
	return _AgenticCommerce.Contract.Unpause(&_AgenticCommerce.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgenticCommerce *AgenticCommerceTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgenticCommerce.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgenticCommerce *AgenticCommerceSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.UpgradeToAndCall(&_AgenticCommerce.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AgenticCommerce *AgenticCommerceTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AgenticCommerce.Contract.UpgradeToAndCall(&_AgenticCommerce.TransactOpts, newImplementation, data)
}

// AgenticCommerceBudgetSetIterator is returned from FilterBudgetSet and is used to iterate over the raw logs and unpacked data for BudgetSet events raised by the AgenticCommerce contract.
type AgenticCommerceBudgetSetIterator struct {
	Event *AgenticCommerceBudgetSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceBudgetSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceBudgetSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceBudgetSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceBudgetSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceBudgetSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceBudgetSet represents a BudgetSet event raised by the AgenticCommerce contract.
type AgenticCommerceBudgetSet struct {
	JobId  *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBudgetSet is a free log retrieval operation binding the contract event 0x869e2577b006bf47ee981cf6fec2e25583548081c14b98deab587f77b5068038.
//
// Solidity: event BudgetSet(uint256 indexed jobId, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterBudgetSet(opts *bind.FilterOpts, jobId []*big.Int) (*AgenticCommerceBudgetSetIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "BudgetSet", jobIdRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceBudgetSetIterator{contract: _AgenticCommerce.contract, event: "BudgetSet", logs: logs, sub: sub}, nil
}

// WatchBudgetSet is a free log subscription operation binding the contract event 0x869e2577b006bf47ee981cf6fec2e25583548081c14b98deab587f77b5068038.
//
// Solidity: event BudgetSet(uint256 indexed jobId, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchBudgetSet(opts *bind.WatchOpts, sink chan<- *AgenticCommerceBudgetSet, jobId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "BudgetSet", jobIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceBudgetSet)
				if err := _AgenticCommerce.contract.UnpackLog(event, "BudgetSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBudgetSet is a log parse operation binding the contract event 0x869e2577b006bf47ee981cf6fec2e25583548081c14b98deab587f77b5068038.
//
// Solidity: event BudgetSet(uint256 indexed jobId, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseBudgetSet(log types.Log) (*AgenticCommerceBudgetSet, error) {
	event := new(AgenticCommerceBudgetSet)
	if err := _AgenticCommerce.contract.UnpackLog(event, "BudgetSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceEvaluatorFeePaidIterator is returned from FilterEvaluatorFeePaid and is used to iterate over the raw logs and unpacked data for EvaluatorFeePaid events raised by the AgenticCommerce contract.
type AgenticCommerceEvaluatorFeePaidIterator struct {
	Event *AgenticCommerceEvaluatorFeePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceEvaluatorFeePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceEvaluatorFeePaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceEvaluatorFeePaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceEvaluatorFeePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceEvaluatorFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceEvaluatorFeePaid represents a EvaluatorFeePaid event raised by the AgenticCommerce contract.
type AgenticCommerceEvaluatorFeePaid struct {
	JobId     *big.Int
	Evaluator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEvaluatorFeePaid is a free log retrieval operation binding the contract event 0x253dd534010ac976fa263caa123bae79b9c50292adf7ce67bdc5ec309f784e61.
//
// Solidity: event EvaluatorFeePaid(uint256 indexed jobId, address indexed evaluator, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterEvaluatorFeePaid(opts *bind.FilterOpts, jobId []*big.Int, evaluator []common.Address) (*AgenticCommerceEvaluatorFeePaidIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "EvaluatorFeePaid", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceEvaluatorFeePaidIterator{contract: _AgenticCommerce.contract, event: "EvaluatorFeePaid", logs: logs, sub: sub}, nil
}

// WatchEvaluatorFeePaid is a free log subscription operation binding the contract event 0x253dd534010ac976fa263caa123bae79b9c50292adf7ce67bdc5ec309f784e61.
//
// Solidity: event EvaluatorFeePaid(uint256 indexed jobId, address indexed evaluator, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchEvaluatorFeePaid(opts *bind.WatchOpts, sink chan<- *AgenticCommerceEvaluatorFeePaid, jobId []*big.Int, evaluator []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "EvaluatorFeePaid", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceEvaluatorFeePaid)
				if err := _AgenticCommerce.contract.UnpackLog(event, "EvaluatorFeePaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEvaluatorFeePaid is a log parse operation binding the contract event 0x253dd534010ac976fa263caa123bae79b9c50292adf7ce67bdc5ec309f784e61.
//
// Solidity: event EvaluatorFeePaid(uint256 indexed jobId, address indexed evaluator, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseEvaluatorFeePaid(log types.Log) (*AgenticCommerceEvaluatorFeePaid, error) {
	event := new(AgenticCommerceEvaluatorFeePaid)
	if err := _AgenticCommerce.contract.UnpackLog(event, "EvaluatorFeePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceHookWhitelistUpdatedIterator is returned from FilterHookWhitelistUpdated and is used to iterate over the raw logs and unpacked data for HookWhitelistUpdated events raised by the AgenticCommerce contract.
type AgenticCommerceHookWhitelistUpdatedIterator struct {
	Event *AgenticCommerceHookWhitelistUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceHookWhitelistUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceHookWhitelistUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceHookWhitelistUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceHookWhitelistUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceHookWhitelistUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceHookWhitelistUpdated represents a HookWhitelistUpdated event raised by the AgenticCommerce contract.
type AgenticCommerceHookWhitelistUpdated struct {
	Hook   common.Address
	Status bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterHookWhitelistUpdated is a free log retrieval operation binding the contract event 0x7ee54953080e392a475a25b6acacb85417ca4e1953293c90934233ca13612510.
//
// Solidity: event HookWhitelistUpdated(address indexed hook, bool status)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterHookWhitelistUpdated(opts *bind.FilterOpts, hook []common.Address) (*AgenticCommerceHookWhitelistUpdatedIterator, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "HookWhitelistUpdated", hookRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceHookWhitelistUpdatedIterator{contract: _AgenticCommerce.contract, event: "HookWhitelistUpdated", logs: logs, sub: sub}, nil
}

// WatchHookWhitelistUpdated is a free log subscription operation binding the contract event 0x7ee54953080e392a475a25b6acacb85417ca4e1953293c90934233ca13612510.
//
// Solidity: event HookWhitelistUpdated(address indexed hook, bool status)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchHookWhitelistUpdated(opts *bind.WatchOpts, sink chan<- *AgenticCommerceHookWhitelistUpdated, hook []common.Address) (event.Subscription, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "HookWhitelistUpdated", hookRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceHookWhitelistUpdated)
				if err := _AgenticCommerce.contract.UnpackLog(event, "HookWhitelistUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHookWhitelistUpdated is a log parse operation binding the contract event 0x7ee54953080e392a475a25b6acacb85417ca4e1953293c90934233ca13612510.
//
// Solidity: event HookWhitelistUpdated(address indexed hook, bool status)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseHookWhitelistUpdated(log types.Log) (*AgenticCommerceHookWhitelistUpdated, error) {
	event := new(AgenticCommerceHookWhitelistUpdated)
	if err := _AgenticCommerce.contract.UnpackLog(event, "HookWhitelistUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AgenticCommerce contract.
type AgenticCommerceInitializedIterator struct {
	Event *AgenticCommerceInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceInitialized represents a Initialized event raised by the AgenticCommerce contract.
type AgenticCommerceInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterInitialized(opts *bind.FilterOpts) (*AgenticCommerceInitializedIterator, error) {

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceInitializedIterator{contract: _AgenticCommerce.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AgenticCommerceInitialized) (event.Subscription, error) {

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceInitialized)
				if err := _AgenticCommerce.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseInitialized(log types.Log) (*AgenticCommerceInitialized, error) {
	event := new(AgenticCommerceInitialized)
	if err := _AgenticCommerce.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobCompletedIterator is returned from FilterJobCompleted and is used to iterate over the raw logs and unpacked data for JobCompleted events raised by the AgenticCommerce contract.
type AgenticCommerceJobCompletedIterator struct {
	Event *AgenticCommerceJobCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobCompleted represents a JobCompleted event raised by the AgenticCommerce contract.
type AgenticCommerceJobCompleted struct {
	JobId     *big.Int
	Evaluator common.Address
	Reason    [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterJobCompleted is a free log retrieval operation binding the contract event 0x0fd54bd364fa9e67f17b091aefe930932c09fe7651cf5ad02c71a418f3341444.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed evaluator, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobCompleted(opts *bind.FilterOpts, jobId []*big.Int, evaluator []common.Address) (*AgenticCommerceJobCompletedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobCompleted", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobCompletedIterator{contract: _AgenticCommerce.contract, event: "JobCompleted", logs: logs, sub: sub}, nil
}

// WatchJobCompleted is a free log subscription operation binding the contract event 0x0fd54bd364fa9e67f17b091aefe930932c09fe7651cf5ad02c71a418f3341444.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed evaluator, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobCompleted(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobCompleted, jobId []*big.Int, evaluator []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobCompleted", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobCompleted)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobCompleted is a log parse operation binding the contract event 0x0fd54bd364fa9e67f17b091aefe930932c09fe7651cf5ad02c71a418f3341444.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed evaluator, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobCompleted(log types.Log) (*AgenticCommerceJobCompleted, error) {
	event := new(AgenticCommerceJobCompleted)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobCreatedIterator is returned from FilterJobCreated and is used to iterate over the raw logs and unpacked data for JobCreated events raised by the AgenticCommerce contract.
type AgenticCommerceJobCreatedIterator struct {
	Event *AgenticCommerceJobCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobCreated represents a JobCreated event raised by the AgenticCommerce contract.
type AgenticCommerceJobCreated struct {
	JobId     *big.Int
	Client    common.Address
	Provider  common.Address
	Evaluator common.Address
	ExpiredAt *big.Int
	Hook      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterJobCreated is a free log retrieval operation binding the contract event 0xb0f0239bfdd96453e24733e18bfc24b70d8fadf123dd977473518dd577ee79b9.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed client, address indexed provider, address evaluator, uint256 expiredAt, address hook)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobCreated(opts *bind.FilterOpts, jobId []*big.Int, client []common.Address, provider []common.Address) (*AgenticCommerceJobCreatedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobCreated", jobIdRule, clientRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobCreatedIterator{contract: _AgenticCommerce.contract, event: "JobCreated", logs: logs, sub: sub}, nil
}

// WatchJobCreated is a free log subscription operation binding the contract event 0xb0f0239bfdd96453e24733e18bfc24b70d8fadf123dd977473518dd577ee79b9.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed client, address indexed provider, address evaluator, uint256 expiredAt, address hook)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobCreated(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobCreated, jobId []*big.Int, client []common.Address, provider []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobCreated", jobIdRule, clientRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobCreated)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobCreated is a log parse operation binding the contract event 0xb0f0239bfdd96453e24733e18bfc24b70d8fadf123dd977473518dd577ee79b9.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed client, address indexed provider, address evaluator, uint256 expiredAt, address hook)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobCreated(log types.Log) (*AgenticCommerceJobCreated, error) {
	event := new(AgenticCommerceJobCreated)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobExpiredIterator is returned from FilterJobExpired and is used to iterate over the raw logs and unpacked data for JobExpired events raised by the AgenticCommerce contract.
type AgenticCommerceJobExpiredIterator struct {
	Event *AgenticCommerceJobExpired // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobExpiredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobExpired)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobExpired)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobExpiredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobExpiredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobExpired represents a JobExpired event raised by the AgenticCommerce contract.
type AgenticCommerceJobExpired struct {
	JobId *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterJobExpired is a free log retrieval operation binding the contract event 0x97237956f8810192811e2c3f273fd02c5d6295206fdd9c62e6fe2bfc19ba9232.
//
// Solidity: event JobExpired(uint256 indexed jobId)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobExpired(opts *bind.FilterOpts, jobId []*big.Int) (*AgenticCommerceJobExpiredIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobExpired", jobIdRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobExpiredIterator{contract: _AgenticCommerce.contract, event: "JobExpired", logs: logs, sub: sub}, nil
}

// WatchJobExpired is a free log subscription operation binding the contract event 0x97237956f8810192811e2c3f273fd02c5d6295206fdd9c62e6fe2bfc19ba9232.
//
// Solidity: event JobExpired(uint256 indexed jobId)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobExpired(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobExpired, jobId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobExpired", jobIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobExpired)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobExpired", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobExpired is a log parse operation binding the contract event 0x97237956f8810192811e2c3f273fd02c5d6295206fdd9c62e6fe2bfc19ba9232.
//
// Solidity: event JobExpired(uint256 indexed jobId)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobExpired(log types.Log) (*AgenticCommerceJobExpired, error) {
	event := new(AgenticCommerceJobExpired)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobExpired", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobFundedIterator is returned from FilterJobFunded and is used to iterate over the raw logs and unpacked data for JobFunded events raised by the AgenticCommerce contract.
type AgenticCommerceJobFundedIterator struct {
	Event *AgenticCommerceJobFunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobFunded represents a JobFunded event raised by the AgenticCommerce contract.
type AgenticCommerceJobFunded struct {
	JobId  *big.Int
	Client common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterJobFunded is a free log retrieval operation binding the contract event 0xe3fbcc1ea1bdc559ec7f0347efde7655e58b5f45a30b0e4470a583c3ef5496b3.
//
// Solidity: event JobFunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobFunded(opts *bind.FilterOpts, jobId []*big.Int, client []common.Address) (*AgenticCommerceJobFundedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobFunded", jobIdRule, clientRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobFundedIterator{contract: _AgenticCommerce.contract, event: "JobFunded", logs: logs, sub: sub}, nil
}

// WatchJobFunded is a free log subscription operation binding the contract event 0xe3fbcc1ea1bdc559ec7f0347efde7655e58b5f45a30b0e4470a583c3ef5496b3.
//
// Solidity: event JobFunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobFunded(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobFunded, jobId []*big.Int, client []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobFunded", jobIdRule, clientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobFunded)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobFunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobFunded is a log parse operation binding the contract event 0xe3fbcc1ea1bdc559ec7f0347efde7655e58b5f45a30b0e4470a583c3ef5496b3.
//
// Solidity: event JobFunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobFunded(log types.Log) (*AgenticCommerceJobFunded, error) {
	event := new(AgenticCommerceJobFunded)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobRejectedIterator is returned from FilterJobRejected and is used to iterate over the raw logs and unpacked data for JobRejected events raised by the AgenticCommerce contract.
type AgenticCommerceJobRejectedIterator struct {
	Event *AgenticCommerceJobRejected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobRejected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobRejected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobRejected represents a JobRejected event raised by the AgenticCommerce contract.
type AgenticCommerceJobRejected struct {
	JobId    *big.Int
	Rejector common.Address
	Reason   [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterJobRejected is a free log retrieval operation binding the contract event 0xae7362b1af91f4492868987b9c73990d780060811551b58728fbe96fd1bab275.
//
// Solidity: event JobRejected(uint256 indexed jobId, address indexed rejector, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobRejected(opts *bind.FilterOpts, jobId []*big.Int, rejector []common.Address) (*AgenticCommerceJobRejectedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var rejectorRule []interface{}
	for _, rejectorItem := range rejector {
		rejectorRule = append(rejectorRule, rejectorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobRejected", jobIdRule, rejectorRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobRejectedIterator{contract: _AgenticCommerce.contract, event: "JobRejected", logs: logs, sub: sub}, nil
}

// WatchJobRejected is a free log subscription operation binding the contract event 0xae7362b1af91f4492868987b9c73990d780060811551b58728fbe96fd1bab275.
//
// Solidity: event JobRejected(uint256 indexed jobId, address indexed rejector, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobRejected(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobRejected, jobId []*big.Int, rejector []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var rejectorRule []interface{}
	for _, rejectorItem := range rejector {
		rejectorRule = append(rejectorRule, rejectorItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobRejected", jobIdRule, rejectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobRejected)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobRejected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobRejected is a log parse operation binding the contract event 0xae7362b1af91f4492868987b9c73990d780060811551b58728fbe96fd1bab275.
//
// Solidity: event JobRejected(uint256 indexed jobId, address indexed rejector, bytes32 reason)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobRejected(log types.Log) (*AgenticCommerceJobRejected, error) {
	event := new(AgenticCommerceJobRejected)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceJobSubmittedIterator is returned from FilterJobSubmitted and is used to iterate over the raw logs and unpacked data for JobSubmitted events raised by the AgenticCommerce contract.
type AgenticCommerceJobSubmittedIterator struct {
	Event *AgenticCommerceJobSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceJobSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceJobSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceJobSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceJobSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceJobSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceJobSubmitted represents a JobSubmitted event raised by the AgenticCommerce contract.
type AgenticCommerceJobSubmitted struct {
	JobId       *big.Int
	Provider    common.Address
	Deliverable [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterJobSubmitted is a free log retrieval operation binding the contract event 0x80c17db79857f338a6a6df68a6883ecc0ce78e2202fe61ed979733573f40538e.
//
// Solidity: event JobSubmitted(uint256 indexed jobId, address indexed provider, bytes32 deliverable)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterJobSubmitted(opts *bind.FilterOpts, jobId []*big.Int, provider []common.Address) (*AgenticCommerceJobSubmittedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "JobSubmitted", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceJobSubmittedIterator{contract: _AgenticCommerce.contract, event: "JobSubmitted", logs: logs, sub: sub}, nil
}

// WatchJobSubmitted is a free log subscription operation binding the contract event 0x80c17db79857f338a6a6df68a6883ecc0ce78e2202fe61ed979733573f40538e.
//
// Solidity: event JobSubmitted(uint256 indexed jobId, address indexed provider, bytes32 deliverable)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchJobSubmitted(opts *bind.WatchOpts, sink chan<- *AgenticCommerceJobSubmitted, jobId []*big.Int, provider []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "JobSubmitted", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceJobSubmitted)
				if err := _AgenticCommerce.contract.UnpackLog(event, "JobSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobSubmitted is a log parse operation binding the contract event 0x80c17db79857f338a6a6df68a6883ecc0ce78e2202fe61ed979733573f40538e.
//
// Solidity: event JobSubmitted(uint256 indexed jobId, address indexed provider, bytes32 deliverable)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseJobSubmitted(log types.Log) (*AgenticCommerceJobSubmitted, error) {
	event := new(AgenticCommerceJobSubmitted)
	if err := _AgenticCommerce.contract.UnpackLog(event, "JobSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommercePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the AgenticCommerce contract.
type AgenticCommercePausedIterator struct {
	Event *AgenticCommercePaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommercePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommercePaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommercePaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommercePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommercePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommercePaused represents a Paused event raised by the AgenticCommerce contract.
type AgenticCommercePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterPaused(opts *bind.FilterOpts) (*AgenticCommercePausedIterator, error) {

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AgenticCommercePausedIterator{contract: _AgenticCommerce.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AgenticCommercePaused) (event.Subscription, error) {

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommercePaused)
				if err := _AgenticCommerce.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) ParsePaused(log types.Log) (*AgenticCommercePaused, error) {
	event := new(AgenticCommercePaused)
	if err := _AgenticCommerce.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommercePaymentReleasedIterator is returned from FilterPaymentReleased and is used to iterate over the raw logs and unpacked data for PaymentReleased events raised by the AgenticCommerce contract.
type AgenticCommercePaymentReleasedIterator struct {
	Event *AgenticCommercePaymentReleased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommercePaymentReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommercePaymentReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommercePaymentReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommercePaymentReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommercePaymentReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommercePaymentReleased represents a PaymentReleased event raised by the AgenticCommerce contract.
type AgenticCommercePaymentReleased struct {
	JobId    *big.Int
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPaymentReleased is a free log retrieval operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 indexed jobId, address indexed provider, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterPaymentReleased(opts *bind.FilterOpts, jobId []*big.Int, provider []common.Address) (*AgenticCommercePaymentReleasedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "PaymentReleased", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommercePaymentReleasedIterator{contract: _AgenticCommerce.contract, event: "PaymentReleased", logs: logs, sub: sub}, nil
}

// WatchPaymentReleased is a free log subscription operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 indexed jobId, address indexed provider, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchPaymentReleased(opts *bind.WatchOpts, sink chan<- *AgenticCommercePaymentReleased, jobId []*big.Int, provider []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "PaymentReleased", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommercePaymentReleased)
				if err := _AgenticCommerce.contract.UnpackLog(event, "PaymentReleased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentReleased is a log parse operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 indexed jobId, address indexed provider, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) ParsePaymentReleased(log types.Log) (*AgenticCommercePaymentReleased, error) {
	event := new(AgenticCommercePaymentReleased)
	if err := _AgenticCommerce.contract.UnpackLog(event, "PaymentReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceProviderSetIterator is returned from FilterProviderSet and is used to iterate over the raw logs and unpacked data for ProviderSet events raised by the AgenticCommerce contract.
type AgenticCommerceProviderSetIterator struct {
	Event *AgenticCommerceProviderSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceProviderSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceProviderSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceProviderSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceProviderSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceProviderSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceProviderSet represents a ProviderSet event raised by the AgenticCommerce contract.
type AgenticCommerceProviderSet struct {
	JobId    *big.Int
	Provider common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProviderSet is a free log retrieval operation binding the contract event 0x9a87df076ea1725aba8ba29d32517ce37c9597d88cbf16ec6707892cc330ab69.
//
// Solidity: event ProviderSet(uint256 indexed jobId, address indexed provider)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterProviderSet(opts *bind.FilterOpts, jobId []*big.Int, provider []common.Address) (*AgenticCommerceProviderSetIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "ProviderSet", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceProviderSetIterator{contract: _AgenticCommerce.contract, event: "ProviderSet", logs: logs, sub: sub}, nil
}

// WatchProviderSet is a free log subscription operation binding the contract event 0x9a87df076ea1725aba8ba29d32517ce37c9597d88cbf16ec6707892cc330ab69.
//
// Solidity: event ProviderSet(uint256 indexed jobId, address indexed provider)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchProviderSet(opts *bind.WatchOpts, sink chan<- *AgenticCommerceProviderSet, jobId []*big.Int, provider []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "ProviderSet", jobIdRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceProviderSet)
				if err := _AgenticCommerce.contract.UnpackLog(event, "ProviderSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProviderSet is a log parse operation binding the contract event 0x9a87df076ea1725aba8ba29d32517ce37c9597d88cbf16ec6707892cc330ab69.
//
// Solidity: event ProviderSet(uint256 indexed jobId, address indexed provider)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseProviderSet(log types.Log) (*AgenticCommerceProviderSet, error) {
	event := new(AgenticCommerceProviderSet)
	if err := _AgenticCommerce.contract.UnpackLog(event, "ProviderSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the AgenticCommerce contract.
type AgenticCommerceRefundedIterator struct {
	Event *AgenticCommerceRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceRefunded represents a Refunded event raised by the AgenticCommerce contract.
type AgenticCommerceRefunded struct {
	JobId  *big.Int
	Client common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterRefunded(opts *bind.FilterOpts, jobId []*big.Int, client []common.Address) (*AgenticCommerceRefundedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "Refunded", jobIdRule, clientRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceRefundedIterator{contract: _AgenticCommerce.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AgenticCommerceRefunded, jobId []*big.Int, client []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "Refunded", jobIdRule, clientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceRefunded)
				if err := _AgenticCommerce.contract.UnpackLog(event, "Refunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefunded is a log parse operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed jobId, address indexed client, uint256 amount)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseRefunded(log types.Log) (*AgenticCommerceRefunded, error) {
	event := new(AgenticCommerceRefunded)
	if err := _AgenticCommerce.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceReputationSignalIterator is returned from FilterReputationSignal and is used to iterate over the raw logs and unpacked data for ReputationSignal events raised by the AgenticCommerce contract.
type AgenticCommerceReputationSignalIterator struct {
	Event *AgenticCommerceReputationSignal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceReputationSignalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceReputationSignal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceReputationSignal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceReputationSignalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceReputationSignalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceReputationSignal represents a ReputationSignal event raised by the AgenticCommerce contract.
type AgenticCommerceReputationSignal struct {
	JobId   *big.Int
	Subject common.Address
	Role    string
	Signal  int8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReputationSignal is a free log retrieval operation binding the contract event 0xc407f1b818177970f0dcd9a17c8a2068a2a07f449bf34cfcee4f2667db8ea4ab.
//
// Solidity: event ReputationSignal(uint256 indexed jobId, address indexed subject, string role, int8 signal)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterReputationSignal(opts *bind.FilterOpts, jobId []*big.Int, subject []common.Address) (*AgenticCommerceReputationSignalIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "ReputationSignal", jobIdRule, subjectRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceReputationSignalIterator{contract: _AgenticCommerce.contract, event: "ReputationSignal", logs: logs, sub: sub}, nil
}

// WatchReputationSignal is a free log subscription operation binding the contract event 0xc407f1b818177970f0dcd9a17c8a2068a2a07f449bf34cfcee4f2667db8ea4ab.
//
// Solidity: event ReputationSignal(uint256 indexed jobId, address indexed subject, string role, int8 signal)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchReputationSignal(opts *bind.WatchOpts, sink chan<- *AgenticCommerceReputationSignal, jobId []*big.Int, subject []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var subjectRule []interface{}
	for _, subjectItem := range subject {
		subjectRule = append(subjectRule, subjectItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "ReputationSignal", jobIdRule, subjectRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceReputationSignal)
				if err := _AgenticCommerce.contract.UnpackLog(event, "ReputationSignal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReputationSignal is a log parse operation binding the contract event 0xc407f1b818177970f0dcd9a17c8a2068a2a07f449bf34cfcee4f2667db8ea4ab.
//
// Solidity: event ReputationSignal(uint256 indexed jobId, address indexed subject, string role, int8 signal)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseReputationSignal(log types.Log) (*AgenticCommerceReputationSignal, error) {
	event := new(AgenticCommerceReputationSignal)
	if err := _AgenticCommerce.contract.UnpackLog(event, "ReputationSignal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the AgenticCommerce contract.
type AgenticCommerceRoleAdminChangedIterator struct {
	Event *AgenticCommerceRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceRoleAdminChanged represents a RoleAdminChanged event raised by the AgenticCommerce contract.
type AgenticCommerceRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AgenticCommerceRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceRoleAdminChangedIterator{contract: _AgenticCommerce.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AgenticCommerceRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceRoleAdminChanged)
				if err := _AgenticCommerce.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseRoleAdminChanged(log types.Log) (*AgenticCommerceRoleAdminChanged, error) {
	event := new(AgenticCommerceRoleAdminChanged)
	if err := _AgenticCommerce.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the AgenticCommerce contract.
type AgenticCommerceRoleGrantedIterator struct {
	Event *AgenticCommerceRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceRoleGranted represents a RoleGranted event raised by the AgenticCommerce contract.
type AgenticCommerceRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgenticCommerceRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceRoleGrantedIterator{contract: _AgenticCommerce.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AgenticCommerceRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceRoleGranted)
				if err := _AgenticCommerce.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseRoleGranted(log types.Log) (*AgenticCommerceRoleGranted, error) {
	event := new(AgenticCommerceRoleGranted)
	if err := _AgenticCommerce.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the AgenticCommerce contract.
type AgenticCommerceRoleRevokedIterator struct {
	Event *AgenticCommerceRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceRoleRevoked represents a RoleRevoked event raised by the AgenticCommerce contract.
type AgenticCommerceRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgenticCommerceRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceRoleRevokedIterator{contract: _AgenticCommerce.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AgenticCommerceRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceRoleRevoked)
				if err := _AgenticCommerce.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseRoleRevoked(log types.Log) (*AgenticCommerceRoleRevoked, error) {
	event := new(AgenticCommerceRoleRevoked)
	if err := _AgenticCommerce.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the AgenticCommerce contract.
type AgenticCommerceUnpausedIterator struct {
	Event *AgenticCommerceUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceUnpaused represents a Unpaused event raised by the AgenticCommerce contract.
type AgenticCommerceUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AgenticCommerceUnpausedIterator, error) {

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceUnpausedIterator{contract: _AgenticCommerce.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AgenticCommerceUnpaused) (event.Subscription, error) {

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceUnpaused)
				if err := _AgenticCommerce.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseUnpaused(log types.Log) (*AgenticCommerceUnpaused, error) {
	event := new(AgenticCommerceUnpaused)
	if err := _AgenticCommerce.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgenticCommerceUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AgenticCommerce contract.
type AgenticCommerceUpgradedIterator struct {
	Event *AgenticCommerceUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgenticCommerceUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgenticCommerceUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgenticCommerceUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgenticCommerceUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgenticCommerceUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgenticCommerceUpgraded represents a Upgraded event raised by the AgenticCommerce contract.
type AgenticCommerceUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AgenticCommerce *AgenticCommerceFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AgenticCommerceUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AgenticCommerce.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AgenticCommerceUpgradedIterator{contract: _AgenticCommerce.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AgenticCommerce *AgenticCommerceFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AgenticCommerceUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AgenticCommerce.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgenticCommerceUpgraded)
				if err := _AgenticCommerce.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AgenticCommerce *AgenticCommerceFilterer) ParseUpgraded(log types.Log) (*AgenticCommerceUpgraded, error) {
	event := new(AgenticCommerceUpgraded)
	if err := _AgenticCommerce.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
