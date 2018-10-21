// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import store from '@/vuex/store'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faLock, faUserSecret } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import App from '@/App'
import router from '@/router'
import * as filter from '@/assets/filter'
import ElementUI from 'element-ui'
import locale from 'element-ui/lib/locale/lang/ja'
import GlobalMixin from './mixin/globalMixin'
import 'element-ui/lib/theme-chalk/index.css'

library.add(faLock)
library.add(faUserSecret)
Vue.component('font-awesome-icon', FontAwesomeIcon)

Vue.use(GlobalMixin)

Vue.use(ElementUI, { locale })

Vue.config.productionTip = false

// カスタムフィルタ登録
for (const key in filter) {
  Vue.filter(key, filter[key])
}

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  components: { App },
  template: '<App/>'
})
