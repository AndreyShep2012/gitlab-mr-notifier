package cron

import (
	"errors"
	"gitlab-mr-notifier/interfaces"
	"strconv"
	"strings"

	"github.com/go-co-op/gocron"

	"time"
)

type cron struct {
	shed *gocron.Scheduler
}

func New() interfaces.Cron {
	return new(cron)
}

func (o *cron) Start(period, tm string, callback interfaces.Callback) error {
	o.shed = gocron.NewScheduler(time.Local).WaitForSchedule()
	var f func() error

	if tm == "" {
		f = func() error {
			_, err := o.shed.Every(period).Do(callback)
			return err
		}
	} else if strings.HasSuffix(period, "d") {
		f = func() error {
			days, err := strconv.Atoi(strings.Split(period, "d")[0])
			if err != nil {
				return err
			}
			_, err = o.shed.Every(days).Day().At(tm).Do(callback)
			return err
		}
	} else if strings.HasSuffix(period, "w") {
		f = func() error {
			weeks, err := strconv.Atoi(strings.Split(period, "w")[0])
			if err != nil {
				return err
			}
			_, err = o.shed.Every(weeks).Week().At(tm).Do(callback)
			return err
		}
	} else {
		return errors.New("not supported period " + period + ", please use 'd' or 'w' as suffix")
	}

	return startScheduler(o.shed, f)
}

func startScheduler(shed *gocron.Scheduler, f func() error) error {
	err := f()
	if err == nil {
		shed.StartAsync()
	}
	return err
}
