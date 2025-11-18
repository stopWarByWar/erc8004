package helper

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type awsCli struct {
	s3SVC      *s3.S3
	bucketName string
}

func (h Helper) UploadAgentProfileToS3(chainId, identityRegistry, agentId string, data []byte) (string, error) {
	key := fmt.Sprintf("/erc8004/agent_profile/%s/%s/%s.jpg", chainId, identityRegistry, agentId)

	_, err := h.s3SVC.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (h Helper) UploadLogoToS3(chainId, identityRegistry, agentId string, data []byte) (string, error) {
	ext := ".png" // 默认后缀
	contentType := "image/png"

	if len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		ext = ".png"
		contentType = "image/png"
	} else if len(data) >= 3 && bytes.Equal(data[:3], []byte{0xFF, 0xD8, 0xFF}) {
		ext = ".jpg"
		contentType = "image/jpeg"
	} else if len(data) >= 2 && bytes.Equal(data[:2], []byte{0x42, 0x4D}) {
		ext = ".bmp"
		contentType = "image/bmp"
	} else if len(data) >= 4 && bytes.Equal(data[:4], []byte{0x47, 0x49, 0x46, 0x38}) {
		ext = ".gif"
		contentType = "image/gif"
	}

	key := fmt.Sprintf("/erc8004/logo/%s/%s/%s%s", chainId, identityRegistry, agentId, ext)

	_, err := h.s3SVC.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(h.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ACL:         aws.String(s3.ObjectCannedACLPublicRead),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (h Helper) UploadFeedbackToS3(chainId, identityRegistry, agentId, userAddress string, indexLimit uint64, data []byte) (string, error) {
	key := fmt.Sprintf("/erc8004/feedback/%s/%s/%s/%s/%d.json", chainId, identityRegistry, agentId, userAddress, indexLimit)
	_, err := h.s3SVC.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}
