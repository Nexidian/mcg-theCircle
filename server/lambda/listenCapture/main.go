package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	guuid "github.com/google/uuid"
)

func genUUID() string {
	return guuid.New().String()
}

func createFileName(queryParameters map[string]string, fileNameId string) string {
	quizId := queryParameters["quizId"]
	year, month, day := time.Now().Date()
	fmt.Println(year, int(month), day)
	return fmt.Sprintf("%s/%d-%d/%s.json", quizId, int(month), year, fileNameId)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")

	for key, value := range request.Headers {
		fmt.Printf("      %s: %s\n", key, value)
	}

	fmt.Println("Parameters")

	for key, value := range request.QueryStringParameters {
		fmt.Printf("     %s: %s\n", key, value)
	}

	fmt.Println("Body:")
	fmt.Printf("%s", request.Body)

	// Moved the UUID generation out of the createFileName function to make it deterministic/testable
	fileNameId := genUUID()
	filename := createFileName(request.QueryStringParameters, fileNameId)
	storeS3(request.Body, filename)

	response := events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200, Headers: make(map[string]string)}
	response.Headers["Access-Control-Allow-Origin"] = "*"
	return response, nil
}

func storeS3(data string, filename string) {
	bucket := os.Getenv("SurveyBucket")
	fmt.Printf("bucket: %s/n", bucket)

	sess, err := session.NewSession()
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(filename),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: strings.NewReader(data),
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

func main() {
	lambda.Start(handleRequest)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
