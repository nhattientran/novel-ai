import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue'
import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from '../pages/RegisterPage.vue'
import StoryListPage from '../pages/creator/StoryListPage.vue'
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
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterPage,
  },
  {
    path: '/creator/stories',
    name: 'StoryList',
    component: StoryListPage,
  },
  {
    path: '/creator/stories/:id/map',
    name: 'StoryMap',
    component: StoryMapPage,
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

export default router
