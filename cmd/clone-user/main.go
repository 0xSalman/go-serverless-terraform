package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/rethesis/backend/dynamo"
	"gitlab.com/rethesis/backend/internal"
)

// TODO
//  1) improve logging
//  2) add user to a cognito/iam group: https://github.com/aws-amplify/amplify-js/issues/1213#issuecomment-432887279

func handler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	log.Printf("Confirmed user: {userName=%s, email=%s}\n", event.UserName, event.Request.UserAttributes["email"])

	// we use `nickname` for `group` value to enforce
	// required field in cognito during sign up
	// this is due to limitation set by cognito on custom attributes
	userGroup := internal.UserGroup(event.Request.UserAttributes["nickname"])
	if userGroup == "" {
		err := errors.New("user group is missing")
		log.Println(err)
		return event, err
	}

	newUser := internal.User{
		ID:          event.UserName,
		Group:       userGroup,
		Email:       event.Request.UserAttributes["email"],
		Created:     time.Now().UTC(),
		LastUpdated: time.Now().UTC(),
	}
	err := userService.Add(newUser)
	if err != nil {
		log.Println(err)
		return event, err
	}

	return event, nil
}

var (
	userService internal.UserService
)

func init() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: *aws.NewConfig().WithRegion("us-east-1"),
	}))
	dynamoClient := dynamodb.New(awsSession)
	userTable := os.Getenv("user_table")
	userService = dynamo.NewService(dynamoClient, userTable)
}

func main() {
	lambda.Start(handler)
}
