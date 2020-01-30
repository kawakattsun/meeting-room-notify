package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kawakattsun/meeting-room-notify/internal/handlers"
	"github.com/kawakattsun/meeting-room-notify/pkg/dynamodb"
)

func init() {
	fmt.Printf("Start lambda function. %s\n", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))
	dynamodb.Connect()
	session := session.Must(session.NewSession())
	svc := cloudformation.New(session, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	out, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(os.Getenv("STACK_NAME")),
	})
	if err != nil {
		fmt.Printf("error: Cloudformation DescribeStacks. %+v", err)
	}
	webSocketURIKey := os.Getenv("WEB_SOCKET_URI_KEY")

EXISTS:
	for _, s := range out.Stacks {
		for _, o := range s.Outputs {
			if *o.OutputKey == webSocketURIKey {
				handlers.SetIoTMessageWebSocketURI(*o.OutputValue)
				break EXISTS
			}
		}
	}
}

func main() {
	lambda.Start(handlers.IoTMessage)
}
