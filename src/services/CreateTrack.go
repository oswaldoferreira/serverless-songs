package services

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	database "github.com/oswaldoferreira/serverless-songs/src/database"
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
	UserID          string `json:"userId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	TrackID         string `json:"trackId"`
	CreatedAt       string `json:"createdAt"`
	TrackURL        string `json:"trackUrl"`
	SignedUploadURL string `json:"signedUploadUrl"`
}

// CreateTrack is the interactor between the caller and DB
// insertion.
func CreateTrack(req *TrackRequest) (*TrackItem, error) {
	uploadItem, err := GenerateUploadURL()
	if err != nil {
		fmt.Println("Got error generating upload item")
		fmt.Println(err.Error())

		return nil, err
	}

	svc := database.NewDynamoDBClient()

	track := TrackItem{
		UserID:      "mock (2)",
		TrackID:     uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		TrackURL:    uploadItem.fileURL,
		CreatedAt:   time.Now().String(),
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

	track.SignedUploadURL = uploadItem.signedPUTUrl

	return &track, nil
}
