import { Button } from "reactstrap"
import React, { useContext } from "react"
import PropTypes from "prop-types"
import { RenderContext } from "../contexts/RenderContext"

const RenderButton = (props) => {
  const { renderState, dispatch } = useContext(RenderContext)
  const question = renderState.questionMap.get(renderState.currentQuestionId)

  const handleClick = (id) => {
    dispatch({
      type: "AnswerQuestion",
      update: { answer: id, nextPage: props.nextPage },
    })
  }

  const button = (answer) => {
    return (
      <div key={answer.id}>
        <Button
          key={answer.id}
          color="primary"
          className="mcg-button-primary"
          onClick={() => {
            handleClick(answer.id)
          }}
        >
          {answer.text}
        </Button>
      </div>
    )
  }

  return (
    <div>
      <div>
        <div>
          <b>{question.title}</b>
        </div>
        <span>
          {question.answers.map((answer) => {
            return button(answer)
          })}
        </span>
      </div>
    </div>
  )
}

RenderButton.propTypes = {
  nextPage: PropTypes.func.isRequired,
}

export default RenderButton
