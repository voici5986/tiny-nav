import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDarkTheme = ref(false)

  const toggleTheme = () => {
    isDarkTheme.value = !isDarkTheme.value
    document.documentElement.classList.toggle('dark', isDarkTheme.value)
  }

  return {
    isDarkTheme,
    toggleTheme,
  }
})
