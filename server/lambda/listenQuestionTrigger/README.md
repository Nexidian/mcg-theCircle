# theCircle - listenQuestionTrigger

The responsibility of this lambda function is to listen to documents being added to s3 and to put the bucket name and s3 key on a FIFO queue for processing.

The reason why we don't process the documents in this function is that s3 triggers could happen at any time and in any order.  Given that we are trying to avoid databases (and thus are outputing our data into s3 in a format that is easy and scalable for reporting) what would happen if two clients sent data at the same time?

## Build

There is a build script (build.sh), that is integrated with the global build.sh file found in the root directory.

The build produces a zip file that can be uploaded as a lambda function into aws.

Part of me thinks that the build could be improved with a make file?  If we made that change how would that effect the global build?
