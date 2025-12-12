import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTerminalStore = defineStore('terminal', () => {
  const commandHistory = ref<string[]>([])
  const historyIndex = ref(-1)

  const addToHistory = (command: string) => {
    if (command.trim() && commandHistory.value[commandHistory.value.length - 1] !== command) {
      commandHistory.value.push(command)
      historyIndex.value = commandHistory.value.length
    }
  }

  const getPreviousCommand = (): string | null => {
    if (commandHistory.value.length === 0) return null
    if (historyIndex.value > 0) {
      historyIndex.value--
    }
    return commandHistory.value[historyIndex.value] || null
  }

  const getNextCommand = (): string | null => {
    if (commandHistory.value.length === 0) return null
    if (historyIndex.value < commandHistory.value.length - 1) {
      historyIndex.value++
      return commandHistory.value[historyIndex.value] || null
    }
    historyIndex.value = commandHistory.value.length
    return null
  }

  const resetHistoryIndex = () => {
    historyIndex.value = commandHistory.value.length
  }

  return {
    commandHistory,
    historyIndex,
    addToHistory,
    getPreviousCommand,
    getNextCommand,
    resetHistoryIndex,
  }
})



