package main

import (
	"testing"
)

func TestCreateMessage(t *testing.T) {
	cases := []struct {
		BucketName string
		BucketKey  string
		Expected   string
	}{
		{
			BucketName: "test-bucket-name-1",
			BucketKey:  "test-bucket-key-1",
			Expected:   "{\"bucketName\": \"test-bucket-name-1\",\"key\": \"test-bucket-key-1\"}",
		},
		{
			BucketName: "test-bucket-name-2",
			BucketKey:  "test-bucket-key-2",
			Expected:   "{\"bucketName\": \"test-bucket-name-2\",\"key\": \"test-bucket-key-2\"}",
		},
	}

	for _, c := range cases {
		result := createMessage(c.BucketName, c.BucketKey)
		if result != c.Expected {
			t.Errorf("Incorrect Message. Expected %s but got %s", c.Expected, result)
		}
	}
}
