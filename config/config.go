package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"gitlab-mr-notifier/utils"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GitlabToken     string `env:"GITLAB_TOKEN" env-required:"true"`
	GitlabGroupID   int    `env:"GITLAB_GROUP_ID" env-required:"true"`
	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL" env-required:"true"`
	CronPeriod      string `env:"CRON_PERIOD"`
	CronTime        string `env:"CRON_TIME"`
}

func Load() (Config, error) {
	var c Config
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return c, fmt.Errorf("load config file %s error %v", path, err.Error())
		}
	}

	if utils.IsAWSLambda() {
		var err error
		c, err = configFromAWS()
		if err != nil {
			return Config{}, err
		}
	} else if err := cleanenv.ReadEnv(&c); err != nil {
		return Config{}, fmt.Errorf("read env error %v", err)
	}

	if err := checkRequred(c); err != nil {
		return Config{}, fmt.Errorf("check required error %v", err)
	}

	return c, nil
}

func configFromAWS() (Config, error) {
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

	return Config{
		GitlabToken:     gitlabToken,
		GitlabGroupID:   groupId,
		SlackWebhookURL: slackWebhookURL,
	}, nil
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

func checkRequred(c Config) error {
	var errs []string
	if c.GitlabToken == "" {
		errs = append(errs, "GITLAB_TOKEN can't be empty")
	}

	if c.GitlabGroupID == 0 {
		errs = append(errs, "GITLAB_GROUP_ID can't be empty")
	}

	if c.SlackWebhookURL == "" {
		errs = append(errs, "SLACK_WEBHOOK_URL can't be empty")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}
