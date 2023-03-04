package main

import (
	"gitlab-mr-notifier/internal/interfaces"
	"gitlab-mr-notifier/internal/runner"
	"gitlab-mr-notifier/internal/utils"
)

func main() {
	var rnr interfaces.Runner
	if utils.IsAWSLambda() {
		rnr = runner.NewAWSLambda()
	} else {
		rnr = runner.NewLocal()
	}
	rnr.Run()
}
