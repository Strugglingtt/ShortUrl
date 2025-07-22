/**
 * 根级 mutations
 * 用于修改多个模块共用的状态
 */
export default {
  SET_LOADING(state, isLoading) {
    state.isLoading = isLoading
  },

  SET_ERROR(state, errorMsg) {
    state.error = errorMsg
  },

  // 重置所有状态（用于退出登录）
  RESET_STATE(state) {
    Object.assign(state, initialState())
  }
}

// 初始状态函数
function initialState() {
  return {
    isLoading: false,
    error: null
  }
}