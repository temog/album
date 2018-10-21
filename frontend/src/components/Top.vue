<template>
  <div>
    <el-row class="myGrid">
      <div v-for="(i, key) in index" :key="key" @click="moveTag(i._id, i.name)" >
        <el-card :body-style="{ padding: '0px' }">
          <div :style="imageQuery(i.thumbnail)" class="img"></div>
          <div class="info">
            <span>{{ i.name }}</span>
            <div class="bottom clearfix">
              <time class="time">{{ i.lastModified }}</time>
            </div>
          </div>
        </el-card>
      </div>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import GlobalMixin from '@/mixin/globalMixin'
export default {
  name: 'Top',
  computed: {
    ...mapGetters(['token', 'index'])
  },
  mounted () { this.$store.dispatch('getIndex') },
  watch: {
    '$route.params.action': function (action) {
      if (action === 'refresh') {
        this.$store.dispatch('getIndex')
      }
    }
  },
  methods: {
    getData () {
    },
    imageQuery (id) {
      return {
        backgroundImage: 'url(' + GlobalMixin.imageQuery(id) + ')'
      }
    },
    moveTag (tagId, tagName) {
      this.$router.push({path: '/tag/' + tagId + '/' + tagName})
    }
  }
}
</script>

<style>

</style>
