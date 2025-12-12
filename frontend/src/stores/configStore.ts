import { defineStore } from 'pinia'
import { ref, onMounted } from 'vue'

export const useConfigStore = defineStore('config', () => {
  const theme = ref('light')
  const language = ref('zh-CN')
  const fontSize = ref(14)

  const setTheme = (newTheme: string) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
  }

  const setLanguage = (lang: string) => {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  const setFontSize = (size: number) => {
    fontSize.value = size
    localStorage.setItem('fontSize', size.toString())
  }

  const init = () => {
    const savedTheme = localStorage.getItem('theme')
    const savedLanguage = localStorage.getItem('language')
    const savedFontSize = localStorage.getItem('fontSize')
    
    if (savedTheme) {
      theme.value = savedTheme
    }
    if (savedLanguage) {
      language.value = savedLanguage
    }
    if (savedFontSize) {
      fontSize.value = parseInt(savedFontSize, 10)
    }
  }

  return {
    theme,
    language,
    fontSize,
    setTheme,
    setLanguage,
    setFontSize,
    init,
  }
})

