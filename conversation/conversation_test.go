package conversation

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/rethesis/backend/common"
)

type testInterface interface {
	doNothing()
}

func (c *Conversation) doNothing() {
}

func (m *Message) doNothing() {
}

func dbClient() *dynamodb.DynamoDB {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "rethesis_personal",
	}))
	return dynamodb.New(awsSession)
}

func createDBItems(items []testInterface) error {
	requestItems := make([]*dynamodb.WriteRequest, len(items))
	for i, item := range items {
		atrValMap, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return err
		}
		requestItems[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: atrValMap},
		}
	}
	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"conversation-dev": requestItems,
		},
	}

	dynamo := dbClient()
	_, err := dynamo.BatchWriteItem(batchInput)
	if err != nil {
		return err
	}

	return nil
}

func TestCreateConversation(t *testing.T) {
	user1 := user{
		ID:   "2350048295201640858",
		Name: "Salman Ahmad",
	}
	user2 := user{
		ID:   "2350048295201640999",
		Name: "Jennin Raffington",
	}
	conversationID := common.GenerateID()

	items := []testInterface{
		&Conversation{
			shared: shared{
				PKey:     user1.ID,
				ID:       conversationID,
				Entity:   ThreadEntity,
				Subject:  "Salman and Jen",
				Sender:   user1,
				Receiver: user2,
			},
			Category:     Application,
			Created:      time.Now(),
			LastRead:     time.Now(),
			LastUpdated:  time.Now(),
			MessageCount: 2,
		},
		&Conversation{
			shared: shared{
				PKey:     user2.ID,
				ID:       conversationID,
				Entity:   ThreadEntity,
				Subject:  "Salman and Jen",
				Sender:   user1,
				Receiver: user2,
			},
			Category:     Application,
			Created:      time.Now().Add(time.Minute),
			LastRead:     time.Now().Add(time.Minute),
			LastUpdated:  time.Now().Add(time.Minute),
			MessageCount: 2,
		},
		&Message{
			shared: shared{
				PKey:     conversationID,
				ID:       common.GenerateID(),
				Entity:   MessageEntity,
				Subject:  "Salman and Jen",
				Sender:   user1,
				Receiver: user2,
			},
			Created: time.Now().Add(2 * time.Minute),
			Content: "hi, this is a test message.",
		},
		&Message{
			shared: shared{
				PKey:     conversationID,
				ID:       common.GenerateID(),
				Entity:   MessageEntity,
				Subject:  "Re: Salman and Jen",
				Sender:   user2,
				Receiver: user1,
			},
			Created: time.Now().Add(time.Hour),
			Content: "hi, this is a reply to your test message.",
		},
	}
	err := createDBItems(items)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindAllConversations(t *testing.T) {
	dynamo := dbClient()
	userID := "2350048295201640858"

	result, err := dynamo.Query(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("pkey = :pkey"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pkey": {S: &userID},
		},
		TableName: aws.String("conversation-dev"),
	})
	if err != nil {
		t.Fatal(err)
	}

	conversations := make([]Conversation, len(result.Items))
	for i, item := range result.Items {
		var conv Conversation
		err := dynamodbattribute.UnmarshalMap(item, &conv)
		if err != nil {
			t.Fatal(err)
		}
		conversations[i] = conv
	}
	log.Println(conversations)
}

func TestFindAllMessages(t *testing.T) {
	dynamo := dbClient()
	convID := "2350411575581577626"

	result, err := dynamo.Query(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("pkey = :pkey"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pkey": {S: &convID},
		},
		TableName: aws.String("conversation-dev"),
	})
	if err != nil {
		t.Fatal(err)
	}

	messages := make([]Message, len(result.Items))
	for i, item := range result.Items {
		var msg Message
		err := dynamodbattribute.UnmarshalMap(item, &msg)
		if err != nil {
			t.Fatal(err)
		}
		messages[i] = msg
	}
	log.Println(messages)
}
