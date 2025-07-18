<template>
  <div class="container">
    <el-card class="box-card">
      <h2>短链生成器</h2>
      <el-input 
        v-model="longUrl" 
        placeholder="粘贴你的长链接"
        style="width: 80%; margin-bottom: 20px">
      </el-input>
      
      <el-button 
        type="primary" 
        @click="generateShortUrl"
        :loading="loading">
        生成短链
      </el-button>

      <div v-if="shortUrl" class="result-box">
        <el-input 
          v-model="fullShortUrl" 
          readonly 
          style="width: 60%">
          <template slot="append">
            <el-button @click="copyUrl">复制</el-button>
          </template>
        </el-input>
        
        <div id="qrcode" style="margin-top: 20px"></div>
        <p>访问次数: {{ stats.total }}</p>
      </div>
    </el-card>
  </div>
</template>

<script>
import QRCode from 'qrcodejs2'
export default {
  data() {
    return {
      longUrl: '',
      shortUrl: '',
      loading: false,
      stats: { total: 0 }
    }
  },
  computed: {
    fullShortUrl() {
      return `${window.location.origin}/${this.shortUrl}`
    }
  },
  methods: {
    async generateShortUrl() {
      this.loading = true
      try {
        const { data } = await this.$axios.post('/api/shorten', {
          url: this.longUrl
        })
        this.shortUrl = data.shortUrl
        this.generateQRCode()
        this.fetchStats()
      } finally {
        this.loading = false
      }
    },
    generateQRCode() {
      this.$nextTick(() => {
        document.getElementById('qrcode').innerHTML = ''
        new QRCode('qrcode', {
          text: this.fullShortUrl,
          width: 150,
          height: 150
        })
      })
    },
    copyUrl() {
      navigator.clipboard.writeText(this.fullShortUrl)
      this.$message.success('已复制到剪贴板')
    },
    async fetchStats() {
      const { data } = await this.$axios.get(`/api/stats/${this.shortUrl}`)
      this.stats = data
    }
  }
}
</script>

<style scoped>
.container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}
.result-box {
  margin-top: 30px;
  text-align: center;
}
</style>