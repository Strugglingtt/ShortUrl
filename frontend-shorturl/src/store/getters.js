const getters = {
  token: state => state.user.token,
  name: state => state.user.name,
  avatar: state => state.user.avatar,
  myUrls: state => state.shorturl.myUrls,
  statsData: state => state.shorturl.statsData
}

export default getters