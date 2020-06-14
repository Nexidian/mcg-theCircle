# theCircle - listenDataStore

The responsibility of this lambda function is to listen to a FIFO message queue, and process completed surveys for reporting.

As we are trying to avoid the use of databases, and make our reporting really scalable, this lambda will need to output the same data in multiple formats.

## Build

There is a build script (build.sh), that is integrated with the global build.sh file found in the root directory.

The build produces a zip file that can be uploaded as a lambda function into aws.

Part of me thinks that the build could be improved with a make file?  If we made that change how would that effect the global build?
