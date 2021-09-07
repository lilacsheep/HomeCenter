<template>
    <a-layout style="padding: 0 12px 12px">
      <a-page-header
        title="Title"
        sub-title="This is a subtitle"
        @back="() => $router.go(-1)"
      >
      </a-page-header>
        <a-layout-content :style="{padding: '10px', margin: 0, minHeight: '280px' }">
          <a-row :gutter="16">
            <a-col :span="24">
              <ve-line :data="chartData" :settings="chartSettings" height="220px"></ve-line>
            </a-col>
            <a-col :span="24">
              <ve-line :data="chartData"  height="220px"></ve-line>
            </a-col>
            <a-col :span="24">
              <ve-line :data="chartData"  height="220px"></ve-line>
            </a-col>
          </a-row>
        </a-layout-content>
      </a-layout>
</template>

<script>
export default {
  data () {
    this.chartSettings = {
      dataType: {
        '使用率': 'percent'
      }
    }
    return {
      info: {
        Name: '',
        Env: [],
        Mounts: [],
        Created: '',
        Config: {
          Hostname: '',
          Image: '',
          Volumes: {}
        },
        State: {
          Pid: '',
          RestartCount: '',
          FinishedAt: '',
          StartedAt: '',
          Error: ''
        },
        NetworkSettings: {
          IPAddress: '',
          MacAddress: '',
          Ports: []
        },
        HostConfig: {
          RestartPolicy: {
            Name: '',
            MaximumRetryCount: 0,
            Privileged: ''
          }
        }
      },
      stats: {
        read: '',
        memory_stats: {
          usage: 0,
          max_usage: 0,
          stats: {}
        },
        cpu_stats: {
          cpu_usage: {
            total_usage: 0,
            percpu_usage: [],
            system_cpu_usage: 0,
            online_cpus: 0
          }
        },
        limit: 0,
        networks: {
          eth0: {
            rx_bytes: 0,
            rx_packets: 0,
            rx_errors: 0,
            rx_dropped: 0,
            tx_bytes: 0,
            tx_packets: 0,
            tx_errors: 0,
            tx_dropped: 0
          }
        }
      },
      chartData: {
        columns: ['日期', '使用率'],
        rows: [
          { '日期': '21:48:05', '使用率': 0.12 }
        ]
      }
    }
  },
  methods: {},
  mounted: function () {
    let that = this
    this.$api.post('/containers/info', {id: this.$route.params.id}).then(function (response) {
      that.$data.info = response.data
    })
    this.$api.post('/containers/stats', {id: this.$route.params.id}).then(function (response) {
      that.$data.stats = JSON.parse(response.data)
    })
  }
}
</script>

<style>
#card {
  text-align: center;
}
.ant-card-head {
  min-height: 20px!important;
}
.ant-card-head .ant-card-head-title {
  padding: 5px 0!important;
}
</style>
