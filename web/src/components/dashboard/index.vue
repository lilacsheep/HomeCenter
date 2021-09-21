<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href="">
        控制面板
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :gutter="16" style="height: 100%">
      <a-col :span="8">
        <a-descriptions bordered :loading="system.loading" style="background: #FFFFFF;">
          <a-descriptions-item label="主机名" :span="3">{{system.info.host.hostname}}</a-descriptions-item>
          <a-descriptions-item label="运行时间" :span="3">{{system.info.host.uptime|formatSeconds}}</a-descriptions-item>
          <a-descriptions-item label="核心数" :span="3">{{cpu_core_num(system)}}</a-descriptions-item>
          <a-descriptions-item label="主频率" :span="3">{{(system.info.cpu_info[0].mhz / 1000).toFixed(2)}} GHz</a-descriptions-item>
          <a-descriptions-item label="总内存" :span="3">{{humanSize(system.info.memory.total)}}</a-descriptions-item>
          <a-descriptions-item label="操作系统" :span="3">{{system.info.host.platform}}</a-descriptions-item>
          <a-descriptions-item label="系统架构" :span="3">{{`${system.info.host.os} ${system.info.host.kernelArch}`}}</a-descriptions-item>
        </a-descriptions>
        <a-card :loading="system.loading" style="margin-top: 5px">
          <h3>性能指标</h3>
          <span>CPU使用率</span>
          <a-progress :percent="system.info.cpu_percent.toFixed(2)" :stroke-color="percent_color(system.info.cpu_percent)"/>
          <span>内存使用率</span>
          <a-progress :percent="system.info.memory.usedPercent.toFixed(2)"  status="active" />
        </a-card>
        <a-card :loading="system.loading" style="margin-top: 5px">
          <h3>磁盘信息</h3>
          <template v-for="(info, i) in system.info.disk">
            <span v-if="filter_fs(info)" :key="i">{{`${info.path} 可用: ${humanSize(info.free)} 共计: ${humanSize(info.total)}`}}</span>
            <a-progress v-if="filter_fs(info)" :key="i" :percent="GetPercent((info.total - info.free), info.total)" :stroke-color="fs_percent_color(info)"/>
          </template>
        </a-card>
      </a-col>
      <a-col :span="16">
        <a-table size="small" :columns="system.process.columns" :row-key="record => record.pid" :data-source="system.processes" style="background: #FFFFFF;">
          <span slot="pid" slot-scope="text">
            <a-button type="link" @click="open_process_info(text)">{{text}}</a-button>
          </span>
          <span slot="cpu" slot-scope="text">{{text.toFixed(1)}}%</span>
          <span slot="mem" slot-scope="text">{{text.toFixed(1)}}%</span>
        </a-table>
      </a-col>
    </a-row>
    <a-drawer :visible="system.process.visiable" :title="system.process.title" @close="close_process_info" width="60%">
      <a-spin :spinning="system.process.spinning">
        <a-descriptions size="small" bordered>
          <a-descriptions-item label="名称" :span="3">
            {{system.process.info.name}}
          </a-descriptions-item>
          <a-descriptions-item label="状态" :span="3">
            <a-tag v-if="system.process.info.status == ''">未知</a-tag>
            <a-badge v-else-if="system.process.info.status == 'R'" status="processing" text="运行" />
            <a-tag v-else-if="system.process.info.status == 'S'">休眠</a-tag>
            <a-tag v-else-if="system.process.info.status == 'T'">停止</a-tag>
            <a-tag v-else-if="system.process.info.status == 'I'">空闲</a-tag>
            <a-tag v-else-if="system.process.info.status == 'Z'">僵死</a-tag>
            <a-tag v-else-if="system.process.info.status == 'W'">等待</a-tag>
            <a-tag v-else-if="system.process.info.status == 'L'">锁</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="内存信息" :span="3">
            {{system.process.info.mem_percent ? system.process.info.mem_percent.toFixed(1): 0}}%
          </a-descriptions-item>
          <a-descriptions-item label="Status" :span="3">
            <a-badge status="processing" text="Running" />
          </a-descriptions-item>
          <a-descriptions-item label="网络连接" :span="3">
            <span v-for="e, index in system.process.info.connections" :key="index">{{e}}<br /></span>
          </a-descriptions-item>
          <a-descriptions-item label="环境变量" :span="3">
            <span v-for="e, index in system.process.info.env" :key="index">{{e}}<br /></span>
          </a-descriptions-item>
        </a-descriptions>
      </a-spin>
    </a-drawer>
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
            uptime: 0,
            os: ''
          },
          disk: [],
          cpu_info: [{cores: 0, hmz: 0,}],
          cpu_percent: 0,
          memory:{
            usedPercent: 0,
            total: 0
          },
        },
        processes: [],
        process: {
          columns: [
            {title: 'PID', dataIndex: 'pid', scopedSlots: { customRender: 'pid' }},
            {title: 'Name', dataIndex: 'name',},
            {title: '状态', dataIndex: 'status',},
            {title: 'CPU', sorter: (a, b) => a.cpu_percent - b.cpu_percent , dataIndex: 'cpu_percent', scopedSlots: { customRender: 'cpu' },},
            {title: '内存', sorter: (a, b) => a.mem_percent - b.mem_percent, dataIndex: 'mem_percent', scopedSlots: { customRender: 'mem' },},
          ],
          visiable: false,
          title: "",
          info: {},
          spinning: true
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
      this.$api.system.info().then(function (response) {
        that.system.info = response.detail
        that.system.loading = false
      }).then(function (response) {
        that.$message.info(response.message)
      })
    },
    refreshProcesses() {
      let that = this;
      this.$api.system.processes().then(function (response) {
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
      if (s == null) {
        return false
      }
      if (s.total == 0) {
        return false
      }
      let v = true
      let l = ["cgroup", "tmpfs", "binfmt_misc", "squashfs", "mqueue", "cgroupfs", "devpts", ""]
      l.forEach(function(item) {
        if (s.fstype == item) {
          v = false
        }
      })
      return v
    },
    GetPercent(num, total) {
      num = parseFloat(num);
      total = parseFloat(total);
      if (isNaN(num) || isNaN(total)) {
          return "-";
      }
      return total <= 0 ? 0 : (Math.round(num / total * 10000) / 100.00);
    },
    percent_color(num) {
      if (num < 60) {
        return "#99FF66"
      } else {
        if (num < 80) {
          return "#CCFF66"
        } else {
          return "#FF3366"
        }
      }
    },
    fs_percent_color(s) {
      let p = this.GetPercent((s.total - s.free), s.total)
      return this.percent_color(p)
    },
    open_process_info(pid) {
      let this_ = this
      this.$api.system.process(pid).then(function(response) {
        this_.system.process.spinning = false
        this_.system.process.visiable = true
        this_.system.process.info = response.detail
      }).catch(function(response) {
        this_.$message.error("获取进程信息错误: "+response.message)
      })
    },
    close_process_info() {
      this.system.process.visiable = false
    },
    cpu_core_num(row) {
      if (row.info.host.os == 'linux') {
        return row.info.cpu_info.length
      } else {
        return row.info.cpu_info[0].cores
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

.ant-descriptions-bordered .ant-descriptions-item-content, .ant-descriptions-bordered .ant-descriptions-item-label {
  padding: 5px;
}

</style>