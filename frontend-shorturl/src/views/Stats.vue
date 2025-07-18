<template>
  <div class="stats-container">
    <el-card>
      <h3>短链分析: {{ $route.params.code }}</h3>
      
      <el-row :gutter="20">
        <el-col :span="8">
          <el-statistic title="总访问量" :value="stats.total" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="今日访问" :value="stats.today" />
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
export default {
  data() {
    return {
      stats: {
        total: 0,
        today: 0,
        logs: []
      }
    }
  },
  async created() {
    const { data } = await this.$axios.get(`/api/stats/${this.$route.params.code}`)
    this.stats = data
  }
}
</script>   