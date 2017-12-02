import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import Index from '@/components/Index'
import SearchPage from '@/components/Tickets/SearchPage'
import Ticket from '@/components/Tickets/Show/Page.vue'
import TicketCreate from '@/components/Tickets/Create'
import Project from '@/components/Projects/Show'
import User from '@/components/Users/Show'
import Login from '@/components/Users/Login'
import Dashboard from '@/components/Dashboard/Index'
import ProjectList from '@/components/Projects/List'

const NotFound = {
  name: 'not-found',
  template: '<h1>Whoops, you\'re off the beaten path!</h1>'
}

export default new Router({
  routes: [
    {
      path: '/queries',
      name: 'Tickets/SearchPage',
      component: SearchPage
    },
    {
      path: '/projects',
      name: 'Projects/List',
      component: ProjectList
    },
    {
      path: '/tickets/create',
      name: 'Tickets/Create',
      component: TicketCreate
    },
    {
      path: '/tickets/:key',
      name: 'Tickets/Ticket',
      component: Ticket
    },
    {
      path: '/projects/:key',
      name: 'Projects/Project',
      component: Project
    },
    {
      path: '/login',
      name: 'Users/Login',
      component: Login
    },
    {
      path: '/register',
      name: 'Users/Register',
      component: Login,
      props: { register: true }
    },
    {
      path: '/users/:username',
      name: 'Users/Profile',
      component: User
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/',
      name: 'Index',
      component: Index
    },
    {
      path: '/404',
      name: 'FourOhFour',
      component: NotFound
    },
    {
      path: '*',
      redirect: '/404'
    }
  ]
})
