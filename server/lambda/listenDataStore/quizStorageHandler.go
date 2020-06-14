/**
A storage handler that creates a consolidated json document for each question where the
answers are batched by question over a month time period.

Here is an example of the output:

{
    "quizId": "1",
    "questionData": [
        {
            "questionId": "0",
            "results": [
                {
                    "answerId": "0",
                    "count": 1,
                    "rawResponseFiles": [
                        "1/6-2020/eecc43d9-b437-4537-984c-8181babfceeb.json"
                    ]
                }
            ]
        },
        {
            "questionId": "1",
            "results": [
                {
                    "answerId": "2",
                    "count": 1,
                    "rawResponseFiles": [
                        "1/6-2020/eecc43d9-b437-4537-984c-8181babfceeb.json"
                    ]
                }
            ]
        }
    ]
}

This json document is stored in s3 in a folder /<quizId>/month-year/<quizId>

This means that a new file will be created every month if this quiz is used!
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type QuizResults struct {
	QuizId       string         `json:"quizId"`
	QuestionData []QuestionData `json:"questionData"`
}

type QuestionData struct {
	QuestionId      string            `json:"questionId"`
	QuestionResults []QuestionResults `json:"results"`
}

type QuestionResults struct {
	AnswerId         string   `json:"answerId"`
	Count            int32    `json:"count"`
	RawResponseFiles []string `json:"rawResponseFiles"`
}

// QuizStorageHandler is a struct that implements the DataStorageHandler Interface
type QuizStorageHandler struct {
	QuizObject QuizObject
}

// NewQuizStorageHandler is a constructor for the QuizStorageHandler struct
func NewQuizStorageHandler(quizData QuizData) *QuizStorageHandler {
	questionHandler := new(QuizStorageHandler)
	questionHandler.QuizObject = quizData.QuizData
	return questionHandler
}

func (data QuizStorageHandler) processData(dataFile string) {
	// Get the bucket name where we will store the processed data
	bucket := os.Getenv("DataStorage")

	fileName := getQuestionFileName(data.QuizObject.QuizID)

	// Moved session creation out of loadData to make it mockable
	sess := session.New()
	svc := s3.New(sess)

	// Try to load an existing file.
	questionFileBytes := loadData(bucket, fileName, svc)

	if questionFileBytes == nil {
		data.createNewQuestionDataFile(data.QuizObject.QuestionAnswers, fileName, dataFile)
	} else {
		data.updateQuestionDataFile(questionFileBytes, data.QuizObject.QuestionAnswers, fileName, dataFile)
	}

}

func getQuestionFileName(quizId string) string {
	year, month, _ := time.Now().Date()
	return fmt.Sprintf("%s/%d-%d/%s.json", quizId, int(month), year, quizId)
}

func (data QuizStorageHandler) createNewQuestionDataFile(answers []QuestionAnswers, fileName string, dataFile string) {
	quizResult := new(QuizResults)
	quizResult.QuizId = data.QuizObject.QuizID
	quizResult.QuestionData = []QuestionData{}

	for _, question := range answers {
		questionData := new(QuestionData)
		questionData.QuestionId = question.QuestionID
		questionData.QuestionResults = []QuestionResults{}

		for _, answer := range question.Answers {
			questionResult := new(QuestionResults)
			questionResult.AnswerId = answer
			questionResult.Count = 1
			questionResult.RawResponseFiles = []string{dataFile}

			questionData.QuestionResults = append(questionData.QuestionResults, *questionResult)
		}

		quizResult.QuestionData = append(quizResult.QuestionData, *questionData)
	}

	jsonData, err := json.Marshal(*quizResult)

	if err == nil {
		fmt.Println("Storing new consolidated file", string(jsonData))
		storeS3(string(jsonData), fileName)
	} else {
		exitErrorf("Unable to marshal json, %v", err)
	}
}

func (data QuizStorageHandler) updateQuestionDataFile(existingQuizResultsBytes []byte, rawQuestions []QuestionAnswers, fileName string, dataFile string) {
	existingQuizResults := parseExistingQuizResults(existingQuizResultsBytes)

	updatedQuizResults := QuizResults{
		QuizId:       existingQuizResults.QuizId,
		QuestionData: existingQuizResults.QuestionData,
	}

	for _, rawQuestion := range rawQuestions {
		// Has this question been answered already?
		hasBeenAnswered, index := hasQuestionBeenAnswered(rawQuestion, existingQuizResults.QuestionData)

		if hasBeenAnswered {
			// Question has been answered before, we need to update the existing entry
			for _, rawAnswer := range rawQuestion.Answers {
				exists, answerIndex := checkAnswerExistsInResult(rawAnswer, existingQuizResults.QuestionData[index].QuestionResults)
				if exists {
					updatedQuizResults.QuestionData[index].QuestionResults[answerIndex].Count = updatedQuizResults.QuestionData[index].QuestionResults[answerIndex].Count + 1
					updatedQuizResults.QuestionData[index].QuestionResults[answerIndex].RawResponseFiles = append(
						updatedQuizResults.QuestionData[index].QuestionResults[answerIndex].RawResponseFiles, dataFile)
				} else {
					questionResult := new(QuestionResults)
					questionResult.AnswerId = rawAnswer
					questionResult.Count = 1
					questionResult.RawResponseFiles = []string{dataFile}
					updatedQuizResults.QuestionData[index].QuestionResults = append(updatedQuizResults.QuestionData[index].QuestionResults, *questionResult)
				}
			}
		} else {
			// Question has not been answered, creating a new entry
			questionData := createNewQuestionData(rawQuestion, dataFile)
			updatedQuizResults.QuestionData = append(updatedQuizResults.QuestionData, questionData)
		}
	}

	jsonData, err := json.Marshal(updatedQuizResults)

	if err == nil {
		fmt.Println("Storing new consolidated file", string(jsonData))
		storeS3(string(jsonData), fileName)
	} else {
		exitErrorf("Unable to marshal json, %v", err)
	}

	fmt.Println("finished updateQuestionDataFile")
}

// Check the data from the existing file to see if an answer entry exists
func checkAnswerExistsInResult(answerId string, questionResults []QuestionResults) (bool, int) {
	// Has we seen this answer before?
	for index, questionResult := range questionResults {
		if answerId == questionResult.AnswerId {
			return true, index
		}
	}
	return false, 0
}

// Check the data from the existing file to see if an unique question is present
func hasQuestionBeenAnswered(rawQuestion QuestionAnswers, existingQuestionData []QuestionData) (bool, int) {
	for index, questionData := range existingQuestionData {
		if rawQuestion.QuestionID == questionData.QuestionId {
			return true, index
		}
	}

	return false, 0
}

func createNewQuestionData(rawQuestion QuestionAnswers, dataFile string) QuestionData {
	questionData := new(QuestionData)
	questionData.QuestionId = rawQuestion.QuestionID
	questionData.QuestionResults = []QuestionResults{}

	for _, answer := range rawQuestion.Answers {
		questionResult := new(QuestionResults)
		questionResult.AnswerId = answer
		questionResult.Count = 1
		questionResult.RawResponseFiles = []string{dataFile}
		questionData.QuestionResults = append(questionData.QuestionResults, *questionResult)
	}

	return *questionData
}

// ParseExistingQuestionData is a function that takes in a byte representation of QuestionData ( data formatted to analyse the responses to a question) and returns a QuestionData object
func parseExistingQuizResults(data []byte) QuizResults {
	var result QuizResults
	json.Unmarshal(data, &result)
	return result
}
