package services

import (
	"fmt"
	"os"
	"time"

	database "github.com/oswaldoferreira/serverless-songs/src"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

// TrackRequest holds the user given data
type TrackRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TrackItem holds the required attributes to create a track.
// JSON definitions are required, otherwise the marshalling is messed up
// either by removing private attributes, or when converting to JSON.
type TrackItem struct {
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TrackID     string `json:"trackId"`
	CreatedAt   string `json:"createdAt"`
	TrackURL    string `json:"trackUrl"`
}

// CreateTrack is the interactor between the caller and DB
// insertion.
func CreateTrack(req *TrackRequest) (*TrackItem, error) {
	svc := database.NewDynamoDBClient()

	track := TrackItem{
		UserID:      "mock (2)",
		Name:        req.Name,
		Description: req.Description,
		TrackID:     uuid.New().String(),
		CreatedAt:   time.Now().String(),
		TrackURL:    "",
	}

	av, err := dynamodbattribute.MarshalMap(track)
	if err != nil {
		fmt.Println("Got error marshalling new track item:")
		fmt.Println(err.Error())

		return nil, err
	}

	tableName := os.Getenv("TRACKS_TABLE")
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())

		return nil, err
	}
	fmt.Println("Successfully added a track")

	return &track, nil
}
