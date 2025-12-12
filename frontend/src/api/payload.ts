import request from './request'

export const payloadApi = {
  getTemplates: (): Promise<string[]> => {
    return request.get('/payload/templates')
  },
  inject: (data: {
    url: string
    password?: string
    encode_type?: string
    template_name?: string
    shell_path?: string
    headers?: Record<string, string>
  }): Promise<any> => {
    return request.post('/payload/inject', data)
  },
}



