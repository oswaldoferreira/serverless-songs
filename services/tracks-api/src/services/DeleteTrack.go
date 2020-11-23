package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	database "github.com/oswaldoferreira/serverless-songs/src/database"
)

// DeleteTrackRequest holds the user given data
type DeleteTrackRequest struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
}

// DeleteTrack deletes a track and returns an empty body.
func DeleteTrack(req *DeleteTrackRequest) error {
	svc := database.NewDynamoDBClient()

	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())

		return err
	}

	tableName := os.Getenv("TRACKS_TABLE")
	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())

		return err
	}

	return nil
}
