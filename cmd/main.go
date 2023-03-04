package main

import (
	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/cron"
	"gitlab-mr-notifier/internal/gitlabapi"
	"gitlab-mr-notifier/internal/slack"
	"gitlab-mr-notifier/internal/utils"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

var conf config.Config

func main() {
	var err error
	conf, err = config.Load()
	if err != nil {
		log.Fatalf("config load error %v", err)
	}

	if utils.IsAWSLambda() {
		log.Println("start in AWS lambda")
		lambda.Start(runWithAWSLambda)
		return
	}

	log.Printf("start service with cron settings: period %s, time: %s", conf.CronPeriod, conf.CronTime)

	if conf.CronPeriod == "" {
		check(conf)
		return
	}

	runWithCron(conf)
}

func runWithCron(config config.Config) {
	cron := cron.New()
	cron.Start(config.CronPeriod, config.CronTime, func() {
		check(config)
	})

	select {}
}

func runWithAWSLambda() {
	check(conf)
}

func check(config config.Config) {
	sl := slack.New()
	gitapi := gitlabapi.New()

	log.Println("start checking")
	mrs, err := gitapi.GetMRList(config.GitlabToken, config.GitlabGroupID)
	if err != nil {
		log.Println("getting mr list error:", err)
		return
	}
	log.Println("found ", len(mrs), " MRs")
	err = sl.Notify(config.SlackWebhookURL, mrs)
	if err != nil {
		log.Println("slack notification error:", err)
	}
}
