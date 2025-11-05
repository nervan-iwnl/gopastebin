// frontend/src/components/PasteForm.jsx
import { useState } from 'react'
import { api } from '../api'

export default function PasteForm({ onCreated }) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [folder, setFolder] = useState('')
  const [slug, setSlug] = useState('')
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setSaving(true)
    try {
      await api.createPaste({
        title,
        content,
        extension: 'txt',
        folder,
        slug,
        is_public: true,
      })
      setTitle('')
      setContent('')
      setFolder('')
      setSlug('')
      onCreated && onCreated()
    } catch (e) {
      setError(e?.error?.message || 'failed')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="card mb-4">
      <div className="card-body">
        <h5 className="card-title">Create paste</h5>
        {error && <div className="alert alert-danger">{error}</div>}
        <form onSubmit={handleSubmit}>
          <div className="mb-2">
            <input className="form-control" placeholder="title" value={title} onChange={e => setTitle(e.target.value)} />
          </div>
          <div className="mb-2">
            <input className="form-control" placeholder="folder (like project-euler/001-100)" value={folder} onChange={e => setFolder(e.target.value)} />
          </div>
          <div className="mb-2">
            <input className="form-control" placeholder="custom slug (optional)" value={slug} onChange={e => setSlug(e.target.value)} />
          </div>
          <div className="mb-2">
            <textarea className="form-control" rows={8} placeholder="content" value={content} onChange={e => setContent(e.target.value)} />
          </div>
          <div className="text-end">
            <button className="btn btn-primary" disabled={saving}>{saving ? 'Creating...' : 'Create'}</button>
          </div>
        </form>
      </div>
    </div>
  )
}
