<template>
  <div id="app">
    <vue-header/>
    <el-main>
      <!--<router-view :key="$route.fullPath" />-->
      <router-view />
    </el-main>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'
import VueHeader from '@/components/Header'

export default {
  name: 'app',
  computed: {
    ...mapGetters(['token', 'userList'])
  },
  components: {
    VueHeader
  },
  mounted () {
    this.$nextTick(function () {
      this.autoTransition()
    })
  },
  watch: {
    token (val) {
      this.autoTransition()
    }
  },
  methods: {
    ...mapMutations(['getUserList']),
    autoTransition () {
      if (!this.token) {
        this.$gm_routerPush('SignIn')
        return
      }

      if (!this.userList.length) {
        this.getUserList()
      }

      if (this.$route.name === 'SignIn') {
        this.$gm_routerPush('Top')
      }
    }
  }
}
</script>

<style>
body,html,
#app,
#app .markdown-body
{
  font-family: 'Avenir', 'M PLUS Rounded 1c', Helvetica, Arial, sans-serif;
}
body,html {
  margin:0;
  padding:0;
  height:100%;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.el-main {
  padding: 10px;
}

.el-card {
  position: relative;
}
.el-card * {
  box-sizing: border-box;
}
.el-card .img {
  display: block;
  width: 100%;
  height: 225px;
  background-size: cover;
  background-position: center;
  cursor: pointer;
}
.el-card time {
  color: #999;
  font-size: 13px;
}
.el-card .info {
  padding: 14px;
}
.el-card .info span {
  display: block;
  width: 100%;
  text-overflow: ellipsis;
}
.el-card .fileName {
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.myGrid {
  display: grid;
  grid-template-columns: 1fr;
}
.myGrid [class*=el-col-] {
  float: none;
  padding: 3px;
}
.myGrid.el-row::before,
.myGrid.el-row::after
 {
    display: none;
}
.myGrid .el-card {
  height: 300px;
}
.el-textarea__inner {
  height: 200px;
}

@media (min-width: 480px){
  .myGrid {
    grid-template-columns: 1fr 1fr;
  }
}
@media (min-width: 768px) {
  .el-main {
    padding: 20px;
  }
  .myGrid {
    grid-template-columns: 1fr 1fr 1fr;
  }
}
@media (min-width: 980px) {
  .myGrid {
    grid-template-columns: 1fr 1fr 1fr 1fr;
  }
}

@media (min-width: 1180px) {
  .myGrid {
    grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  }
}
</style>
