// eslint-disable-next-line consistent-return
const setRender = (question) => {
  if (question.type === "multiSelect") {
    return { renderButton: false, renderMultiSelect: true, renderText: false }
  }
  if (question.type === "button") {
    return { renderButton: true, renderMultiSelect: false, renderText: false }
  }
  if (question.type === "text") {
    return { renderButton: false, renderMultiSelect: false, renderText: true }
  }
}

exports.setRender = setRender
