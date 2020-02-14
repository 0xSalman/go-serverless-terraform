package model

import "time"

type UserService interface {
	Save(user User) error
	Update(id string, values map[string]interface{}) error
	FindByID(id string, fields ...string) (User, error)
}

type group string

const (
	ProfessorGroup group = "professor"
	StudentGroup   group = "student"
)

func UserGroup(source string) group {
	switch source {
	case "student":
		return StudentGroup
	case "professor":
		return ProfessorGroup
	}
	return ""
}

type User struct {
	ID          string    `json:"id"`
	IdentityID  string    `json:"identityId,omitempty"`
	Group       group     `json:"group"`
	FirstName   string    `json:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	PictureKey  string    `json:"pictureKey,omitempty"`
	Created     time.Time `json:"created" dynamodbav:",unixtime"`
	LastUpdated time.Time `json:"lastUpdated" dynamodbav:",unixtime"`

	Student
	Professor
}

type Student struct {
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

type Professor struct {
	Title               string `json:"title,omitempty"`
	School              string `json:"school,omitempty"`
	Department          string `json:"department,omitempty"`
	AcceptingApplicants bool   `json:"acceptingApplicants,omitempty"`
	StudentCriteria     string `json:"studentCriteria,omitempty"`
}
