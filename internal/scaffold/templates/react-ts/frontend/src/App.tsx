import { useState } from 'react'
// @ts-ignore
import { System } from './api'

function App() {
  const [status, setStatus] = useState("Waiting for Go backend...")
  const [echo, setEcho] = useState("")

  const checkStatus = async () => {
    try {
      const res = await System.GetStatus();
      setStatus(res);
      const e = await System.Echo("Hello from React!");
      setEcho(e);
    } catch (e) {
      setStatus("Error communicating with Go backend.");
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', fontFamily: 'sans-serif', color: '#fafafa', background: '#09090b' }}>
      <div style={{ background: 'rgba(255,255,255,0.03)', padding: '3rem', borderRadius: '16px', textAlign: 'center', border: '1px solid rgba(255,255,255,0.1)' }}>
        <h1 style={{ margin: '0 0 1rem 0' }}>Glyra UI</h1>
        <div style={{ marginBottom: '0.5rem', color: '#10b981' }}>{status}</div>
        <div style={{ marginBottom: '2rem', color: '#6366f1' }}>{echo}</div>
        <button onClick={checkStatus} style={{ background: '#3b82f6', color: 'white', padding: '1rem 2rem', border: 'none', borderRadius: '8px', cursor: 'pointer' }}>
          Ping Go Backend
        </button>
      </div>
    </div>
  )
}

export default App
