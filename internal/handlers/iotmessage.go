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

var webSocketURI string

func SetIoTMessageWebSocketURI(uri string) {
	webSocketURI = uri
}

type iotMessageBody struct {
	State State `json:"state"`
}

type State struct {
	Reported Reported `json:"reported"`
}

type Reported struct {
	Sensor string `json:"sensor"`
}

// IoTMessage Lambda handler function.
func IoTMessage(event events.KinesisEvent) error {
	msg := sensorOff
	for _, r := range event.Records {
		fmt.Printf("Kinesis SequenceNumber: %+v\n", r.Kinesis.SequenceNumber)
		fmt.Printf("Kinesis Data: %s\n", string(r.Kinesis.Data))
		body := new(iotMessageBody)
		if err := json.Unmarshal(r.Kinesis.Data, body); err != nil {
			fmt.Print("error: Umnmarshal Kinesis Data\n")
			continue
		}
		if body.State.Reported.Sensor == sensorOn {
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
