import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: () => import('@/views/Dashboard.vue'),
    },
    {
      path: '/shells',
      name: 'ShellManager',
      component: () => import('@/views/ShellManager.vue'),
    },
    {
      path: '/shell/:shellId/:tab?',
      name: 'ShellDetail',
      component: () => import('@/views/ShellDetail.vue'),
    },
    {
      path: '/payload',
      name: 'PayloadGen',
      component: () => import('@/views/PayloadGen.vue'),
    },
    {
      path: '/proxy',
      name: 'ProxyManager',
      component: () => import('@/views/ProxyManager.vue'),
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('@/views/Settings.vue'),
    },
  ],
})

export default router

