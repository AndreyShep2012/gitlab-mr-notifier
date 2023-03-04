package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func LoadAWSLambda() (Config, error) {
	session, err := session.NewSession()
	if err != nil {
		return Config{}, err
	}

	kmsClient := kms.New(session)
	encryptionContext := aws.StringMap(map[string]string{"LambdaFunctionName": os.Getenv("AWS_LAMBDA_FUNCTION_NAME")})

	gitlabToken, err := decryptAWSString(os.Getenv("GITLAB_TOKEN"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	gitlabGroupId, err := decryptAWSString(os.Getenv("GITLAB_GROUP_ID"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	groupId, err := strconv.Atoi(gitlabGroupId)
	if err != nil {
		return Config{}, err
	}

	slackWebhookURL, err := decryptAWSString(os.Getenv("SLACK_WEBHOOK_URL"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	c := Config{
		GitlabToken:     gitlabToken,
		GitlabGroupID:   groupId,
		SlackWebhookURL: slackWebhookURL,
	}

	if err := checkRequred(c); err != nil {
		return Config{}, fmt.Errorf("check required error %v", err)
	}

	return c, nil
}

func decryptAWSString(encrypted string, kmsClient *kms.KMS, encryptionContext map[string]*string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	input := &kms.DecryptInput{
		CiphertextBlob:    decodedBytes,
		EncryptionContext: encryptionContext,
	}

	response, err := kmsClient.Decrypt(input)
	if err != nil {
		return "", err
	}

	return string(response.Plaintext), nil
}
