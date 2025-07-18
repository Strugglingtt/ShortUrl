<template>
  <el-card shadow="hover" class="stats-card">
    <div slot="header" class="clearfix">
      <span>{{ title }}</span>
      <el-tag :type="tagType" style="float: right">{{ tagText }}</el-tag>
    </div>
    <div class="card-content">
      <div class="stat-value">{{ value }}</div>
      <div class="stat-compare">
        <span :class="trendClass">
          <i :class="trendIcon"></i> {{ trendText }}
        </span>
      </div>
    </div>
  </el-card>
</template>

<script>
export default {
  props: {
    title: String,
    value: [String, Number],
    trend: {
      type: Number, // 正数表示上升，负数表示下降
      default: 0
    }
  },
  computed: {
    trendClass() {
      return this.trend >= 0 ? 'positive' : 'negative'
    },
    trendIcon() {
      return this.trend >= 0 ? 'el-icon-top' : 'el-icon-bottom'
    },
    trendText() {
      return Math.abs(this.trend) + '%'
    },
    tagType() {
      return this.trend >= 0 ? 'success' : 'danger'
    },
    tagText() {
      return this.trend >= 0 ? '上升' : '下降'
    }
  }
}
</script>

<style scoped>
.stats-card {
  margin-bottom: 20px;
}
.card-content {
  text-align: center;
}
.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin: 10px 0;
}
.positive {
  color: #67C23A;
}
.negative {
  color: #F56C6C;
}
</style>