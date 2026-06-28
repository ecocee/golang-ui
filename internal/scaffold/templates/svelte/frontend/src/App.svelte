<script>
  let status = 'Fetching...'

  async function fetchStatus() {
    // window.System_GetStatus is bound by Glyra's Go ↔ JS bridge.
    status = await window.System_GetStatus()
  }

  fetchStatus()
</script>

<main>
  <h1>{{.Title}}</h1>
  <p>Your Svelte frontend is talking to a Go backend.</p>
  <p><strong>Backend says:</strong> {status}</p>
  <button on:click={fetchStatus}>Refresh</button>
</main>

<style>
  :global(body) {
    font-family: system-ui, sans-serif;
    background: #0f172a;
    color: #e2e8f0;
    margin: 0;
  }
  main {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }
  h1 { color: #22d3ee; }
  button {
    margin-top: 1rem;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    background: #22d3ee;
    color: #0f172a;
    cursor: pointer;
  }
</style>
