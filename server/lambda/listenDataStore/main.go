package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Message is a struct that we use to parse SQS messages
type Message struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

// QuizData is a struct that we use to parse the responses from a feedback quiz.
type QuizData struct {
	QuizData QuizObject `json:"quizData"`
}

// QuizObject is a struct that stores the QuizID and answers
type QuizObject struct {
	QuizID          string            `json:"quizId"`
	QuestionAnswers []QuestionAnswers `json:"questionAnswers"`
}

// QuestionAnswers is a struct that stores the answers for a Question
type QuestionAnswers struct {
	QuestionID string   `json:"id"`
	Answers    []string `json:"answers"`
}

/**
The handler is called as a SQS event from a FIFO queue.
*/
func handler(ctx context.Context, event events.SQSEvent) {
	fmt.Println("Received new sqs event - updated")

	// Process each of the records within the queue.
	for _, record := range event.Records {
		processRecords(record)
	}
}

// Process a record from an SQS event.  The record will contain a bucket and key of a raw feedback file.
func processRecords(record events.SQSMessage) {

	body := record.Body
	message := parseSQSMessage(body)

	sess := session.New()
	svc := s3.New(sess)

	// Load the raw feedback file
	data := loadData(message.BucketName, message.Key, svc)

	if data != nil {
		// Parse the feedback data
		quizData := ParseQuizData(data)

		/**
		This code is a little more complex than it needs to be at the moment as we only have one way of procesing feedback.
		However as we write more reports, it is pretty likely that we will need to format our data in new ways.  Hence the interface.
		*/

		quizHandler := NewQuizStorageHandler(quizData)
		quizHandlers := []DataStorageHandler{quizHandler}

		for _, handler := range quizHandlers {
			handler.processData(message.Key)
		}
	} else {
		fmt.Println("This is not good - no data in process records")
	}
}

func parseSQSMessage(input string) Message {
	var message Message
	json.Unmarshal([]byte(input), &message)
	return message
}

func loadData(bucketName string, key string, svc s3iface.S3API) []byte {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return nil
	}

	s3objectBytes, err := ioutil.ReadAll(result.Body)

	return s3objectBytes
}

func storeS3(data string, filename string) {
	bucket := os.Getenv("DataStorage")

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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// ParseQuizData is a function that takes in an byte representation of the response to a quiz and turns it into a QuizData object
func ParseQuizData(data []byte) QuizData {
	var result QuizData
	json.Unmarshal(data, &result)
	return result
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
