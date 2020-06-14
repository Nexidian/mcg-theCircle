package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func main() {
	svc := cloudformation.New(nil)

	params := &cloudformation.CreateStackInput{
		StackName: aws.String("my-go-stack"),
		Parameters: []*cloudformation.Parameter{
			&cloudformation.Parameter{
				ParameterKey:   aws.String("BucketName"),
				ParameterValue: aws.String("my-go-bucket-1710"),
			},
		},
		TemplateBody: aws.String("Template body"),
	}
}
