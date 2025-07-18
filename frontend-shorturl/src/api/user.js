import request from '@/utils/request'  // 假设你已经封装了axios

/**
 * 用户登录
 * @param {string} username - 用户名
 * @param {string} password - 密码
 */
export function login(username, password) {
  return request({
    url: '/auth/login',
    method: 'post',
    data: {
      username,
      password
    }
  })
}

/**
 * 获取用户信息
 * @param {string} token - 认证令牌
 */
export function getInfo(token) {
  return request({
    url: '/user/info',
    method: 'get',
    params: { token }
  })
}

/**
 * 用户注销
 */
export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

/**
 * 注册新用户
 * @param {Object} userData - 用户数据
 */
export function register(userData) {
  return request({
    url: '/user/register',
    method: 'post',
    data: userData
  })
}

/**
 * 刷新Token
 * @param {string} refreshToken 
 */
export function refreshToken(refreshToken) {
  return request({
    url: '/auth/refresh',
    method: 'post',
    data: { refreshToken }
  })
}