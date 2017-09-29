// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import resource from 'vue-resource'
import bootstrap from 'bootstrap-vue'
import store from './store'

import Axios from 'axios'
// TODO: Set this during build time.
Axios.defaults.baseURL = 'http://localhost:8080'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false

Vue.use(resource)
Vue.use(bootstrap)

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  template: '<App/>',
  components: { App }
})
