import { createRouter, createWebHistory } from 'vue-router'
import PromQLWorkbench from '@/views/PromQLWorkbench.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: PromQLWorkbench
    }
  ]
})

export default router
