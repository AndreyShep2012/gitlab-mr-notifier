package utils

import "os"

func IsAWSLambda() bool {
	return os.Getenv("AWS_LAMBDA_RUNTIME_API") != ""
}
