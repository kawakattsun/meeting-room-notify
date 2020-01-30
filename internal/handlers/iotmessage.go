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
)

const (
	sensorOn  = "on"
	sensorOff = "off"
)

type iotMessageBody struct {
	Sensor string `json:"sensor"`
}

// IoTMessage Lambda handler function.
func IoTMessage(event events.KinesisEvent) error {
	msg := sensorOff
	for _, r := range event.Records {
		fmt.Printf("Kinesis SequenceNumber: %+v\n", r.Kinesis.SequenceNumber)
		fmt.Printf("Kinesis Data: %+v\n", r.Kinesis.Data)
		body := new(iotMessageBody)
		if err := json.Unmarshal(r.Kinesis.Data, body); err != nil {
			fmt.Print("error: Umnmarshal Kinesis Data\n")
			continue
		}
		if body.Sensor == sensorOn {
			msg = sensorOn
			break
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
		return nil
	}

	svc := apigatewaymanagementapi.New(newSession)
	svc.Endpoint = os.Getenv("WS_ENDPOINT")

	connections, err := repositories.GetAllConnection()
	if err != nil {
		fmt.Print("error: DynamoDB GetAllConnection.\n")
		return nil
	}

	for _, connection := range connections {
		connectionID := connection.ConnectionID
		svc.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte(msg),
		})
	}

	return nil
}
