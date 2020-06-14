import React from "react"
import RenderContextProvider from "./contexts/RenderContext"
import Quiz from "./components/Quiz"

import "./App.css"

function App() {
  return (
    <div className="App">
      <RenderContextProvider>
        <Quiz />
      </RenderContextProvider>
    </div>
  )
}

export default App
