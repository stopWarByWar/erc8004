// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "./interfaces/IReputationRegistry.sol";
import "./interfaces/IIdentityRegistry.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/**
 * @title ReputationRegistry
 * @dev Implementation of the Reputation Registry for ERC-8004 Trustless Agents
 * @notice Lightweight entry point for task feedback between agents
 * @author ChaosChain Labs
 */
contract ReputationRegistry is IReputationRegistry,Initializable,OwnableUpgradeable {
    // ============ Constants ============
    
    /// @dev Contract version for tracking implementation changes
    string public constant VERSION = "1.0.0";
    
    // ============ State Variables ============
    
    /// @dev Reference to the IdentityRegistry for agent validation
    IIdentityRegistry public  _identityRegistry;
    
    /// @dev Mapping from feedback auth ID to whether it exists
    mapping(bytes32 => bool) private _feedbackAuthorizations;
    
    /// @dev Mapping from client-server pair to feedback auth ID
    mapping(uint256 => mapping(uint256 => bytes32)) private _clientServerToAuthId;


    function initialize(IIdentityRegistry identityRegistry) public initializer{
        __Ownable_init();
       _identityRegistry =  identityRegistry;
    }

    // ============ Write Functions ============
    
    /**
     * @inheritdoc IReputationRegistry
     */
    function acceptFeedback(uint256 agentClientId, uint256 agentServerId) external {
        // Validate that both agents exist
        if (!_identityRegistry.agentExists(agentClientId)) {
            revert AgentNotFound();
        }
        if (!_identityRegistry.agentExists(agentServerId)) {
            revert AgentNotFound();
        }
        
        // Get server agent info to check authorization
        IIdentityRegistry.AgentInfo memory serverAgent = _identityRegistry.getAgent(agentServerId);
        
        // Only the server agent can authorize feedback
        if (msg.sender != serverAgent.agentAddress) {
            revert UnauthorizedFeedback();
        }
    
    // SECURITY: Prevent self-feedback to maintain integrity
        if (agentClientId == agentServerId) {
            revert SelfFeedbackNotAllowed();
        }
        
        // Check if feedback is already authorized
        bytes32 existingAuthId = _clientServerToAuthId[agentClientId][agentServerId];
        if (existingAuthId != bytes32(0)) {
            revert FeedbackAlreadyAuthorized();
        }
        
        // Generate unique feedback authorization ID
        bytes32 feedbackAuthId = _generateFeedbackAuthId(agentClientId, agentServerId);
        
        // Store the authorization
        _feedbackAuthorizations[feedbackAuthId] = true;
        _clientServerToAuthId[agentClientId][agentServerId] = feedbackAuthId;
        
        emit AuthFeedback(agentClientId, agentServerId, feedbackAuthId);
    }

    // ============ Read Functions ============
    
    /**
     * @inheritdoc IReputationRegistry
     */
    function isFeedbackAuthorized(
        uint256 agentClientId,
        uint256 agentServerId
    ) external view returns (bool isAuthorized, bytes32 feedbackAuthId) {
        feedbackAuthId = _clientServerToAuthId[agentClientId][agentServerId];
        isAuthorized = feedbackAuthId != bytes32(0) && _feedbackAuthorizations[feedbackAuthId];
    }
    
    /**
     * @inheritdoc IReputationRegistry
     */
    function getFeedbackAuthId(
        uint256 agentClientId,
        uint256 agentServerId
    ) external view returns (bytes32 feedbackAuthId) {
        feedbackAuthId = _clientServerToAuthId[agentClientId][agentServerId];
    }

    // ============ Internal Functions ============
    
    /**
     * @dev Generates a unique feedback authorization ID
     * @param agentClientId The client agent ID
     * @param agentServerId The server agent ID
     * @return feedbackAuthId The unique authorization ID
     */
    function _generateFeedbackAuthId(
        uint256 agentClientId,
        uint256 agentServerId
    ) private view returns (bytes32 feedbackAuthId) {
        // Include block timestamp, prevrandao, and sender for uniqueness
        feedbackAuthId = keccak256(
            abi.encodePacked(
                agentClientId,
                agentServerId,
                block.timestamp,
                block.prevrandao, // Use block.prevrandao for additional entropy (post-merge)
                msg.sender // Use msg.sender for meta-transaction compatibility
            )
        );
    }
}