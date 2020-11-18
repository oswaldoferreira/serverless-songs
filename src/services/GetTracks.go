package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// GetTracksFromUser queries the DB for tracks from a given user.
func GetTracksFromUser(userID string) (*[]TrackItem, error) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	filter := expression.Name("userId").Equal(expression.Value(userID))

	// What I want to read from the result:
	projection := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("trackId"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())

		return nil, err
	}

	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("Tracks-dev"),
		IndexName:                 aws.String("TracksIdIndex-dev"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Query(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))

		return nil, err
	}

	var tracks []TrackItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &tracks)
	if err != nil {
		fmt.Println("Error unmarshaling DynamoDB query result")
		fmt.Println((err.Error()))

		return nil, err
	}

	return &tracks, nil
}
