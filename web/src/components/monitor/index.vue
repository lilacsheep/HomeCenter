<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <span style="float: right;">
        刷新间隔: 
        <el-select size="small" v-model="timesleep" placeholder="请选择" style="margin-left: 10px" @change="monitor_interval">
          <el-option label="1秒" value="1000"></el-option>
          <el-option label="10秒" value="10000"></el-option>
          <el-option label="30秒" value="30000"></el-option>
          <el-option label="1分钟" value="60000"></el-option>
        </el-select>
      </span>
      
    </el-col>
    <el-col :span="12" style="margin-top: 10px;">
      <el-card>
        <div slot="header" class="clearfix">
          <span>应用内存</span>
        </div>
        <ve-line height="300px" :settings="ChartSettings" :extend="extend" :data="memoryChart" :legend-visible="false"></ve-line>
      </el-card>
    </el-col>
    <el-col :span="12" style="margin-top: 10px;">
      <el-card>
        <div slot="header" class="clearfix">
          <span>CPU使用率</span>
        </div>
        <ve-line height="300px" :extend="extend" :data="cpuChart" :legend-visible="false"></ve-line>
      </el-card>
    </el-col>
    <el-col :span="12" style="margin-top: 10px;">
      <el-card>
        <div slot="header" class="clearfix">
          <span>网络连接</span>
        </div>
        <ve-line height="300px" :extend="extend" :data="connectionsChart" :legend-visible="false"></ve-line>
      </el-card>
    </el-col>
    <el-col :span="12" style="margin-top: 10px;">
      <el-card>
        <div slot="header" class="clearfix">
          <span>网络流量</span>
        </div>
        <ve-line height="300px" :extend="extend" :data="netChart" :settings="ChartSettings" :legend-visible="false"></ve-line>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      timesleep: '1000',
      ChartSettings: {
        yAxisType: ['0.0 b']
      },
      memoryChart: {
          columns: ['create_at', 'memory_size'],
          rows: []
      },
      cpuChart: {
          columns: ['create_at', 'cpu_percent'],
          rows: []
      },
      connectionsChart: {
          columns: ['create_at', 'connections'],
          rows: []
      },
      netChart: {
          columns: ['create_at', 'bytes_sent', 'bytes_recv'],
          rows: []
      },
      extend: {
        'xAxis.0.axisLabel.rotate': 45
      }
    }
  },
  methods: {
    filter_status(value, row) {
      return row.status === value;
    },
    get_info () {
      let that = this
      this.$api.get("/proxy/server/monitor").then(function (response) {
        if (that.memoryChart.rows.length > 120) {
          that.memoryChart.rows.shift()
          that.cpuChart.rows.shift()
          that.connectionsChart.rows.shift()
          that.netChart.rows.shift()
        }
        let item = response.detail.pop()
        let time = item.create_at.split(" ")[1]

        that.memoryChart.rows.push({"create_at": time, "memory_size": item.memory_size})
        that.cpuChart.rows.push({"create_at": time, "cpu_percent": item.cpu_percent})
        that.connectionsChart.rows.push({"create_at": time, "connections": item.connections})
        that.netChart.rows.push({"create_at": time, "bytes_sent": item.bytes_sent, 'bytes_recv': item.bytes_recv})
      })
    },
    monitor_interval: function (value) {
      clearInterval(this.timer)
      this.timer = setInterval(this.get_info, value)
    },
    refresh_info () {
      let that = this
      this.$api.get("/proxy/server/monitor").then(function (response) {
        response.detail.forEach(function (item) {
          let time = item.create_at.split(" ")[1]
          that.memoryChart.rows.push({"create_at": time, "memory_size": item.memory_size})
          that.cpuChart.rows.push({"create_at": time, "cpu_percent": item.cpu_percent})
          that.connectionsChart.rows.push({"create_at": time, "connections": item.connections})
          that.netChart.rows.push({"create_at": time, "bytes_sent": item.bytes_sent, 'bytes_recv': item.bytes_recv})
        })
      })
    },
    get_instance: function (id) {
      name = "默认转发"
      this.instanceData.forEach(function (item) {
        if (item.id === id) {
          name = item.address
        }
      })
      return name
    }
  },
  created: function () {
    this.refresh_info()
    this.timer = setInterval(this.get_info, 1000)
  },
  beforeDestroy () {
    clearInterval(this.timer)
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}

.el-card__body {
  padding: 0;
}

.el-dialog__header {
  padding: 10px 10px 5px;
  border-bottom: 1px solid whitesmoke;
}

.el-dialog__headerbtn {
  top: 12px;
}
.el-dialog__body {
  padding: 15px 10px;
}
.el-dialog__footer {
  border-top: 1px solid whitesmoke;
  padding: 5px 10px 10px;
}
</style>