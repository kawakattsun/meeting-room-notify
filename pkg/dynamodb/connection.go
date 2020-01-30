package dynamodb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var db *dynamo.DB

// Connect connect dynamoDB.
func Connect() {
	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	db = dynamo.New(session.New(), config)
}

// Put dynamoDB put.
func Put(tableName string, item interface{}) error {
	return db.Table(tableName).Put(item).Run()
}
