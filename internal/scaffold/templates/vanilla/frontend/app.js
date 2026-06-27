// JS counter — pure client-side state.
let count = 0
const jsCounter = document.getElementById('jsCounter')
const jsCounterBtn = document.getElementById('jsCounterBtn')
jsCounterBtn.addEventListener('click', () => {
  count += 1
  jsCounter.textContent = count
})

// Echo — sends a message to Go and renders the reply.
const echoInput = document.getElementById('echoInput')
const echoBtn = document.getElementById('echoBtn')
const echoReply = document.getElementById('echoReply')
echoBtn.addEventListener('click', async () => {
  const msg = echoInput.value || '(empty)'
  try {
    const reply = await window.System_Echo(msg)
    echoReply.textContent = reply
  } catch (err) {
    echoReply.textContent = 'Error: ' + err.message
    echoReply.style.color = '#ef4444'
  }
})

// Ping — Go backend health check.
const statusBox = document.getElementById('statusBox')
const pingButton = document.getElementById('pingButton')
pingButton.addEventListener('click', async () => {
  statusBox.textContent = 'Pinging…'
  try {
    const res = await window.System_GetStatus()
    statusBox.textContent = res
  } catch (err) {
    statusBox.textContent = 'Error: Cannot reach Go backend.'
    statusBox.style.color = '#ef4444'
  }
})

// Bridge health indicator.
const bridgeDot = document.getElementById('bridgeDot')
const bridgeLabel = document.getElementById('bridgeLabel')

function setBridgeConnected(ok) {
  bridgeDot.classList.toggle('connected', ok)
  bridgeLabel.textContent = ok ? 'Go bridge connected' : 'Go bridge disconnected'
}

// Probe the bridge on load. window.Echo / window.GetSystemStatus are
// injected by the Go backend via w.Bind(…).
function probeBridge() {
  const ok = typeof window.System_Echo === 'function' && typeof window.System_GetStatus === 'function'
  setBridgeConnected(ok)
  if (!ok) {
    statusBox.textContent = '⚠ Go bridge not available'
    statusBox.style.color = '#f59e0b'
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', probeBridge)
} else {
  probeBridge()
}
