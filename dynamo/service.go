package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type service struct {
	dynamo    *dynamodb.DynamoDB
	tableName string
}
