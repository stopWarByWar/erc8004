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

// ReputationRegistryMetaData contains all meta data concerning the ReputationRegistry contract.
var ReputationRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_identityRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"}],\"name\":\"FeedbackRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"score\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"tag1\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"tag2\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"feedbackUri\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"feedbackHash\",\"type\":\"bytes32\"}],\"name\":\"NewFeedback\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"responder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"responseUri\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"responseHash\",\"type\":\"bytes32\"}],\"name\":\"ResponseAppended\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"responseUri\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"responseHash\",\"type\":\"bytes32\"}],\"name\":\"appendResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"}],\"name\":\"getClients\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getIdentityRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"}],\"name\":\"getLastIndex\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"responders\",\"type\":\"address[]\"}],\"name\":\"getResponseCount\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"count\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"clientAddresses\",\"type\":\"address[]\"},{\"internalType\":\"bytes32\",\"name\":\"tag1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tag2\",\"type\":\"bytes32\"}],\"name\":\"getSummary\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"count\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"averageScore\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"score\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"tag1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tag2\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"feedbackUri\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"feedbackHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"feedbackAuth\",\"type\":\"bytes\"}],\"name\":\"giveFeedback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"clientAddresses\",\"type\":\"address[]\"},{\"internalType\":\"bytes32\",\"name\":\"tag1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tag2\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"includeRevoked\",\"type\":\"bool\"}],\"name\":\"readAllFeedback\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"clients\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"scores\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"tag1s\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"tag2s\",\"type\":\"bytes32[]\"},{\"internalType\":\"bool[]\",\"name\":\"revokedStatuses\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"clientAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"}],\"name\":\"readFeedback\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"score\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"tag1\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"tag2\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isRevoked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentId\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"feedbackIndex\",\"type\":\"uint64\"}],\"name\":\"revokeFeedback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ReputationRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ReputationRegistryMetaData.ABI instead.
var ReputationRegistryABI = ReputationRegistryMetaData.ABI

// ReputationRegistry is an auto generated Go binding around an Ethereum contract.
type ReputationRegistry struct {
	ReputationRegistryCaller     // Read-only binding to the contract
	ReputationRegistryTransactor // Write-only binding to the contract
	ReputationRegistryFilterer   // Log filterer for contract events
}

// ReputationRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReputationRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReputationRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReputationRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReputationRegistrySession struct {
	Contract     *ReputationRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ReputationRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReputationRegistryCallerSession struct {
	Contract *ReputationRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ReputationRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReputationRegistryTransactorSession struct {
	Contract     *ReputationRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ReputationRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReputationRegistryRaw struct {
	Contract *ReputationRegistry // Generic contract binding to access the raw methods on
}

// ReputationRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReputationRegistryCallerRaw struct {
	Contract *ReputationRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ReputationRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReputationRegistryTransactorRaw struct {
	Contract *ReputationRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReputationRegistry creates a new instance of ReputationRegistry, bound to a specific deployed contract.
func NewReputationRegistry(address common.Address, backend bind.ContractBackend) (*ReputationRegistry, error) {
	contract, err := bindReputationRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistry{ReputationRegistryCaller: ReputationRegistryCaller{contract: contract}, ReputationRegistryTransactor: ReputationRegistryTransactor{contract: contract}, ReputationRegistryFilterer: ReputationRegistryFilterer{contract: contract}}, nil
}

// NewReputationRegistryCaller creates a new read-only instance of ReputationRegistry, bound to a specific deployed contract.
func NewReputationRegistryCaller(address common.Address, caller bind.ContractCaller) (*ReputationRegistryCaller, error) {
	contract, err := bindReputationRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryCaller{contract: contract}, nil
}

// NewReputationRegistryTransactor creates a new write-only instance of ReputationRegistry, bound to a specific deployed contract.
func NewReputationRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ReputationRegistryTransactor, error) {
	contract, err := bindReputationRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryTransactor{contract: contract}, nil
}

// NewReputationRegistryFilterer creates a new log filterer instance of ReputationRegistry, bound to a specific deployed contract.
func NewReputationRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ReputationRegistryFilterer, error) {
	contract, err := bindReputationRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryFilterer{contract: contract}, nil
}

// bindReputationRegistry binds a generic wrapper to an already deployed contract.
func bindReputationRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReputationRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReputationRegistry *ReputationRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReputationRegistry.Contract.ReputationRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReputationRegistry *ReputationRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.ReputationRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReputationRegistry *ReputationRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.ReputationRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReputationRegistry *ReputationRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReputationRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReputationRegistry *ReputationRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReputationRegistry *ReputationRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetClients is a free data retrieval call binding the contract method 0x42dd519c.
//
// Solidity: function getClients(uint256 agentId) view returns(address[])
func (_ReputationRegistry *ReputationRegistryCaller) GetClients(opts *bind.CallOpts, agentId *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "getClients", agentId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetClients is a free data retrieval call binding the contract method 0x42dd519c.
//
// Solidity: function getClients(uint256 agentId) view returns(address[])
func (_ReputationRegistry *ReputationRegistrySession) GetClients(agentId *big.Int) ([]common.Address, error) {
	return _ReputationRegistry.Contract.GetClients(&_ReputationRegistry.CallOpts, agentId)
}

// GetClients is a free data retrieval call binding the contract method 0x42dd519c.
//
// Solidity: function getClients(uint256 agentId) view returns(address[])
func (_ReputationRegistry *ReputationRegistryCallerSession) GetClients(agentId *big.Int) ([]common.Address, error) {
	return _ReputationRegistry.Contract.GetClients(&_ReputationRegistry.CallOpts, agentId)
}

// GetIdentityRegistry is a free data retrieval call binding the contract method 0xbc4d861b.
//
// Solidity: function getIdentityRegistry() view returns(address)
func (_ReputationRegistry *ReputationRegistryCaller) GetIdentityRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "getIdentityRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIdentityRegistry is a free data retrieval call binding the contract method 0xbc4d861b.
//
// Solidity: function getIdentityRegistry() view returns(address)
func (_ReputationRegistry *ReputationRegistrySession) GetIdentityRegistry() (common.Address, error) {
	return _ReputationRegistry.Contract.GetIdentityRegistry(&_ReputationRegistry.CallOpts)
}

// GetIdentityRegistry is a free data retrieval call binding the contract method 0xbc4d861b.
//
// Solidity: function getIdentityRegistry() view returns(address)
func (_ReputationRegistry *ReputationRegistryCallerSession) GetIdentityRegistry() (common.Address, error) {
	return _ReputationRegistry.Contract.GetIdentityRegistry(&_ReputationRegistry.CallOpts)
}

// GetLastIndex is a free data retrieval call binding the contract method 0xf2d81759.
//
// Solidity: function getLastIndex(uint256 agentId, address clientAddress) view returns(uint64)
func (_ReputationRegistry *ReputationRegistryCaller) GetLastIndex(opts *bind.CallOpts, agentId *big.Int, clientAddress common.Address) (uint64, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "getLastIndex", agentId, clientAddress)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastIndex is a free data retrieval call binding the contract method 0xf2d81759.
//
// Solidity: function getLastIndex(uint256 agentId, address clientAddress) view returns(uint64)
func (_ReputationRegistry *ReputationRegistrySession) GetLastIndex(agentId *big.Int, clientAddress common.Address) (uint64, error) {
	return _ReputationRegistry.Contract.GetLastIndex(&_ReputationRegistry.CallOpts, agentId, clientAddress)
}

// GetLastIndex is a free data retrieval call binding the contract method 0xf2d81759.
//
// Solidity: function getLastIndex(uint256 agentId, address clientAddress) view returns(uint64)
func (_ReputationRegistry *ReputationRegistryCallerSession) GetLastIndex(agentId *big.Int, clientAddress common.Address) (uint64, error) {
	return _ReputationRegistry.Contract.GetLastIndex(&_ReputationRegistry.CallOpts, agentId, clientAddress)
}

// GetResponseCount is a free data retrieval call binding the contract method 0x6e04cacd.
//
// Solidity: function getResponseCount(uint256 agentId, address clientAddress, uint64 feedbackIndex, address[] responders) view returns(uint64 count)
func (_ReputationRegistry *ReputationRegistryCaller) GetResponseCount(opts *bind.CallOpts, agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responders []common.Address) (uint64, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "getResponseCount", agentId, clientAddress, feedbackIndex, responders)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetResponseCount is a free data retrieval call binding the contract method 0x6e04cacd.
//
// Solidity: function getResponseCount(uint256 agentId, address clientAddress, uint64 feedbackIndex, address[] responders) view returns(uint64 count)
func (_ReputationRegistry *ReputationRegistrySession) GetResponseCount(agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responders []common.Address) (uint64, error) {
	return _ReputationRegistry.Contract.GetResponseCount(&_ReputationRegistry.CallOpts, agentId, clientAddress, feedbackIndex, responders)
}

// GetResponseCount is a free data retrieval call binding the contract method 0x6e04cacd.
//
// Solidity: function getResponseCount(uint256 agentId, address clientAddress, uint64 feedbackIndex, address[] responders) view returns(uint64 count)
func (_ReputationRegistry *ReputationRegistryCallerSession) GetResponseCount(agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responders []common.Address) (uint64, error) {
	return _ReputationRegistry.Contract.GetResponseCount(&_ReputationRegistry.CallOpts, agentId, clientAddress, feedbackIndex, responders)
}

// GetSummary is a free data retrieval call binding the contract method 0x31259cff.
//
// Solidity: function getSummary(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2) view returns(uint64 count, uint8 averageScore)
func (_ReputationRegistry *ReputationRegistryCaller) GetSummary(opts *bind.CallOpts, agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte) (struct {
	Count        uint64
	AverageScore uint8
}, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "getSummary", agentId, clientAddresses, tag1, tag2)

	outstruct := new(struct {
		Count        uint64
		AverageScore uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Count = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.AverageScore = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetSummary is a free data retrieval call binding the contract method 0x31259cff.
//
// Solidity: function getSummary(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2) view returns(uint64 count, uint8 averageScore)
func (_ReputationRegistry *ReputationRegistrySession) GetSummary(agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte) (struct {
	Count        uint64
	AverageScore uint8
}, error) {
	return _ReputationRegistry.Contract.GetSummary(&_ReputationRegistry.CallOpts, agentId, clientAddresses, tag1, tag2)
}

// GetSummary is a free data retrieval call binding the contract method 0x31259cff.
//
// Solidity: function getSummary(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2) view returns(uint64 count, uint8 averageScore)
func (_ReputationRegistry *ReputationRegistryCallerSession) GetSummary(agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte) (struct {
	Count        uint64
	AverageScore uint8
}, error) {
	return _ReputationRegistry.Contract.GetSummary(&_ReputationRegistry.CallOpts, agentId, clientAddresses, tag1, tag2)
}

// ReadAllFeedback is a free data retrieval call binding the contract method 0x80b87594.
//
// Solidity: function readAllFeedback(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2, bool includeRevoked) view returns(address[] clients, uint8[] scores, bytes32[] tag1s, bytes32[] tag2s, bool[] revokedStatuses)
func (_ReputationRegistry *ReputationRegistryCaller) ReadAllFeedback(opts *bind.CallOpts, agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte, includeRevoked bool) (struct {
	Clients         []common.Address
	Scores          []uint8
	Tag1s           [][32]byte
	Tag2s           [][32]byte
	RevokedStatuses []bool
}, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "readAllFeedback", agentId, clientAddresses, tag1, tag2, includeRevoked)

	outstruct := new(struct {
		Clients         []common.Address
		Scores          []uint8
		Tag1s           [][32]byte
		Tag2s           [][32]byte
		RevokedStatuses []bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Clients = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Scores = *abi.ConvertType(out[1], new([]uint8)).(*[]uint8)
	outstruct.Tag1s = *abi.ConvertType(out[2], new([][32]byte)).(*[][32]byte)
	outstruct.Tag2s = *abi.ConvertType(out[3], new([][32]byte)).(*[][32]byte)
	outstruct.RevokedStatuses = *abi.ConvertType(out[4], new([]bool)).(*[]bool)

	return *outstruct, err

}

// ReadAllFeedback is a free data retrieval call binding the contract method 0x80b87594.
//
// Solidity: function readAllFeedback(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2, bool includeRevoked) view returns(address[] clients, uint8[] scores, bytes32[] tag1s, bytes32[] tag2s, bool[] revokedStatuses)
func (_ReputationRegistry *ReputationRegistrySession) ReadAllFeedback(agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte, includeRevoked bool) (struct {
	Clients         []common.Address
	Scores          []uint8
	Tag1s           [][32]byte
	Tag2s           [][32]byte
	RevokedStatuses []bool
}, error) {
	return _ReputationRegistry.Contract.ReadAllFeedback(&_ReputationRegistry.CallOpts, agentId, clientAddresses, tag1, tag2, includeRevoked)
}

// ReadAllFeedback is a free data retrieval call binding the contract method 0x80b87594.
//
// Solidity: function readAllFeedback(uint256 agentId, address[] clientAddresses, bytes32 tag1, bytes32 tag2, bool includeRevoked) view returns(address[] clients, uint8[] scores, bytes32[] tag1s, bytes32[] tag2s, bool[] revokedStatuses)
func (_ReputationRegistry *ReputationRegistryCallerSession) ReadAllFeedback(agentId *big.Int, clientAddresses []common.Address, tag1 [32]byte, tag2 [32]byte, includeRevoked bool) (struct {
	Clients         []common.Address
	Scores          []uint8
	Tag1s           [][32]byte
	Tag2s           [][32]byte
	RevokedStatuses []bool
}, error) {
	return _ReputationRegistry.Contract.ReadAllFeedback(&_ReputationRegistry.CallOpts, agentId, clientAddresses, tag1, tag2, includeRevoked)
}

// ReadFeedback is a free data retrieval call binding the contract method 0x232b0810.
//
// Solidity: function readFeedback(uint256 agentId, address clientAddress, uint64 index) view returns(uint8 score, bytes32 tag1, bytes32 tag2, bool isRevoked)
func (_ReputationRegistry *ReputationRegistryCaller) ReadFeedback(opts *bind.CallOpts, agentId *big.Int, clientAddress common.Address, index uint64) (struct {
	Score     uint8
	Tag1      [32]byte
	Tag2      [32]byte
	IsRevoked bool
}, error) {
	var out []interface{}
	err := _ReputationRegistry.contract.Call(opts, &out, "readFeedback", agentId, clientAddress, index)

	outstruct := new(struct {
		Score     uint8
		Tag1      [32]byte
		Tag2      [32]byte
		IsRevoked bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Score = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Tag1 = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Tag2 = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.IsRevoked = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// ReadFeedback is a free data retrieval call binding the contract method 0x232b0810.
//
// Solidity: function readFeedback(uint256 agentId, address clientAddress, uint64 index) view returns(uint8 score, bytes32 tag1, bytes32 tag2, bool isRevoked)
func (_ReputationRegistry *ReputationRegistrySession) ReadFeedback(agentId *big.Int, clientAddress common.Address, index uint64) (struct {
	Score     uint8
	Tag1      [32]byte
	Tag2      [32]byte
	IsRevoked bool
}, error) {
	return _ReputationRegistry.Contract.ReadFeedback(&_ReputationRegistry.CallOpts, agentId, clientAddress, index)
}

// ReadFeedback is a free data retrieval call binding the contract method 0x232b0810.
//
// Solidity: function readFeedback(uint256 agentId, address clientAddress, uint64 index) view returns(uint8 score, bytes32 tag1, bytes32 tag2, bool isRevoked)
func (_ReputationRegistry *ReputationRegistryCallerSession) ReadFeedback(agentId *big.Int, clientAddress common.Address, index uint64) (struct {
	Score     uint8
	Tag1      [32]byte
	Tag2      [32]byte
	IsRevoked bool
}, error) {
	return _ReputationRegistry.Contract.ReadFeedback(&_ReputationRegistry.CallOpts, agentId, clientAddress, index)
}

// AppendResponse is a paid mutator transaction binding the contract method 0xc2349ab2.
//
// Solidity: function appendResponse(uint256 agentId, address clientAddress, uint64 feedbackIndex, string responseUri, bytes32 responseHash) returns()
func (_ReputationRegistry *ReputationRegistryTransactor) AppendResponse(opts *bind.TransactOpts, agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responseUri string, responseHash [32]byte) (*types.Transaction, error) {
	return _ReputationRegistry.contract.Transact(opts, "appendResponse", agentId, clientAddress, feedbackIndex, responseUri, responseHash)
}

// AppendResponse is a paid mutator transaction binding the contract method 0xc2349ab2.
//
// Solidity: function appendResponse(uint256 agentId, address clientAddress, uint64 feedbackIndex, string responseUri, bytes32 responseHash) returns()
func (_ReputationRegistry *ReputationRegistrySession) AppendResponse(agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responseUri string, responseHash [32]byte) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.AppendResponse(&_ReputationRegistry.TransactOpts, agentId, clientAddress, feedbackIndex, responseUri, responseHash)
}

// AppendResponse is a paid mutator transaction binding the contract method 0xc2349ab2.
//
// Solidity: function appendResponse(uint256 agentId, address clientAddress, uint64 feedbackIndex, string responseUri, bytes32 responseHash) returns()
func (_ReputationRegistry *ReputationRegistryTransactorSession) AppendResponse(agentId *big.Int, clientAddress common.Address, feedbackIndex uint64, responseUri string, responseHash [32]byte) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.AppendResponse(&_ReputationRegistry.TransactOpts, agentId, clientAddress, feedbackIndex, responseUri, responseHash)
}

// GiveFeedback is a paid mutator transaction binding the contract method 0x155e5bbd.
//
// Solidity: function giveFeedback(uint256 agentId, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash, bytes feedbackAuth) returns()
func (_ReputationRegistry *ReputationRegistryTransactor) GiveFeedback(opts *bind.TransactOpts, agentId *big.Int, score uint8, tag1 [32]byte, tag2 [32]byte, feedbackUri string, feedbackHash [32]byte, feedbackAuth []byte) (*types.Transaction, error) {
	return _ReputationRegistry.contract.Transact(opts, "giveFeedback", agentId, score, tag1, tag2, feedbackUri, feedbackHash, feedbackAuth)
}

// GiveFeedback is a paid mutator transaction binding the contract method 0x155e5bbd.
//
// Solidity: function giveFeedback(uint256 agentId, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash, bytes feedbackAuth) returns()
func (_ReputationRegistry *ReputationRegistrySession) GiveFeedback(agentId *big.Int, score uint8, tag1 [32]byte, tag2 [32]byte, feedbackUri string, feedbackHash [32]byte, feedbackAuth []byte) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.GiveFeedback(&_ReputationRegistry.TransactOpts, agentId, score, tag1, tag2, feedbackUri, feedbackHash, feedbackAuth)
}

// GiveFeedback is a paid mutator transaction binding the contract method 0x155e5bbd.
//
// Solidity: function giveFeedback(uint256 agentId, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash, bytes feedbackAuth) returns()
func (_ReputationRegistry *ReputationRegistryTransactorSession) GiveFeedback(agentId *big.Int, score uint8, tag1 [32]byte, tag2 [32]byte, feedbackUri string, feedbackHash [32]byte, feedbackAuth []byte) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.GiveFeedback(&_ReputationRegistry.TransactOpts, agentId, score, tag1, tag2, feedbackUri, feedbackHash, feedbackAuth)
}

// RevokeFeedback is a paid mutator transaction binding the contract method 0x4ab3ca99.
//
// Solidity: function revokeFeedback(uint256 agentId, uint64 feedbackIndex) returns()
func (_ReputationRegistry *ReputationRegistryTransactor) RevokeFeedback(opts *bind.TransactOpts, agentId *big.Int, feedbackIndex uint64) (*types.Transaction, error) {
	return _ReputationRegistry.contract.Transact(opts, "revokeFeedback", agentId, feedbackIndex)
}

// RevokeFeedback is a paid mutator transaction binding the contract method 0x4ab3ca99.
//
// Solidity: function revokeFeedback(uint256 agentId, uint64 feedbackIndex) returns()
func (_ReputationRegistry *ReputationRegistrySession) RevokeFeedback(agentId *big.Int, feedbackIndex uint64) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.RevokeFeedback(&_ReputationRegistry.TransactOpts, agentId, feedbackIndex)
}

// RevokeFeedback is a paid mutator transaction binding the contract method 0x4ab3ca99.
//
// Solidity: function revokeFeedback(uint256 agentId, uint64 feedbackIndex) returns()
func (_ReputationRegistry *ReputationRegistryTransactorSession) RevokeFeedback(agentId *big.Int, feedbackIndex uint64) (*types.Transaction, error) {
	return _ReputationRegistry.Contract.RevokeFeedback(&_ReputationRegistry.TransactOpts, agentId, feedbackIndex)
}

// ReputationRegistryFeedbackRevokedIterator is returned from FilterFeedbackRevoked and is used to iterate over the raw logs and unpacked data for FeedbackRevoked events raised by the ReputationRegistry contract.
type ReputationRegistryFeedbackRevokedIterator struct {
	Event *ReputationRegistryFeedbackRevoked // Event containing the contract specifics and raw log

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
func (it *ReputationRegistryFeedbackRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationRegistryFeedbackRevoked)
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
		it.Event = new(ReputationRegistryFeedbackRevoked)
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
func (it *ReputationRegistryFeedbackRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationRegistryFeedbackRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationRegistryFeedbackRevoked represents a FeedbackRevoked event raised by the ReputationRegistry contract.
type ReputationRegistryFeedbackRevoked struct {
	AgentId       *big.Int
	ClientAddress common.Address
	FeedbackIndex uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeedbackRevoked is a free log retrieval operation binding the contract event 0x25156fd3288212246d8b008d5921fde376c71ed14ac2e072a506eb06fde6d09d.
//
// Solidity: event FeedbackRevoked(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex)
func (_ReputationRegistry *ReputationRegistryFilterer) FilterFeedbackRevoked(opts *bind.FilterOpts, agentId []*big.Int, clientAddress []common.Address, feedbackIndex []uint64) (*ReputationRegistryFeedbackRevokedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}
	var feedbackIndexRule []interface{}
	for _, feedbackIndexItem := range feedbackIndex {
		feedbackIndexRule = append(feedbackIndexRule, feedbackIndexItem)
	}

	logs, sub, err := _ReputationRegistry.contract.FilterLogs(opts, "FeedbackRevoked", agentIdRule, clientAddressRule, feedbackIndexRule)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryFeedbackRevokedIterator{contract: _ReputationRegistry.contract, event: "FeedbackRevoked", logs: logs, sub: sub}, nil
}

// WatchFeedbackRevoked is a free log subscription operation binding the contract event 0x25156fd3288212246d8b008d5921fde376c71ed14ac2e072a506eb06fde6d09d.
//
// Solidity: event FeedbackRevoked(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex)
func (_ReputationRegistry *ReputationRegistryFilterer) WatchFeedbackRevoked(opts *bind.WatchOpts, sink chan<- *ReputationRegistryFeedbackRevoked, agentId []*big.Int, clientAddress []common.Address, feedbackIndex []uint64) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}
	var feedbackIndexRule []interface{}
	for _, feedbackIndexItem := range feedbackIndex {
		feedbackIndexRule = append(feedbackIndexRule, feedbackIndexItem)
	}

	logs, sub, err := _ReputationRegistry.contract.WatchLogs(opts, "FeedbackRevoked", agentIdRule, clientAddressRule, feedbackIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationRegistryFeedbackRevoked)
				if err := _ReputationRegistry.contract.UnpackLog(event, "FeedbackRevoked", log); err != nil {
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

// ParseFeedbackRevoked is a log parse operation binding the contract event 0x25156fd3288212246d8b008d5921fde376c71ed14ac2e072a506eb06fde6d09d.
//
// Solidity: event FeedbackRevoked(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex)
func (_ReputationRegistry *ReputationRegistryFilterer) ParseFeedbackRevoked(log types.Log) (*ReputationRegistryFeedbackRevoked, error) {
	event := new(ReputationRegistryFeedbackRevoked)
	if err := _ReputationRegistry.contract.UnpackLog(event, "FeedbackRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationRegistryNewFeedbackIterator is returned from FilterNewFeedback and is used to iterate over the raw logs and unpacked data for NewFeedback events raised by the ReputationRegistry contract.
type ReputationRegistryNewFeedbackIterator struct {
	Event *ReputationRegistryNewFeedback // Event containing the contract specifics and raw log

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
func (it *ReputationRegistryNewFeedbackIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationRegistryNewFeedback)
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
		it.Event = new(ReputationRegistryNewFeedback)
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
func (it *ReputationRegistryNewFeedbackIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationRegistryNewFeedbackIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationRegistryNewFeedback represents a NewFeedback event raised by the ReputationRegistry contract.
type ReputationRegistryNewFeedback struct {
	AgentId       *big.Int
	ClientAddress common.Address
	FeedbackIndex uint64
	Score         uint8
	Tag1          [32]byte
	Tag2          [32]byte
	FeedbackUri   string
	FeedbackHash  [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewFeedback is a free log retrieval operation binding the contract event 0xb655ce21b319053e24bad48a8f38fa1a42101e27866f559ca10f597d2bb584a1.
//
// Solidity: event NewFeedback(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash)
func (_ReputationRegistry *ReputationRegistryFilterer) FilterNewFeedback(opts *bind.FilterOpts, agentId []*big.Int, clientAddress []common.Address, feedbackIndex []uint64) (*ReputationRegistryNewFeedbackIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}
	var feedbackIndexRule []interface{}
	for _, feedbackIndexItem := range feedbackIndex {
		feedbackIndexRule = append(feedbackIndexRule, feedbackIndexItem)
	}

	logs, sub, err := _ReputationRegistry.contract.FilterLogs(opts, "NewFeedback", agentIdRule, clientAddressRule, feedbackIndexRule)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryNewFeedbackIterator{contract: _ReputationRegistry.contract, event: "NewFeedback", logs: logs, sub: sub}, nil
}

// WatchNewFeedback is a free log subscription operation binding the contract event 0xb655ce21b319053e24bad48a8f38fa1a42101e27866f559ca10f597d2bb584a1.
//
// Solidity: event NewFeedback(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash)
func (_ReputationRegistry *ReputationRegistryFilterer) WatchNewFeedback(opts *bind.WatchOpts, sink chan<- *ReputationRegistryNewFeedback, agentId []*big.Int, clientAddress []common.Address, feedbackIndex []uint64) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}
	var feedbackIndexRule []interface{}
	for _, feedbackIndexItem := range feedbackIndex {
		feedbackIndexRule = append(feedbackIndexRule, feedbackIndexItem)
	}

	logs, sub, err := _ReputationRegistry.contract.WatchLogs(opts, "NewFeedback", agentIdRule, clientAddressRule, feedbackIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationRegistryNewFeedback)
				if err := _ReputationRegistry.contract.UnpackLog(event, "NewFeedback", log); err != nil {
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

// ParseNewFeedback is a log parse operation binding the contract event 0xb655ce21b319053e24bad48a8f38fa1a42101e27866f559ca10f597d2bb584a1.
//
// Solidity: event NewFeedback(uint256 indexed agentId, address indexed clientAddress, uint64 indexed feedbackIndex, uint8 score, bytes32 tag1, bytes32 tag2, string feedbackUri, bytes32 feedbackHash)
func (_ReputationRegistry *ReputationRegistryFilterer) ParseNewFeedback(log types.Log) (*ReputationRegistryNewFeedback, error) {
	event := new(ReputationRegistryNewFeedback)
	if err := _ReputationRegistry.contract.UnpackLog(event, "NewFeedback", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationRegistryResponseAppendedIterator is returned from FilterResponseAppended and is used to iterate over the raw logs and unpacked data for ResponseAppended events raised by the ReputationRegistry contract.
type ReputationRegistryResponseAppendedIterator struct {
	Event *ReputationRegistryResponseAppended // Event containing the contract specifics and raw log

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
func (it *ReputationRegistryResponseAppendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationRegistryResponseAppended)
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
		it.Event = new(ReputationRegistryResponseAppended)
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
func (it *ReputationRegistryResponseAppendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationRegistryResponseAppendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationRegistryResponseAppended represents a ResponseAppended event raised by the ReputationRegistry contract.
type ReputationRegistryResponseAppended struct {
	AgentId       *big.Int
	ClientAddress common.Address
	FeedbackIndex uint64
	Responder     common.Address
	ResponseUri   string
	ResponseHash  [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterResponseAppended is a free log retrieval operation binding the contract event 0xb1c6be0b5b8aef6539e2fac0fd131a2faa7b49edf8e505b5eb0ad487d56051d4.
//
// Solidity: event ResponseAppended(uint256 indexed agentId, address indexed clientAddress, uint64 feedbackIndex, address indexed responder, string responseUri, bytes32 responseHash)
func (_ReputationRegistry *ReputationRegistryFilterer) FilterResponseAppended(opts *bind.FilterOpts, agentId []*big.Int, clientAddress []common.Address, responder []common.Address) (*ReputationRegistryResponseAppendedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}

	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}

	logs, sub, err := _ReputationRegistry.contract.FilterLogs(opts, "ResponseAppended", agentIdRule, clientAddressRule, responderRule)
	if err != nil {
		return nil, err
	}
	return &ReputationRegistryResponseAppendedIterator{contract: _ReputationRegistry.contract, event: "ResponseAppended", logs: logs, sub: sub}, nil
}

// WatchResponseAppended is a free log subscription operation binding the contract event 0xb1c6be0b5b8aef6539e2fac0fd131a2faa7b49edf8e505b5eb0ad487d56051d4.
//
// Solidity: event ResponseAppended(uint256 indexed agentId, address indexed clientAddress, uint64 feedbackIndex, address indexed responder, string responseUri, bytes32 responseHash)
func (_ReputationRegistry *ReputationRegistryFilterer) WatchResponseAppended(opts *bind.WatchOpts, sink chan<- *ReputationRegistryResponseAppended, agentId []*big.Int, clientAddress []common.Address, responder []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var clientAddressRule []interface{}
	for _, clientAddressItem := range clientAddress {
		clientAddressRule = append(clientAddressRule, clientAddressItem)
	}

	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}

	logs, sub, err := _ReputationRegistry.contract.WatchLogs(opts, "ResponseAppended", agentIdRule, clientAddressRule, responderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationRegistryResponseAppended)
				if err := _ReputationRegistry.contract.UnpackLog(event, "ResponseAppended", log); err != nil {
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

// ParseResponseAppended is a log parse operation binding the contract event 0xb1c6be0b5b8aef6539e2fac0fd131a2faa7b49edf8e505b5eb0ad487d56051d4.
//
// Solidity: event ResponseAppended(uint256 indexed agentId, address indexed clientAddress, uint64 feedbackIndex, address indexed responder, string responseUri, bytes32 responseHash)
func (_ReputationRegistry *ReputationRegistryFilterer) ParseResponseAppended(log types.Log) (*ReputationRegistryResponseAppended, error) {
	event := new(ReputationRegistryResponseAppended)
	if err := _ReputationRegistry.contract.UnpackLog(event, "ResponseAppended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
