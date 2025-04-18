export interface Link {
  name: string
  url: string
  icon: string
  category: string
  sortIndex: number
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface SortIndexUpdate {
  index: number
  sortIndex: number
  category?: string
}

export interface Config {
  enableNoAuth: boolean
  enableNoAuthView: boolean
}
