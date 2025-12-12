import { defineStore } from 'pinia'
import { ref } from 'vue'
import { shellApi } from '@/api/shell'
import type { Shell } from '@/types'

export const useShellStore = defineStore('shell', () => {
  const shells = ref<Shell[]>([])

  const fetchShells = async () => {
    shells.value = await shellApi.list()
  }

  const createShell = async (shell: Partial<Shell>) => {
    const newShell = await shellApi.create(shell)
    await fetchShells()
    return newShell
  }

  const updateShell = async (id: number, shell: Partial<Shell>) => {
    const updated = await shellApi.update(id, shell)
    await fetchShells()
    return updated
  }

  const deleteShell = async (id: number) => {
    await shellApi.delete(id)
    await fetchShells()
  }

  const testShell = async (id: number) => {
    return await shellApi.test(id)
  }

  return {
    shells,
    fetchShells,
    createShell,
    updateShell,
    deleteShell,
    testShell,
  }
})



