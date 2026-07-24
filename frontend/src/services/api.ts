import type { Note, CreateNoteData, UpdateNoteData } from '../types/note'

const BASE_URL = '/api'

class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const url = `${BASE_URL}${path}`
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  if (!res.ok) {
    let message = `Request failed with status ${res.status}`
    try {
      const body = await res.json()
      if (body.error) {
        message = body.error
      }
    } catch {
      // ignore JSON parse errors in error responses
    }
    throw new ApiError(message, res.status)
  }

  return res.json()
}

export async function fetchNotes(search?: string, tag?: string): Promise<Note[]> {
  const params = new URLSearchParams()
  if (search) params.set('search', search)
  if (tag) params.set('tag', tag)
  const query = params.toString()
  const path = query ? `/notes?${query}` : '/notes'
  return request<Note[]>(path)
}

export async function fetchNote(id: string): Promise<Note> {
  return request<Note>(`/notes/${id}`)
}

export async function createNote(data: CreateNoteData): Promise<Note> {
  return request<Note>('/notes', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateNote(id: string, data: UpdateNoteData): Promise<Note> {
  return request<Note>(`/notes/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteNote(id: string): Promise<void> {
  await request<void>(`/notes/${id}`, { method: 'DELETE' })
}

export async function checkHealth(): Promise<{ status: string }> {
  return request<{ status: string }>('/health')
}
