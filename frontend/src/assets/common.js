import router from '../router'
import store from '../vuex/store'
import { Loading } from 'element-ui'

class Common {
  // 可変長引数
  static config (...key) {
    let conf = null
    key.forEach(function (value) {
      if (!conf) {
        conf = process.env
      }
      conf = conf[value]
    })
    return conf
  }

  static createFormData (data) {
    let formData = new FormData()
    for (let key in data) {
      formData.append(key, data[key])
    }
    return formData
  }

  // token必須api用
  static api (url, data) {
    return new Promise(function (resolve, reject) {
      if (!data.token) {
        Common.routerPush('signIn')
        return
      }

      Common.post(url, data).then(json => {
        if (json.status === 'authError') {
          store.dispatch('signOut')
          Common.routerPush('signIn')
          return
        }
        resolve(json)
      })
    })
  }

  static post (url, data) {
    return new Promise(function (resolve, reject) {
      let loadingInstance = Loading.service({ fullscreen: true })
      let formData = Common.createFormData(data)
      fetch(url, {
        method: 'POST',
        body: formData
      })
      .then(res => res.json())
      .then(json => {
        loadingInstance.close()
        resolve(json)
      })
    })
  }

  static routerPush (name, param = null) {
    if (param) {
      router.push({name: name, params: param})
    } else {
      router.push({name: name})
    }
  }

  static resizeTextbox (e) {
    let height = 0
    if (e.target) {
      height = e.target.innerHeight
    } else {
      height = e.innerHeight
    }

    // 面倒なので footer 分きめうち・・
    let footerHeight = 80
    let bodyBlock = document.getElementById('bodyBlock')
    let boxHeight = height - bodyBlock.getBoundingClientRect().top - footerHeight
    bodyBlock.style.height = boxHeight + 'px'
    document.getElementById('previewBlock').style.height = boxHeight + 'px'
  }
}
export default Common
