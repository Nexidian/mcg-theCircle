import { Alert, Button } from "reactstrap"
import React, { useState, useContext } from "react"
import PropTypes from "prop-types"
import { RenderContext } from "../contexts/RenderContext"

const RenderMultiSelect = (props) => {
  const [state, setState] = useState({ selected: [] })

  const { renderState, dispatch } = useContext(RenderContext)
  const question = renderState.questionMap.get(renderState.currentQuestionId)

  const renderOption = (answer, onChange) => {
    return (
      <div key={answer.id} className="checkbox">
        <label htmlFor={answer.id}>
          <input
            id={answer.id}
            className="question-answer"
            type="checkbox"
            value={answer.text}
            onChange={onChange}
          />
          <span className="checkbox-text">{answer.text}</span>
        </label>
      </div>
    )
  }

  const renderAlert = () => {
    return <Alert>Please select one or more answers from the select.</Alert>
  }

  const handleChange = () => {
    const value = []
    const checkboxes = document.getElementsByClassName("question-answer")
    // eslint-disable-next-line no-restricted-syntax
    for (const checkbox of checkboxes) {
      if (checkbox.checked) {
        value.push(checkbox.id)
      }
    }

    setState({ selected: value })
  }

  const handleClick = () => {
    if (state.selected.length > 0) {
      dispatch({
        type: "AnswerQuestion",
        update: { answer: state.selected, nextPage: props.nextPage },
      })
    } else {
      setState({ ...state, viewAlert: true })
    }
  }

  return (
    <div>
      <div>
        <b>{question.title}</b>
      </div>
      <div>
        <div id="question-answers-container">
          {question.answers.map((answer) => {
            return renderOption(answer, handleChange)
          })}
        </div>
      </div>

      <div>
        <Button
          className="mcg-button-primary"
          color="primary"
          type="button"
          onClick={() => handleClick()}
        >
          Continue
        </Button>
      </div>

      {state.viewAlert && renderAlert()}
    </div>
  )
}

RenderMultiSelect.propTypes = {
  nextPage: PropTypes.func.isRequired,
}

export default RenderMultiSelect
