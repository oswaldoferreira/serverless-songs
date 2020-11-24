package main

import (
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
	trackID := request.PathParameters["trackId"]
	if trackID == "" {
		return Response{StatusCode: 404}, nil
	}

	userID := utils.GetUserID(request.RequestContext.Authorizer)
	var req = services.TrackRequest{
		TrackID: trackID,
		UserID:  userID,
	}
	track, err := services.GetTrack(&req)
	if err != nil {
		fmt.Println("Got error fetching the track")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	if track.TrackID == "" {
		return Response{StatusCode: 404}, err
	}

	err = services.DeleteTrack(&req)
	if err != nil {
		fmt.Println("Got error deleting the track")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers:         utils.JSONHeader,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
