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

// IValidationRegistryRequest is an auto generated low-level Go binding around an user-defined struct.
type IValidationRegistryRequest struct {
	AgentValidatorId *big.Int
	AgentServerId    *big.Int
	DataHash         [32]byte
	Timestamp        *big.Int
	Responded        bool
}

// ValidationRegistryMetaData contains all meta data concerning the ValidationRegistry contract.
var ValidationRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AgentNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataHash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidResponse\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RequestExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SelfValidationNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedValidator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValidationAlreadyResponded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValidationRequestNotFound\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentValidatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"ValidationRequestEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentValidatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"response\",\"type\":\"uint8\"}],\"name\":\"ValidationResponseEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EXPIRATION_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_identityRegistry\",\"outputs\":[{\"internalType\":\"contractIIdentityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpirationSlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"slots\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getValidationRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"agentValidatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"responded\",\"type\":\"bool\"}],\"internalType\":\"structIValidationRegistry.Request\",\"name\":\"request\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getValidationResponse\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"hasResponse\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"response\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIIdentityRegistry\",\"name\":\"identityRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"isValidationPending\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"pending\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentValidatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"validationRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"response\",\"type\":\"uint8\"}],\"name\":\"validationResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ValidationRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidationRegistryMetaData.ABI instead.
var ValidationRegistryABI = ValidationRegistryMetaData.ABI

// ValidationRegistry is an auto generated Go binding around an Ethereum contract.
type ValidationRegistry struct {
	ValidationRegistryCaller     // Read-only binding to the contract
	ValidationRegistryTransactor // Write-only binding to the contract
	ValidationRegistryFilterer   // Log filterer for contract events
}

// ValidationRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidationRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidationRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidationRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidationRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidationRegistrySession struct {
	Contract     *ValidationRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ValidationRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidationRegistryCallerSession struct {
	Contract *ValidationRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ValidationRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidationRegistryTransactorSession struct {
	Contract     *ValidationRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ValidationRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidationRegistryRaw struct {
	Contract *ValidationRegistry // Generic contract binding to access the raw methods on
}

// ValidationRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidationRegistryCallerRaw struct {
	Contract *ValidationRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ValidationRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidationRegistryTransactorRaw struct {
	Contract *ValidationRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidationRegistry creates a new instance of ValidationRegistry, bound to a specific deployed contract.
func NewValidationRegistry(address common.Address, backend bind.ContractBackend) (*ValidationRegistry, error) {
	contract, err := bindValidationRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistry{ValidationRegistryCaller: ValidationRegistryCaller{contract: contract}, ValidationRegistryTransactor: ValidationRegistryTransactor{contract: contract}, ValidationRegistryFilterer: ValidationRegistryFilterer{contract: contract}}, nil
}

// NewValidationRegistryCaller creates a new read-only instance of ValidationRegistry, bound to a specific deployed contract.
func NewValidationRegistryCaller(address common.Address, caller bind.ContractCaller) (*ValidationRegistryCaller, error) {
	contract, err := bindValidationRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryCaller{contract: contract}, nil
}

// NewValidationRegistryTransactor creates a new write-only instance of ValidationRegistry, bound to a specific deployed contract.
func NewValidationRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidationRegistryTransactor, error) {
	contract, err := bindValidationRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryTransactor{contract: contract}, nil
}

// NewValidationRegistryFilterer creates a new log filterer instance of ValidationRegistry, bound to a specific deployed contract.
func NewValidationRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidationRegistryFilterer, error) {
	contract, err := bindValidationRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryFilterer{contract: contract}, nil
}

// bindValidationRegistry binds a generic wrapper to an already deployed contract.
func bindValidationRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidationRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidationRegistry *ValidationRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidationRegistry.Contract.ValidationRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidationRegistry *ValidationRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidationRegistry *ValidationRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidationRegistry *ValidationRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidationRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidationRegistry *ValidationRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidationRegistry *ValidationRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.contract.Transact(opts, method, params...)
}

// EXPIRATIONTIME is a free data retrieval call binding the contract method 0x4a5c7348.
//
// Solidity: function EXPIRATION_TIME() view returns(uint256)
func (_ValidationRegistry *ValidationRegistryCaller) EXPIRATIONTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "EXPIRATION_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EXPIRATIONTIME is a free data retrieval call binding the contract method 0x4a5c7348.
//
// Solidity: function EXPIRATION_TIME() view returns(uint256)
func (_ValidationRegistry *ValidationRegistrySession) EXPIRATIONTIME() (*big.Int, error) {
	return _ValidationRegistry.Contract.EXPIRATIONTIME(&_ValidationRegistry.CallOpts)
}

// EXPIRATIONTIME is a free data retrieval call binding the contract method 0x4a5c7348.
//
// Solidity: function EXPIRATION_TIME() view returns(uint256)
func (_ValidationRegistry *ValidationRegistryCallerSession) EXPIRATIONTIME() (*big.Int, error) {
	return _ValidationRegistry.Contract.EXPIRATIONTIME(&_ValidationRegistry.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ValidationRegistry *ValidationRegistryCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ValidationRegistry *ValidationRegistrySession) VERSION() (string, error) {
	return _ValidationRegistry.Contract.VERSION(&_ValidationRegistry.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ValidationRegistry *ValidationRegistryCallerSession) VERSION() (string, error) {
	return _ValidationRegistry.Contract.VERSION(&_ValidationRegistry.CallOpts)
}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_ValidationRegistry *ValidationRegistryCaller) IdentityRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "_identityRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_ValidationRegistry *ValidationRegistrySession) IdentityRegistry() (common.Address, error) {
	return _ValidationRegistry.Contract.IdentityRegistry(&_ValidationRegistry.CallOpts)
}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_ValidationRegistry *ValidationRegistryCallerSession) IdentityRegistry() (common.Address, error) {
	return _ValidationRegistry.Contract.IdentityRegistry(&_ValidationRegistry.CallOpts)
}

// GetExpirationSlots is a free data retrieval call binding the contract method 0x6cc19969.
//
// Solidity: function getExpirationSlots() pure returns(uint256 slots)
func (_ValidationRegistry *ValidationRegistryCaller) GetExpirationSlots(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "getExpirationSlots")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExpirationSlots is a free data retrieval call binding the contract method 0x6cc19969.
//
// Solidity: function getExpirationSlots() pure returns(uint256 slots)
func (_ValidationRegistry *ValidationRegistrySession) GetExpirationSlots() (*big.Int, error) {
	return _ValidationRegistry.Contract.GetExpirationSlots(&_ValidationRegistry.CallOpts)
}

// GetExpirationSlots is a free data retrieval call binding the contract method 0x6cc19969.
//
// Solidity: function getExpirationSlots() pure returns(uint256 slots)
func (_ValidationRegistry *ValidationRegistryCallerSession) GetExpirationSlots() (*big.Int, error) {
	return _ValidationRegistry.Contract.GetExpirationSlots(&_ValidationRegistry.CallOpts)
}

// GetValidationRequest is a free data retrieval call binding the contract method 0xd097037a.
//
// Solidity: function getValidationRequest(bytes32 dataHash) view returns((uint256,uint256,bytes32,uint256,bool) request)
func (_ValidationRegistry *ValidationRegistryCaller) GetValidationRequest(opts *bind.CallOpts, dataHash [32]byte) (IValidationRegistryRequest, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "getValidationRequest", dataHash)

	if err != nil {
		return *new(IValidationRegistryRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(IValidationRegistryRequest)).(*IValidationRegistryRequest)

	return out0, err

}

// GetValidationRequest is a free data retrieval call binding the contract method 0xd097037a.
//
// Solidity: function getValidationRequest(bytes32 dataHash) view returns((uint256,uint256,bytes32,uint256,bool) request)
func (_ValidationRegistry *ValidationRegistrySession) GetValidationRequest(dataHash [32]byte) (IValidationRegistryRequest, error) {
	return _ValidationRegistry.Contract.GetValidationRequest(&_ValidationRegistry.CallOpts, dataHash)
}

// GetValidationRequest is a free data retrieval call binding the contract method 0xd097037a.
//
// Solidity: function getValidationRequest(bytes32 dataHash) view returns((uint256,uint256,bytes32,uint256,bool) request)
func (_ValidationRegistry *ValidationRegistryCallerSession) GetValidationRequest(dataHash [32]byte) (IValidationRegistryRequest, error) {
	return _ValidationRegistry.Contract.GetValidationRequest(&_ValidationRegistry.CallOpts, dataHash)
}

// GetValidationResponse is a free data retrieval call binding the contract method 0xc941b161.
//
// Solidity: function getValidationResponse(bytes32 dataHash) view returns(bool hasResponse, uint8 response)
func (_ValidationRegistry *ValidationRegistryCaller) GetValidationResponse(opts *bind.CallOpts, dataHash [32]byte) (struct {
	HasResponse bool
	Response    uint8
}, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "getValidationResponse", dataHash)

	outstruct := new(struct {
		HasResponse bool
		Response    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.HasResponse = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Response = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetValidationResponse is a free data retrieval call binding the contract method 0xc941b161.
//
// Solidity: function getValidationResponse(bytes32 dataHash) view returns(bool hasResponse, uint8 response)
func (_ValidationRegistry *ValidationRegistrySession) GetValidationResponse(dataHash [32]byte) (struct {
	HasResponse bool
	Response    uint8
}, error) {
	return _ValidationRegistry.Contract.GetValidationResponse(&_ValidationRegistry.CallOpts, dataHash)
}

// GetValidationResponse is a free data retrieval call binding the contract method 0xc941b161.
//
// Solidity: function getValidationResponse(bytes32 dataHash) view returns(bool hasResponse, uint8 response)
func (_ValidationRegistry *ValidationRegistryCallerSession) GetValidationResponse(dataHash [32]byte) (struct {
	HasResponse bool
	Response    uint8
}, error) {
	return _ValidationRegistry.Contract.GetValidationResponse(&_ValidationRegistry.CallOpts, dataHash)
}

// IsValidationPending is a free data retrieval call binding the contract method 0x625123da.
//
// Solidity: function isValidationPending(bytes32 dataHash) view returns(bool exists, bool pending)
func (_ValidationRegistry *ValidationRegistryCaller) IsValidationPending(opts *bind.CallOpts, dataHash [32]byte) (struct {
	Exists  bool
	Pending bool
}, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "isValidationPending", dataHash)

	outstruct := new(struct {
		Exists  bool
		Pending bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Pending = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// IsValidationPending is a free data retrieval call binding the contract method 0x625123da.
//
// Solidity: function isValidationPending(bytes32 dataHash) view returns(bool exists, bool pending)
func (_ValidationRegistry *ValidationRegistrySession) IsValidationPending(dataHash [32]byte) (struct {
	Exists  bool
	Pending bool
}, error) {
	return _ValidationRegistry.Contract.IsValidationPending(&_ValidationRegistry.CallOpts, dataHash)
}

// IsValidationPending is a free data retrieval call binding the contract method 0x625123da.
//
// Solidity: function isValidationPending(bytes32 dataHash) view returns(bool exists, bool pending)
func (_ValidationRegistry *ValidationRegistryCallerSession) IsValidationPending(dataHash [32]byte) (struct {
	Exists  bool
	Pending bool
}, error) {
	return _ValidationRegistry.Contract.IsValidationPending(&_ValidationRegistry.CallOpts, dataHash)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidationRegistry *ValidationRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidationRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidationRegistry *ValidationRegistrySession) Owner() (common.Address, error) {
	return _ValidationRegistry.Contract.Owner(&_ValidationRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidationRegistry *ValidationRegistryCallerSession) Owner() (common.Address, error) {
	return _ValidationRegistry.Contract.Owner(&_ValidationRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address identityRegistry) returns()
func (_ValidationRegistry *ValidationRegistryTransactor) Initialize(opts *bind.TransactOpts, identityRegistry common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.contract.Transact(opts, "initialize", identityRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address identityRegistry) returns()
func (_ValidationRegistry *ValidationRegistrySession) Initialize(identityRegistry common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.Initialize(&_ValidationRegistry.TransactOpts, identityRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address identityRegistry) returns()
func (_ValidationRegistry *ValidationRegistryTransactorSession) Initialize(identityRegistry common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.Initialize(&_ValidationRegistry.TransactOpts, identityRegistry)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidationRegistry *ValidationRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidationRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidationRegistry *ValidationRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidationRegistry.Contract.RenounceOwnership(&_ValidationRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidationRegistry *ValidationRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidationRegistry.Contract.RenounceOwnership(&_ValidationRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidationRegistry *ValidationRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidationRegistry *ValidationRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.TransferOwnership(&_ValidationRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidationRegistry *ValidationRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.TransferOwnership(&_ValidationRegistry.TransactOpts, newOwner)
}

// ValidationRequest is a paid mutator transaction binding the contract method 0x9748c8cb.
//
// Solidity: function validationRequest(uint256 agentValidatorId, uint256 agentServerId, bytes32 dataHash) returns()
func (_ValidationRegistry *ValidationRegistryTransactor) ValidationRequest(opts *bind.TransactOpts, agentValidatorId *big.Int, agentServerId *big.Int, dataHash [32]byte) (*types.Transaction, error) {
	return _ValidationRegistry.contract.Transact(opts, "validationRequest", agentValidatorId, agentServerId, dataHash)
}

// ValidationRequest is a paid mutator transaction binding the contract method 0x9748c8cb.
//
// Solidity: function validationRequest(uint256 agentValidatorId, uint256 agentServerId, bytes32 dataHash) returns()
func (_ValidationRegistry *ValidationRegistrySession) ValidationRequest(agentValidatorId *big.Int, agentServerId *big.Int, dataHash [32]byte) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationRequest(&_ValidationRegistry.TransactOpts, agentValidatorId, agentServerId, dataHash)
}

// ValidationRequest is a paid mutator transaction binding the contract method 0x9748c8cb.
//
// Solidity: function validationRequest(uint256 agentValidatorId, uint256 agentServerId, bytes32 dataHash) returns()
func (_ValidationRegistry *ValidationRegistryTransactorSession) ValidationRequest(agentValidatorId *big.Int, agentServerId *big.Int, dataHash [32]byte) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationRequest(&_ValidationRegistry.TransactOpts, agentValidatorId, agentServerId, dataHash)
}

// ValidationResponse is a paid mutator transaction binding the contract method 0x4146624a.
//
// Solidity: function validationResponse(bytes32 dataHash, uint8 response) returns()
func (_ValidationRegistry *ValidationRegistryTransactor) ValidationResponse(opts *bind.TransactOpts, dataHash [32]byte, response uint8) (*types.Transaction, error) {
	return _ValidationRegistry.contract.Transact(opts, "validationResponse", dataHash, response)
}

// ValidationResponse is a paid mutator transaction binding the contract method 0x4146624a.
//
// Solidity: function validationResponse(bytes32 dataHash, uint8 response) returns()
func (_ValidationRegistry *ValidationRegistrySession) ValidationResponse(dataHash [32]byte, response uint8) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationResponse(&_ValidationRegistry.TransactOpts, dataHash, response)
}

// ValidationResponse is a paid mutator transaction binding the contract method 0x4146624a.
//
// Solidity: function validationResponse(bytes32 dataHash, uint8 response) returns()
func (_ValidationRegistry *ValidationRegistryTransactorSession) ValidationResponse(dataHash [32]byte, response uint8) (*types.Transaction, error) {
	return _ValidationRegistry.Contract.ValidationResponse(&_ValidationRegistry.TransactOpts, dataHash, response)
}

// ValidationRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ValidationRegistry contract.
type ValidationRegistryInitializedIterator struct {
	Event *ValidationRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ValidationRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationRegistryInitialized)
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
		it.Event = new(ValidationRegistryInitialized)
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
func (it *ValidationRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationRegistryInitialized represents a Initialized event raised by the ValidationRegistry contract.
type ValidationRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ValidationRegistry *ValidationRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ValidationRegistryInitializedIterator, error) {

	logs, sub, err := _ValidationRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryInitializedIterator{contract: _ValidationRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ValidationRegistry *ValidationRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ValidationRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ValidationRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationRegistryInitialized)
				if err := _ValidationRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ValidationRegistry *ValidationRegistryFilterer) ParseInitialized(log types.Log) (*ValidationRegistryInitialized, error) {
	event := new(ValidationRegistryInitialized)
	if err := _ValidationRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidationRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ValidationRegistry contract.
type ValidationRegistryOwnershipTransferredIterator struct {
	Event *ValidationRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidationRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationRegistryOwnershipTransferred)
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
		it.Event = new(ValidationRegistryOwnershipTransferred)
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
func (it *ValidationRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ValidationRegistry contract.
type ValidationRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidationRegistry *ValidationRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidationRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidationRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryOwnershipTransferredIterator{contract: _ValidationRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidationRegistry *ValidationRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidationRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidationRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationRegistryOwnershipTransferred)
				if err := _ValidationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidationRegistry *ValidationRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ValidationRegistryOwnershipTransferred, error) {
	event := new(ValidationRegistryOwnershipTransferred)
	if err := _ValidationRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidationRegistryValidationRequestEventIterator is returned from FilterValidationRequestEvent and is used to iterate over the raw logs and unpacked data for ValidationRequestEvent events raised by the ValidationRegistry contract.
type ValidationRegistryValidationRequestEventIterator struct {
	Event *ValidationRegistryValidationRequestEvent // Event containing the contract specifics and raw log

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
func (it *ValidationRegistryValidationRequestEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationRegistryValidationRequestEvent)
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
		it.Event = new(ValidationRegistryValidationRequestEvent)
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
func (it *ValidationRegistryValidationRequestEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationRegistryValidationRequestEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationRegistryValidationRequestEvent represents a ValidationRequestEvent event raised by the ValidationRegistry contract.
type ValidationRegistryValidationRequestEvent struct {
	AgentValidatorId *big.Int
	AgentServerId    *big.Int
	DataHash         [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterValidationRequestEvent is a free log retrieval operation binding the contract event 0x3fcaa5f72d6fb1deb76744afac4fbc92c823e52ffcf39ad2cdd9da66a80b5196.
//
// Solidity: event ValidationRequestEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash)
func (_ValidationRegistry *ValidationRegistryFilterer) FilterValidationRequestEvent(opts *bind.FilterOpts, agentValidatorId []*big.Int, agentServerId []*big.Int, dataHash [][32]byte) (*ValidationRegistryValidationRequestEventIterator, error) {

	var agentValidatorIdRule []interface{}
	for _, agentValidatorIdItem := range agentValidatorId {
		agentValidatorIdRule = append(agentValidatorIdRule, agentValidatorIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _ValidationRegistry.contract.FilterLogs(opts, "ValidationRequestEvent", agentValidatorIdRule, agentServerIdRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryValidationRequestEventIterator{contract: _ValidationRegistry.contract, event: "ValidationRequestEvent", logs: logs, sub: sub}, nil
}

// WatchValidationRequestEvent is a free log subscription operation binding the contract event 0x3fcaa5f72d6fb1deb76744afac4fbc92c823e52ffcf39ad2cdd9da66a80b5196.
//
// Solidity: event ValidationRequestEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash)
func (_ValidationRegistry *ValidationRegistryFilterer) WatchValidationRequestEvent(opts *bind.WatchOpts, sink chan<- *ValidationRegistryValidationRequestEvent, agentValidatorId []*big.Int, agentServerId []*big.Int, dataHash [][32]byte) (event.Subscription, error) {

	var agentValidatorIdRule []interface{}
	for _, agentValidatorIdItem := range agentValidatorId {
		agentValidatorIdRule = append(agentValidatorIdRule, agentValidatorIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _ValidationRegistry.contract.WatchLogs(opts, "ValidationRequestEvent", agentValidatorIdRule, agentServerIdRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationRegistryValidationRequestEvent)
				if err := _ValidationRegistry.contract.UnpackLog(event, "ValidationRequestEvent", log); err != nil {
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

// ParseValidationRequestEvent is a log parse operation binding the contract event 0x3fcaa5f72d6fb1deb76744afac4fbc92c823e52ffcf39ad2cdd9da66a80b5196.
//
// Solidity: event ValidationRequestEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash)
func (_ValidationRegistry *ValidationRegistryFilterer) ParseValidationRequestEvent(log types.Log) (*ValidationRegistryValidationRequestEvent, error) {
	event := new(ValidationRegistryValidationRequestEvent)
	if err := _ValidationRegistry.contract.UnpackLog(event, "ValidationRequestEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidationRegistryValidationResponseEventIterator is returned from FilterValidationResponseEvent and is used to iterate over the raw logs and unpacked data for ValidationResponseEvent events raised by the ValidationRegistry contract.
type ValidationRegistryValidationResponseEventIterator struct {
	Event *ValidationRegistryValidationResponseEvent // Event containing the contract specifics and raw log

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
func (it *ValidationRegistryValidationResponseEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidationRegistryValidationResponseEvent)
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
		it.Event = new(ValidationRegistryValidationResponseEvent)
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
func (it *ValidationRegistryValidationResponseEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidationRegistryValidationResponseEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidationRegistryValidationResponseEvent represents a ValidationResponseEvent event raised by the ValidationRegistry contract.
type ValidationRegistryValidationResponseEvent struct {
	AgentValidatorId *big.Int
	AgentServerId    *big.Int
	DataHash         [32]byte
	Response         uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterValidationResponseEvent is a free log retrieval operation binding the contract event 0x2f1d7a4e773d3a02a2db8d4ba16df1476ec97974bbc93c41427689f236664f1b.
//
// Solidity: event ValidationResponseEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash, uint8 response)
func (_ValidationRegistry *ValidationRegistryFilterer) FilterValidationResponseEvent(opts *bind.FilterOpts, agentValidatorId []*big.Int, agentServerId []*big.Int, dataHash [][32]byte) (*ValidationRegistryValidationResponseEventIterator, error) {

	var agentValidatorIdRule []interface{}
	for _, agentValidatorIdItem := range agentValidatorId {
		agentValidatorIdRule = append(agentValidatorIdRule, agentValidatorIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _ValidationRegistry.contract.FilterLogs(opts, "ValidationResponseEvent", agentValidatorIdRule, agentServerIdRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return &ValidationRegistryValidationResponseEventIterator{contract: _ValidationRegistry.contract, event: "ValidationResponseEvent", logs: logs, sub: sub}, nil
}

// WatchValidationResponseEvent is a free log subscription operation binding the contract event 0x2f1d7a4e773d3a02a2db8d4ba16df1476ec97974bbc93c41427689f236664f1b.
//
// Solidity: event ValidationResponseEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash, uint8 response)
func (_ValidationRegistry *ValidationRegistryFilterer) WatchValidationResponseEvent(opts *bind.WatchOpts, sink chan<- *ValidationRegistryValidationResponseEvent, agentValidatorId []*big.Int, agentServerId []*big.Int, dataHash [][32]byte) (event.Subscription, error) {

	var agentValidatorIdRule []interface{}
	for _, agentValidatorIdItem := range agentValidatorId {
		agentValidatorIdRule = append(agentValidatorIdRule, agentValidatorIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _ValidationRegistry.contract.WatchLogs(opts, "ValidationResponseEvent", agentValidatorIdRule, agentServerIdRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidationRegistryValidationResponseEvent)
				if err := _ValidationRegistry.contract.UnpackLog(event, "ValidationResponseEvent", log); err != nil {
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

// ParseValidationResponseEvent is a log parse operation binding the contract event 0x2f1d7a4e773d3a02a2db8d4ba16df1476ec97974bbc93c41427689f236664f1b.
//
// Solidity: event ValidationResponseEvent(uint256 indexed agentValidatorId, uint256 indexed agentServerId, bytes32 indexed dataHash, uint8 response)
func (_ValidationRegistry *ValidationRegistryFilterer) ParseValidationResponseEvent(log types.Log) (*ValidationRegistryValidationResponseEvent, error) {
	event := new(ValidationRegistryValidationResponseEvent)
	if err := _ValidationRegistry.contract.UnpackLog(event, "ValidationResponseEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
