import { defineStore } from 'pinia'
import type { Link, Config } from '@/api/types'
import { api } from '@/api'

interface State {
  token: string | null
  links: Link[]
  categories: string[]
  lastModified: number
  config: Config
}

export const useMainStore = defineStore('main', {
  state: (): State => ({
    token: null,
    links: [],
    categories: [],
    lastModified: 0,
    config: {
      enableNoAuth: false,
      enableNoAuthView: false,
    },
  }),

  actions: {
    setToken(token: string | null) {
      this.token = token
    },
    async fetchConfig(): Promise<Config> {
      try {
        this.config = await api.getConfig()
      } catch (error) {
        console.error('获取配置失败:', error)
      }
      return this.config
    },
    logout() {
      this.token = null
      this.links = []
    },
    // 验证 token 是否有效
    async validateToken(): Promise<boolean> {
      if (!this.token) return false

      try {
        await api.validateToken()
        return true
      } catch (error) {
        // 如果token无效，清除存储的token
        this.logout()
        return false
      }
    },
    async getNavigation(): Promise<{ links: Link[], categories: string[] }> {
      const lastModifiedResponse: { lastModified: number } = await api.getLastModified()
      if (this.lastModified === lastModifiedResponse.lastModified) {
        return {
          links: this.links,
          categories: this.categories,
        }
      }
      const { links, categories, lastModified } = await api.getNavigation()
      this.links = links
      this.categories = categories
      this.lastModified = lastModified
      return { links, categories }
    },
  },

  persist: true
})
