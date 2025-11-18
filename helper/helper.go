package helper

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var helper *Helper

func GetHelper() *Helper {
	return helper
}

type Helper struct {
	awsCli
}

func InitHelper(s3Region, _bucketName, _awsKey, _awsSecret string) *Helper {
	_ = os.Setenv("AWS_ACCESS_KEY_ID", _awsKey)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", _awsSecret)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region), // 请根据实际情况修改区域
	})
	if err != nil {
		panic(err)
	}

	svc := s3.New(sess)

	helper = &Helper{
		awsCli: awsCli{
			s3SVC:      svc,
			bucketName: _bucketName,
		},
	}

	return helper
}
