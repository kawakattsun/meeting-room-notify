package handlers

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kawakattsun/meeting-room-notify/internal/repositories"
	"github.com/kawakattsun/meeting-room-notify/pkg/responder"
)

// OnConnect Lambda handler function.
func OnConnect(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	fmt.Printf("connectionId: %s\n", connectionID)

	if err := repositories.PutConnection(connectionID); err != nil {
		fmt.Print("error: DynamoDB PutConnection.\n")
		return responder.Response500(err), nil
	}

	return responder.Response200("ok"), nil
}