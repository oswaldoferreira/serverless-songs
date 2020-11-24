package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	database "github.com/oswaldoferreira/serverless-songs/src/database"
)

// DeleteTrack deletes a track and returns an empty body.
func GetTrack(req *TrackRequest) (*TrackItem, error) {
	svc := database.NewDynamoDBClient()

	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())

		return nil, err
	}

	tableName := os.Getenv("TRACKS_TABLE")
	getItemInput := &dynamodb.GetItemInput{
		Key:       av,
		TableName: aws.String(tableName),
	}

	result, err := svc.GetItem(getItemInput)
	if err != nil {
		fmt.Println("Got error calling GetItem")
		fmt.Println(err.Error())

		return nil, err
	}
	var track TrackItem
	err = dynamodbattribute.UnmarshalMap(result.Item, &track)
	if err != nil {
		fmt.Println("Error unmarshaling DynamoDB GetItem result")
		fmt.Println((err.Error()))

		return nil, err
	}
	fmt.Printf("TrackID: %s", track.TrackID)

	return &track, err
}
