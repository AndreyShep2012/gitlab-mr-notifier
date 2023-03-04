package runner

import (
	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/cron"
	"gitlab-mr-notifier/internal/interfaces"
	"log"
)

func NewLocal() interfaces.Runner {
	return new(local)
}

type local struct {
}

func (r local) Run() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("config load error %v", err)
	}

	log.Printf("start service with cron settings: period %s, time: %s", config.CronPeriod, config.CronTime)

	if config.CronPeriod == "" {
		check(config)
		return
	}

	r.runWithCron(config)
}

func (r local) runWithCron(config config.Config) {
	cron := cron.New()
	cron.Start(config.CronPeriod, config.CronTime, func() {
		check(config)
	})

	select {}
}
