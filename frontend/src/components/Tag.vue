<template>
  <div>
    <h1 class="h1">{{ $route.params.tagName }}</h1>

    <el-row class="myGrid">
      <div class="el-card" v-for="(i, key) in taggedImage" :key="key" :class="{active: i.active}">
        <i v-if="i.active" @click="inactiveCard(key)" class="cardClose el-icon-circle-close"></i>
        <div @click="activeCard(key)" :style="imageQuery(i._id)" class="img">
          <div class="nickname">{{ getNickname(i.userId) }}</div>
        </div>
        <div class="info">

          <div class="actions">

            <el-tooltip v-if="i.secret" class="item" effect="dark" content="secret" placement="top">
              <font-awesome-icon icon="lock" class="secret" />
            </el-tooltip>

            <el-tooltip v-if="i.memo" class="item" effect="dark" content="メモあり" placement="top">
              <i class="el-icon-document memoIcon"></i>
            </el-tooltip>

            <el-tooltip class="item" effect="dark" content="メモ追記" placement="top">
              <i @click="openEditMemo(key, i)" class="el-icon-edit editIcon"></i>
            </el-tooltip>

          </div>

          <div class="fileName">{{ i.name }}</div>
          <div class="bottom clearfix">
            <time class="time">{{ i.updated_at | date }}</time>
          </div>
        </div>

        <div v-if="!mode && i.memo" class="memo markdown-body" v-html="markdown(i.memo)"></div>

        <div v-if="mode === 'edit'" class="editArea">

          <el-form status-icon label-position="top" class="add">
            <el-form-item label="Tag" prop="tag">
              <el-input v-model="editTag" style="width:50%" />
              <el-select v-model="editTagCandidate" placeholder="タグ一覧">
                <el-option v-for="tag in tags"
                  :key="tag.name"
                  :label="tag.name"
                  :value="tag.name" />
              </el-select>
            </el-form-item>

            <div class="control">
              <el-switch v-if="isAdmin" inactive-text="secret" v-model="editSecret" class="switchSecret"/>
              <el-switch inactive-text="preview markdown" v-model="editMarkdown" />
            </div>

            <el-input v-if="!editMarkdown" type="textarea" v-model="memo"></el-input>
            <div v-if="editMarkdown" class="previewMarkdown markdown-body" v-html="markdown(memo)"></div>

            <el-button @click="closeEditMemo" icon="el-icon-close">Cancel</el-button>
            <el-button v-if="isAdmin" @click="deleteMemo(i, key)" type="danger" icon="el-icon-close">Delete</el-button>
            <el-button @click="saveMemo(i)" type="primary" icon="el-icon-edit">Save</el-button>
          </el-form>
        </div>
      </div>
    </el-row>
  </div>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex'
import GlobalMixin from '@/mixin/globalMixin'
import Marked from 'marked'
import 'github-markdown-css/github-markdown.css'
export default {
  data () {
    return {
      page: 1,
      limit: 20,
      getImageAutoTimer: null,
      mode: null,
      memo: null,
      editTag: null,
      editMarkdown: false,
      editSecret: false
    }
  },
  computed: {
    ...mapGetters(['token', 'isAdmin', 'index', 'taggedImage', 'tags', 'taggedImageLoading', 'getNickname', 'userList']),
    editTagCandidate: {
      get () {
        return null
      },
      set (value) {
        this.editTag += ' ' + value
      }
    }
  },
  mounted () {
    this.$nextTick(function () {
      this.initTaggedImage()
      const tagId = this.$route.params.tagId
      if (tagId) {
        this.$store.dispatch('getTaggedImage', {
          tag: tagId,
          page: this.page,
          limit: this.limit
        })
      }

      this.$store.dispatch('getTagAll')
      this.getImageScroll()
    })
  },
  methods: {
    ...mapMutations(['initTaggedImage', 'spliceTaggedImage']),
    nickname (userId) {
      console.log(userId)
      for (let i = 0; i < this.userList.length; i++) {
        const user = this.userList[i]
        console.log(user)
        if (user._id === userId) {
          return user.nickname
        }
      }

      return null
    },
    getImageScroll () {
      window.addEventListener('scroll', () => {
        clearTimeout(this.getImageAutoTimer)
        this.getImageAutoTimer = setTimeout(() => {
          this.getImageAuto()
        }, 500)
      })
    },
    getImageAuto () {
      if (this.taggedImageLoading) {
        return
      }

      const innerHeight = window.innerHeight
      const height = document.documentElement.scrollHeight - 50
      const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
      const current = scrollTop + innerHeight

      if (current < height) {
        return
      }

      this.page++
      const tagId = this.$route.params.tagId
      this.$store.dispatch('getTaggedImage', {
        tag: tagId,
        page: this.page,
        limit: this.limit
      })
    },
    getData () {
    },
    imageQuery (id) {
      return {
        backgroundImage: 'url(' + GlobalMixin.imageQuery(id) + ')'
      }
    },
    activeCard (index) {
      this.$store.dispatch('activeImage', index)
    },
    inactiveCard (index) {
      this.closeEditMemo()
      this.$store.dispatch('inactiveImage', index)
    },
    markdown (str) {
      return Marked(str)
    },
    openEditMemo (key, image) {
      const tags = this.tags
      // 既存のタグをセット
      const currentTag = []
      image.tag.forEach(tag => {
        for (let i = 0; i < tags.length; i++) {
          if (tags[i]._id === tag) { currentTag.push(tags[i].name) }
        }
      })
      this.editTag = currentTag.join(' ')

      this.mode = 'edit'
      this.memo = image.memo
      this.editMarkdown = false
      this.editSecret = image.secret
      this.activeCard(key)
    },
    closeEditMemo () {
      this.mode = null
    },
    saveMemo (image) {
      const data = {
        id: image._id,
        tag: this.editTag,
        secret: this.editSecret,
        memo: this.memo
      }

      const url = GlobalMixin.config('api', 'updateImage')
      GlobalMixin.api(url, data).then((resp) => {
        if (!resp.status) {
          this.$notify({
            title: '更新失敗',
            message: resp.error,
            type: 'warning'
          })
          return
        }

        this.$notify({
          title: 'Success',
          message: '更新しました',
          type: 'success'
        })
      })
    },
    deleteMemo (image, key) {
      this.$confirm('削除しますか?', 'Warning', {
        confirmButtonText: '削除',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.deleteMemoExec(image, key)
      }).catch(() => {
      })
    },
    deleteMemoExec (image, key) {
      const url = GlobalMixin.config('api', 'deleteImage')
      const data = {
        id: image._id
      }
      GlobalMixin.api(url, data).then((resp) => {
        if (!resp.status) {
          this.$notify({
            title: '削除失敗',
            message: resp.error,
            type: 'warning'
          })
          return
        }

        this.$notify({
          title: 'Success',
          message: '削除しました',
          type: 'success'
        })
        this.spliceTaggedImage(key)
        this.closeEditMemo()
      })
    }
  }
}
</script>

<style>
.el-card .memo {
  padding: 15px;
  color: #444;
}
.el-card .memoIcon,
.el-card .editIcon
{
  color: #e68e20;
  font-size: 26px;
}
.el-card .editIcon {
  cursor: pointer;
}
.el-card .actions {
  position: absolute;
  bottom: 5px;
  right: 5px;
  text-aligin: right;
}
.el-card.active {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  z-index: 999;
}
.el-card .img {
  position: relative;
}
.el-card.active .img {
  height: 100%;
}
.el-card.active .info {
  position: absolute;
  top: 45%;
  background-color: #fff;
  min-width: 200px;
}
.el-card.active .memo {
  position: absolute;
  bottom: 3%;
  width: 80%;
  max-height: 38%;
  background-color: #fff;
  margin-left: 10%;
  border-radius: 10px;
  opacity: 0.8;
  overflow-y: auto;
}

.cardClose {
  position: absolute;
  right: 5px;
  top: 5px;
  font-size: 60px;
  text-shadow: 0 0 12px #fff;
  cursor: pointer;
  z-index: 1;
}
.editArea {
  position: fixed;
  padding: 15px;
  top: 20%;
  left: 20%;
  width: 60%;
  background-color: #fff;
}
.control {
  margin-bottom: 10px;
}
.secret {
  color: #666;
}
.previewMarkdown {
  margin: 30px 0;
}
.switchSecret {
  margin-right: 20px;
}
textarea {
  margin-bottom: 15px;
  height: 200px;
}
.nickname {
  position: absolute;
  bottom: 0;
  left: 5px;
  color: #fff;
  font-weight: bold;
}
</style>
