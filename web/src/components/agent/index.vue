<template>
    <a-layout-content :style="{ padding: '12px'}">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          节点操作
        </a-breadcrumb-item>
      </a-breadcrumb>
      <a-row :gutter="16" style="background: #fbfbfb;border: 1px solid #f4f4f4;height: 100%">
        
      </a-row>
    </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      system: {}
    }
  },
  computed: {
    count() {
      return this.expand ? 11 : 7;
    },
  },
  created: function () {
    var that = this;
    this.$api.systemInfo().then(function (response) {
      that.data.system = response.details
    }).then(function (response) {
      that.$message.info(response.message)
    })
  },
  methods: {},
};
</script>

<style>

</style>