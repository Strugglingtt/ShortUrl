/**
 * 根级 actions
 * 用于处理多个模块共用的异步操作
 */
export default {
  // 示例：初始化应用数据
  initApp({ dispatch }) {
    return Promise.all([
      dispatch('user/getInfo'),      // 获取用户信息
      dispatch('shorturl/fetchAll')  // 获取短链列表
    ])
  },

  // 全局错误处理
  handleError({ commit }, error) {
    commit('SET_ERROR', error.response?.data?.message || error.message)
    return Promise.reject(error)
  }
}