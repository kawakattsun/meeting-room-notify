package repositories

import (
	"os"

	"github.com/kawakattsun/meeting-room-notify/pkg/dynamodb"
)

var connectionTable string

func init() {
	connectionTable = os.Getenv("CONNECTION_TABLE_NAME")
}

// ConnectionItem dynamo struct.
type ConnectionItem struct {
	ConnectionID string `dynamo:"connectionId,hash"`
}

// GetAllConnection dynamodb get all connections
func GetAllConnection() ([]*ConnectionItem, error) {
	var items []*ConnectionItem
	err := dynamodb.ScanAll(connectionTable, &items)
	return items, err
}

// PutConnection dynamodb put to connections
func PutConnection(connectionID string) error {
	item := ConnectionItem{ConnectionID: connectionID}
	return dynamodb.Put(connectionTable, item)
}

// DeleteConnection dynamodb get all connections
func DeleteConnection(connectionID string) error {
	return dynamodb.Delete(connectionTable, "connectionId", connectionID)
}
