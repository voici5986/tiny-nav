import { useMainStore } from '@/stores'
import type { Link, LoginCredentials, SortIndexUpdate } from './types'

const apiBase = import.meta.env.VITE_API_BASE

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

const apiFetch = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<{ data: T; headers: Headers }> => {
  const store = useMainStore()
  const token = store.token

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: token } : {}),
    ...options.headers
  }

  const response = await fetch(`${apiBase}${endpoint}`, {
    ...options,
    headers
  })

  // 如果返回 401，清除 token
  if (response.status === 401) {
    store.logout()
    throw new ApiError(401, 'Unauthorized')
  }

  if (!response.ok) {
    throw new ApiError(response.status, `HTTP error! status: ${response.status}`)
  }

  const responseHeaders = response.headers
  let data: T
  const contentType = response.headers.get('content-type')
  if (contentType?.includes('application/json')) {
    data = await response.json()
  } else {
    data = {} as T
  }

  return { data, headers: responseHeaders }
}

export const api = {
  async login(credentials: LoginCredentials): Promise<string> {
    const { headers } = await apiFetch('/login', {
      method: 'POST',
      body: JSON.stringify(credentials)
    })

    // 从响应头获取 token
    const token = headers.get('Authorization')
    if (!token) {
      throw new Error('No token received')
    }

    return token
  },

  async getNavigation(): Promise<{ links: Link[], categories: string[] }> {
    const { data } = await apiFetch<{ links: Link[], categories: string[] }>('/navigation')
    return data
  },

  async addLink(link: Link): Promise<void> {
    apiFetch('/navigation/add', {
      method: 'POST',
      body: JSON.stringify(link)
    })
  },

  async updateLink(index: number, link: Link): Promise<void> {
    apiFetch(`/navigation/update/${index}`, {
      method: 'PUT',
      body: JSON.stringify(link)
    })
  },

  async deleteLink(index: number): Promise<void> {
    apiFetch(`/navigation/delete/${index}`, {
      method: 'DELETE'
    })
  },

  async getWebsiteIcon(url: string): Promise<string> {
    const { data } = await apiFetch<{ iconData: string }>(`/get-icon?url=${encodeURIComponent(url)}`)
    return data.iconData
  },

  async updateSortIndices(updates: SortIndexUpdate[]): Promise<void> {
    apiFetch('/navigation/sort', {
      method: 'PUT',
      body: JSON.stringify({ updates })
    })
  },

  async updateCategories(categories: string[]): Promise<void> {
    apiFetch('/navigation/categories', {
      method: 'PUT',
      body: JSON.stringify({ categories })
    })
  }
}
