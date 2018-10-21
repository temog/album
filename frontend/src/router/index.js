import Vue from 'vue'
import Router from 'vue-router'
import Top from '@/components/Top'
import SignIn from '@/components/SignIn'
import CreateUser from '@/components/CreateUser'
import Add from '@/components/Add'
import Tag from '@/components/Tag'
import Admin from '@/components/Admin'
import EditUser from '@/components/EditUser'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: '/',
  routes: [
    {
      path: '/',
      name: 'Top',
      component: Top
    },
    {
      path: '/signIn',
      name: 'SignIn',
      component: SignIn
    },
    {
      path: '/createUser',
      name: 'CreateUser',
      component: CreateUser
    },
    {
      path: '/add',
      name: 'Add',
      component: Add
    },
    {
      path: '/tag/:tagId/:tagName',
      name: 'Tag',
      component: Tag
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin
    },
    {
      path: '/editUser/:id',
      name: 'EditUser',
      component: EditUser
    }
  ]
})
