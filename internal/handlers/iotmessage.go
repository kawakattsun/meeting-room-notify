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
	doSendMessage := false
	for _, r := range event.Records {
		fmt.Printf("eventID: %s, eventName: %s, eventSourceARN: %s\n",
			r.EventID,
			r.EventName,
			r.EventSourceArn,
		)
		switch r.EventName {
		case "INSERT":
			fmt.Print("Event execute.\n")
			doSendMessage = true
			item := r.Change.NewImage
			fmt.Printf("item: %+v\n", item)
			if v, ok := item[detectedAtKey]; ok {
				fmt.Printf("delete dynamodb record. table: %s, detected_at: %s\n", iotMessageTableName, v.String())
				if err := dynamodb.Delete(iotMessageTableName, detectedAtKey, v.String()); err != nil {
					fmt.Printf("error: delete dynamodb record. %+v\n", err)
				}

			}
			if v, ok := item[sensorKey]; ok {
				sensor := v.Map()
				fmt.Printf("sensor: %+v\n", sensor["sensor"].String())
				if msg != sensorOn && sensor["sensor"].String() == sensorOn {
					fmt.Print("Detected sensor.\n")
					msg = sensorOn
				}
			}
		default:
			fmt.Print("Not executable event.\n")
		}
	}

	if doSendMessage {
		if err := sendMessage(msg); err != nil {
			fmt.Printf("error: sendMessage %s. %+v\n", msg, err)
		}
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
