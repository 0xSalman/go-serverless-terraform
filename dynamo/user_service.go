package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/rethesis/backend/internal"
)

// TODO add logging

type serviceImpl struct {
	dynamo    *dynamodb.DynamoDB
	tableName string
}

func NewService(dynamo *dynamodb.DynamoDB, tableName string) internal.UserService {
	return serviceImpl{dynamo: dynamo, tableName: tableName}
}

func (s serviceImpl) Add(user internal.User) error {
	atrValue, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      atrValue,
		TableName: &s.tableName,
	}
	_, err = s.dynamo.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
