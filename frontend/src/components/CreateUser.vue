<template>
  <div>
    <div v-if="!complete">
      <h1 v-if="!roleOption">Create Administrator</h1>
      <h1 v-else>Create Account</h1>

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

        <el-form-item v-if="roleOption" label="Role" prop="role">
          <el-select v-model="form.role">
            <el-option v-for="(value, key) in roleOption" :key="key"
              :label="value"
              :value="value" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="createUser('form')">Create Account</el-button>
        </el-form-item>

      </el-form>
    </div>

    <div v-else>
      <el-alert title="ユーザを作成しました" type="success" />
      <br>
      <router-link :to="{name: 'SignIn'}">
        <el-button type="primary">Sign In</el-button>
      </router-link>
    </div>

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  data () {
    return {
      complete: false,
      form: {
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
          { required: true, message: 'Please input Password', trigger: 'blur' }
        ],
        role: [
          { required: true, message: 'Please input Role', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters(['isAdmin'])
  },
  mounted () {
    const url = this.$gm_config('api', 'userCount')
    this.$gm_get(url, null).then((resp) => {
      if (!resp.status) {
        this.$alert('invalid api response (userCount)', 'システムエラー', {
          type: 'warning',
          confirmButtonText: 'Close'
        })
        return
      }
      if (resp.count === 0) {
        this.roleOption = null
      } else if (resp.count !== 0 && !this.isAdmin) {
        this.$router.push({name: 'Top'})
      }
    })
  },
  methods: {
    createUser (formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          return
        }

        const url = this.$gm_config('api', 'createUser')
        console.warn(url)

        if (!this.roleOption) {
          // 初回ユーザ
          this.$gm_post(url, this.form).then((resp) => {
            if (!resp.status) {
              this.$alert('invalid api response (userCount)', 'システムエラー', {
                type: 'warning',
                confirmButtonText: 'Close'
              })
              return
            }
            this.complete = true
          })
        } else {
          // 初回以降
          this.$gm_api(url, this.form).then((resp) => {
            if (!resp.status) {
              this.$alert('invalid api response (userCount)', 'システムエラー', {
                type: 'warning',
                confirmButtonText: 'Close'
              })
              return
            }
            this.complete = true
          })
        }
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
