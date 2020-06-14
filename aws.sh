BUCKETNAME="luke-mcg-test-bucket-01"
TEMPLATE_LOCATION="cfn/bucketCreation.yaml"
MAIN_TEMPLATE_LOCATION="cfn/main.yaml"
RAW_ANSWERS_BUCKET_NAME="mcg-raw-answers-01"

STACKNAME=`aws cloudformation create-stack \
--stack-name myteststack \
--template-body file://$TEMPLATE_LOCATION \
--parameters ParameterKey=BucketName,ParameterValue=$BUCKETNAME \
--query StackId \
--output text`

echo "Waiting for stack ${STACKNAME} to complete..."
OUTPUT=`aws cloudformation wait stack-create-complete --stack-name $STACKNAME`
if [[ $? != 0 ]]; then
    echo "Stack failed to complete."
    exit 0
else 
    echo "Stack completed successfully"
fi

echo "Attempting to upload YAML files..."
aws s3 cp cfn/ s3://$BUCKETNAME/cfn --recursive --exclude "*" --include "*.yaml"

echo "Attempting to upload source files..."
aws s3 cp dist/ s3://$BUCKETNAME/src --recursive --exclude "*" --include "*.zip"


echo "Attempting to create main stack"
MAINSTACK=`aws cloudformation create-stack \
--stack-name theCircle \
--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
--template-body file://$MAIN_TEMPLATE_LOCATION \
--parameters ParameterKey=CodeBucket,ParameterValue=$BUCKETNAME ParameterKey=AnswersBucketName,ParameterValue=$RAW_ANSWERS_BUCKET_NAME \
--query StackId \
--output text \
--debug`

echo "Waiting for stack ${MAINSTACK} to complete..."
OUTPUT=`aws cloudformation wait stack-create-complete --stack-name $MAINSTACK`
if [[ $? != 0 ]]; then
    echo "Stack failed to complete."
    exit 0
else 
    echo "Stack completed successfully"
fi
