package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kawakattsun/meeting-room-notify/responder"
)

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	fmt.Printf("connectionId : %s ¥n", connectionID)

	err := dynamodb.Put(connectionID)
	if err != nil {
		fmt.Print("error: dynamodb not puted.\n")
		return responder.Response500(err)
	}

	fmt.Println("end on_connect")
	return responder.Response200("ok")
}