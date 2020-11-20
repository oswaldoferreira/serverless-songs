package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	database "github.com/oswaldoferreira/serverless-songs/src/database"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// GetTracksFromUser queries the DB for tracks from a given user.
// Is that possible to unit-test this using local dynamodb?
func GetTracksFromUser(userID string) (*[]TrackItem, error) {
	svc := database.NewDynamoDBClient()

	filter := expression.Name("userId").Equal(expression.Value(userID))

	// What I want to read from the result:
	projection := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("trackId"), expression.Name("trackUrl"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())

		return nil, err
	}
	tableName := os.Getenv("TRACKS_TABLE")
	indexName := os.Getenv("TRACKS_ID_INDEX")

	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		ScanIndexForward:          aws.Bool(false),
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
