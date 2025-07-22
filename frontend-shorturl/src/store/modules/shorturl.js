import { createShortUrl, getStats } from '@/api/shorturl'

const state = {
  myUrls: [],
  statsData: null
}

const mutations = {
  SET_MY_URLS: (state, urls) => {
    state.myUrls = urls
  },
  SET_STATS_DATA: (state, data) => {
    state.statsData = data
  }
}

const actions = {
  // 创建短链
  createShortUrl({ commit }, longUrl) {
    return new Promise((resolve, reject) => {
      createShortUrl({ long_url: longUrl })
        .then(response => {
          commit('SET_MY_URLS', [...state.myUrls, response.data])
          resolve(response.data)
        })
        .catch(error => {
          reject(error)
        })
    })
  },

  // 获取统计信息
  fetchStats({ commit }, code) {
    return new Promise((resolve, reject) => {
      getStats(code)
        .then(response => {
          commit('SET_STATS_DATA', response.data)
          resolve(response.data)
        })
        .catch(error => {
          reject(error)
        })
    })
  }
}

const getters = {
  totalClicks: state => {
    return state.statsData?.total || 0
  },
  todayClicks: state => {
    return state.statsData?.today || 0
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}