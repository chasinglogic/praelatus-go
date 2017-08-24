import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import SearchPage from '@/components/Tickets/SearchPage'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Index',
      component: Index
    },
    {
      path: '/tickets',
      name: 'Tickets/SearchPage',
      component: SearchPage
    }
  ]
})
