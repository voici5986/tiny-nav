import { createRouter, createWebHashHistory } from 'vue-router'
import { useMainStore } from '@/stores'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/',
      name: 'nav',
      component: () => import('@/views/Nav.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach(async (to, _) => {
  // 如果路由需要认证
  if (to.meta.requiresAuth) {
    const store = useMainStore()
    let needAuth = true
    if (to.name === 'nav') {
      const config = store.config
      if (config.enableNoAuth || config.enableNoAuthView) {
        needAuth = false
      }
    }
    if (needAuth) {
      // 验证 token
      const isValid = await store.validateToken()

      if (!isValid) {
        // token 无效，重定向到登录页
        return { name: 'login' }
      }
    }
    console.log('需要认证的路由:', to.name)
  }
})

export default router
