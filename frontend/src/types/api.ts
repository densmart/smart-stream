// Auth types
export interface SignInRequest {
  login: string
  password: string
}

export interface SignInResponse {
  result: {
    access_token: string
  }
}

// Client types
export interface Client {
  id: string
  login: string
  email?: string | null
  is_active: boolean
  current_version?: string | null
  created_at: string
  last_login_at: string | null
}

export interface CreateClientRequest {
  login: string
  password: string
  email?: string
  is_active?: boolean
}

export interface UpdateClientRequest {
  login?: string
  password?: string
  email?: string
  is_active?: boolean
}

// Media types
export interface Media {
  id: string
  playlist_id?: string | null
  name: string
  format: string
  path: string
  poster: string | null
  duration: number
  size: number
  order: number
  created_at: string
  updated_at: string
}

export interface CreateMediaRequest {
  playlist_id?: string
  name: string
  format: string
  path: string
  poster?: string
  duration?: number
  size?: number
  order?: number
}

export interface UpdateMediaRequest {
  playlist_id?: string
  name?: string
  format?: string
  path?: string
  poster?: string
  duration?: number
  size?: number
  order?: number
}

// Playlist types
export enum PlaylistType {
  SERIES = 'series',
  FRANCHISE = 'franchise',
  SEASON = 'season'
}

export interface Playlist {
  id: string
  name: string
  type: PlaylistType
  poster: string | null
  parent_id?: string | null
  has_children: boolean
  created_at: string
  updated_at: string
}

export interface CreatePlaylistRequest {
  name: string
  type: PlaylistType
  poster?: string
  parent_id?: string
}

export interface UpdatePlaylistRequest {
  name?: string
  type?: PlaylistType
  poster?: string
  parent_id?: string
}

// API Response types
export interface ApiResponse<T> {
  result: T
  error?: string
}

export interface PaginatedResponse<T> {
  result: T[]
  pagination: {
    total: number
    pages: number
  }
  error?: string
}

// Error types
export interface ApiError {
  error: string
  message?: string
}

// File browser types
export type FileSystemItemType = 'file' | 'directory'

export interface FileSystemItem {
  name: string
  type: FileSystemItemType
  path: string
  size?: number
  format?: string
  modified_at?: string
}

export interface BrowseMediaFilesResponse {
  current_path: string
  parent_path?: string | null
  items: FileSystemItem[]
}