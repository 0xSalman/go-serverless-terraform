package model

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestAddUser(t *testing.T) {
	table := "user-dev"
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "rethesis_personal",
	}))
	dynamo := dynamodb.New(awsSession)

	users := []User{
		{
			ID:          "1111",
			Group:       Student,
			FirstName:   "Salman",
			LastName:    "Ahmad",
			Email:       "m.salman86@gmail.com",
			PhoneNumber: "647-832-839",
			Created:     time.Now(),
			LastUpdated: time.Now(),
			student: student{
				ActivelyLooking: true,
			},
		},
		{
			ID:          "2222",
			Group:       Professor,
			FirstName:   "Jennin",
			LastName:    "Raffington",
			Email:       "jennin@rethesis.com",
			PhoneNumber: "416-111-1111",
			professor: professor{
				AcceptingApplicants: true,
				Title:               "Professor",
				School:              "University of Guelph",
				Department:          "Science",
			},
			Created:     time.Now(),
			LastUpdated: time.Now(),
		},
	}

	for _, usr := range users {
		result, err := dynamodbattribute.MarshalMap(usr)
		if err != nil {
			t.Fatal(err)
		}
		input := &dynamodb.PutItemInput{
			Item:      result,
			TableName: &table,
		}
		_, err = dynamo.PutItem(input)
		if err != nil {
			t.Fatal(err)
		}
	}
}
