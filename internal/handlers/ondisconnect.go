package handlers

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kawakattsun/meeting-room-notify/internal/repositories"
	"github.com/kawakattsun/meeting-room-notify/pkg/responder"
)

// OnDisconnect Lambda handler function.
func OnDisconnect(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	fmt.Printf("connectionId: %s\n", connectionID)

	if err := repositories.DeleteConnection(connectionID); err != nil {
		fmt.Print("error: DynamoDB DeleteConnection.\n")
		return responder.Response500(err), nil
	}

	return responder.Response200("ok"), nil
}
