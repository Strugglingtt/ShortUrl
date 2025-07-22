<template>
  <div class="stats-container">
    <el-card>
      <h3>短链分析: {{ $route.params.code }}</h3>
      
      <el-row :gutter="20" style="align-items: stretch; min-height: 220px;">
          <el-col :span="8" style="text-align: center; display: flex; flex-direction: column; justify-content: center; height: 100%;">
            <el-statistic title="总访问量" :value="stats.totalClicks" style="font-size: 18px;" />
          </el-col>
          <el-col :span="8" style="text-align: center; display: flex; flex-direction: column; justify-content: center; height: 100%;">
            <el-statistic title="今日访问" :value="stats.today" style="font-size: 18px;" />
          </el-col>
          <el-col :span="8">
            <div style="height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; border: 1px solid #e8e8e8; border-radius: 4px; padding: 15px;">
              <h3 style="margin-bottom: 15px; font-size: 16px;">短链接二维码</h3>
              <QRCode :text="`${baseUrl}/s/${$route.params.code}`" :size="120"/>
              <p style="margin-top: 15px; word-break: break-all; max-width: 200px; text-align: center; font-size: 14px;">{{ baseUrl }}/{{ $route.params.code }}</p>
            </div>
          </el-col>
        </el-row>

      <el-divider />

      <el-table :data="stats.logs" style="width: 100%">
        <el-table-column prop="time" label="时间" width="180" />
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="device" label="设备" />
      </el-table>
    </el-card>
  </div>
</template>

<script>
import QRCode from '@/components/QRCode.vue'
export default {
  components: {
    QRCode
  },
  data() {
    return {
        baseUrl: process.env.VUE_APP_BASE_URL,
        stats: {
          totalClicks: 0,
          today: 0,
          logs: []
        }
      }
  },
  async created() {
    console.log(this.$route.params.code)
    const { data } = await this.$axios.get(`/api/stats/${this.$route.params.code}`)
    this.stats = data
    console.log(data)
    console.log(this.stats)
  }
}
</script>