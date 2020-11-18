package main

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	// var req services.TrackRequest

	// json.Unmarshal([]byte(request.Body), &req)

	tracks, err := services.GetTracksFromUser("mock (2)")
	if err != nil {
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	body, _ := json.Marshal(tracks)
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
