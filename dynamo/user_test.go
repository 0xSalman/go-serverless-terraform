package dynamo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bxcodec/faker/v3"

	"gitlab.com/rethesis/backend/errors"
	"gitlab.com/rethesis/backend/model"
)

type fake struct {
	FirstName   string `faker:"first_name"`
	LastName    string `faker:"last_name"`
	PhoneNumber string `faker:"phone_number"`
}

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

func TestCreateFakeUser(t *testing.T) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "rethesis_personal",
	}))
	cognitoClient := cognito.New(awsSession)
	dynamoClient := dynamodb.New(awsSession)
	service := NewUserService(dynamoClient, "user-dev")
	titles := []string{
		"Professor",
		"Associate Professor",
		"Assistant Professor",
		"Adjunct Professor",
		"Emeritus Professor",
	}
	schools := []string{
		"McMaster University",
		"Queen's University",
		"Ryerson University",
		"University of Guelph",
		"University of Toronto",
		"University of Waterloo",
		"Western University",
		"Wilfrid Laurier University",
		"York University",
	}
	departments := []string{
		"Biology",
		"Chemistry",
		"Computer Science",
		"Engineering",
		"Math",
		"Physics",
	}

	for i := 1; i <= 100; i++ {
		studentEmail := fmt.Sprintf("student%v@rethesis.com", i)
		profEmail := fmt.Sprintf("professor%v@rethesis.com", i)

		studentID, err := signup(cognitoClient, studentEmail, string(model.StudentGroup))
		if err != nil {
			t.Fatal(err)
		}
		profID, err := signup(cognitoClient, profEmail, string(model.ProfessorGroup))
		if err != nil {
			t.Fatal(err)
		}

		// follow https://benincosa.com/?p=3714 to fetch identityID

		student := fake{}
		err = faker.FakeData(&student)
		if err != nil {
			fmt.Println(err)
		}
		prof := fake{}
		err = faker.FakeData(&prof)
		if err != nil {
			fmt.Println(err)
		}

		titleRandIndex := randomNumber(0, 4)
		schoolRandIndex := randomNumber(0, 8)
		departRandIndex := randomNumber(0, 5)
		users := []model.User{
			{
				ID:          studentID,
				Group:       model.StudentGroup,
				FirstName:   student.FirstName,
				LastName:    student.LastName,
				Email:       studentEmail,
				PhoneNumber: student.PhoneNumber,
				Created:     time.Now(),
				LastUpdated: time.Now(),
				Student: model.Student{
					ActivelyLooking: i < 73,
					Country:         "Canada",
					City:            randomdata.City(),
					PostalCode:      randomdata.PostalCode("CA"),
				},
			},
			{
				ID:          profID,
				Group:       model.ProfessorGroup,
				FirstName:   prof.FirstName,
				LastName:    prof.LastName,
				Email:       profEmail,
				PhoneNumber: prof.PhoneNumber,
				Professor: model.Professor{
					AcceptingApplicants: i < 73,
					Title:               titles[titleRandIndex],
					School:              schools[schoolRandIndex],
					Department:          departments[departRandIndex],
				},
				Created:     time.Now(),
				LastUpdated: time.Now(),
			},
		}

		for _, usr := range users {
			err := service.Save(usr)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func signup(cognitoClient *cognito.CognitoIdentityProvider, email, group string) (string, error) {
	user := &cognito.SignUpInput{
		Username: aws.String(email),
		Password: aws.String("Test@1234"),
		ClientId: aws.String("3vn2mu5pphi2sq9busqadnvtg4"),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("nickname"),
				Value: aws.String(group),
			},
		},
	}
	ouput, err := cognitoClient.SignUp(user)
	if err != nil {
		return "", err
	}
	return *ouput.UserSub, nil
}

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func TestDeleteCognitoUsers(t *testing.T) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "rethesis_personal",
	}))
	cognitoClient := cognito.New(awsSession)

	// has pagination, need to run this multiple times for a large dataset
	output, err := cognitoClient.ListUsers(&cognito.ListUsersInput{
		UserPoolId: aws.String("us-east-1_6KyAY7Eb9"),
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, user := range output.Users {
		_, err = cognitoClient.AdminDeleteUser(&cognito.AdminDeleteUserInput{
			Username:   user.Username,
			UserPoolId: aws.String("us-east-1_6KyAY7Eb9"),
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
