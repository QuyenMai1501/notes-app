import { useState, useEffect, useReducer } from 'react'
import type { Note } from './types/note'
import { fetchNotes, createNote, updateNote, deleteNote } from './services/api'
import SearchBar from './components/SearchBar/SearchBar'
import NoteList from './components/NoteList/NoteList'
import NoteForm from './components/NoteForm/NoteForm'
import ConfirmDialog from './components/ConfirmDialog/ConfirmDialog'
import styles from './App.module.css'

type FetchState = {
  notes: Note[]
  loading: boolean
  error: string | null
}

type FetchAction =
  | { type: 'start' }
  | { type: 'success'; notes: Note[] }
  | { type: 'error'; error: string }

function fetchReducer(state: FetchState, action: FetchAction): FetchState {
  switch (action.type) {
    case 'start':
      return { notes: state.notes, loading: true, error: null }
    case 'success':
      return { notes: action.notes, loading: false, error: null }
    case 'error':
      return { notes: state.notes, loading: false, error: action.error }
  }
}

export default function App() {
  const [{ notes, loading, error }, dispatch] = useReducer(fetchReducer, {
    notes: [],
    loading: true,
    error: null,
  })

  const [search, setSearch] = useState('')
  const [tag, setTag] = useState('')
  const [editingNote, setEditingNote] = useState<Note | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [deletingNote, setDeletingNote] = useState<Note | null>(null)
  const [refreshKey, setRefreshKey] = useState(0)

  useEffect(() => {
    let ignore = false

    dispatch({ type: 'start' })

    fetchNotes(search, tag)
      .then((data) => {
        if (!ignore) dispatch({ type: 'success', notes: data })
      })
      .catch((err) => {
        if (!ignore) dispatch({ type: 'error', error: err instanceof Error ? err.message : 'Failed to load notes' })
      })

    return () => { ignore = true }
  }, [search, tag, refreshKey])

  async function handleFormSubmit(title: string, content: string, tags: string[]) {
    if (editingNote) {
      await updateNote(editingNote.id, { title, content, tags })
    } else {
      await createNote({ title, content, tags })
    }
    setShowForm(false)
    setEditingNote(null)
    setRefreshKey((k) => k + 1)
  }

  async function handleDeleteConfirm() {
    if (!deletingNote) return
    const id = deletingNote.id
    setDeletingNote(null)
    try {
      await deleteNote(id)
      setRefreshKey((k) => k + 1)
    } catch (err) {
      dispatch({ type: 'error', error: err instanceof Error ? err.message : 'Failed to delete note' })
    }
  }

  const allTags = [...new Set(notes.flatMap((n) => n.tags))].sort()

  return (
    <div className={styles.app}>
      <header className={styles.header}>
        <h1>Notes</h1>
        <button className={styles.createBtn} onClick={() => { setEditingNote(null); setShowForm(true) }}>
          + Create Note
        </button>
      </header>

      <SearchBar
        search={search}
        tag={tag}
        allTags={allTags}
        onSearchChange={setSearch}
        onTagChange={setTag}
      />

      <NoteList
        notes={notes}
        loading={loading}
        error={error}
        onEdit={(note) => { setEditingNote(note); setShowForm(true) }}
        onDelete={setDeletingNote}
      />

      {showForm && (
        <NoteForm
          note={editingNote}
          onSubmit={handleFormSubmit}
          onCancel={() => { setShowForm(false); setEditingNote(null) }}
        />
      )}

      {deletingNote && (
        <ConfirmDialog
          title="Delete Note"
          message={`Are you sure you want to delete "${deletingNote.title}"? This action cannot be undone.`}
          onConfirm={handleDeleteConfirm}
          onCancel={() => setDeletingNote(null)}
        />
      )}
    </div>
  )
}
