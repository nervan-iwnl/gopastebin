// frontend/src/pages/AdminStorage.jsx
import { useEffect, useState } from 'react'
import { api } from '../api'

export default function AdminStorage() {
  const [current, setCurrent] = useState('')
  const [err, setErr] = useState('')

  async function load() {
    try {
      const data = await api.getStorage()
      setCurrent(data.storage)
    } catch (e) {
      setErr('for admin only or not logged in')
    }
  }

  useEffect(() => {
    load()
  }, [])

  async function changeTo(v) {
    try {
      await api.setStorage(v)
      setCurrent(v)
    } catch (e) {
      setErr('cannot change')
    }
  }

  return (
    <div>
      <h2>Storage settings</h2>
      {err && <p style={{ color: 'red' }}>{err}</p>}
      <p>Current: <b>{current || 'unknown'}</b></p>
      <button onClick={() => changeTo('local')}>Use local</button>
      <button onClick={() => changeTo('firebase')} style={{ marginLeft: 10 }}>Use firebase</button>
    </div>
  )
}
