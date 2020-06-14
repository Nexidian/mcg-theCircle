# theCircle - CloudFormation

theCircle is a serverless app that runs on AWS.  It is an experiment in trying to build a system that does not use databases, instead data is stored within s3 and manupulated via lambda functions.  

When working with AWS it is important to automate installations as you go along, as building automation in hindsight is often painful.

The current automation is based within three scripts.

* main.yaml - Is the main cloudFormation script.  It creates a nested stack and then calls out to the other scripts to do the work.

* listenCaptureLambda.yaml - Creates the lamda functions and all of the necessary roles, s3 buckets and queues.  It is really important to make sure that roles conform to the principle of least priviledge as security is job number ONE!!!  It feels to me like there is too much going on in this template so it should be refactored to break it down a bit.  However, a note to the wise it easy easy to get into a cloudFormation race condition when creating s3 triggers, as resources have to be created in the right order.  That's why the ListenAnswersBucket has a hard coded name, which is something that I would normally avoid.

* listenCaptureApiGateway.yaml - Sets up the API Gateway to run the listenCapture lambda function.  Maybe we should think about adding some security here.  Cognito would be a good idea?