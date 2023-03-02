package interfaces

type Callback func()

type Cron interface {
	// peroid is defined as number with postfix, like 1s, 1m, 1h, 1d, 1w
	// time is used if period is in days or weeks, format 10:30
	Start(period, time string, callback Callback) error
}
