package model

import "time"

type ConversationService interface {
	Save(conversation Conversation) error
}

type category string
type entity string

const (
	Application   category = "application"
	Direct        category = "direct"
	ThreadEntity  entity   = "thread"
	MessageEntity entity   = "message"
)

type Conversation struct {
	shared
	Category     category  `json:"category"`
	LastRead     time.Time `json:"lastRead" dynamodbav:",unixtime"`
	LastUpdated  time.Time `json:"lastUpdated" dynamodbav:"skey,unixtime"`
	Created      time.Time `json:"created" dynamodbav:",unixtime"`
	MessageCount int64     `json:"messageCount"`
}

type Message struct {
	shared
	Content string    `json:"content"`
	Created time.Time `json:"created" dynamodbav:"skey,unixtime"`
}

type shared struct {
	PKey     string `json:"_" dynamodbav:"pkey"`
	ID       string `json:"id"`
	Entity   entity `json:"entity"`
	Subject  string `json:"subject"`
	Sender   user   `json:"sender"`
	Receiver user   `json:"receiver"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
