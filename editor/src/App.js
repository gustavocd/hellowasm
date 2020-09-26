import React, { useState } from 'react'
import LiveEditor from './components/Editor'
import runGoWasm from './wasm'
import './App.css'

function App() {
  const [code, setCode] = useState('')
  const [success, setSuccess] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = value => setCode(value)
  const handleReset = () => setCode('')

  async function handleRun() {
    if (code === '') {
      setError('Please, write some code in the editor 😃, then hit Run code!')
      return
    }
    setLoading(true)

    const blob = await fetch('http://localhost:8080/execute', {
      method: 'POST',
      body: JSON.stringify({ code })
    })

    if (!blob.ok) {
      setError(await blob.text())
      return
    }

    runGoWasm(await blob.arrayBuffer())
      .then(response => {
        setSuccess(response)
        setError('')
        setLoading(false)
      })
      .catch(error => {
        setError(error)
        setSuccess([])
        setLoading(false)
      })
  }

  return (
    <div className="app">
      <LiveEditor onChange={handleChange} value={code} className="pb-1" />
      <div className="output">
        <div className="controls">
          <button onClick={handleReset} className="run" type="button">
            Reset <span role="img" aria-label="click on this to clear the editor">🗑</span>
          </button>
          <button onClick={handleRun} className="run" type="button">
            {!loading && (
              <span>
              Run code <span role="img" aria-label="click on this to run your code">▶️</span>
              </span>
            )}
            {loading && <span>Loading...</span>}
          </button>
        </div>
        {success.length > 0 && success.map(s => <p key={s}>{s}</p>)}
        {error !== '' && <p>{error}</p>}
      </div>
    </div>
  );
}

export default App;
