package user

import "time"

type group string

const (
	Professor group = "professor"
	Student   group = "student"
)

type User struct {
	ID                  string    `json:"id"`
	Group               group     `json:"group"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	Email               string    `json:"email"`
	PhoneNumber         string    `json:"phoneNumber"`
	JobStatus           string    `json:"jobStatus,omitempty"`
	Title               string    `json:"title,omitempty"`
	School              string    `json:"school,omitempty"`
	Department          string    `json:"department,omitempty"`
	AcceptingApplicants bool      `json:"acceptingApplicants,omitempty"`
	Created             time.Time `json:"created" dynamodbav:",unixtime"`
	LastUpdated         time.Time `json:"lastUpdated" dynamodbav:",unixtime"`
}
