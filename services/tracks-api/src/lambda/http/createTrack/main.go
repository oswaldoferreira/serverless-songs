package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	utils "github.com/oswaldoferreira/serverless-songs/src"
	"github.com/oswaldoferreira/serverless-songs/src/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request follows the same rule described above
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	var buf bytes.Buffer
	var req services.CreateTrackRequest

	json.Unmarshal([]byte(request.Body), &req)

	userID := utils.GetUserID(request.RequestContext.Authorizer)
	req.UserID = userID

	track, err := services.CreateTrack(&req)
	if err != nil {
		fmt.Println("Got error creating the track")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	body, _ := json.Marshal(track)
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      201,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers:         utils.JSONHeader,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
