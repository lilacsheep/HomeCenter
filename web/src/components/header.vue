<template>
  <el-container>
    <el-header>
      <el-menu class="menu" v-if="login" :default-active="activeIndex" mode="horizontal" @select="handleSelect">
        <el-menu-item index="dashboard"><i class="el-icon-s-home"></i>首页</el-menu-item>
        <el-menu-item index="download"><i class="el-icon-download"></i>资源下载</el-menu-item>
        <el-menu-item index="filesystem"><i class="el-icon-folder"></i>资源管理</el-menu-item>
        <el-menu-item index="monitor"><i class="el-icon-data-analysis"></i>应用监控</el-menu-item>
        <el-menu-item index="other"><i class="el-icon-more"></i>其他功能</el-menu-item>
        <el-submenu style="float: right">
          <template slot="title">用户操作</template>
          <el-menu-item index="user">信息操作</el-menu-item>
          <el-menu-item index="logout">注销</el-menu-item>
        </el-submenu>
      </el-menu>
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
    <el-dialog title="提示" :visible.sync="logout" width="300px">
      <span>确认是否登出?</span>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="logout = false">取 消</el-button>
        <el-button size="small" type="primary" @click="logout_complete">确 定</el-button>
      </span>
    </el-dialog>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      activeIndex: this.$route.name,
      logout: false,
      login: false
    };
  },
  methods: {
    handleSelect: function (index, path) {
      if (index != "logout") {
        this.$router.push({name: index})
      } else {
        this.logout = true
      }
    },
    logout_complete: function () {
      let that = this
      this.$api.get("/auth/logout").then(function (response) {
        window.location.reload()
      }).catch(function (response) {
        window.location.reload()
      })
    }
  },
  mounted: function () {
    if (this.$route.name != "login") {
      var that = this
      this.$api.post("/auth/self").then(function (response) {
        that.login = true
      }).catch(function (response) {})
    }
  }
};
</script>
<style>
html body #app{
  height: 100%;
  background-color: #DDDDDD;
  overflow: auto;
}

.el-container .el-header {
  height: 40px!important;
  background-color: white;
}

.menu {
  width: 1200px;
  margin: 0 auto;
}

.el-main {
  height: 100%;
  width: 1200px;
  margin: 0 auto;
}

.el-menu {
  height: 40px;
}

.el-menu--horizontal>.el-menu-item {
  height: 40px;
  line-height: 40px
}

.el-menu--horizontal>.el-submenu .el-submenu__title {
  height: 40px;
  line-height: 40px
}

.el-menu--collapse .el-menu .el-submenu, .el-menu--popup {
  min-width: 150px;
}
</style>
