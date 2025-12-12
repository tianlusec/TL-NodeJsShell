import request from './request'
import type { Shell } from '@/types'

export const shellApi = {
  list: (): Promise<Shell[]> => {
    return request.get('/shells')
  },
  get: (id: number): Promise<Shell> => {
    return request.get(`/shells/${id}`)
  },
  create: (shell: Partial<Shell>): Promise<Shell> => {
    return request.post('/shells', shell)
  },
  update: (id: number, shell: Partial<Shell>): Promise<Shell> => {
    return request.put(`/shells/${id}`, shell)
  },
  delete: (id: number): Promise<void> => {
    return request.delete(`/shells/${id}`)
  },
  test: (id: number): Promise<{ success: boolean; latency: number }> => {
    return request.post(`/shells/${id}/test`)
  },
  execute: (id: number, command: string): Promise<any> => {
    return request.post(`/shells/${id}/execute`, { command })
  },
  getInfo: (id: number, params?: any): Promise<any> => {
    return request.get(`/shells/${id}/info`, { params })
  },
}

