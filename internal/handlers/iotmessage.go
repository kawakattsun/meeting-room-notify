package handlers

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/kawakattsun/meeting-room-notify/internal/repositories"
	"github.com/kawakattsun/meeting-room-notify/pkg/dynamodb"
)

const (
	sensorKey     = "sensor"
	detectedAtKey = "detected_at"
	sensorOn      = "on"
	sensorOff     = "off"
)

var webSocketURI string

// SetIoTMessageWebSocketURI set IoT message web socket uri.
func SetIoTMessageWebSocketURI(uri string) {
	webSocketURI = uri
}

var iotMessageTableName string

// SetIoTMessageTableName set IoT message table name.
func SetIoTMessageTableName(name string) {
	iotMessageTableName = name
}

// IoTMessage Lambda handler function.
func IoTMessage(event events.DynamoDBEvent) error {
	msg := sensorOff
	for _, r := range event.Records {
		fmt.Printf("eventID: %s, eventName: %s, eventSourceARN: %s\n",
			r.EventID,
			r.EventName,
			r.EventSourceArn,
		)
		switch r.EventName {
		case "INSERT":
			fmt.Print("Event execute.\n")
			item := r.Change.NewImage
			if v, ok := item[detectedAtKey]; ok {
				dynamodb.Delete(iotMessageTableName, detectedAtKey, v.String())
			}
			if v, ok := item[sensorKey]; ok {
				if v.String() == sensorOn {
					msg = sensorOn
					break
				}
			}
		default:
			fmt.Print("Not executable event.\n")
		}
	}

	if err := sendMessage(msg); err != nil {
		fmt.Printf("error: sendMessage %s. %+v\n", msg, err)
	}

	return nil
}

func sendMessage(msg string) error {
	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	newSession, err := session.NewSession(config)
	if err != nil {
		fmt.Print("error: New aws session.\n")
		return err
	}

	svc := apigatewaymanagementapi.New(newSession)
	svc.Endpoint = webSocketURI

	connections, err := repositories.GetAllConnection()
	if err != nil {
		fmt.Print("error: DynamoDB GetAllConnection.\n")
		return err
	}

	for _, connection := range connections {
		connectionID := connection.ConnectionID
		_, err := svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte(fmt.Sprintf(`{"message": "%s"}`, msg)),
		})
		if err != nil {
			fmt.Printf("error: PostToConnection. %+v\n", err)
		}
	}

	return nil
}
