<template>
  <el-container>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      collapsed: false,
    };
  },
  methods: {}
};
</script>
<style>
html body {
  height: 100%;
}
.el-container {
  height: 100%;
  background-color: #DDDDDD;
}

.el-main {
  height: 100%;
}
</style>
