package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/go-cmp/cmp"
)

type mockS3Client struct {
	s3iface.S3API
}

type MockError struct{}

func (m *MockError) Error() string {
	return "An error occured"
}

func (c *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	// List of bucket/key values to simulate 'finding' a entry
	items := map[string]bool{
		"test-bucket1/test-key1": true,
		"test-bucket2/test-key2": true,
	}

	if items[*input.Bucket+"/"+*input.Key] {
		r := s3.GetObjectOutput{
			// Create a new IoUtils.ReadClose.
			// It does not matter what we response with here, the fact that we have something to return satisfies the test
			Body: ioutil.NopCloser(strings.NewReader("{\"quizData\":{\"quizId\":\"1\",\"questionAnswers\":[{\"id\":\"0\",\"answers\":[\"0\"]},{\"id\":\"1\",\"answers\":[\"2\",\"3\",\"4\"]}]}}")),
		}

		return &r, nil
	}

	return &s3.GetObjectOutput{}, &MockError{}
}

func TestParseSQSMessage(t *testing.T) {
	testSQSMessage := "{\"bucketName\": \"test-bucket-name\",\"key\": \"test/file/name.json\"}"

	cases := []struct {
		Resp     string
		Expected Message
	}{
		{
			Resp: testSQSMessage,
			Expected: Message{
				BucketName: "test-bucket-name",
				Key:        "test/file/name.json",
			},
		},
	}

	for _, c := range cases {
		resp := parseSQSMessage(c.Resp)
		if resp != c.Expected {
			t.Errorf(
				"Expected {%s,%s} but got {%s,%s}",
				c.Expected.BucketName,
				c.Expected.Key,
				resp.BucketName,
				resp.Key,
			)
		}
	}
}

func TestParseQuizData(t *testing.T) {
	testS3ObjectBytes := []byte("{\"quizData\":{\"quizId\":\"1\",\"questionAnswers\":[{\"id\":\"0\",\"answers\":[\"0\"]},{\"id\":\"1\",\"answers\":[\"2\",\"3\",\"4\"]}]}}")

	cases := []struct {
		Resp     []byte
		Expected QuizData
	}{
		{
			Resp: testS3ObjectBytes,
			Expected: QuizData{
				QuizObject{
					QuizID: "1",
					QuestionAnswers: []QuestionAnswers{
						{
							QuestionID: "0",
							Answers:    []string{"0"},
						},
						{
							QuestionID: "1",
							Answers:    []string{"2", "3", "4"},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		resp := ParseQuizData(c.Resp)
		if !cmp.Equal(resp, c.Expected) {
			t.Error("Did not match")
		}
	}
}

func TestLoadData(t *testing.T) {
	mockSvc := &mockS3Client{}

	cases := []struct {
		BucketName  string
		Key         string
		ExpectedNil bool
	}{
		{
			BucketName:  "test-bucket1",
			Key:         "test-key1",
			ExpectedNil: false,
		},
		{
			BucketName:  "test-bucket2",
			Key:         "test-key2",
			ExpectedNil: false,
		},
		{
			BucketName:  "incorrect-name",
			Key:         "incorrect-key",
			ExpectedNil: true,
		},
	}

	for _, c := range cases {
		resp := loadData(c.BucketName, c.Key, mockSvc)

		if c.ExpectedNil && resp != nil {
			t.Errorf("Expected NIL response got %s", resp)
		} else if !c.ExpectedNil && resp == nil {
			t.Errorf("Expected []byte response got %s", resp)
		}
	}
}
