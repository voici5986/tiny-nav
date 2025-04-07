export interface Link {
  name: string
  url: string
  icon: string
  category: string
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}
