package dynamo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"gitlab.com/rethesis/backend/errors"
	"gitlab.com/rethesis/backend/model"
)

type userService struct {
	service
}

func NewUserService(dynamo *dynamodb.DynamoDB, tableName string) model.UserService {
	return userService{service{dynamo: dynamo, tableName: tableName}}
}

func (s userService) Save(user model.User) error {
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

func (s userService) expression(key string) string {
	expression := strings.Builder{}
	expression.WriteString(", ")
	expression.WriteString(key)
	expression.WriteString(" = ")
	return expression.String()
}

func (s userService) Update(id string, values map[string]interface{}) error {
	exprValues := map[string]*dynamodb.AttributeValue{
		":lastUpdated": {
			N: aws.String(strconv.FormatInt(time.Now().UTC().Unix(), 10)),
		},
	}
	expression := strings.Builder{}
	expression.WriteString("set lastUpdated = :lastUpdated")
	var removeExpression string

	for key, v := range values {
		exprKey := ":" + key
		switch val := v.(type) {
		case string:
			expression.WriteString(s.expression(key))
			expression.WriteString(exprKey)
			exprValues[exprKey] = &dynamodb.AttributeValue{S: &val}
		case bool:
			expression.WriteString(s.expression(key))
			expression.WriteString(exprKey)
			exprValues[exprKey] = &dynamodb.AttributeValue{BOOL: &val}
		case map[string]interface{}:
			if index, ok := val["index"]; ok {
				// remove key[index]
				removeExpression = fmt.Sprintf("remove %s[%v]", key, index)
				break
			}

			// key = list_append(if_not_exists(key, :empty_list), :key)
			expression.WriteString(s.expression(key))
			expression.WriteString("list_append(if_not_exists(")
			expression.WriteString(key)
			expression.WriteString(", :empty_list), ")
			expression.WriteString(exprKey)
			expression.WriteString(")")

			exprValues[":empty_list"] = &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{},
			}
			atrValues, _ := dynamodbattribute.MarshalMap(val)
			exprValues[exprKey] = &dynamodb.AttributeValue{
				L: []*dynamodb.AttributeValue{{M: atrValues}},
			}
		}
	}
	expression.WriteString(removeExpression)
	log.Printf("update expression: %s\n", expression.String())

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: exprValues,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: &id,
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        &s.tableName,
		UpdateExpression: aws.String(expression.String()),
	}

	_, err := s.dynamo.UpdateItem(input)
	if err != nil {
		return errors.Request{
			StatusCode:   500,
			Err:          fmt.Errorf("failed to update user %s. %v", id, err),
			UserFriendly: fmt.Errorf("failed to update user data"),
		}
	}

	return nil
}

func (s userService) FindByID(id string, fields ...string) (model.User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: &id,
			},
		},
		TableName: &s.tableName,
	}
	if len(fields) > 0 {
		attrs := make([]*string, len(fields))
		for i, field := range fields {
			attrs[i] = &field
		}
		input.AttributesToGet = attrs
	}

	result, err := s.dynamo.GetItem(input)
	if err != nil {
		reqError := errors.Request{
			StatusCode:   404,
			Err:          fmt.Errorf("failed to find user %s. %v", id, err),
			UserFriendly: fmt.Errorf("failed to find user"),
		}
		return model.User{}, reqError
	}

	var user model.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		reqError := errors.Request{
			StatusCode:   500,
			Err:          fmt.Errorf("failed to unmarshall db response for user %s. %v", id, err),
			UserFriendly: fmt.Errorf("failed to parse data"),
		}
		return model.User{}, reqError
	}

	return user, nil
}
