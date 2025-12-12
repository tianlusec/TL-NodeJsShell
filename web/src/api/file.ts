import request from './request'
import type { FileItem } from '@/types'

export const fileApi = {
  list: (shellId: number, path: string): Promise<{ path: string; files: FileItem[] }> => {
    return request.get(`/shells/${shellId}/files`, { params: { path } })
  },
  read: (shellId: number, path: string): Promise<{ path: string; content: string }> => {
    return request.get(`/shells/${shellId}/files/read`, { params: { path } })
  },
  upload: (shellId: number, remotePath: string, content: string, chunkIndex?: number, totalChunks?: number): Promise<any> => {
    const data: any = { remote_path: remotePath, content }
    if (chunkIndex !== undefined && totalChunks !== undefined) {
      data.chunk_index = chunkIndex
      data.total_chunks = totalChunks
    }
    return request.post(`/shells/${shellId}/files/upload`, data)
  },
  download: (shellId: number, path: string, chunkIndex?: number, chunkSize?: number): Promise<any> => {
    const params: any = { path }
    if (chunkIndex !== undefined && chunkSize !== undefined) {
      params.chunk_index = chunkIndex
      params.chunk_size = chunkSize
    }
    return request.get(`/shells/${shellId}/files/download`, { params, responseType: chunkIndex !== undefined ? 'json' : 'blob' })
  },
  update: (shellId: number, path: string, content: string): Promise<any> => {
    return request.put(`/shells/${shellId}/files`, { path, content })
  },
  delete: (shellId: number, path: string): Promise<any> => {
    return request.delete(`/shells/${shellId}/files`, { params: { path } })
  },
  mkdir: (shellId: number, path: string): Promise<any> => {
    return request.post(`/shells/${shellId}/files/mkdir`, { path })
  },
}


