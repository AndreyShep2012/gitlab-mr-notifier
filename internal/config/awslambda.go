package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	gitlabGroupIds, err := decryptAWSString(os.Getenv("GITLAB_GROUP_IDS"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	var groupIds []int
	data := strings.Split(gitlabGroupIds, ",")

	for _, v := range data {
		id, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return Config{}, err
		}

		groupIds = append(groupIds, id)
	}

	gitlabProjectIds, err := decryptAWSString(os.Getenv("GITLAB_PROJECT_IDS"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	var projectIds []int
	data = strings.Split(gitlabProjectIds, ",")

	for _, v := range data {
		id, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return Config{}, err
		}

		projectIds = append(projectIds, id)
	}

	slackWebhookURL, err := decryptAWSString(os.Getenv("SLACK_WEBHOOK_URL"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	c := Config{
		GitlabToken:      gitlabToken,
		GitlabGroupIDS:   groupIds,
		GitlabProjectIDS: projectIds,
		SlackWebhookURL:  slackWebhookURL,
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
