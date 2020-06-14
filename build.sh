#!/bin/bash
# A shell script to build all of the source files and copy them to the dist directory
rm -rf dist
mkdir dist

# Build the listenCaptureLambda - This is the lambda that captures responses to questions
(cd ./server/lambda/listenCapture/ && make)
cp ./server/lambda/listenCapture/dist/listenCapture.zip ./dist

# Build the listenDataStoreLambda - This is the Lambda function that listens to the FIFO queue and stores the data into s3
(cd ./server/lambda/listenDataStore/ && make)
cp ./server/lambda/listenDataStore/dist/listenDataStore.zip ./dist

# Build the listenQuestionTrigger 
(cd ./server/lambda/listenQuestionTrigger/ && make)
cp ./server/lambda/listenQuestionTrigger/dist/listenQuestionTrigger.zip ./dist