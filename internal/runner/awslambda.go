package runner

import (
	"gitlab-mr-notifier/internal/config"
	"gitlab-mr-notifier/internal/interfaces"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func NewAWSLambda() interfaces.Runner {
	return new(awslambda)
}

type awslambda struct {
}

func (r awslambda) Run() {
	lambda.Start(r.run)
}

func (r awslambda) run() {
	config, err := config.LoadAWSLambda()
	if err != nil {
		log.Fatalf("config load error %v", err)
	}

	check(config)
}
