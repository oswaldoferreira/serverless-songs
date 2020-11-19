package services

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// UploadItem holds the data that will be returned to the caller of
// the service. It has both the signed URL and the future file URL.
type UploadItem struct {
	signedPUTUrl string
	fileURL      string
}

// GenerateUploadURL returns a PUT sign URL from S3 along a future uploaded
// file URL.
func GenerateUploadURL() (*UploadItem, error) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	bucketName := os.Getenv("TRACKS_S3_BUCKET")
	imageID := uuid.New().String()

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageID),
	})

	signedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("Failed to sign request", err)
		return nil, err
	}

	uploadItem := UploadItem{
		signedPUTUrl: signedURL,
		fileURL:      fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, imageID),
	}

	return &uploadItem, nil
}
