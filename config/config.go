package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GitlabToken     string `env:"GITLAB_TOKEN"`
	GitlabGroupID   int    `env:"GITLAB_GROUP_ID"`
	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL"`
	CronPeriod      string `env:"CRON_PERIOD" env-default:"1d"`
	CronTime        string `env:"CRON_TIME" env-default:"10:30"`
}

func Load() Config {
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		if err := godotenv.Load(path); err != nil {
			log.Fatalf("load config file %s error %v", path, err.Error())
		}
	}

	var c Config
	if err := cleanenv.ReadEnv(&c); err != nil {
		log.Fatalf("read env error %v", err)
	}

	if err := checkRequred(c); err != nil {
		log.Fatalf("check required error %v", err)
	}

	return c
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

	if c.CronPeriod == "" {
		errs = append(errs, "CRON_PERIOD can't be empty")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}
