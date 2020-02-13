package model

import "time"

type UserService interface {
	Save(user User) error
	Update(id string, values map[string]interface{}) error
	FindByID(id string, fields ...string) (User, error)
}

type group string

const (
	Professor group = "professor"
	Student   group = "student"
)

func UserGroup(source string) group {
	switch source {
	case "student":
		return Student
	case "professor":
		return Professor
	}
	return ""
}

type User struct {
	ID          string    `json:"id"`
	IdentityID  string    `json:"identityId"`
	Group       group     `json:"group"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Created     time.Time `json:"created" dynamodbav:",unixtime"`
	LastUpdated time.Time `json:"lastUpdated" dynamodbav:",unixtime"`

	student
	professor
}

type student struct {
	Country         string `json:"country,omitempty"`
	City            string `json:"city,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	ActivelyLooking bool   `json:"activelyLooking,omitempty"`
	Resumes         []file `json:"resumes,omitempty"`
	CoverLetters    []file `json:"coverLetters,omitempty"`
	Transcripts     []file `json:"transcripts,omitempty"`
}

type file struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Primary  bool   `json:"primary"`
	Uploaded int64  `json:"uploaded"`
}

type professor struct {
	Title               string `json:"title,omitempty"`
	School              string `json:"school,omitempty"`
	Department          string `json:"department,omitempty"`
	AcceptingApplicants bool   `json:"acceptingApplicants,omitempty"`
	StudentCriteria     string `json:"studentCriteria,omitempty"`
}
