package main

/**
An interface that allows us to flexibly process all incoming feedback.

In the POC we are attempting to avoid using databases.  This means that we are storing data in JSON documents
that are sympathetic to the needs of reporting.

As we write more reports, it is pretty likely that we will need to format our data in new ways.

This is why all of the StorageHandlers (who are responsible for writing data to s3 in the correct format for reporting)
will implement this interface.

Currently there is one implementation (questionStorageHandler) which batches answers by question over a month time period.

We hope to write an answerStorageHandler soon, which will store the data in a format where we can write a report to show us
how users have answwered other questions.
*/

// DataStorageHandler is an interface that allows us to output feedback data in a flexible way.
type DataStorageHandler interface {
	processData(dataFile string)
}
