/**
 * A hardcoded feedback survey.
 *
 * One of the next tasks is to build a back-end that loads this data from a lambda via the api gateway.
 *
 * However, this gives a good example of what a feedback survey could look like.
 *
 * Each question has a type, text, optional next question and one or more answers that themselves can lead on to another question.
 */
exports.RenderData = {
  id: 1,
  title: "Daily Engagement Service",
  questions: [
    {
      id: 0,
      title: "Did you enjoy yourself today?",
      type: "button",
      answers: [
        {
          id: 0,
          text: "Hell Yes!",
          nextQuestionId: 1,
        },
        {
          id: 1,
          text: "Nope... :-(",
          nextQuestionId: 2,
        },
      ],
    },
    {
      id: 1,
      title: "Great what contributed to your success today?",
      type: "multiSelect",
      nextQuestionId: 4,
      answers: [
        {
          id: 2,
          text: "My colleagues.",
        },
        {
          id: 3,
          text: "I felt empowered to get my job done.",
        },
        {
          id: 4,
          text: "I had access to all of the tools I needed to smash it.",
        },
        {
          id: 5,
          text: "I could focus on my most important tasks.",
        },
        {
          id: 6,
          text: "I knew exactly what I was doing so just nailed it.",
        },
        {
          id: 7,
          text: "I don't know.  I was just on a role.",
        },
      ],
    },
    {
      id: 2,
      title: "Sorry to hear that.  Why didn't work out today?",
      type: "multiSelect",
      nextQuestionId: 3,
      answers: [
        {
          id: 8,
          text: "Too much to do.",
        },
        {
          id: 9,
          text: "Our processes got in the way.",
        },
        {
          id: 10,
          text: "Too many meetings.",
        },
        {
          id: 11,
          text: "I'm missing some skills.",
        },
        {
          id: 12,
          text: "I found my work boring.",
        },
        {
          id: 13,
          text: "I don't know I just had a bad day.",
        },
      ],
    },
    {
      id: 3,
      title: "Do you need some help?",
      type: "button",
      answers: [
        {
          id: 14,
          text: "Yes Please",
          nextQuestionId: 5,
        },
        {
          id: 15,
          text: "No thanks",
          nextQuestionId: 4,
        },
      ],
    },
    {
      id: 4,
      title: "Thanks for your input.",
      type: "text",
    },
    {
      id: 5,
      title: "Thanks for your input.  Help is on its way!",
      type: "text",
    },
  ],
}
