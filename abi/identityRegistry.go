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

// IIdentityRegistryAgentInfo is an auto generated low-level Go binding around an user-defined struct.
type IIdentityRegistryAgentInfo struct {
	AgentId      *big.Int
	AgentDomain  string
	AgentAddress common.Address
}

// IdentityRegistryMetaData contains all meta data concerning the IdentityRegistry contract.
var IdentityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AgentNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DomainAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDomain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedRegistration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedUpdate\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"name\":\"AgentRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"name\":\"AgentUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"agentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getAgent\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"internalType\":\"structIIdentityRegistry.AgentInfo\",\"name\":\"agentInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAgentCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"name\":\"newAgent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"name\":\"resolveByAddress\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"internalType\":\"structIIdentityRegistry.AgentInfo\",\"name\":\"agentInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"}],\"name\":\"resolveByDomain\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"agentDomain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"}],\"internalType\":\"structIIdentityRegistry.AgentInfo\",\"name\":\"agentInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"newAgentDomain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"newAgentAddress\",\"type\":\"address\"}],\"name\":\"updateAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IdentityRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IdentityRegistryMetaData.ABI instead.
var IdentityRegistryABI = IdentityRegistryMetaData.ABI

// IdentityRegistry is an auto generated Go binding around an Ethereum contract.
type IdentityRegistry struct {
	IdentityRegistryCaller     // Read-only binding to the contract
	IdentityRegistryTransactor // Write-only binding to the contract
	IdentityRegistryFilterer   // Log filterer for contract events
}

// IdentityRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentityRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentityRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentityRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentityRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentityRegistrySession struct {
	Contract     *IdentityRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IdentityRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentityRegistryCallerSession struct {
	Contract *IdentityRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IdentityRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentityRegistryTransactorSession struct {
	Contract     *IdentityRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IdentityRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentityRegistryRaw struct {
	Contract *IdentityRegistry // Generic contract binding to access the raw methods on
}

// IdentityRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentityRegistryCallerRaw struct {
	Contract *IdentityRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IdentityRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentityRegistryTransactorRaw struct {
	Contract *IdentityRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentityRegistry creates a new instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistry(address common.Address, backend bind.ContractBackend) (*IdentityRegistry, error) {
	contract, err := bindIdentityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistry{IdentityRegistryCaller: IdentityRegistryCaller{contract: contract}, IdentityRegistryTransactor: IdentityRegistryTransactor{contract: contract}, IdentityRegistryFilterer: IdentityRegistryFilterer{contract: contract}}, nil
}

// NewIdentityRegistryCaller creates a new read-only instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryCaller(address common.Address, caller bind.ContractCaller) (*IdentityRegistryCaller, error) {
	contract, err := bindIdentityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryCaller{contract: contract}, nil
}

// NewIdentityRegistryTransactor creates a new write-only instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentityRegistryTransactor, error) {
	contract, err := bindIdentityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryTransactor{contract: contract}, nil
}

// NewIdentityRegistryFilterer creates a new log filterer instance of IdentityRegistry, bound to a specific deployed contract.
func NewIdentityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentityRegistryFilterer, error) {
	contract, err := bindIdentityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryFilterer{contract: contract}, nil
}

// bindIdentityRegistry binds a generic wrapper to an already deployed contract.
func bindIdentityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IdentityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistry *IdentityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistry.Contract.IdentityRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistry *IdentityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.IdentityRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistry *IdentityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.IdentityRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentityRegistry *IdentityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IdentityRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentityRegistry *IdentityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentityRegistry *IdentityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IdentityRegistry *IdentityRegistryCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IdentityRegistry *IdentityRegistrySession) VERSION() (string, error) {
	return _IdentityRegistry.Contract.VERSION(&_IdentityRegistry.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IdentityRegistry *IdentityRegistryCallerSession) VERSION() (string, error) {
	return _IdentityRegistry.Contract.VERSION(&_IdentityRegistry.CallOpts)
}

// AgentExists is a free data retrieval call binding the contract method 0xde99f157.
//
// Solidity: function agentExists(uint256 agentId) view returns(bool exists)
func (_IdentityRegistry *IdentityRegistryCaller) AgentExists(opts *bind.CallOpts, agentId *big.Int) (bool, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "agentExists", agentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentExists is a free data retrieval call binding the contract method 0xde99f157.
//
// Solidity: function agentExists(uint256 agentId) view returns(bool exists)
func (_IdentityRegistry *IdentityRegistrySession) AgentExists(agentId *big.Int) (bool, error) {
	return _IdentityRegistry.Contract.AgentExists(&_IdentityRegistry.CallOpts, agentId)
}

// AgentExists is a free data retrieval call binding the contract method 0xde99f157.
//
// Solidity: function agentExists(uint256 agentId) view returns(bool exists)
func (_IdentityRegistry *IdentityRegistryCallerSession) AgentExists(agentId *big.Int) (bool, error) {
	return _IdentityRegistry.Contract.AgentExists(&_IdentityRegistry.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCaller) GetAgent(opts *bind.CallOpts, agentId *big.Int) (IIdentityRegistryAgentInfo, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "getAgent", agentId)

	if err != nil {
		return *new(IIdentityRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IIdentityRegistryAgentInfo)).(*IIdentityRegistryAgentInfo)

	return out0, err

}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistrySession) GetAgent(agentId *big.Int) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.GetAgent(&_IdentityRegistry.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCallerSession) GetAgent(agentId *big.Int) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.GetAgent(&_IdentityRegistry.CallOpts, agentId)
}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256 count)
func (_IdentityRegistry *IdentityRegistryCaller) GetAgentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "getAgentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256 count)
func (_IdentityRegistry *IdentityRegistrySession) GetAgentCount() (*big.Int, error) {
	return _IdentityRegistry.Contract.GetAgentCount(&_IdentityRegistry.CallOpts)
}

// GetAgentCount is a free data retrieval call binding the contract method 0x91cab63e.
//
// Solidity: function getAgentCount() view returns(uint256 count)
func (_IdentityRegistry *IdentityRegistryCallerSession) GetAgentCount() (*big.Int, error) {
	return _IdentityRegistry.Contract.GetAgentCount(&_IdentityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistrySession) Owner() (common.Address, error) {
	return _IdentityRegistry.Contract.Owner(&_IdentityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IdentityRegistry *IdentityRegistryCallerSession) Owner() (common.Address, error) {
	return _IdentityRegistry.Contract.Owner(&_IdentityRegistry.CallOpts)
}

// ResolveByAddress is a free data retrieval call binding the contract method 0x0c638107.
//
// Solidity: function resolveByAddress(address agentAddress) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCaller) ResolveByAddress(opts *bind.CallOpts, agentAddress common.Address) (IIdentityRegistryAgentInfo, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "resolveByAddress", agentAddress)

	if err != nil {
		return *new(IIdentityRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IIdentityRegistryAgentInfo)).(*IIdentityRegistryAgentInfo)

	return out0, err

}

// ResolveByAddress is a free data retrieval call binding the contract method 0x0c638107.
//
// Solidity: function resolveByAddress(address agentAddress) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistrySession) ResolveByAddress(agentAddress common.Address) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.ResolveByAddress(&_IdentityRegistry.CallOpts, agentAddress)
}

// ResolveByAddress is a free data retrieval call binding the contract method 0x0c638107.
//
// Solidity: function resolveByAddress(address agentAddress) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCallerSession) ResolveByAddress(agentAddress common.Address) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.ResolveByAddress(&_IdentityRegistry.CallOpts, agentAddress)
}

// ResolveByDomain is a free data retrieval call binding the contract method 0x78ddad12.
//
// Solidity: function resolveByDomain(string agentDomain) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCaller) ResolveByDomain(opts *bind.CallOpts, agentDomain string) (IIdentityRegistryAgentInfo, error) {
	var out []interface{}
	err := _IdentityRegistry.contract.Call(opts, &out, "resolveByDomain", agentDomain)

	if err != nil {
		return *new(IIdentityRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IIdentityRegistryAgentInfo)).(*IIdentityRegistryAgentInfo)

	return out0, err

}

// ResolveByDomain is a free data retrieval call binding the contract method 0x78ddad12.
//
// Solidity: function resolveByDomain(string agentDomain) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistrySession) ResolveByDomain(agentDomain string) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.ResolveByDomain(&_IdentityRegistry.CallOpts, agentDomain)
}

// ResolveByDomain is a free data retrieval call binding the contract method 0x78ddad12.
//
// Solidity: function resolveByDomain(string agentDomain) view returns((uint256,string,address) agentInfo)
func (_IdentityRegistry *IdentityRegistryCallerSession) ResolveByDomain(agentDomain string) (IIdentityRegistryAgentInfo, error) {
	return _IdentityRegistry.Contract.ResolveByDomain(&_IdentityRegistry.CallOpts, agentDomain)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityRegistry *IdentityRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityRegistry *IdentityRegistrySession) Initialize() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.Initialize(&_IdentityRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.Initialize(&_IdentityRegistry.TransactOpts)
}

// NewAgent is a paid mutator transaction binding the contract method 0x4750d0fa.
//
// Solidity: function newAgent(string agentDomain, address agentAddress) returns(uint256 agentId)
func (_IdentityRegistry *IdentityRegistryTransactor) NewAgent(opts *bind.TransactOpts, agentDomain string, agentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "newAgent", agentDomain, agentAddress)
}

// NewAgent is a paid mutator transaction binding the contract method 0x4750d0fa.
//
// Solidity: function newAgent(string agentDomain, address agentAddress) returns(uint256 agentId)
func (_IdentityRegistry *IdentityRegistrySession) NewAgent(agentDomain string, agentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.NewAgent(&_IdentityRegistry.TransactOpts, agentDomain, agentAddress)
}

// NewAgent is a paid mutator transaction binding the contract method 0x4750d0fa.
//
// Solidity: function newAgent(string agentDomain, address agentAddress) returns(uint256 agentId)
func (_IdentityRegistry *IdentityRegistryTransactorSession) NewAgent(agentDomain string, agentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.NewAgent(&_IdentityRegistry.TransactOpts, agentDomain, agentAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RenounceOwnership(&_IdentityRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IdentityRegistry.Contract.RenounceOwnership(&_IdentityRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.TransferOwnership(&_IdentityRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IdentityRegistry *IdentityRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.TransferOwnership(&_IdentityRegistry.TransactOpts, newOwner)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x3bfbbad7.
//
// Solidity: function updateAgent(uint256 agentId, string newAgentDomain, address newAgentAddress) returns(bool success)
func (_IdentityRegistry *IdentityRegistryTransactor) UpdateAgent(opts *bind.TransactOpts, agentId *big.Int, newAgentDomain string, newAgentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.contract.Transact(opts, "updateAgent", agentId, newAgentDomain, newAgentAddress)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x3bfbbad7.
//
// Solidity: function updateAgent(uint256 agentId, string newAgentDomain, address newAgentAddress) returns(bool success)
func (_IdentityRegistry *IdentityRegistrySession) UpdateAgent(agentId *big.Int, newAgentDomain string, newAgentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateAgent(&_IdentityRegistry.TransactOpts, agentId, newAgentDomain, newAgentAddress)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0x3bfbbad7.
//
// Solidity: function updateAgent(uint256 agentId, string newAgentDomain, address newAgentAddress) returns(bool success)
func (_IdentityRegistry *IdentityRegistryTransactorSession) UpdateAgent(agentId *big.Int, newAgentDomain string, newAgentAddress common.Address) (*types.Transaction, error) {
	return _IdentityRegistry.Contract.UpdateAgent(&_IdentityRegistry.TransactOpts, agentId, newAgentDomain, newAgentAddress)
}

// IdentityRegistryAgentRegisteredIterator is returned from FilterAgentRegistered and is used to iterate over the raw logs and unpacked data for AgentRegistered events raised by the IdentityRegistry contract.
type IdentityRegistryAgentRegisteredIterator struct {
	Event *IdentityRegistryAgentRegistered // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryAgentRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryAgentRegistered)
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
		it.Event = new(IdentityRegistryAgentRegistered)
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
func (it *IdentityRegistryAgentRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryAgentRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryAgentRegistered represents a AgentRegistered event raised by the IdentityRegistry contract.
type IdentityRegistryAgentRegistered struct {
	AgentId      *big.Int
	AgentDomain  string
	AgentAddress common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistered is a free log retrieval operation binding the contract event 0xaef1fcdf962a03943b677ae3b92ef1b34c972aeef794ed6e9ba5f599b461883b.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterAgentRegistered(opts *bind.FilterOpts, agentId []*big.Int) (*IdentityRegistryAgentRegisteredIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "AgentRegistered", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryAgentRegisteredIterator{contract: _IdentityRegistry.contract, event: "AgentRegistered", logs: logs, sub: sub}, nil
}

// WatchAgentRegistered is a free log subscription operation binding the contract event 0xaef1fcdf962a03943b677ae3b92ef1b34c972aeef794ed6e9ba5f599b461883b.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchAgentRegistered(opts *bind.WatchOpts, sink chan<- *IdentityRegistryAgentRegistered, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "AgentRegistered", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryAgentRegistered)
				if err := _IdentityRegistry.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
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

// ParseAgentRegistered is a log parse operation binding the contract event 0xaef1fcdf962a03943b677ae3b92ef1b34c972aeef794ed6e9ba5f599b461883b.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseAgentRegistered(log types.Log) (*IdentityRegistryAgentRegistered, error) {
	event := new(IdentityRegistryAgentRegistered)
	if err := _IdentityRegistry.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryAgentUpdatedIterator is returned from FilterAgentUpdated and is used to iterate over the raw logs and unpacked data for AgentUpdated events raised by the IdentityRegistry contract.
type IdentityRegistryAgentUpdatedIterator struct {
	Event *IdentityRegistryAgentUpdated // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryAgentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryAgentUpdated)
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
		it.Event = new(IdentityRegistryAgentUpdated)
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
func (it *IdentityRegistryAgentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryAgentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryAgentUpdated represents a AgentUpdated event raised by the IdentityRegistry contract.
type IdentityRegistryAgentUpdated struct {
	AgentId      *big.Int
	AgentDomain  string
	AgentAddress common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0x1b37db56c92268433e33fda6783791cb9e5f8ffababd3952ab40716443f1bf89.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterAgentUpdated(opts *bind.FilterOpts, agentId []*big.Int) (*IdentityRegistryAgentUpdatedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "AgentUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryAgentUpdatedIterator{contract: _IdentityRegistry.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0x1b37db56c92268433e33fda6783791cb9e5f8ffababd3952ab40716443f1bf89.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *IdentityRegistryAgentUpdated, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "AgentUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryAgentUpdated)
				if err := _IdentityRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
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

// ParseAgentUpdated is a log parse operation binding the contract event 0x1b37db56c92268433e33fda6783791cb9e5f8ffababd3952ab40716443f1bf89.
//
// Solidity: event AgentUpdated(uint256 indexed agentId, string agentDomain, address agentAddress)
func (_IdentityRegistry *IdentityRegistryFilterer) ParseAgentUpdated(log types.Log) (*IdentityRegistryAgentUpdated, error) {
	event := new(IdentityRegistryAgentUpdated)
	if err := _IdentityRegistry.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IdentityRegistry contract.
type IdentityRegistryInitializedIterator struct {
	Event *IdentityRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryInitialized)
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
		it.Event = new(IdentityRegistryInitialized)
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
func (it *IdentityRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryInitialized represents a Initialized event raised by the IdentityRegistry contract.
type IdentityRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*IdentityRegistryInitializedIterator, error) {

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryInitializedIterator{contract: _IdentityRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IdentityRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryInitialized)
				if err := _IdentityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseInitialized(log types.Log) (*IdentityRegistryInitialized, error) {
	event := new(IdentityRegistryInitialized)
	if err := _IdentityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IdentityRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IdentityRegistry contract.
type IdentityRegistryOwnershipTransferredIterator struct {
	Event *IdentityRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IdentityRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IdentityRegistryOwnershipTransferred)
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
		it.Event = new(IdentityRegistryOwnershipTransferred)
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
func (it *IdentityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IdentityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IdentityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the IdentityRegistry contract.
type IdentityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistry *IdentityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IdentityRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IdentityRegistryOwnershipTransferredIterator{contract: _IdentityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IdentityRegistry *IdentityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IdentityRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IdentityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IdentityRegistryOwnershipTransferred)
				if err := _IdentityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IdentityRegistry *IdentityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*IdentityRegistryOwnershipTransferred, error) {
	event := new(IdentityRegistryOwnershipTransferred)
	if err := _IdentityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
