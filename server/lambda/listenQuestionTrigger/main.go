package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	fmt.Printf("Triggering Question Trigger")
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		addToQueue(s3.Bucket.Name, s3.Object.Key)
	}
}

func createMessage(bucketName string, bucketKey string) string {
	return fmt.Sprintf("{\"bucketName\": \"%s\",\"key\": \"%s\"}", bucketName, bucketKey)
}

func addToQueue(bucketName string, bucketKey string) {
	queueURL := os.Getenv("AnswersFifoQueue")
	sess, err := session.NewSession()

	message := createMessage(bucketName, bucketKey)
	svc := sqs.New(sess)

	sendParams := &sqs.SendMessageInput{
		MessageBody:    aws.String(message),       // Required
		QueueUrl:       aws.String(queueURL),      // Required
		MessageGroupId: aws.String("answerGroup"), // For now we want to single thread.
	}

	sendResp, err := svc.SendMessage(sendParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", sendResp)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
