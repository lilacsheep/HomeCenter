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
.ant-list-bordered .ant-list-header {
  padding-left: 0;
}
.ant-list-header, .ant-list-footer {
  padding: 0;
}

.ant-list-item {
  padding: 5px;
}
.card-container {
  background: #fbfbfb;
  overflow: hidden;
  /* padding: 24px; */
}
.card-container > .ant-tabs-card > .ant-tabs-content {
  margin-top: -16px;
}

.card-container > .ant-tabs-card > .ant-tabs-content > .ant-tabs-tabpane {
  background: #fff;
  padding: 0;
}

.card-container > .ant-tabs-card > .ant-tabs-bar {
  border-color: #fff;
}

.card-container > .ant-tabs-card > .ant-tabs-bar .ant-tabs-tab {
  border-color: transparent;
  background: transparent;
}

.card-container > .ant-tabs-card > .ant-tabs-bar .ant-tabs-tab-active {
  /* border-color: #fff; */
  background: #fff;
}

.ant-collapse > .ant-collapse-item > .ant-collapse-header {
  padding: 4px 8px;
  padding-left: 40px;
}
</style>