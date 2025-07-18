<template>
  <div id="app">
    <!-- 顶部导航栏 -->
    <el-header v-if="showHeader">
      <el-menu
        :default-active="activeMenu"
        mode="horizontal"
        @select="handleMenuSelect"
      >
        <el-menu-item index="/">短链生成</el-menu-item>
        <el-menu-item index="/stats">数据统计</el-menu-item>
        <el-menu-item v-if="!isAuthenticated" index="/login">登录</el-menu-item>
        <el-submenu v-else index="user-menu">
          <template #title>{{ userName }}</template>
          <el-menu-item @click="handleLogout">退出登录</el-menu-item>
        </el-submenu>
      </el-menu>
    </el-header>

    <!-- 主内容区 -->
    <el-main>
      <router-view/>
    </el-main>

    <!-- 全局加载动画 -->
    <el-dialog
      :visible.sync="globalLoading"
      :show-close="false"
      width="30%"
      top="45vh"
      custom-class="loading-dialog"
    >
      <div class="loading-content">
        <el-progress type="circle" :percentage="50" :status="loadingStatus"/>
        <p>{{ loadingText }}</p>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'App',
  data() {
    return {
      loadingText: '加载中...',
      loadingStatus: null
    }
  },
  computed: {
    ...mapGetters('user', ['isAuthenticated', 'userName']),
    
    // 控制是否显示顶部导航
    showHeader() {
      return this.$route.meta.showHeader !== false
    },
    
    // 当前激活的菜单项
    activeMenu() {
      return this.$route.path
    },
    
    // 全局加载状态
    globalLoading: {
      get() {
        return this.$store.state.isLoading
      },
      set(val) {
        this.$store.commit('SET_LOADING', val)
      }
    }
  },
  methods: {
    ...mapActions('user', ['logout']),
    
    // 菜单选择处理
    handleMenuSelect(index) {
      if (index === 'user-menu') return
      if (index === '/login') {
        this.$router.push('/login')
      } else {
        this.$router.push(index)
      }
    },
    
    // 退出登录
    async handleLogout() {
      this.loadingText = '正在退出...'
      this.loadingStatus = null
      this.globalLoading = true
      
      try {
        await this.logout()
        this.$router.push('/login')
      } catch (error) {
        this.loadingStatus = 'exception'
        this.loadingText = '退出失败'
        setTimeout(() => {
          this.globalLoading = false
        }, 1000)
      } finally {
        setTimeout(() => {
          this.globalLoading = false
        }, 500)
      }
    }
  }
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.el-header {
  padding: 0;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.el-main {
  flex: 1;
  padding: 20px;
}

.loading-dialog {
  background-color: transparent !important;
  box-shadow: none !important;
  
  .el-dialog__header {
    display: none;
  }
  
  .loading-content {
    text-align: center;
    
    p {
      margin-top: 10px;
      color: #fff;
    }
  }
}
</style>