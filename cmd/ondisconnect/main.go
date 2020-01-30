package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kawakattsun/meeting-room-notify/internal/handlers"
	"github.com/kawakattsun/meeting-room-notify/pkg/dynamodb"
)

func init() {
	fmt.Printf("Start lambda function. %s\n", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))
	dynamodb.Connect()
}

func main() {
	lambda.Start(handlers.OnDisconnect)
}
