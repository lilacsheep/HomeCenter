<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href="">
        控制面板
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :gutter="16" style="height: 100%">
      <a-col span="8">
        <a-descriptions bordered :loading="system.loading" style="background: #FFFFFF;">
          <a-descriptions-item label="主机名" span="3">{{system.info.host.hostname}}</a-descriptions-item>
          <a-descriptions-item label="运行时间" span="3">{{system.info.host.uptime|formatSeconds}}</a-descriptions-item>
          <a-descriptions-item label="核心数" span="3">{{system.info.cpu_info[0].cores}}</a-descriptions-item>
          <a-descriptions-item label="主频率" span="3">{{(system.info.cpu_info[0].mhz / 1000).toFixed(2)}} GHz</a-descriptions-item>
          <a-descriptions-item label="总内存" span="3">{{humanSize(system.info.memory.total)}}</a-descriptions-item>
          <a-descriptions-item label="操作系统" span="3">{{system.info.host.platform}}</a-descriptions-item>
          <a-descriptions-item label="系统架构" span="3">{{`${system.info.host.os} ${system.info.host.kernelArch}`}}</a-descriptions-item>
        </a-descriptions>
        <a-card :loading="system.loading" style="margin-top: 5px">
          <h3>性能指标</h3>
          <span>CPU使用率</span>
          <a-progress :percent="90" />
          <span>内存使用率</span>
          <a-progress :percent="system.info.memory.usedPercent.toFixed(2)"  status="active" />
        </a-card>
        <a-card :loading="system.loading" style="margin-top: 5px">
          <h3>磁盘信息</h3>
          <template v-for="(info, i) in system.info.disk">
            <span v-if="filter_fs(info)" :key="i">{{`${info.path} 可用: ${humanSize(info.free)} 共计: ${humanSize(info.total)}`}}</span>
            <a-progress v-if="filter_fs(info)" :key="i" :percent="GetPercent((info.total - info.free), info.total)" :stroke-color="percent_color(info)"/>
          </template>
        </a-card>
      </a-col>
      <a-col span="16">
        <a-table size="small" :columns="system.process.columns" :row-key="record => record.pid" :data-source="system.processes" style="background: #FFFFFF;">
          <span slot="cpu" slot-scope="text">{{text.toFixed(1)}}%</span>
          <span slot="mem" slot-scope="text">{{text.toFixed(1)}}%</span>
        </a-table>
      </a-col>
    </a-row>
  </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      system: {
        info: {
          host: {
            hostname: "",
            uptime: 0
          },
          disk: [],
          cpu_info: [{cores: 0, hmz: 0,}],
          memory:{
            usedPercent: 0
          }
        },
        processes: [],
        process: {
          columns: [
            {title: 'PID', dataIndex: 'pid',},
            {title: 'Name', dataIndex: 'name',},
            {title: '状态', dataIndex: 'status',},
            {title: 'CPU', sorter: (a, b) => a.cpu_percent - b.cpu_percent , dataIndex: 'cpu_percent', scopedSlots: { customRender: 'cpu' },},
            {title: '内存', sorter: (a, b) => a.mem_percent - b.mem_percent, dataIndex: 'mem_percent', scopedSlots: { customRender: 'mem' },},
          ]
        },
        loading: true,
      }
    };
  },
  created: function () {
    this.refreshSystem()
    this.refreshProcesses()
  },
  methods: {
    refreshSystem() {
      let that = this;
      this.$api.systemInfo().then(function (response) {
        that.system.info = response.detail
        that.system.loading = false
      }).then(function (response) {
        that.$message.info(response.message)
      })
    },
    refreshProcesses() {
      let that = this;
      this.$api.systemProcesses().then(function (response) {
        that.system.processes = response.detail
      })
    },
    humanSize(num, func) {
      if (num === 0) return 0, 'B';
      let k = 1024; //设定基础容量大小
      let sizeStr = ['B','KB','MB','GB','TB','PB','EB','ZB','YB']; //容量单位
      let i = 0; //单位下标和次幂
      for(let l=0;l<8;l++){   //因为只有8个单位所以循环八次
          if(num / Math.pow(k, l) < 1){ //判断传入数值 除以 基础大小的次幂 是否小于1，这里小于1 就代表已经当前下标的单位已经不合适了所以跳出循环
              break; //小于1跳出循环
          }
          i = l; //不小于1的话这个单位就合适或者还要大于这个单位 接着循环
      }
      if (func) {
        func((num / Math.pow(k, i)).toFixed(2), sizeStr[i])
      } else {
        return (num / Math.pow(k, i)).toFixed(2) + " " +sizeStr[i]
      }
    },
    filter_fs(s) {
      if (s.total == 0) {
        return false
      }
      let l = ["cgroup", "tmpfs", "binfmt_misc", "squashfs", "mqueue", "cgroupfs", "devpts"]
      l.forEach(function(item) {
        if (s.type == item) {
          return false
        }
      })
      return true
    },
    GetPercent(num, total) {
      num = parseFloat(num);
      total = parseFloat(total);
      if (isNaN(num) || isNaN(total)) {
          return "-";
      }
      return total <= 0 ? 0 : (Math.round(num / total * 10000) / 100.00);
    },
    percent_color(s) {
      let p = this.GetPercent((s.total - s.free), s.total)
      if (s < 60) {
        return "#99FF66"
      } else {
        if (s < 80) {
          return "#CCFF66"
        } else {
          return "#FF3366"
        }
      }
    }
  },
};
</script>

<style>
.ant-drawer-header {
  padding: 8px 12px;
}

.ant-drawer-body {
  padding: 12px;
}
.ant-list-bordered .ant-list-header {
  padding-left: 5px;
}
.ant-list-header {
  min-height: 20px;
  padding-top: 5px;
  padding-bottom: 5px;
}
.ant-list-bordered .ant-list-item {
  padding: 5px;
}

.ant-timeline-item {
  padding: 0;
}

.ant-card-head {
  background-color: #ddf3f5;
}

.ant-timeline-item-last > .ant-timeline-item-content {
  min-height: 0;
}

.ant-descriptions-bordered .ant-descriptions-item-content, .ant-descriptions-bordered .ant-descriptions-item-label {
  padding: 5px;
}

</style>