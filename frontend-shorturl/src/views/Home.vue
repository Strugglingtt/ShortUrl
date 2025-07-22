<template>
  <div class="container">
  <el-card class="box-card">
      <h2 style="text-align: center; margin-bottom: 30px; font-size: 24px;">短链生成器</h2>
      <el-input 
  v-model="longUrl" 
  placeholder="粘贴你的长链接"
  style="width: 100%; margin-bottom: 20px">
</el-input>
      
      <el-button 
        type="primary" 
        @click="generateShortUrl"
        :loading="loading"
        style="width: 100%; height: 44px; font-size: 16px;">
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
        <p>访问次数: {{ stats.totalClicks }}</p>
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
      shortCode: '',
      loading: false,
      stats: { 
        shortCode:"",
        originalUrl:"",
        totalClicks: 0,
       }
    }
  },
  computed: {
    fullShortUrl() {
      console.log(process.env.VUE_APP_BASE_API)
      const baseUrl = process.env.VUE_APP_BASE_API ;
      return `${baseUrl}/${this.shortUrl}`;
    }
  },
  methods: {
    async generateShortUrl() {
      this.loading = true
      try {
        const { data } = await this.$axios.post('/api/shorten', {
          long_url: this.longUrl,
          expire_time: '2025-10-01T00:00:00Z',
        })
        console.log(data)
        this.shortUrl = data.shortUrl
        this.shortCode =data.shortCode
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
      const { data } = await this.$axios.get(`/api/stats/${this.shortCode}`)
      this.stats = data
    }
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.box-card {
  width: 100%;
  max-width: 550px;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border: none !important;
}
.result-box {
  margin-top: 30px;
  text-align: center;
}
</style>