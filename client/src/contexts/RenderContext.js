import React, { createContext, useReducer, useEffect } from "react"
import PropTypes from "prop-types"
import { RenderData } from "../API/RenderData"
import { renderReducer } from "../reducers/RenderReducer"

export const RenderContext = createContext()

// The questions are currently hardcoded within
const createQuestionMap = (questions) => {
  return questions.reduce(
    (acc, question) => {
      acc.questions.set(question.id, question)

      if (question.answers) {
        acc.answers = question.answers.reduce((acc2, answer) => {
          return acc2.set(answer.id, answer)
        }, acc.answers)
      }

      return acc
    },
    { questions: new Map(), answers: new Map() }
  )
}

const RenderContextProvider = (props) => {
  useEffect(() => {
    console.log("We will call out to do a quiz ICL here")
  }, [])

  const questionData = createQuestionMap(RenderData.questions)

  console.log(questionData)

  const [renderState, dispatch] = useReducer(renderReducer, {
    currentQuestionId: 0,
    answers: new Map(),
    questionMap: questionData.questions,
    answerMap: questionData.answers,
    renderButton: true,
    renderDropDown: true,
    completed: false,
    quizId: RenderData.id,
  })

  const { children } = props

  return (
    <RenderContext.Provider value={{ renderState, dispatch }}>
      {children}
    </RenderContext.Provider>
  )
}

RenderContextProvider.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]).isRequired,
}

export default RenderContextProvider
