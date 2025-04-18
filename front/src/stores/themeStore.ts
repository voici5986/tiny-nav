import { defineStore } from 'pinia'
import { ref, onMounted } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDarkTheme = ref(false)

  const toggleTheme = () => {
    isDarkTheme.value = !isDarkTheme.value
    document.documentElement.classList.toggle('dark', isDarkTheme.value)
  }

  const applyTheme = () => {
    document.documentElement.classList.toggle('dark', isDarkTheme.value)
  }

  // Check the saved theme when the store is created
  onMounted(() => {
    applyTheme()
  })

  return {
    isDarkTheme,
    toggleTheme,
    applyTheme
  }
}, {
  persist: true
})

