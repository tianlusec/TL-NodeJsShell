export interface Shell {
  id: number
  url: string
  password: string
  encode_type: string
  protocol: string
  method: string
  group?: string
  name: string
  status: string
  last_active?: string
  latency?: number
  system_info?: string
  custom_headers?: string
  proxy_id?: number
  created_at: string
  updated_at: string
}

export interface Proxy {
  id: number
  name: string
  type: string
  host: string
  port: number
  username?: string
  password?: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface FileItem {
  name: string
  path: string
  type: string
  size: string
  mode: string
  time: string
}



