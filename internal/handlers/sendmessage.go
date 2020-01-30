package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/kawakattsun/meeting-room-notify/internal/repositories"
	"github.com/kawakattsun/meeting-room-notify/pkg/responder"
)

type sendMessageRequestBody struct {
	Message string `json:"message"`
	Data    string `data:"message"`
}

// SendMessage Lambda handler function.
func SendMessage(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := new(sendMessageRequestBody)
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		fmt.Printf("error: Umnmarshal request body. body: %+v\n", request.Body)
		return responder.Response500(err), nil
	}

	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	newSession, err := session.NewSession(config)
	if err != nil {
		fmt.Print("error: New aws session.\n")
		return responder.Response500(err), nil
	}

	svc := apigatewaymanagementapi.New(newSession)
	svc.Endpoint = fmt.Sprintf("https://%s/%s", request.RequestContext.DomainName, request.RequestContext.Stage)

	connections, err := repositories.GetAllConnection()
	if err != nil {
		fmt.Print("error: DynamoDB GetAllConnection.\n")
		return responder.Response500(err), nil
	}

	for _, connection := range connections {
		connectionID := connection.ConnectionID
		svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte(body.Data),
		})
	}

	return responder.Response200("ok"), nil
}
