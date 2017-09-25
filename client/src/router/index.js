import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import SearchPage from '@/components/Tickets/SearchPage'
import Ticket from '@/components/Tickets/Ticket'

Vue.use(Router)

const NotFound = {
  name: 'not-found',
  template: '<h1>Nothing to see here!</h1>'
}

export default new Router({
  routes: [
    {
      path: '/queries',
      name: 'Tickets/SearchPage',
      component: SearchPage
    },
    {
      path: '/ticket/:key',
      name: 'Tickets/Ticket',
      component: Ticket
    },
    {
      path: '/',
      name: 'Index',
      component: Index
    },
    {
      path: '*',
      name: 'FourOhFour',
      component: NotFound
    }
  ]
})
