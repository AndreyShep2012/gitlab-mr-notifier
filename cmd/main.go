package main

import (
	"gitlab-mr-notifier/config"
	"gitlab-mr-notifier/cron"
	"gitlab-mr-notifier/gitapi"
	"gitlab-mr-notifier/slack"
	"log"
)

func main() {
	config := config.Load()

	log.Printf("start service with cron settings: period %s, time: %s", config.CronPeriod, config.CronTime)

	if config.CronPeriod == "" {
		check(config)
		return
	}

	runWithCron(config)
}

func runWithCron(config config.Config) {
	cron := cron.New()
	cron.Start(config.CronPeriod, config.CronTime, func() {
		check(config)
	})

	select {}
}

func check(config config.Config) {
	sl := slack.New()
	gitapi := gitapi.New()

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
