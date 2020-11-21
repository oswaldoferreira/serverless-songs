package main

import (
	"fmt"

	headers "github.com/oswaldoferreira/serverless-songs/src"
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
	// TODO: Check if both are present, otherwise fail
	var req = services.DeleteTrackRequest{
		TrackID: request.PathParameters["trackId"],
		UserID:  "mock (2)",
	}
	fmt.Println("TrackID: " + req.TrackID)

	err := services.DeleteTrack(&req)
	if err != nil {
		fmt.Println("Got error creating the track")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers:         headers.JSONHeader,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
