import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedState'
import Axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    sidebarWidth: '250px',
    currentUser: null,
    token: '',
    showSidebar: false
  },

  getters: {
    showSidebar: function (state) {
      return state.showSidebar
    },

    sidebarWidth: function (state) {
      return state.sidebarWidth
    },

    currentUser: function (state) {
      return state.currentUser
    },

    token: function (state) {
      return state.token
    }
  },

  mutations: {
    setSidebarWidth: function (state, width) {
      if (typeof width === 'number') {
        state.sidebarWidth = width.toString() + 'px'
      } else if (width) {
        state.sidebarWidth = width
      } else {
        state.sidebarWidth = '250px'
      }
    },

    login: function (state, { token, user }) {
      state.currentUser = user
      state.token = token
      Axios.defaults.headers.common['Authorization'] = 'Bearer ' + token
    },

    logout: function (state) {
      state.currentUser = null
      state.token = null
      Axios.defaults.headers.common['Authorization'] = null
    },

    sidebarShown: function (state, show) {
      state.showSidebar = show
    }
  },

  plugins: [createPersistedState({
    key: 'praelatusSession',
    getState: (key, storage) => {
      let value = storage.getItem(key)

      try {
        let s = value && value !== 'undefined' ? JSON.parse(value) : undefined
        if (s && s['token']) {
          Axios.defaults.headers.common['Authorization'] = 'Bearer ' + s['token']
        }
        return s
      } catch (err) {
        return undefined
      }
    }
  })]
})
