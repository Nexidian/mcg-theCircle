import React, { useState, useContext } from "react"
import { RenderContext } from "../contexts/RenderContext"

import RenderButton from "./RenderButton"
import RenderMultiSelect from "./RenderMultiSelect"
import RenderText from "./RenderText"

import { setRender } from "../sharedCode"

import "../App.css"

function Quiz() {
  const { renderState } = useContext(RenderContext)
  const question = renderState.questionMap.get(renderState.currentQuestionId)
  const [state, setState] = useState(setRender(question))

  return (
    <div>
      {state.renderButton && <RenderButton nextPage={setState} />}
      {state.renderMultiSelect && <RenderMultiSelect nextPage={setState} />}
      {state.renderText && <RenderText />}
    </div>
  )
}

export default Quiz
