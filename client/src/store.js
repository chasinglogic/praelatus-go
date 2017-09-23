import Vue from 'vue'
import Vuex from 'vuex'
import Axios from 'axios'

// TODO: Set this during build time.
Axios.defaults.baseURL = 'http://localhost:8080'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    sidebarWidth: '250px',
    showSidebar: false,
    tickets: [],
    errors: []
  },

  getters: {
    tickets: function (state) {
      return state.tickets
    },

    errors: function (state) {
      return state.errors
    },

    showSidebar: function (state) {
      return state.showSidebar
    },

    sidebarWidth: function (state) {
      return state.sidebarWidth
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

    setSidebarShown: function (state, show) {
      state.sidebarShown = show
    },

    API_SUCCESS: function (state, payload) {
      state[payload.key] = payload.data
    },

    API_FAILURE: function (state, payload) {
      if (state[payload.key]) {
        state[payload.key] = Array.isArray(state[payload.key]) ? [] : {}
      }

      state.errors.push(payload.data)
    }
  },

  actions: {
    request: function (context, { url, key }) {
      Axios.get(url).then(res => {
        console.log(res)
        context.commit('API_SUCCESS', { key: key, data: res.data })
      })
      .catch(err => {
        context.commit('API_FAILURE', { key: key, data: err.response.data })
      })
    }
  }
})
