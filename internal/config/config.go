package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GitlabToken             string `env:"GITLAB_TOKEN" env-required:"true"`
	GitlabGroupID           int    `env:"GITLAB_GROUP_ID" env-required:"true"`
	SlackWebhookURL         string `env:"SLACK_WEBHOOK_URL" env-required:"true"`
	CronPeriod              string `env:"CRON_PERIOD"`
	CronTime                string `env:"CRON_TIME"`
	MessageDescriptionLimit int    `env:"MESSAGE_DESCRIPTION_LIMIT" env-default:"500"`
	Notifier                string `env:"NOTIFIER" env-default:"slack"`
	ShortMsgAuthors         string `env:"SHORT_MESSAGE_AUTHORS"`
}

func Load() (Config, error) {
	var c Config
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return c, fmt.Errorf("load config file %s error %v", path, err.Error())
		}
	}

	if err := cleanenv.ReadEnv(&c); err != nil {
		return Config{}, fmt.Errorf("read env error %v", err)
	}

	if err := checkRequred(c); err != nil {
		return Config{}, fmt.Errorf("check required error %v", err)
	}

	return c, nil
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
