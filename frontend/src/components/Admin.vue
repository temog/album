<template>
  <div>
    <h1>
      ユーザ一覧
      <el-button type="primary" style="float:right">
        <router-link to="/createUser">
          Create Account
        </router-link>
      </el-button>
    </h1>

    <el-card v-for="(user, key) in userList" class="box-card" :key="key">
      <div slot="header" class="clearfix">
        <span>{{ user.nickname }}</span>
        <el-button style="float: right; padding: 3px 0" type="text">
          <router-link :to="'editUser/' + user._id">
            Edit
          </router-link>
        </el-button>
      </div>
      <div class="text item">
        Account: {{ user.account }}
      </div>
      <div class="text item">
        Role: {{ user.role }}
      </div>
    </el-card>

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters(['userList', 'isAdmin'])
  },
  mounted () {
    if (!this.isAdmin) {
      this.$router.push({name: 'Top'})
    }
  }
}
</script>
<style>
a {
  color: inherit;
  text-decoration: none;
}
.el-card {
  margin-bottom: 20px;
}
</style>
