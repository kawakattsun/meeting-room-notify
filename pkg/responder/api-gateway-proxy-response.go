package responder

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func ResponseOK() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    commonHeaders(),
		Body:       `{"message":"ok"}`,
	}
}

func Response200(data string) events.APIGatewayProxyResponse {
	fmt.Print("body: " + data)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    commonHeaders(),
		Body:       data,
	}
}

func Response500(err error) events.APIGatewayProxyResponse {
	fmt.Printf("%+v\n", err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    commonHeaders(),
		Body:       `{"message":"サーバエラーが発生しました。"}`,
	}
}
