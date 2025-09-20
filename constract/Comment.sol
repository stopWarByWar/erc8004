// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "../IBAS.sol";
import "./interfaces/IReputationRegistry.sol";
import "./interfaces/IIdentityRegistry.sol";
import "../ISchemaRegistry.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol"; 
import {EMPTY_UID} from "../Common.sol";

struct CommentReplyData {
    address commenter;
    bytes32 commentAttestationId;
    string commentText;
}

struct CommentData {
    address commenter;
    uint256 agentClientId;
    uint256 agentServerId;
    uint16 score;
    string comment;
}


error InvalidScore();
error UnauthorizedFeedback();
error UnauthorizedCommentReply();
error InvalidSchema();
error NotFound();


//todo: add schema to the contract
bytes32 constant CommentSchema = 0xa803780ede4c8797ec0b1e0de04f23f3795c035e523b43e0a7145c557824dbdb;
bytes32 constant CommentReplySchema = 0x24c87dcea7b00cf6a33a983d615079a22380a1ad58f82da279ab6f5a1e01c5fe;

event CommentEvent(uint256 indexed agentClientId, uint256 indexed agentServerId, bytes32 indexed commentAttestationId);
event CommentReplyEvent(bytes32 indexed commentAttestationId, bytes32 indexed commentReplyAttestationId);

contract Comment is Initializable,OwnableUpgradeable {
    IReputationRegistry public _reputationRegistry;
    IIdentityRegistry public _identityRegistry;
    IBAS public _bas;

    function initialize(IReputationRegistry reputationRegistry, IIdentityRegistry identityRegistry, IBAS bas) public initializer {
        __Ownable_init();
        require(address(reputationRegistry) != address(0), "reputationRegistry is zero");
        require(address(identityRegistry) != address(0), "identityRegistry is zero");
        require(address(bas) != address(0), "bas is zero");

        _reputationRegistry = reputationRegistry;
        _identityRegistry = identityRegistry;
        _bas = bas;
    }

    function comment(uint256 agentClientId, uint256 agentServerId, string calldata commentText,uint16 score) external {
        if (score > 100) {
            revert InvalidScore();
        }

        (bool isFeedbackAuthorized,bytes32 feedbackAuthId) = _reputationRegistry.isFeedbackAuthorized(agentClientId, agentServerId);

        if (!isFeedbackAuthorized) {
            revert UnauthorizedFeedback();
        }


        IIdentityRegistry.AgentInfo memory agentClientInfo = _identityRegistry.getAgent(agentClientId);
        if (msg.sender != agentClientInfo.agentAddress) {
            revert UnauthorizedFeedback();
        }

        IIdentityRegistry.AgentInfo memory agentServerInfo = _identityRegistry.getAgent(agentServerId);

        bytes32 commentAttestationId = _bas.attest(AttestationRequest({
            schema: CommentSchema,
            data: AttestationRequestData({
                recipient: agentServerInfo.agentAddress,
                expirationTime: 0,
                revocable: false,
                refUID: feedbackAuthId,
                data: abi.encode(CommentData({
                    commenter: msg.sender,
                    agentClientId: agentClientId,
                    agentServerId: agentServerId,
                    score: score,
                    comment: commentText
                })),
                value: 0
            })
        }));

        emit CommentEvent(agentClientId, agentServerId, commentAttestationId);
    }

    function commentReply(bytes32 commentAttestationId,string calldata commentText) external {
        Attestation memory attestation = _bas.getAttestation(commentAttestationId);
        if (attestation.uid == EMPTY_UID) {
            revert NotFound();
        }
        if (attestation.schema != CommentReplySchema && attestation.schema != CommentSchema) {
            revert InvalidSchema();
        }

        if (attestation.recipient != msg.sender) {
            revert UnauthorizedCommentReply();
        }
        address recipient;

        if (attestation.schema == CommentSchema) {
            CommentData memory data = abi.decode(attestation.data, (CommentData));
            recipient = data.commenter;
        }

        if (attestation.schema == CommentReplySchema) {
            CommentReplyData memory data = abi.decode(attestation.data, (CommentReplyData));
            recipient = data.commenter;
        }
        
        bytes32 commentReplyAttestationId = _bas.attest(AttestationRequest({
            schema: CommentReplySchema,
            data: AttestationRequestData({
                recipient: recipient,
                expirationTime: 0,
                revocable: false,
                refUID: commentAttestationId,
                data: abi.encode(CommentReplyData({
                    commenter: msg.sender,
                    commentAttestationId: commentAttestationId,
                    commentText: commentText
                })),
                value: 0
            })
        }));

        emit CommentReplyEvent(commentAttestationId, commentReplyAttestationId);
    }

    function setReputationRegistry(IReputationRegistry reputationRegistry) external onlyOwner {
        require(address(reputationRegistry) != address(0), "reputationRegistry is zero");
        _reputationRegistry = reputationRegistry;
    }

    function setIdentityRegistry(IIdentityRegistry identityRegistry) external onlyOwner {
        require(address(identityRegistry) != address(0), "identityRegistry is zero");
        _identityRegistry = identityRegistry;
    }

    function setBas(IBAS bas) external onlyOwner {
        require(address(bas) != address(0), "bas is zero");
        _bas = bas;
    }
}