package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/rethesis/backend/errors"

	"gitlab.com/rethesis/backend/model"
	"gitlab.com/rethesis/backend/transport"
)

// TODO
//  1) add logging
//  2) improve outgoing error messages

type Endpoints struct {
	service model.UserService
}

func NewEndpoints(service model.UserService) Endpoints {
	return Endpoints{service: service}
}

func (e Endpoints) loggedInID(authProvider string) string {
	i := strings.LastIndex(authProvider, ":")
	return authProvider[i+1:]
}

func (e Endpoints) GetByID(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "me" {
		id = e.loggedInID(req.RequestContext.Identity.CognitoAuthenticationProvider)
	}
	log.Printf("Get user profile for %s\n", id)

	user, err := e.service.FindByID(id)
	if err != nil {
		return transport.Error(err)
	}

	// update user identity id in the database
	// it is expected this will happen only when
	// user logs in the first time
	if user.IdentityID != req.RequestContext.Identity.CognitoIdentityID {
		user.IdentityID = req.RequestContext.Identity.CognitoIdentityID
		user.LastUpdated = time.Now().UTC()
		values := map[string]interface{}{"identityId": user.IdentityID}
		err := e.service.Update(user.ID, values)
		if err != nil {
			return transport.Error(err)
		}
		log.Printf("updated user identityId for %s\n", user.ID)
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		return transport.Error(err)
	}
	log.Printf("Found user: %s\n", jsonResp)

	return transport.Ok(string(jsonResp))
}

func (e Endpoints) Update(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	loggedInID := e.loggedInID(req.RequestContext.Identity.CognitoAuthenticationProvider)
	if id != "me" && id != loggedInID {
		return transport.Error(errors.Request{
			StatusCode:   403,
			Err:          fmt.Errorf("user %s attempted to update user %s's profile", loggedInID, id),
			UserFriendly: fmt.Errorf("you cannot update someone else's profile"),
		})
	}
	log.Printf("Update user profile for %s\n", loggedInID)

	var values map[string]interface{}
	err := json.Unmarshal([]byte(req.Body), &values)
	if err != nil {
		return transport.Error(errors.Request{
			StatusCode:   400,
			Err:          fmt.Errorf("bad request: failed to unmarshall request body. %v", req.Body),
			UserFriendly: fmt.Errorf("bad request: could not parse data"),
		})
	}
	if len(values) == 0 {
		return transport.Error(errors.Request{
			StatusCode:   400,
			Err:          fmt.Errorf("bad request: missing data in request body. %v", req.Body),
			UserFriendly: fmt.Errorf("bad request: missing data"),
		})
	}

	err = e.service.Update(loggedInID, values)
	if err != nil {
		return transport.Error(err)
	}
	log.Printf("Updated user profile for %s\n", loggedInID)

	return transport.Ok("")
}
