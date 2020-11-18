package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// TrackItem holds the required attributes to create a track.
// JSON definitions are required, otherwise the marshalling is messed up
// either by removing private attributes, or when converting to JSON.
type TrackItem struct {
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TrackID     string `json:"trackId"`
	CreatedAt   string
	TrackURL    string `json:"trackUrl"`
}

type TrackRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request follows the same rule described above
type Request events.APIGatewayProxyRequest

// func CreateTrack() {
// 	// await this.docClient.put({
// 	// 	TableName: this.todosTable,
// 	// 	Item: todoItem
// 	//   }).promise()

// 	//   return todoItem
// }

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	var buf bytes.Buffer

	fmt.Printf("request.Body: %s", request.Body)

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	var req TrackRequest

	json.Unmarshal([]byte(request.Body), &req)

	// TODO:
	// - generate UUID
	track := TrackItem{
		UserID:      "mock (2)",
		Name:        req.Name,
		Description: req.Description,
		TrackID:     "<uuid>",
		CreatedAt:   "date-test",
		TrackURL:    "",
	}

	fmt.Printf("TrackItem: %s", track)

	av, err := dynamodbattribute.MarshalMap(track)
	if err != nil {
		fmt.Println("Got error marshalling new track item:")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}

	// Create item in table Tracks
	tableName := "Tracks-dev"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())

		return Response{StatusCode: 400}, err
	}
	fmt.Println("Successfully added a track")

	body, _ := json.Marshal(track)
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "create-track-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
