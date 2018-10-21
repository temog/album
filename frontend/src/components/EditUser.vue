<template>
  <div>

    <el-form status-icon ref="form" :model="form" :rules="rules">

      <el-form-item label="Nickname" prop="nickname">
        <el-input v-model="form.nickname"/>
      </el-form-item>

      <el-form-item label="Account" prop="account">
        <el-input v-model="form.account"/>
      </el-form-item>

      <el-form-item label="Password" prop="password">
        <el-input type="password" v-model="form.password"/>
      </el-form-item>

      <el-form-item v-if="isAdmin && roleOption" label="Role" prop="role">
        <el-select v-model="form.role">
          <el-option v-for="(value, key) in roleOption" :key="key"
            :label="value"
            :value="value" />
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="editUser('form')">Edit Account</el-button>
      </el-form-item>

    </el-form>

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  data () {
    return {
      complete: false,
      form: {
        id: '',
        nickname: '',
        account: '',
        password: '',
        role: ''
      },
      roleOption: ['guest', 'admin'],
      rules: {
        nickname: [
          { required: true, message: 'Please input Nickname', trigger: 'blur' }
        ],
        account: [
          { required: true, message: 'Please input Account', trigger: 'blur' }
        ],
        password: [
          { required: false, message: 'Please input Password', trigger: 'blur' }
        ],
        role: [
          { required: true, message: 'Please input Role', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters(['userId', 'isAdmin', 'userList', 'getUser'])
  },
  watch: {
    userList (list) {
      if (list.length) {
        this.setUser()
      }
    }
  },
  mounted () {
    this.setUser()
  },
  methods: {
    setUser () {
      console.log(this.$route.params.id)
      const userId = this.$route.params.id

      if (!this.isAdmin && userId !== this.userId) {
        this.$router.push({name: 'Top'})
        return
      }

      const user = this.getUser(userId)
      this.form.id = user._id
      this.form.nickname = user.nickname
      this.form.account = user.account
      this.form.role = user.role
    },
    editUser (formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          return
        }

        const url = this.$gm_config('api', 'editUser')
        this.$gm_api(url, this.form).then((resp) => {
          if (!resp.status) {
            this.$alert('invalid api response (userCount)', 'システムエラー', {
              type: 'warning',
              confirmButtonText: 'Close'
            })
            return
          }

          this.$notify({
            title: 'Success',
            message: '更新しました',
            type: 'success'
          })
        })
      })
    }
  }
}
</script>

<style>
.el-form label {
  display: block;
  float: none;
  text-align: left;
}
</style>
