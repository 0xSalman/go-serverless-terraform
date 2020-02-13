package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/rethesis/backend/errors"
)

func TestDeleteResume(t *testing.T) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "rethesis_personal",
	}))
	dynamoClient := dynamodb.New(awsSession)
	service := NewUserService(dynamoClient, "user-dev")
	userID := "45a3f9e3-0e84-496f-9b89-dd502b8ee680"
	values := map[string]interface{}{
		"resumes": map[string]interface{}{
			"index": 0,
		},
	}

	err := service.Update(userID, values)
	if err != nil {
		t.Error((err.(errors.Request)).Err)
	}
}
