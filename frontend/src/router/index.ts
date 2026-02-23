import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores'
import HomePage from '../pages/HomePage.vue'
import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from '../pages/RegisterPage.vue'
import StoryListPage from '../pages/creator/StoryListPage.vue'
import StoryEditPage from '../pages/creator/StoryEditPage.vue'
import StoryMapPage from '../pages/creator/StoryMapPage.vue'
import ReadPage from '../pages/reader/ReadPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: { guestOnly: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterPage,
    meta: { guestOnly: true },
  },
  {
    path: '/creator/stories',
    name: 'StoryList',
    component: StoryListPage,
    meta: { requiresAuth: true, requiresCreator: true },
  },
  {
    path: '/creator/stories/:storyId/edit',
    name: 'StoryEdit',
    component: StoryEditPage,
    meta: { requiresAuth: true, requiresCreator: true },
  },
  {
    path: '/creator/stories/:storyId/map',
    name: 'StoryMap',
    component: StoryMapPage,
    meta: { requiresAuth: true, requiresCreator: true },
  },
  {
    path: '/read/:storyId',
    name: 'Read',
    component: ReadPage,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guards
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Fetch current user if not already loaded
  if (!authStore.user && !authStore.isLoading) {
    await authStore.fetchMe()
  }

  const isAuthenticated = authStore.isAuthenticated
  const isCreator = authStore.isCreator

  // Redirect authenticated users away from guest-only pages
  if (to.meta.guestOnly && isAuthenticated) {
    if (isCreator) {
      return next({ name: 'StoryList' })
    }
    return next({ name: 'Home' })
  }

  // Require authentication
  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  // Require creator role
  if (to.meta.requiresCreator && !isCreator) {
    return next({ name: 'Home' })
  }

  next()
})

export default router
