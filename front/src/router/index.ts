import { createRouter, createWebHashHistory } from 'vue-router'
import { useMainStore } from '@/stores'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/nav',
      name: 'nav',
      component: () => import('@/views/Nav.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach(async (to, _) => {
  const store = useMainStore()

  // 如果路由需要认证
  if (to.meta.requiresAuth) {
    // 验证 token
    const isValid = await store.validateToken()

    if (!isValid) {
      // token 无效，重定向到登录页
      return { name: 'login' }
    }
  } else if (to.name === 'login') {
    // 如果要去登录页，先检查是否有有效token
    const isValid = await store.validateToken()

    if (isValid) {
      // 如果 token 有效，直接跳转到导航页
      return { name: 'nav' }
    }
  }
})

export default router
