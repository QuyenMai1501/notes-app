import type { Note } from '../../types/note'
import NoteCard from '../NoteCard/NoteCard'
import LoadingSpinner from '../LoadingSpinner/LoadingSpinner'
import styles from './NoteList.module.css'

interface NoteListProps {
  notes: Note[]
  loading: boolean
  error: string | null
  onEdit: (note: Note) => void
  onDelete: (note: Note) => void
}

export default function NoteList({ notes, loading, error, onEdit, onDelete }: NoteListProps) {
  if (loading) {
    return <LoadingSpinner />
  }

  if (error) {
    return (
      <div className={styles.error}>
        <p>Failed to load notes: {error}</p>
      </div>
    )
  }

  if (notes.length === 0) {
    return (
      <div className={styles.empty}>
        <p>No notes yet. Create your first note!</p>
      </div>
    )
  }

  return (
    <div className={styles.list}>
      {notes.map((note) => (
        <NoteCard key={note.id} note={note} onEdit={onEdit} onDelete={onDelete} />
      ))}
    </div>
  )
}
