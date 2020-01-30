package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	fmt.Printf("Start lambda function. %s\n", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))
}

func main() {
	lambda.Start(handler)
}
