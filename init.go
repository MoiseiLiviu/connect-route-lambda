package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func Init(connectionID string, userId string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	tableName := os.Getenv("DYNAMO_TABLE_NAME")
	if tableName == "" {
		return fmt.Errorf("DYNAMO_TABLE_NAME environment variable not set")
	}

	dynamoSaver := NewDynamoConnectionSaver(tableName, svc)

	return dynamoSaver.SaveConnectionDetails(connectionID, userId)
}
