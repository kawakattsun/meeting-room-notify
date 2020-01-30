package dynamodb

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var db *dynamo.DB
var once sync.Once

// Connect connect dynamoDB.
func Connect() {
	once.Do(func() {
		config := &aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		}
		db = dynamo.New(session.New(), config)
	})
}

// ScanAll dynamoDB Scan.
func ScanAll(tableName string, items interface{}) error {
	return db.Table(tableName).Scan().All(items)
}

// Put dynamoDB put.
func Put(tableName string, item interface{}) error {
	return db.Table(tableName).Put(item).Run()
}

// Delete dynamoDB delete.
func Delete(tableName string, name string, item interface{}) error {
	return db.Table(tableName).Delete(name, item).Run()
}
