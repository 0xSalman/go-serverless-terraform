package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/rethesis/backend/dynamo"
	"gitlab.com/rethesis/backend/errors"
	"gitlab.com/rethesis/backend/transport"
	"gitlab.com/rethesis/backend/transport/user"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqJson, _ := json.Marshal(req)
	log.Printf("Received request: %s\n", string(reqJson))

	switch req.HTTPMethod {
	case "GET":
		return endpoints.GetByID(req)
	case "PUT":
		return endpoints.Update(req)
	default:
		return transport.Error(errors.Request{
			StatusCode:   405,
			Err:          fmt.Errorf("bad request: method not allowed"),
			UserFriendly: fmt.Errorf("bad request"),
		})
	}

	return events.APIGatewayProxyResponse{}, nil
}

var (
	endpoints user.Endpoints
)

func init() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: *aws.NewConfig().WithRegion("us-east-1"),
	}))
	dynamoClient := dynamodb.New(awsSession)
	userTable := os.Getenv("user_table")
	service := dynamo.NewUserService(dynamoClient, userTable)
	endpoints = user.NewEndpoints(service)
}

func main() {
	lambda.Start(handler)
}
