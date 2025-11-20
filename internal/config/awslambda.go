package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func LoadAWSLambda() (Config, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return Config{}, err
	}

	kmsClient := kms.NewFromConfig(cfg)
	encryptionContext := map[string]string{"LambdaFunctionName": os.Getenv("AWS_LAMBDA_FUNCTION_NAME")}

	gitlabToken, err := decryptAWSString(ctx, os.Getenv("GITLAB_TOKEN"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	gitlabGroupIds, err := decryptAWSString(ctx, os.Getenv("GITLAB_GROUP_IDS"), kmsClient, encryptionContext)
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

	var projectIds []int
	gitlabProjectIdsEnv := os.Getenv("GITLAB_PROJECT_IDS")
	if gitlabProjectIdsEnv != "" {
		gitlabProjectIds, err := decryptAWSString(ctx, gitlabProjectIdsEnv, kmsClient, encryptionContext)
		if err != nil {
			return Config{}, err
		}

		data = strings.Split(gitlabProjectIds, ",")

		for _, v := range data {
			id, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				return Config{}, err
			}

			projectIds = append(projectIds, id)
		}
	}

	slackWebhookURL, err := decryptAWSString(ctx, os.Getenv("SLACK_WEBHOOK_URL"), kmsClient, encryptionContext)
	if err != nil {
		return Config{}, err
	}

	c := Config{
		GitlabToken:      gitlabToken,
		GitlabGroupIDS:   groupIds,
		GitlabProjectIDS: projectIds,
		SlackWebhookURL:  slackWebhookURL,
		ShortMsgAuthors:  os.Getenv("SHORT_MESSAGE_AUTHORS"),
	}

	if err := checkRequired(c); err != nil {
		return Config{}, fmt.Errorf("check required error %v", err)
	}

	return c, nil
}

func decryptAWSString(ctx context.Context, encrypted string, kmsClient *kms.Client, encryptionContext map[string]string) (string, error) {
	if encrypted == "" {
		return "", fmt.Errorf("encrypted string is empty")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	if len(decodedBytes) == 0 {
		return "", fmt.Errorf("decoded bytes are empty")
	}

	input := &kms.DecryptInput{
		CiphertextBlob:    decodedBytes,
		EncryptionContext: encryptionContext,
	}

	response, err := kmsClient.Decrypt(ctx, input)
	if err != nil {
		return "", fmt.Errorf("KMS decrypt error: %w", err)
	}

	return string(response.Plaintext), nil
}
