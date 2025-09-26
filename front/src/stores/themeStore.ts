import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

type ThemeMode = 'system' | 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  // 当前主题模式
  const mode = ref<ThemeMode>('system')
  // 实际是否为深色（根据系统 + 模式计算出来）
  const isDarkTheme = ref(false)

  // 应用主题
  const applyTheme = () => {
    if (mode.value === 'system') {
      isDarkTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      isDarkTheme.value = mode.value === 'dark'
    }
    document.documentElement.classList.toggle('dark', isDarkTheme.value)
  }

  // 切换主题模式（system → light → dark → system）
  const toggleTheme = () => {
    if (mode.value === 'system') {
      mode.value = 'light'
    } else if (mode.value === 'light') {
      mode.value = 'dark'
    } else {
      mode.value = 'system'
    }
    applyTheme()
  }

  // 监听系统主题变化
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', () => {
    if (mode.value === 'system') {
      applyTheme()
    }
  })

  // 每次 mode 改变都应用
  watch(mode, () => applyTheme(), { immediate: true })

  return {
    mode,
    isDarkTheme,
    toggleTheme,
    applyTheme,
  }
}, {
  persist: true, // 依旧支持持久化
})

