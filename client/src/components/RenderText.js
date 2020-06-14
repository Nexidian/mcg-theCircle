import React, { useEffect, useContext } from "react"
import { RenderContext } from "../contexts/RenderContext"

const RenderText = () => {
  const { renderState, dispatch } = useContext(RenderContext)
  const question = renderState.questionMap.get(renderState.currentQuestionId)

  useEffect(() => {
    if (!renderState.completed) {
      console.log("Store Quiz")
      dispatch({ type: "StoreQuiz", update: { completed: true } })
    }
  }, [dispatch, renderState])

  return (
    <div>
      <b>{question.title}</b>
    </div>
  )
}

export default RenderText
