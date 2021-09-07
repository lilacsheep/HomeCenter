<template>
    <a-layout style="padding: 0 12px 12px">
        <a-breadcrumb style="margin: 16px 0">
          <a-breadcrumb-item>Home</a-breadcrumb-item>
          <a-breadcrumb-item>仪表盘</a-breadcrumb-item>
        </a-breadcrumb>
        <a-layout-content :style="{padding: '10px', margin: 0, minHeight: '280px' }">
          <a-row :gutter="16">
            <a-col :span="6">
              <a-card>
                <a-statistic title="容器总数" :value="info.Containers">
                  <template #suffix>
                    <a-icon type="cloud-server" />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="运行中" :value="info.ContainersRunning">
                  <template #suffix>
                    <a-icon type="play-circle" />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="暂停" :value="info.ContainersPaused">
                  <template #suffix>
                    <a-icon type="pause-circle" />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card>
                <a-statistic title="停止" :value="info.ContainersStopped">
                  <template #suffix>
                    <a-icon type="stop" />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>

            <a-col :span="24" style="margin-top: 10px;" >
              <a-card :bodyStyle="{padding: '0px'}">
                <a-descriptions bordered >
                  <a-descriptions-item label="版本">
                    {{ info.ServerVersion}}
                  </a-descriptions-item>
                  <a-descriptions-item label="CPU个数">
                    {{ info.NCPU}} 个
                  </a-descriptions-item>
                  <a-descriptions-item label="内存">
                    {{ info.MemTotal | diskSize}}
                  </a-descriptions-item>
                  <a-descriptions-item label="Status" :span="3">
                    <a-badge status="processing" text="Running" />
                  </a-descriptions-item>
                  <a-descriptions-item label="操作系统">
                    {{ info.OSType }}
                  </a-descriptions-item>
                  <a-descriptions-item label="主机名">
                    {{ info.Name }}
                  </a-descriptions-item>
                  <a-descriptions-item label="内核版本">
                    {{ info.KernelVersion }}
                  </a-descriptions-item>
                  <a-descriptions-item label="镜像加速">
                    <div v-for="url in info.RegistryConfig.Mirrors" :key="url">
                      {{url}}
                    </div>
                  </a-descriptions-item>
                </a-descriptions>
              </a-card>
            </a-col>
          </a-row>
        </a-layout-content>
      </a-layout>
</template>

<script>
export default {
  data () {
    return {
      info: {
        RegistryConfig: {
          Mirrors: []
        }
      }
    }
  },
  created: function () {
    let that = this
    this.$api.post('/info/system').then(function (response) {
      that.$data.info = response.data
    })
  }
}
</script>

<style>
#card {
  float: none;
  text-align: right;
}
</style>
