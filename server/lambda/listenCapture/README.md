# theCircle - listenCapture

The responsibility of the listenCapture function is to store the feedback requests created by the front end.  It is a simple lambda function, written in golang that is triggered from an API proxy event and stores the data in an s3 bucket.

## Build

There is a build script (build.sh), that is integrated with the global build.sh file found in the root directory.

The build produces a zip file that can be uploaded as a lambda function into aws.

Part of me thinks that the build could be improved with a make file?  If we made that change how would that effect the global build?
