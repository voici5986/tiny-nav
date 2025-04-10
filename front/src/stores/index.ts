import { defineStore } from 'pinia'
import type { Link } from '@/api/types'
import { api } from '@/api'

interface State {
  token: string | null
  links: Link[]
  categories: string[]
}

export const useMainStore = defineStore('main', {
  state: (): State => ({
    token: null,
    links: [],
    categories: [],
  }),

  actions: {
    setToken(token: string | null) {
      this.token = token
    },
    setLinks(links: Link[]) {
      this.links = Array.isArray(links) ? links : []
    },
    setCategories(categories: string[]) {
      this.categories = Array.isArray(categories) ? categories : []
    },
    logout() {
      this.token = null
      this.links = []
    },
    // 验证 token 是否有效
    async validateToken(): Promise<boolean> {
      if (!this.token) return false

      try {
        // 使用 getNavigation 接口来验证 token
        await api.getNavigation()
        return true
      } catch (error) {
        // 如果token无效，清除存储的token
        this.logout()
        return false
      }
    }
  },

  persist: true
})
