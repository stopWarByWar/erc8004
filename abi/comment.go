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

// CommentMetaData contains all meta data concerning the Comment contract.
var CommentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSchema\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidScore\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedCommentReply\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedFeedback\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentClientId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commentAttestationId\",\"type\":\"bytes32\"}],\"name\":\"CommentEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commentAttestationId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commentReplyAttestationId\",\"type\":\"bytes32\"}],\"name\":\"CommentReplyEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_bas\",\"outputs\":[{\"internalType\":\"contractIBAS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_identityRegistry\",\"outputs\":[{\"internalType\":\"contractIIdentityRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_reputationRegistry\",\"outputs\":[{\"internalType\":\"contractIReputationRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"agentClientId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"agentServerId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"commentText\",\"type\":\"string\"},{\"internalType\":\"uint16\",\"name\":\"score\",\"type\":\"uint16\"}],\"name\":\"comment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commentAttestationId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"commentText\",\"type\":\"string\"}],\"name\":\"commentReply\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIReputationRegistry\",\"name\":\"reputationRegistry\",\"type\":\"address\"},{\"internalType\":\"contractIIdentityRegistry\",\"name\":\"identityRegistry\",\"type\":\"address\"},{\"internalType\":\"contractIBAS\",\"name\":\"bas\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIBAS\",\"name\":\"bas\",\"type\":\"address\"}],\"name\":\"setBas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIIdentityRegistry\",\"name\":\"identityRegistry\",\"type\":\"address\"}],\"name\":\"setIdentityRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIReputationRegistry\",\"name\":\"reputationRegistry\",\"type\":\"address\"}],\"name\":\"setReputationRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CommentABI is the input ABI used to generate the binding from.
// Deprecated: Use CommentMetaData.ABI instead.
var CommentABI = CommentMetaData.ABI

// Comment is an auto generated Go binding around an Ethereum contract.
type Comment struct {
	CommentCaller     // Read-only binding to the contract
	CommentTransactor // Write-only binding to the contract
	CommentFilterer   // Log filterer for contract events
}

// CommentCaller is an auto generated read-only Go binding around an Ethereum contract.
type CommentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CommentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommentSession struct {
	Contract     *Comment          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommentCallerSession struct {
	Contract *CommentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// CommentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommentTransactorSession struct {
	Contract     *CommentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CommentRaw is an auto generated low-level Go binding around an Ethereum contract.
type CommentRaw struct {
	Contract *Comment // Generic contract binding to access the raw methods on
}

// CommentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommentCallerRaw struct {
	Contract *CommentCaller // Generic read-only contract binding to access the raw methods on
}

// CommentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommentTransactorRaw struct {
	Contract *CommentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewComment creates a new instance of Comment, bound to a specific deployed contract.
func NewComment(address common.Address, backend bind.ContractBackend) (*Comment, error) {
	contract, err := bindComment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Comment{CommentCaller: CommentCaller{contract: contract}, CommentTransactor: CommentTransactor{contract: contract}, CommentFilterer: CommentFilterer{contract: contract}}, nil
}

// NewCommentCaller creates a new read-only instance of Comment, bound to a specific deployed contract.
func NewCommentCaller(address common.Address, caller bind.ContractCaller) (*CommentCaller, error) {
	contract, err := bindComment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommentCaller{contract: contract}, nil
}

// NewCommentTransactor creates a new write-only instance of Comment, bound to a specific deployed contract.
func NewCommentTransactor(address common.Address, transactor bind.ContractTransactor) (*CommentTransactor, error) {
	contract, err := bindComment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommentTransactor{contract: contract}, nil
}

// NewCommentFilterer creates a new log filterer instance of Comment, bound to a specific deployed contract.
func NewCommentFilterer(address common.Address, filterer bind.ContractFilterer) (*CommentFilterer, error) {
	contract, err := bindComment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommentFilterer{contract: contract}, nil
}

// bindComment binds a generic wrapper to an already deployed contract.
func bindComment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CommentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Comment *CommentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Comment.Contract.CommentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Comment *CommentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Comment.Contract.CommentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Comment *CommentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Comment.Contract.CommentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Comment *CommentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Comment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Comment *CommentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Comment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Comment *CommentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Comment.Contract.contract.Transact(opts, method, params...)
}

// Bas is a free data retrieval call binding the contract method 0x20d171a4.
//
// Solidity: function _bas() view returns(address)
func (_Comment *CommentCaller) Bas(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Comment.contract.Call(opts, &out, "_bas")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bas is a free data retrieval call binding the contract method 0x20d171a4.
//
// Solidity: function _bas() view returns(address)
func (_Comment *CommentSession) Bas() (common.Address, error) {
	return _Comment.Contract.Bas(&_Comment.CallOpts)
}

// Bas is a free data retrieval call binding the contract method 0x20d171a4.
//
// Solidity: function _bas() view returns(address)
func (_Comment *CommentCallerSession) Bas() (common.Address, error) {
	return _Comment.Contract.Bas(&_Comment.CallOpts)
}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_Comment *CommentCaller) IdentityRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Comment.contract.Call(opts, &out, "_identityRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_Comment *CommentSession) IdentityRegistry() (common.Address, error) {
	return _Comment.Contract.IdentityRegistry(&_Comment.CallOpts)
}

// IdentityRegistry is a free data retrieval call binding the contract method 0x91993b61.
//
// Solidity: function _identityRegistry() view returns(address)
func (_Comment *CommentCallerSession) IdentityRegistry() (common.Address, error) {
	return _Comment.Contract.IdentityRegistry(&_Comment.CallOpts)
}

// ReputationRegistry is a free data retrieval call binding the contract method 0x156b912f.
//
// Solidity: function _reputationRegistry() view returns(address)
func (_Comment *CommentCaller) ReputationRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Comment.contract.Call(opts, &out, "_reputationRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReputationRegistry is a free data retrieval call binding the contract method 0x156b912f.
//
// Solidity: function _reputationRegistry() view returns(address)
func (_Comment *CommentSession) ReputationRegistry() (common.Address, error) {
	return _Comment.Contract.ReputationRegistry(&_Comment.CallOpts)
}

// ReputationRegistry is a free data retrieval call binding the contract method 0x156b912f.
//
// Solidity: function _reputationRegistry() view returns(address)
func (_Comment *CommentCallerSession) ReputationRegistry() (common.Address, error) {
	return _Comment.Contract.ReputationRegistry(&_Comment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Comment *CommentCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Comment.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Comment *CommentSession) Owner() (common.Address, error) {
	return _Comment.Contract.Owner(&_Comment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Comment *CommentCallerSession) Owner() (common.Address, error) {
	return _Comment.Contract.Owner(&_Comment.CallOpts)
}

// Comment is a paid mutator transaction binding the contract method 0xc8d871d2.
//
// Solidity: function comment(uint256 agentClientId, uint256 agentServerId, string commentText, uint16 score) returns()
func (_Comment *CommentTransactor) Comment(opts *bind.TransactOpts, agentClientId *big.Int, agentServerId *big.Int, commentText string, score uint16) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "comment", agentClientId, agentServerId, commentText, score)
}

// Comment is a paid mutator transaction binding the contract method 0xc8d871d2.
//
// Solidity: function comment(uint256 agentClientId, uint256 agentServerId, string commentText, uint16 score) returns()
func (_Comment *CommentSession) Comment(agentClientId *big.Int, agentServerId *big.Int, commentText string, score uint16) (*types.Transaction, error) {
	return _Comment.Contract.Comment(&_Comment.TransactOpts, agentClientId, agentServerId, commentText, score)
}

// Comment is a paid mutator transaction binding the contract method 0xc8d871d2.
//
// Solidity: function comment(uint256 agentClientId, uint256 agentServerId, string commentText, uint16 score) returns()
func (_Comment *CommentTransactorSession) Comment(agentClientId *big.Int, agentServerId *big.Int, commentText string, score uint16) (*types.Transaction, error) {
	return _Comment.Contract.Comment(&_Comment.TransactOpts, agentClientId, agentServerId, commentText, score)
}

// CommentReply is a paid mutator transaction binding the contract method 0xd35d4116.
//
// Solidity: function commentReply(bytes32 commentAttestationId, string commentText) returns()
func (_Comment *CommentTransactor) CommentReply(opts *bind.TransactOpts, commentAttestationId [32]byte, commentText string) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "commentReply", commentAttestationId, commentText)
}

// CommentReply is a paid mutator transaction binding the contract method 0xd35d4116.
//
// Solidity: function commentReply(bytes32 commentAttestationId, string commentText) returns()
func (_Comment *CommentSession) CommentReply(commentAttestationId [32]byte, commentText string) (*types.Transaction, error) {
	return _Comment.Contract.CommentReply(&_Comment.TransactOpts, commentAttestationId, commentText)
}

// CommentReply is a paid mutator transaction binding the contract method 0xd35d4116.
//
// Solidity: function commentReply(bytes32 commentAttestationId, string commentText) returns()
func (_Comment *CommentTransactorSession) CommentReply(commentAttestationId [32]byte, commentText string) (*types.Transaction, error) {
	return _Comment.Contract.CommentReply(&_Comment.TransactOpts, commentAttestationId, commentText)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address reputationRegistry, address identityRegistry, address bas) returns()
func (_Comment *CommentTransactor) Initialize(opts *bind.TransactOpts, reputationRegistry common.Address, identityRegistry common.Address, bas common.Address) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "initialize", reputationRegistry, identityRegistry, bas)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address reputationRegistry, address identityRegistry, address bas) returns()
func (_Comment *CommentSession) Initialize(reputationRegistry common.Address, identityRegistry common.Address, bas common.Address) (*types.Transaction, error) {
	return _Comment.Contract.Initialize(&_Comment.TransactOpts, reputationRegistry, identityRegistry, bas)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address reputationRegistry, address identityRegistry, address bas) returns()
func (_Comment *CommentTransactorSession) Initialize(reputationRegistry common.Address, identityRegistry common.Address, bas common.Address) (*types.Transaction, error) {
	return _Comment.Contract.Initialize(&_Comment.TransactOpts, reputationRegistry, identityRegistry, bas)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Comment *CommentTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Comment *CommentSession) RenounceOwnership() (*types.Transaction, error) {
	return _Comment.Contract.RenounceOwnership(&_Comment.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Comment *CommentTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Comment.Contract.RenounceOwnership(&_Comment.TransactOpts)
}

// SetBas is a paid mutator transaction binding the contract method 0x273886a8.
//
// Solidity: function setBas(address bas) returns()
func (_Comment *CommentTransactor) SetBas(opts *bind.TransactOpts, bas common.Address) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "setBas", bas)
}

// SetBas is a paid mutator transaction binding the contract method 0x273886a8.
//
// Solidity: function setBas(address bas) returns()
func (_Comment *CommentSession) SetBas(bas common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetBas(&_Comment.TransactOpts, bas)
}

// SetBas is a paid mutator transaction binding the contract method 0x273886a8.
//
// Solidity: function setBas(address bas) returns()
func (_Comment *CommentTransactorSession) SetBas(bas common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetBas(&_Comment.TransactOpts, bas)
}

// SetIdentityRegistry is a paid mutator transaction binding the contract method 0xcbf3f861.
//
// Solidity: function setIdentityRegistry(address identityRegistry) returns()
func (_Comment *CommentTransactor) SetIdentityRegistry(opts *bind.TransactOpts, identityRegistry common.Address) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "setIdentityRegistry", identityRegistry)
}

// SetIdentityRegistry is a paid mutator transaction binding the contract method 0xcbf3f861.
//
// Solidity: function setIdentityRegistry(address identityRegistry) returns()
func (_Comment *CommentSession) SetIdentityRegistry(identityRegistry common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetIdentityRegistry(&_Comment.TransactOpts, identityRegistry)
}

// SetIdentityRegistry is a paid mutator transaction binding the contract method 0xcbf3f861.
//
// Solidity: function setIdentityRegistry(address identityRegistry) returns()
func (_Comment *CommentTransactorSession) SetIdentityRegistry(identityRegistry common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetIdentityRegistry(&_Comment.TransactOpts, identityRegistry)
}

// SetReputationRegistry is a paid mutator transaction binding the contract method 0x3f0ce714.
//
// Solidity: function setReputationRegistry(address reputationRegistry) returns()
func (_Comment *CommentTransactor) SetReputationRegistry(opts *bind.TransactOpts, reputationRegistry common.Address) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "setReputationRegistry", reputationRegistry)
}

// SetReputationRegistry is a paid mutator transaction binding the contract method 0x3f0ce714.
//
// Solidity: function setReputationRegistry(address reputationRegistry) returns()
func (_Comment *CommentSession) SetReputationRegistry(reputationRegistry common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetReputationRegistry(&_Comment.TransactOpts, reputationRegistry)
}

// SetReputationRegistry is a paid mutator transaction binding the contract method 0x3f0ce714.
//
// Solidity: function setReputationRegistry(address reputationRegistry) returns()
func (_Comment *CommentTransactorSession) SetReputationRegistry(reputationRegistry common.Address) (*types.Transaction, error) {
	return _Comment.Contract.SetReputationRegistry(&_Comment.TransactOpts, reputationRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Comment *CommentTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Comment.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Comment *CommentSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Comment.Contract.TransferOwnership(&_Comment.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Comment *CommentTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Comment.Contract.TransferOwnership(&_Comment.TransactOpts, newOwner)
}

// CommentCommentEventIterator is returned from FilterCommentEvent and is used to iterate over the raw logs and unpacked data for CommentEvent events raised by the Comment contract.
type CommentCommentEventIterator struct {
	Event *CommentCommentEvent // Event containing the contract specifics and raw log

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
func (it *CommentCommentEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommentCommentEvent)
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
		it.Event = new(CommentCommentEvent)
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
func (it *CommentCommentEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommentCommentEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommentCommentEvent represents a CommentEvent event raised by the Comment contract.
type CommentCommentEvent struct {
	AgentClientId        *big.Int
	AgentServerId        *big.Int
	CommentAttestationId [32]byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterCommentEvent is a free log retrieval operation binding the contract event 0x1827f7cf157e8fff05e90ee57d7a5e67fd914e29c7113c4fa065f3f830cb5e4f.
//
// Solidity: event CommentEvent(uint256 indexed agentClientId, uint256 indexed agentServerId, bytes32 indexed commentAttestationId)
func (_Comment *CommentFilterer) FilterCommentEvent(opts *bind.FilterOpts, agentClientId []*big.Int, agentServerId []*big.Int, commentAttestationId [][32]byte) (*CommentCommentEventIterator, error) {

	var agentClientIdRule []interface{}
	for _, agentClientIdItem := range agentClientId {
		agentClientIdRule = append(agentClientIdRule, agentClientIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var commentAttestationIdRule []interface{}
	for _, commentAttestationIdItem := range commentAttestationId {
		commentAttestationIdRule = append(commentAttestationIdRule, commentAttestationIdItem)
	}

	logs, sub, err := _Comment.contract.FilterLogs(opts, "CommentEvent", agentClientIdRule, agentServerIdRule, commentAttestationIdRule)
	if err != nil {
		return nil, err
	}
	return &CommentCommentEventIterator{contract: _Comment.contract, event: "CommentEvent", logs: logs, sub: sub}, nil
}

// WatchCommentEvent is a free log subscription operation binding the contract event 0x1827f7cf157e8fff05e90ee57d7a5e67fd914e29c7113c4fa065f3f830cb5e4f.
//
// Solidity: event CommentEvent(uint256 indexed agentClientId, uint256 indexed agentServerId, bytes32 indexed commentAttestationId)
func (_Comment *CommentFilterer) WatchCommentEvent(opts *bind.WatchOpts, sink chan<- *CommentCommentEvent, agentClientId []*big.Int, agentServerId []*big.Int, commentAttestationId [][32]byte) (event.Subscription, error) {

	var agentClientIdRule []interface{}
	for _, agentClientIdItem := range agentClientId {
		agentClientIdRule = append(agentClientIdRule, agentClientIdItem)
	}
	var agentServerIdRule []interface{}
	for _, agentServerIdItem := range agentServerId {
		agentServerIdRule = append(agentServerIdRule, agentServerIdItem)
	}
	var commentAttestationIdRule []interface{}
	for _, commentAttestationIdItem := range commentAttestationId {
		commentAttestationIdRule = append(commentAttestationIdRule, commentAttestationIdItem)
	}

	logs, sub, err := _Comment.contract.WatchLogs(opts, "CommentEvent", agentClientIdRule, agentServerIdRule, commentAttestationIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommentCommentEvent)
				if err := _Comment.contract.UnpackLog(event, "CommentEvent", log); err != nil {
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

// ParseCommentEvent is a log parse operation binding the contract event 0x1827f7cf157e8fff05e90ee57d7a5e67fd914e29c7113c4fa065f3f830cb5e4f.
//
// Solidity: event CommentEvent(uint256 indexed agentClientId, uint256 indexed agentServerId, bytes32 indexed commentAttestationId)
func (_Comment *CommentFilterer) ParseCommentEvent(log types.Log) (*CommentCommentEvent, error) {
	event := new(CommentCommentEvent)
	if err := _Comment.contract.UnpackLog(event, "CommentEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommentCommentReplyEventIterator is returned from FilterCommentReplyEvent and is used to iterate over the raw logs and unpacked data for CommentReplyEvent events raised by the Comment contract.
type CommentCommentReplyEventIterator struct {
	Event *CommentCommentReplyEvent // Event containing the contract specifics and raw log

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
func (it *CommentCommentReplyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommentCommentReplyEvent)
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
		it.Event = new(CommentCommentReplyEvent)
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
func (it *CommentCommentReplyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommentCommentReplyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommentCommentReplyEvent represents a CommentReplyEvent event raised by the Comment contract.
type CommentCommentReplyEvent struct {
	CommentAttestationId      [32]byte
	CommentReplyAttestationId [32]byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterCommentReplyEvent is a free log retrieval operation binding the contract event 0xfa31746d0264cfd80a2b654a98bcc01cdaac6c48275a79928055244e37c0f87a.
//
// Solidity: event CommentReplyEvent(bytes32 indexed commentAttestationId, bytes32 indexed commentReplyAttestationId)
func (_Comment *CommentFilterer) FilterCommentReplyEvent(opts *bind.FilterOpts, commentAttestationId [][32]byte, commentReplyAttestationId [][32]byte) (*CommentCommentReplyEventIterator, error) {

	var commentAttestationIdRule []interface{}
	for _, commentAttestationIdItem := range commentAttestationId {
		commentAttestationIdRule = append(commentAttestationIdRule, commentAttestationIdItem)
	}
	var commentReplyAttestationIdRule []interface{}
	for _, commentReplyAttestationIdItem := range commentReplyAttestationId {
		commentReplyAttestationIdRule = append(commentReplyAttestationIdRule, commentReplyAttestationIdItem)
	}

	logs, sub, err := _Comment.contract.FilterLogs(opts, "CommentReplyEvent", commentAttestationIdRule, commentReplyAttestationIdRule)
	if err != nil {
		return nil, err
	}
	return &CommentCommentReplyEventIterator{contract: _Comment.contract, event: "CommentReplyEvent", logs: logs, sub: sub}, nil
}

// WatchCommentReplyEvent is a free log subscription operation binding the contract event 0xfa31746d0264cfd80a2b654a98bcc01cdaac6c48275a79928055244e37c0f87a.
//
// Solidity: event CommentReplyEvent(bytes32 indexed commentAttestationId, bytes32 indexed commentReplyAttestationId)
func (_Comment *CommentFilterer) WatchCommentReplyEvent(opts *bind.WatchOpts, sink chan<- *CommentCommentReplyEvent, commentAttestationId [][32]byte, commentReplyAttestationId [][32]byte) (event.Subscription, error) {

	var commentAttestationIdRule []interface{}
	for _, commentAttestationIdItem := range commentAttestationId {
		commentAttestationIdRule = append(commentAttestationIdRule, commentAttestationIdItem)
	}
	var commentReplyAttestationIdRule []interface{}
	for _, commentReplyAttestationIdItem := range commentReplyAttestationId {
		commentReplyAttestationIdRule = append(commentReplyAttestationIdRule, commentReplyAttestationIdItem)
	}

	logs, sub, err := _Comment.contract.WatchLogs(opts, "CommentReplyEvent", commentAttestationIdRule, commentReplyAttestationIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommentCommentReplyEvent)
				if err := _Comment.contract.UnpackLog(event, "CommentReplyEvent", log); err != nil {
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

// ParseCommentReplyEvent is a log parse operation binding the contract event 0xfa31746d0264cfd80a2b654a98bcc01cdaac6c48275a79928055244e37c0f87a.
//
// Solidity: event CommentReplyEvent(bytes32 indexed commentAttestationId, bytes32 indexed commentReplyAttestationId)
func (_Comment *CommentFilterer) ParseCommentReplyEvent(log types.Log) (*CommentCommentReplyEvent, error) {
	event := new(CommentCommentReplyEvent)
	if err := _Comment.contract.UnpackLog(event, "CommentReplyEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommentInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Comment contract.
type CommentInitializedIterator struct {
	Event *CommentInitialized // Event containing the contract specifics and raw log

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
func (it *CommentInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommentInitialized)
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
		it.Event = new(CommentInitialized)
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
func (it *CommentInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommentInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommentInitialized represents a Initialized event raised by the Comment contract.
type CommentInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Comment *CommentFilterer) FilterInitialized(opts *bind.FilterOpts) (*CommentInitializedIterator, error) {

	logs, sub, err := _Comment.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CommentInitializedIterator{contract: _Comment.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Comment *CommentFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CommentInitialized) (event.Subscription, error) {

	logs, sub, err := _Comment.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommentInitialized)
				if err := _Comment.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Comment *CommentFilterer) ParseInitialized(log types.Log) (*CommentInitialized, error) {
	event := new(CommentInitialized)
	if err := _Comment.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommentOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Comment contract.
type CommentOwnershipTransferredIterator struct {
	Event *CommentOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CommentOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommentOwnershipTransferred)
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
		it.Event = new(CommentOwnershipTransferred)
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
func (it *CommentOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommentOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommentOwnershipTransferred represents a OwnershipTransferred event raised by the Comment contract.
type CommentOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Comment *CommentFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CommentOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Comment.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CommentOwnershipTransferredIterator{contract: _Comment.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Comment *CommentFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CommentOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Comment.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommentOwnershipTransferred)
				if err := _Comment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Comment *CommentFilterer) ParseOwnershipTransferred(log types.Log) (*CommentOwnershipTransferred, error) {
	event := new(CommentOwnershipTransferred)
	if err := _Comment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
