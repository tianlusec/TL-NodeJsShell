import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/api/request'
import type { Proxy } from '@/types'

export const useProxyStore = defineStore('proxy', () => {
  const proxies = ref<Proxy[]>([])

  const fetchProxies = async () => {
    proxies.value = await request.get('/proxies')
  }

  const getEnabledProxies = () => {
    return computed(() => proxies.value.filter(p => p.enabled)).value
  }

  return {
    proxies,
    fetchProxies,
    getEnabledProxies,
  }
})



