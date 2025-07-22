<template>
  <div v-if="error">
    <el-result
      icon="error"
      title="获取数据失败"
      subTitle="无法加载统计信息，请稍后重试">
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </el-result>
  </div>

  <div v-else-if="isList">
    <el-card>
      <h2>所有短链统计数据</h2>
      <el-table :data="allStats" border>
        <el-table-column prop="shortCode" label="短码" width="100"></el-table-column>
        <el-table-column prop="originalUrl" label="原始URL" min-width="300"></el-table-column>
        <el-table-column prop="totalClicks" label="点击量" width="100"></el-table-column>
        <el-table-column label="操作" width="120">
          <template slot-scope="scope">
            <el-button @click="$router.push(`/stats/${scope.row.shortCode}`)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
          v-model:current-page="currentPage"
          :page-sizes="[5, 10, 20, 50]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
        </el-pagination>
      </div>
    </el-card>
  </div>

  <div v-else-if="statsData">
    <el-card>
      <h2>短链统计详情</h2>
      <el-descriptions title="基本信息" column="2">
        <el-descriptions-item label="短码">{{ statsData.shortCode }}</el-descriptions-item>
        <el-descriptions-item label="原始URL">{{ statsData.originalUrl }}</el-descriptions-item>
        <el-descriptions-item label="总点击量">{{ statsData.totalClicks }}</el-descriptions-item>
      </el-descriptions>

    </el-card>
  </div>

  <div v-else>
    <el-result icon="loading" title="加载中..." subTitle="正在获取统计数据，请稍候"></el-result>
  </div>
</template>

<script>
import { getStats, fetchAllStats } from '@/api/shorturl'

export default {
  data() {
      return {
        statsData: null,
        allStats: [],
        isList: false,
        currentPage: 1,
        pageSize: 10,
        total: 0,
        totalPages: 0,
        error: false
      }
    },
    watch: {
        currentPage:function(newVal){
          // console.log(newVal)
        }
    },
  async mounted() {
      this.fetchStatsData()
    },
    methods: {
        async fetchStatsData(page) {
          try {
            const code = this.$route.params.code
            if (code) {
              // 获取单个短链统计
              const { data } = await getStats(code)
              this.statsData = data.data
              this.isList = false
            } else {
              // 获取所有短链统计（带分页）
              const res = await fetchAllStats(this.currentPage, this.pageSize)
              const data =res.data
              this.allStats = data
              this.total = res.total
              this.totalPages = res.TotalPages
              this.isList = true
            }
          } catch (err) {
            this.error = true
            console.error('获取统计数据失败:', err)
          }
        },
        handleSizeChange(size) {
        this.pageSize = size
         // 重置到第一页
         this.currentPage=1;
        this.fetchStatsData()
      },
        handlePageChange(page){
        this.currentPage = page
        this.fetchStatsData()
        }
  }
}
</script>