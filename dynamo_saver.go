package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
	"time"
)

type DynamoItem struct {
	ConnectionID string `dynamodbav:"connection_id"`
	UserID       string `dynamodbav:"user_id"`
	TTL          int64  `dynamodbav:"ttl"`
}

type DynamoConnectionSaver interface {
	SaveConnectionDetails(connectionID string, userId string) error
}

type dynamoConnectionSaver struct {
	tableName string
	dynamoSvc *dynamodb.DynamoDB
}

func NewDynamoConnectionSaver(tableName string, dynamoSvc *dynamodb.DynamoDB) DynamoConnectionSaver {
	return &dynamoConnectionSaver{
		tableName: tableName,
		dynamoSvc: dynamoSvc,
	}
}

func (s *dynamoConnectionSaver) SaveConnectionDetails(connectionID string, userId string) error {
	item := DynamoItem{
		ConnectionID: connectionID,
		UserID:       userId,
		TTL:          time.Now().Add(15 * time.Minute).Unix(),
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Error().
			Err(err).
			Interface("item", item).
			Msg("Failed to marshal item for DynamoDB")
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.tableName),
	}

	_, err = s.dynamoSvc.PutItem(input)
	if err != nil {
		log.Error().
			Err(err).
			Interface("input", input).
			Msg("Failed to publish item to DynamoDB")
		return err
	}

	return nil
}
