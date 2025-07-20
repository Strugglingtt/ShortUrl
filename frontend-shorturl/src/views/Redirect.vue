<template>
  <el-result
    icon="success"
    title="正在跳转..."
    :subTitle="`即将访问: ${originalUrl}`">
    <template slot="extra">
      <el-button v-if="error" type="primary" @click="$router.push('/')">
        返回首页
      </el-button>
    </template>
  </el-result>
</template>

<script>
export default {
  client() {
    return {
      originalUrl: '',
      error: false
    }
  },
  async mounted() {
    try {
      const { client } = await this.$axios.get(`/api/${this.$route.params.code}`)
      this.originalUrl = client.url
      setTimeout(() => {
        window.location.href = client.url
      }, 1500)
    } catch {
      this.error = true
    }
  }
}
</script>