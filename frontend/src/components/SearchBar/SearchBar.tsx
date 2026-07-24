import { useState, type FormEvent } from 'react'
import styles from './SearchBar.module.css'

interface SearchBarProps {
  search: string
  tag: string
  allTags: string[]
  onSearchChange: (search: string) => void
  onTagChange: (tag: string) => void
}

export default function SearchBar({ search, tag, allTags, onSearchChange, onTagChange }: SearchBarProps) {
  const [input, setInput] = useState(search)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()
    onSearchChange(input)
  }

  function handleClear() {
    setInput('')
    onSearchChange('')
  }

  return (
    <form className={styles.searchBar} onSubmit={handleSubmit}>
      <div className={styles.searchInput}>
        <input
          type="text"
          placeholder="Search notes..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        {input && (
          <button type="button" className={styles.clearBtn} onClick={handleClear}>
            &times;
          </button>
        )}
        <button type="submit" className={styles.searchBtn}>Search</button>
      </div>
      <div className={styles.tagFilter}>
        <select value={tag} onChange={(e) => onTagChange(e.target.value)}>
          <option value="">All tags</option>
          {allTags.map((t) => (
            <option key={t} value={t}>{t}</option>
          ))}
        </select>
      </div>
    </form>
  )
}
