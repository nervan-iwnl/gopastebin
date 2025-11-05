// frontend/src/pages/Register.jsx
import { useState } from 'react'
import { api } from '../api'

export default function Register() {
  const [email, setEmail] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [msg, setMsg] = useState('')
  const [err, setErr] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    setErr('')
    setMsg('')
    try {
      await api.register(email, username, password)
      setMsg('registered, check email (if enabled)')
    } catch (e) {
      setErr(e?.error?.message || 'failed')
    }
  }

  return (
    <div className="card mx-auto" style={{ maxWidth: 480 }}>
      <div className="card-body">
        <h3 className="card-title mb-3">Register</h3>
        {msg && <div className="alert alert-success">{msg}</div>}
        {err && <div className="alert alert-danger">{err}</div>}
        <form onSubmit={handleSubmit}>
          <div className="mb-2">
            <input className="form-control" placeholder="email" value={email} onChange={e => setEmail(e.target.value)} />
          </div>
          <div className="mb-2">
            <input className="form-control" placeholder="username" value={username} onChange={e => setUsername(e.target.value)} />
          </div>
          <div className="mb-3">
            <input className="form-control" placeholder="password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
          </div>
          <div className="text-end">
            <button className="btn btn-primary">Sign up</button>
          </div>
        </form>
      </div>
    </div>
  )
}
