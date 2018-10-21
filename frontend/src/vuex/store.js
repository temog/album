import Vue from 'vue'
import Vuex from 'vuex'
import GlobalMixin from '../mixin/globalMixin'
Vue.use(Vuex)

const state = {
  token: GlobalMixin.getStorage('token'), // アクセストークン
  userId: GlobalMixin.getStorage('userId'),
  nickname: GlobalMixin.getStorage('nickname'),
  role: GlobalMixin.getStorage('role'),
  tags: null,
  index: null,
  taggedImage: [],
  taggedImageLoading: false,
  userList: []
}

const actions = {
  // signIn
  signIn ({ commit }, data) {
    const url = GlobalMixin.config('api', 'signIn')
    GlobalMixin.post(url, {
      account: data.account,
      password: data.password
    }).then((resp) => {
      commit('signIn', resp)
    })
  },
  // signOut
  signOut ({ commit }) {
    const url = GlobalMixin.config('api', 'signOut')
    GlobalMixin.api(url, {
      token: state.token
    }).then((resp) => {
      commit('signOut', resp)
    })
  },
  // getTagAll
  getTagAll ({ commit }) {
    const url = GlobalMixin.config('api', 'getTagAll')
    GlobalMixin.api(url).then((resp) => {
      commit('getTagAll', resp)
    })
  },
  getIndex ({ commit }) {
    if (!state.token) {
      return
    }
    const url = GlobalMixin.config('api', 'getIndex')
    GlobalMixin.api(url).then((resp) => {
      commit('getIndex', resp)
    })
  },
  getTaggedImage ({ commit }, data) {
    commit('getTaggedImage', data)
  },
  activeImage ({ commit }, index) {
    commit('activeImage', index)
  },
  inactiveImage ({ commit }, index) {
    commit('inactiveImage', index)
  }
}

const mutations = {
  signIn (state, data) {
    if (data.status) {
      state.token = data.token
      state.userId = data.userId
      state.nickname = data.nickname
      state.role = data.role
      GlobalMixin.setStorage('token', data.token)
      GlobalMixin.setStorage('userId', data.userId)
      GlobalMixin.setStorage('nickname', data.nickname)
      GlobalMixin.setStorage('role', data.role)
    }
  },
  signOut (state, resp) {
    if (resp.status) {
      state.token = null
      state.userId = null
      state.nickname = null
      state.role = null
      GlobalMixin.removeStorage('token')
      GlobalMixin.removeStorage('userId')
      GlobalMixin.removeStorage('nickname')
      GlobalMixin.removeStorage('role')
    }
  },
  getUserList (state) {
    const url = GlobalMixin.config('api', 'userList')
    console.warn(url)
    GlobalMixin.api(url).then((resp) => {
      if (!resp.status) {
        return
      }
      state.userList = resp.userList
    })
  },
  getTagAll (state, resp) {
    if (resp.status) {
      state.tags = resp.tags
    }
  },
  getIndex (state, resp) {
    if (resp.status) {
      state.index = resp.tags
    }
  },
  initTaggedImage (state) {
    state.taggedImage = []
  },
  getTaggedImage (state, data) {
    state.taggedImageLoading = true
    const url = GlobalMixin.config('api', 'getTaggedImage')
    GlobalMixin.api(url, data).then((resp) => {
      if (resp.status) {
        resp.images.forEach(i => {
          i.active = false
        })
        state.taggedImage = state.taggedImage.concat(resp.images)
      }
      state.taggedImageLoading = false
    })
  },
  spliceTaggedImage (state, index) {
    state.taggedImage.splice(index, 1)
  },
  activeImage (state, index) {
    state.taggedImage[index].active = true
  },
  inactiveImage (state, index) {
    state.taggedImage[index].active = false
  }
}

const getters = {
  token: state => state.token,
  userId: state => state.userId,
  nickname: state => state.nickname,
  role: state => state.role,
  isAdmin: state => state.role === 'admin',
  tags: state => state.tags,
  index: state => state.index,
  taggedImage: state => state.taggedImage,
  taggedImageLoading: state => state.taggedImageLoading,
  userList: state => state.userList,
  getUser: state => (userId) => {
    if (!state.userList) {
      return false
    }

    for (let i = 0; i < state.userList.length; i++) {
      const u = state.userList[i]
      if (u._id === userId) {
        return u
      }
    }
    return false
  },
  getNickname: (state, getters) => (userId) => {
    if (!state.userList) {
      return null
    }

    const user = getters.getUser(userId)
    if (!user) {
      return null
    }

    return user.nickname
  }
}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})
