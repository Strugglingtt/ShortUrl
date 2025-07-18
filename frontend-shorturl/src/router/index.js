import Vue from 'vue'
import Router from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/:code',
      name: 'redirect',
      component: () => import('../views/Redirect.vue')
    },
    {
      path: '/stats/:code',
      name: 'stats',
      component: () => import('../views/Stats.vue')
    }
  ]
})