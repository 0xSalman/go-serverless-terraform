package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// TODO
//  1) add logging
//  2) use email template
func handler(event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	if event.TriggerSource == "CustomMessage_SignUp" {
		codeParameter := event.Request.CodeParameter
		userID := event.UserName
		link := fmt.Sprintf("<a href=\"%s?userId=%s&code=%s\" target=\"_blank\">here</a>", verificationURL, userID, codeParameter)
		event.Response.EmailSubject = "Your verification link"
		event.Response.EmailMessage = fmt.Sprintf("Thank you for signing up. Click %s to verify your email.", link)
	}
	return event, nil
}

var (
	verificationURL string
)

func init() {
	verificationURL = os.Getenv("verification_url")
}

func main() {
	lambda.Start(handler)
}
