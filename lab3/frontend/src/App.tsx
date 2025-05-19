import { useEffect, useState } from 'react'
import './App.css'
import { api } from './api/client'

function App() {
  const [health, setHealth] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    // api.getHealth()
    api.getTest()
      .then(response => {
        console.log(response.headers)
        setHealth(response.data);
      })
      .catch(err => {
        setError(err.message);
        console.error('API Error:', err);
      });
  }, []);

  return (
    <>
      <div>
        <h1>Proxy Status</h1>
        {health ? (
          <pre>{JSON.stringify(health, null, 2)}</pre>
        ) : error ? (
          <p style={{ color: 'red' }}>Error: {error}</p>
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </>
  )
}

export default App
