<template>
  <div>

  <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal"
    background-color="#545c64"
    text-color="#fff"
    active-text-color="#ffd04b">

    <el-menu-item index="1" @click="toggleSideMenu" class="bars" />
    <el-menu-item @click="move('/')" index="2" class="logo"><img src="../assets/logo.png"></el-menu-item>

    <el-menu-item v-if="token" v-for="(value, key) in pageLink" :index="String(key + 100)" :key="key" @click="move(value.path)" class="desktop">
      <i :class="value.icon"></i>{{ value.label }}
    </el-menu-item>

    <el-menu-item v-if="token && isAdmin" index="3" @click="move('/admin')" class="desktop">
      <i class="el-icon-setting"></i>
      Admin
    </el-menu-item>

    <el-menu-item v-if="token && isAdmin" index="4" @click="confirmSecret" class="desktop">
      <i class="el-icon-question"></i>
    </el-menu-item>

    <div class="right desktop" v-if="token">
      <el-menu-item index="5" @click="move('/editUser/' + userId)">
        Profile
      </el-menu-item>
      <el-menu-item index="6" @click="signOut">
        Sign Out
      </el-menu-item>
    </div>
  </el-menu>

  <div id="sideMenuWrapper">
    <div class="bg" @click="sideMenuOff"></div>
    <el-menu
      default-active="1"
      id="sideMenu"
      class="sideMenu"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b">

      <router-link v-if="token" v-for="(value, key) in pageLink" :to="value.path" :key="key" @click.native="sideMenuOff">
        <el-menu-item :index="String(key + 200)">
          <i :class="value.icon"></i>{{ value.label }}
        </el-menu-item>
      </router-link>

      <el-menu-item v-if="token && isAdmin" index="3" @click="confirmSecret">
        <i class="el-icon-question"></i>
      </el-menu-item>

      <el-menu-item v-if="token" index="5" @click="move('/editUser/' + userId)">
        Profile
      </el-menu-item>

      <el-menu-item v-if="token" index="3" @click="signOut">
        Sign Out
      </el-menu-item>
    </el-menu>
  </div>

  </div>

</template>

<script>
import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters(['token', 'userId', 'isAdmin'])
  },
  mounted () {
    this.$nextTick(function () {
    })
  },
  methods: {
    confirmSecret () {
      const self = this
      this.$prompt('合言葉！', '!?', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        inputErrorMessage: 'Invalid',
        inputType: 'password'
      }).then(data => {
        self.sendSecret(data.value)
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Input canceled'
        })
      })
    },
    sendSecret (password) {
      if (!password) {
        return
      }
      const self = this
      this.$gm_api(this.$gm_config('api', 'secret'), {password: password}).then((resp) => {
        if (!resp.status) {
          self.$notify.error({
            title: '失敗',
            message: 'だめです'
          })
          return false
        }

        self.$notify({
          title: 'OK',
          message: 'secret画像も表示します',
          type: 'success'
        })

        self.$router.push({path: '/top/refresh'})

        /*
        if (resp.status) {
          self.addSuccess = true
          self.$notify({
            title: '成功',
            message: '画像の登録に成功しました',
            type: 'success'
          })
        }
        */
      })
    },
    move (path) {
      if (!path) {
        return
      }
      path = path || '/'
      this.$router.push({path: path})
    },
    signOut () {
      this.sideMenuOff()
      this.$store.dispatch('signOut')
    },
    toggleSideMenu () {
      let menuWrapper = document.getElementById('sideMenuWrapper')
      menuWrapper.classList.toggle('on')

      setTimeout(function () {
        let menu = document.getElementById('sideMenu')
        menu.classList.toggle('on')
      }, 100)
    },
    sideMenuOff () {
      let menu = document.getElementById('sideMenu')
      menu.classList.remove('on')

      setTimeout(function () {
        let menuWrapper = document.getElementById('sideMenuWrapper')
        menuWrapper.classList.remove('on')
      }, 500)
    }
  },
  data () {
    return {
      activeIndex: '',
      pageLink: [
        {
          path: '/add',
          icon: 'el-icon-circle-plus',
          label: '登録'
        }
        /*
        {
          path: '/search',
          icon: 'el-icon-search',
          label: '検索'
        }
        */
      ]
    }
  }
}
</script>

<style scoped>
.right {
  float:right;
}
.right .el-menu-item {
  float: left;
}

.logo img {
  width: 48px;
  margin-top: -1px;
  filter: invert(100%);
  transition: all 500ms 0s ease;
}
.logo img:hover {
  filter: invert(0%);
}
#sideMenuWrapper {
  display: none;
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 100;
}
#sideMenuWrapper.on {
  display: block;
}
#sideMenuWrapper .bg {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  background-color: #000;
  opacity: 0.5;
}
.sideMenu {
  position: fixed;
  top: 0;
  width: 200px;
  height: 100%;
  left: -201px;
  z-index: 1;
  transition: all 500ms 0s ease;
  border: 0;
}
.sideMenu.on {
  left: 0;
}
.sideMenu a {
  text-decoration: none;
}
.bars {
  display: none;
}
.bars:before {
  content: "";
  display: block;
  width: 16px;
  height: 4px;
  border-top: 2px solid #fff;
  border-bottom: 2px solid #fff;
  margin-top: 22px;
}
.bars:after {
  content: "";
  display: block;
  width: 16px;
  border-top: 2px solid #fff;
  margin-top: 4px;
}
@media (max-width: 767px) {
  .desktop {
    display: none;
  }
  .bars {
    display: block;
  }
  .logo {
    width: 40%;
    position: absolute !important;
    left: 30%;
    text-align: center;
  }
  .sideMenu {
    display:block;
  }
}
</style>
