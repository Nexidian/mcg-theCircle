package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGenUUID(t *testing.T) {
	result := genUUID()
	if result == "" {
		t.Errorf("Expected GUUID. Got empty string")
	}
}

func TestCreateFileName(t *testing.T) {
	queryParameters := map[string]string{"quizId": "1"}
	year, month, _ := time.Now().Date()
	guuid := genUUID()
	fileName := fmt.Sprintf("%s/%d-%d/%s.json", queryParameters["quizId"], int(month), year, guuid)

	result := createFileName(queryParameters, guuid)
	if result != fileName {
		t.Errorf("Incorrect filename. Expected %s but got %s", fileName, result)
	}
}
