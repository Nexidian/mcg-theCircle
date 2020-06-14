# Introduction
Just wanted to say up top that I have a lot of fun with this technical assessment. I haven't had one like this before and it was interesting to see how my work flow changed because of it. I enjoyed the open ended aspect, it allowed me to stretch out a bit and not have to follow a strict brief.

This really did feel like it could have a task assigned to me in a sprint planning meeting. If day to day work is a similar scenario to this I would be very happy to come work for you. 

This exercise has been a great opportunity to learn some new technologies. Before this task I had not had any exposure to AWS lambdas or unit tests with Go. This process has allowed me to have a deeper understanding of serverless processes, and has given me many ideas for my own future projects!

Thank you for taking the time to read through my thoughts and please do let me know if you would like clarification on any points.

# My thoughts
## App concept
I think the idea of this POC is fairly solid and throughout even at this stage. As mentioned in TechnicalChallange.md, scaling databases can be a problem and is often hard to get right the first time. I'm sure most engineers still have nightmares about that one "simple database scale" ticket they did out of hours once.

I've been reading a lot recently about database-less applications but have not worked on a solution myself, so I was pleased to have the chance to see how it could work in production.

I can see this being a useful tool for employers, but only if it was embedded into a system that they are already using, or came with a suite of other tools. There are so many services out there that allow employers to create surveys and generate reports from the data. I don't think this POC offers anything to entice users in its current offering. 

## Documentation

The provided READMEs were well written and informative. I only had to read through the initial AWS setup guide once to get a working environment. The smaller README files were also a good reference point to application specifics, and the improvement ideas were useful.

The only issue I ran into was when creating the initial bucket, containing the CFN/Lamda files, my AWS ui did not offer me the option of choosing the region. This meant the bucket was created in eu-west-2, however the root cloudformation template is hardcoded to reference eu-west-1.

```bash
Fn::Sub: https://${CodeBucket}.s3-eu-west-1.amazonaws.com/cfn/listenCaptureApiGateway.yaml
```

The README should either reflect this requirement, or `cfn/main.yaml` should be updated to be dynamic.


# My changes
I've made some changes to the existing code, and have also added a handful of improvements/extra features.
## Client side

### **Added eslint**
I've added the airbnb eslint configuration for the client javascript, and have fixed the issues it raised. Linting files has many benefits, but primarily it helps enforce unified coding standards across projects, files and developers. It also helps eliminate certain classes of bugs that could otherwise go undetected.

### **Improved UI/UX**
I have made minor tweaks to the frontend which creates a better user experience. These changes include:
- Set a `min-width` on the primary buttons
- Added some `margin` to the primary buttons so they are no longer stacked on top of eachother with no space
- Ensured all buttons were using the same style

I also refactored the multi-selected question renderer to use a series of checkboxes instead of a multi input field. This makes it easier for users to select multiple answers, and does not require the use of a keyboard to select multiple options.

```js
  const renderOption = (answer, onChange) => {
    return (
      <div key={answer.id} className="checkbox">
        <label htmlFor={answer.id}>
          <input
            id={answer.id}
            className="question-answer"
            type="checkbox"
            value={answer.text}
            onChange={onChange}
          />
          <span className="checkbox-text">{answer.text}</span>
        </label>
      </div>
    )
  }
```

## Backend

### **listenDataStore** 
I have refactored the listenDataStore lambda to consolidate data by quizId instead of questionId.

The current application process does not seem to keep track of the unique quizID except for in the path of the responseFile, which will cause issues when users want to generate reports based on a specific quiz. 

I believe a typical use case for an end user would be to want to sum answers to a specific quiz, perhaps in a graph view, to see at a glance if there are any predominant answers. For example if 80/100 employees answered with "Too many meetings" that clearly shows a problem point. 

With this in mind I think a better approach would be to adjust the process to consolidate answers keyed by the unique quizID, perhaps with a folder structure of 'bucket-name/quizID/mm-yyyy/results.json'. An example of this JSON could be similar to:

```json
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
                        "1/6-2020/0dfede7b-ca35-4ce9-a460-b3f39546a718.json"
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
                        "1/6-2020/0dfede7b-ca35-4ce9-a460-b3f39546a718.json"
                    ]
                }
            ]
        }
    ]
}
```

I think this format better captures the responses to a quiz at a glance, and would mean reading a single file when we want to generate a report, instead of the multiple reads it would take to generate currently.

### **Added Unit tests**
I have added a handful of unit tests to all of the lambdas. I have not worked with unit tests with Go before so it was fun to learn how they work and how I could integrate them. 

I have made some adjustments to some of the function definitions in order for them to be testable. 

Take the following for example:

```golang
func createFileName(queryParameters map[string]string) string {
	quizId := queryParameters["quizId"]
	year, month, day := time.Now().Date()
	fmt.Println(year, int(month), day)
	return fmt.Sprintf("%s/%d-%d/%s.json", quizId, int(month), year, genUUID())
}
```
This function is not deterministic as it relies on a random element with `genUUID`. I adjusted this function to take fileNameId which allows us to generate the UUID before calling `createFileName`. This makes the function pure and deterministic and can be tested like the following:

```go
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
```

I also had to mock some AWS functions in order to test certain functionality:

```go
type mockS3Client struct {
	s3iface.S3API
}
func (c *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
    // Testing functionality
}
```

### **Added Makefiles**
I've added Makefiles for each of the lambdas, which can be run navigating a lambda root directory and running `make`. This will automatically run the unit tests for that lambda.


## **General**
I've Added a small build script for setting up AWS infrastructure. I first set up the infrastructure manually and then then automated it. I then tore down the infrastructure and used my build script to recreate it. This script does rely on having the aws cli setup with the correct credentials so I have left it as a standalone script rather than including it within a single use build script.

# Further Improvements
- If I was able to spend more time on this project I would have set up a simple CI/CD pipeline to test, build and deploy lambda changes. 
- I would have also liked to work towards getting 100% test coverage. Testing the vital functions like `updateQuestionDataFile` and `createNewQuestionDataFile` would be very beneficial.
- I would have liked to benchmark performance metrics between Go and Nodejs lambdas
- I would have liked to spent more time developing an elegant and beautiful frontend, however I was aware that this being a POC, it does not need to have something polished.
- S3 and Dynamodb are both key-value stores. I don't have a lot of experience with dynamo, but from my research it would have a lower latency for searching and updating. Which is important for what we are doing. I would like to spend some time benchmarking performance between the two.

