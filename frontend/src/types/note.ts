export interface Note {
  id: string
  title: string
  content: string
  tags: string[]
  created_at: string
  updated_at: string
}

export interface CreateNoteData {
  title: string
  content: string
  tags: string[]
}

export interface UpdateNoteData {
  title: string
  content: string
  tags: string[]
}
