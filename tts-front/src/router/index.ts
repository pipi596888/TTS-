import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login/index.vue'),
    },
    {
      path: '/',
      redirect: '/generate',
    },
    {
      path: '/generate',
      name: 'GenerateAudio',
      component: () => import('@/views/GenerateAudio/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/works',
      name: 'Works',
      component: () => import('@/views/Works/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/voice',
      name: 'VoiceManage',
      component: () => import('@/views/VoiceManage/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/custom-voice',
      name: 'CustomVoice',
      component: () => import('@/views/CustomVoice/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/feedback',
      name: 'Feedback',
      component: () => import('@/views/Feedback/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/system',
      name: 'SystemManage',
      component: () => import('@/views/SystemManage/index.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

function getToken() {
  return localStorage.getItem('token')
}

router.beforeEach((to, _from, next) => {
  const token = getToken()
  
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/generate')
  } else {
    next()
  }
})

export default router

