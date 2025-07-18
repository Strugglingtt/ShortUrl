import request from '@/utils/request'

export function createShortUrl(data) {
  return request({
    url: '/api/shorten',
    method: 'post',
    data
  })
}

export function getStats(code) {
  return request({
    url: `/api/stats/${code}`,
    method: 'get'
  })
}

export function fetchAllUrls() {
  return request({
    url: '/api/urls',
    method: 'get'
  })
}