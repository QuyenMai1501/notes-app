import type { Note } from '../../types/note'
import styles from './NoteCard.module.css'

interface NoteCardProps {
  note: Note
  onEdit: (note: Note) => void
  onDelete: (note: Note) => void
}

export default function NoteCard({ note, onEdit, onDelete }: NoteCardProps) {
  const preview = note.content.length > 150
    ? note.content.slice(0, 150) + '...'
    : note.content

  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <h3 className={styles.title}>{note.title}</h3>
        <div className={styles.actions}>
          <button className={styles.editBtn} onClick={() => onEdit(note)}>
            Edit
          </button>
          <button className={styles.deleteBtn} onClick={() => onDelete(note)}>
            Delete
          </button>
        </div>
      </div>
      {preview && <p className={styles.content}>{preview}</p>}
      {note.tags.length > 0 && (
        <div className={styles.tags}>
          {note.tags.map((tag) => (
            <span key={tag} className={styles.tag}>{tag}</span>
          ))}
        </div>
      )}
      <div className={styles.date}>
        {new Date(note.updated_at).toLocaleDateString()}
      </div>
    </div>
  )
}
