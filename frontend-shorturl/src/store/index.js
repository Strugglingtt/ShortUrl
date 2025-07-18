import Vue from 'vue'
import Vuex from 'vuex'
import actions from './actions'
import mutations from './mutations'
import getters from './getters'
import user from './modules/user'
import shorturl from './modules/shorturl'

Vue.use(Vuex)

// 初始化根级状态
const state = {
  isLoading: false,
  error: null
}

export default new Vuex.Store({
  state,
  actions,
  mutations,
  getters,
  modules: {
    user,
    shorturl
  },
  strict: process.env.NODE_ENV !== 'production' // 开发环境开启严格模式
})