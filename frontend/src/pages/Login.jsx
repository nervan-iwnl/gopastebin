// frontend/src/pages/Login.jsx
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth.jsx'

export default function Login() {
  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useAuth()

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    try {
      await login(identifier, password)
      navigate('/')
    } catch (err) {
      setError(err?.error?.message || 'login failed')
    }
  }

  return (
    <div className="card mx-auto" style={{ maxWidth: 420 }}>
      <div className="card-body">
        <h3 className="card-title mb-3">Sign in</h3>
        {error && <div className="alert alert-danger">{error}</div>}
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <input
              className="form-control"
              placeholder="email or username"
              value={identifier}
              onChange={e => setIdentifier(e.target.value)}
            />
          </div>
          <div className="mb-3">
            <input
              className="form-control"
              placeholder="password"
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
            />
          </div>
          <div className="d-flex justify-content-end">
            <button className="btn btn-primary">Sign in</button>
          </div>
        </form>
      </div>
    </div>
  )
}
