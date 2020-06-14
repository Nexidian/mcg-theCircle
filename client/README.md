# theCircle - Client - A tool to help us listen to our employees

theCircle is a POC project where we are experimenting with tools and techniques to gather and then analyse feedback from our teams.

The client is written in React.  Currently the UX is very basic as this is a POC!

The team have decided to use React Hooks, a technology that we are quite new to.

## Architecture

### Contexts - contexts dir

The render context can be thought of as the global state for the application.  It creates and exposes the render reducer.

### Reducers - reducers dir

The RenderReducer manipulates the application state and ultimately stores the feedback.  It is responsible for storing and answering a question, moving on to the next question and ultimately posting the results of the survey to the back-end.

### Components - components dir

* Quiz - The main component that is used to draw the other quiz components.

* RenderButton - Renders a button form where a user can answer a question by clicking on buttons.

* RenderMultiSelect - Renders a multi-select where a user can select multiple answers from a list of options.

* RenderText - Normally used when the user gets to the end of a feedback quiz.

## Build and run

In order to run the client through to completion the AWS back-end will need to be configured.  Please see the instructions in the main README file.

It is important to make sure that the settngs.js file is configured correctly!

Once everything is set:

```
npm install
npm run start
```
