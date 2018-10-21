'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  api: {
    userCount: '"http://local.album.temo.xyz:3000/api/user/count"',
    createUser: '"http://local.album.temo.xyz:3000/api/user/create"',
    editUser: '"http://local.album.temo.xyz:3000/api/user/edit"',
    signIn: '"http://local.album.temo.xyz:3000/api/user/signIn"',
    signOut: '"http://local.album.temo.xyz:3000/api/user/signOut"',
    secret: '"http://local.album.temo.xyz:3000/api/user/secret"',
    userList: '"http://local.album.temo.xyz:3000/api/user/list"',
    add: '"http://local.album.temo.xyz:3000/api/image/add"',
    getImage: '"http://local.album.temo.xyz:3000/api/image/get"',
    updateImage: '"http://local.album.temo.xyz:3000/api/image/update"',
    deleteImage: '"http://local.album.temo.xyz:3000/api/image/delete"',
    getTaggedImage: '"http://local.album.temo.xyz:3000/api/image/getTagged"',
    getTagAll: '"http://local.album.temo.xyz:3000/api/tag/getAll"',
    getIndex: '"http://local.album.temo.xyz:3000/api/tag/getIndex"'
  }
})
