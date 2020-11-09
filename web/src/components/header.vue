<template>
  <el-container>
    <el-header>
      <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
        <el-menu-item index="dashboard"><i class="el-icon-s-home"></i>首页</el-menu-item>
        <el-menu-item index="download"><i class="el-icon-download"></i>资源下载</el-menu-item>
        <el-menu-item index="filesystem"><i class="el-icon-folder"></i>资源管理</el-menu-item>
        <el-menu-item index="monitor"><i class="el-icon-data-analysis"></i>应用监控</el-menu-item>
        <!-- <el-menu-item index="oss"><i class="el-icon-document-checked"></i>文件管理</el-menu-item> -->
        <!-- <el-menu-item style="float: right" index="3">消息中心</el-menu-item> -->
      </el-menu>
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      activeIndex: this.$route.name
    };
  },
  methods: {
    handleSelect: function (index, path) {
      this.$router.push({name: index})
    },
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

.el-menu {
  width: 960px;
  margin: 0 auto;
}

.el-main {
  height: 100%;
  width: 960px;
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
</style>
