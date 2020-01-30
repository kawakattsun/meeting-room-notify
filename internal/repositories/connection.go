package repositories

import (
	"os"

	"github.com/kawakattsun/meeting-room-notify/pkg/dynamodb"
)

var connectionTable string

func init() {
	connectionTable = os.Getenv("CONNECTION_TABLE_NAME")
}

type connectionItem struct {
	ConnectionID string `dynamo:"connectionId,hash"`
}

// GetAllConnection dynamodb get all connections
func GetAllConnection() ([]string, error) {
	return nil, nil
}

// PutConnection dynamodb put to connections
func PutConnection(connectionID string) error {
	item := connectionItem{ConnectionID: connectionID}
	return dynamodb.Put(connectionTable, item)
}

// DeleteConnection dynamodb get all connections
func DeleteConnection(connectionID string) error {
	return nil
}
