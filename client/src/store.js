import Vue from 'vue'
import Vuex from 'vuex'
import Axios from 'axios'

Axios.defaults.baseURL = 'http://localhost:8080'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    tickets: [],
    errors: []
  },
  getters: {
    tickets: function (state) {
      return state.tickets
    },

    errors: function (state) {
      return state.errors
    }
  },
  mutations: {
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
