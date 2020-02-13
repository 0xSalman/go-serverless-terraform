package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"

	"gitlab.com/rethesis/backend/errors"
)

func serverError(err errors.Request) (events.APIGatewayProxyResponse, error) {
	errorResponse := map[string]interface{}{
		"status":  err.StatusCode,
		"message": "Something went wrong",
	}
	if err.Err != nil {
		errorLogger := log.New(os.Stderr, "ERROR ", log.Llongfile)
		errorResponse["message"] = err.Error()
		errorLogger.Println(err.Err)
	}
	jsonResp, _ := json.Marshal(errorResponse)
	return apiResponse(string(jsonResp), err.StatusCode), nil
}

func apiResponse(body string, status int) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}
	if body != "" {
		resp.Body = body
	}
	return resp
}

func Error(err error) (events.APIGatewayProxyResponse, error) {
	reqError := err.(errors.Request)
	return serverError(reqError)
}

func Ok(body string) (events.APIGatewayProxyResponse, error) {
	return apiResponse(body, http.StatusOK), nil
}
