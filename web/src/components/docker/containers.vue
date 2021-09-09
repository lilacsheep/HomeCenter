<template>
    <a-layout style="padding: 0 12px 12px">
      <a-breadcrumb style="margin: 16px 0">
        <a-breadcrumb-item>主页</a-breadcrumb-item>
        <a-breadcrumb-item>容器相关</a-breadcrumb-item>
        <a-breadcrumb-item>容器列表</a-breadcrumb-item>
      </a-breadcrumb>
      <a-layout-content :style="{ background: '#fff', padding: '12px', margin: 0, minHeight: '280px' }">
        <a-row :gutter="16">
          <a-col :span="18">
            <a-button-group>
              <a-button type="primary" @click="create_container">新增容器</a-button>
            </a-button-group>
          </a-col>
          <a-col :span="6">
            <!--  @search="onSearch" -->
            <a-input-search placeholder="input search text" enter-button />
          </a-col>
          <a-col :span="24" style="margin-top: 10px;">
            <a-table :columns="columns" :data-source="data" :loading="tableloading" size="small" style="background: #FFFFFF">
              <a-button slot="id" slot-scope="text" type="link" size="small" @click="container_detail(text)">
                {{ text.slice(0,12) }}
              </a-button>

              <span slot="Name" slot-scope="text">{{ text[0].slice(1) }}</span>
              <span slot="Created" slot-scope="text">{{ text * 1000 | dateformat }}</span>
              <span slot="Ports" slot-scope="Ports">{{ Ports.length > 0 ? `${Ports[0].IP}:${Ports[0].PublicPort}/${Ports[0].Type}` : "Null"}}</span>
              <span slot="action" slot-scope="row">
                <a-button-group size='small'>
                  <a-button icon='edit'></a-button>
                  <a-button icon='area-chart' type="primary" @click="container_stats(row)"></a-button>
                  <a-popconfirm v-if="row.State != 'running'" title="确认启动?" @confirm="() => container_start(row)">
                    <a-button icon='play-square' type="primary"></a-button>
                  </a-popconfirm>
                  <a-popconfirm v-else title="确认停止?" @confirm="() => container_stop(row)">
                    <a-button icon='stop' type="primary"></a-button>
                  </a-popconfirm>
                </a-button-group>
              </span>
            </a-table>
          </a-col>
        </a-row>
        <a-drawer
            :title="`${container.Name.slice(1)} 信息`"
            placement="right"
            :closable="false"
            :visible="visible"
            :after-visible-change="afterVisibleChange"
            @close="onClose"
            :width="520"
            :get-container="false"
          >
          <a-spin :spinning="drawerLoading">
            <div class="spin-content">
              <a-descriptions bordered layout="vertical" size='small'>
                <a-descriptions-item label="容器名">
                  {{ container.Name.slice(1) }}
                </a-descriptions-item>
                <a-descriptions-item label="状态">
                  <a-badge v-if="container.State.Running" status="processing" text="运行中" />
                  <a-badge v-else status="warning" text="停止" ></a-badge>
                </a-descriptions-item>
                <a-descriptions-item label="创建时间">
                  {{ container.Created | dateformat}}
                </a-descriptions-item>
                <a-descriptions-item label="启动时间">
                  {{ container.State.StartedAt | dateformat}}
                </a-descriptions-item>
                <a-descriptions-item label="主机名" :span="2">
                  {{ container.Config.Hostname}}
                </a-descriptions-item>
                <a-descriptions-item label="镜像" :span="3">
                  {{ container.Config.Image}}
                </a-descriptions-item>
                <a-descriptions-item label="ENTRYPOINT" :span="1">
                  {{container.Config.Cmd ? container.Config.Cmd.join(" ") : "无"}}
                </a-descriptions-item>
                <a-descriptions-item label="CMD" :span="2">
                  {{ container.Config.Entrypoint ? container.Config.Entrypoint.join(" ") : "无" }}
                </a-descriptions-item>
                <a-descriptions-item label="地址">
                  {{ container.NetworkSettings.IPAddress}}
                </a-descriptions-item>
                <a-descriptions-item label="网关">
                  {{ container.NetworkSettings.Gateway }}
                </a-descriptions-item>
                <a-descriptions-item label="MAC地址">
                  {{ container.NetworkSettings.MacAddress}}
                </a-descriptions-item>
                <a-descriptions-item label="端口" :span="3">
                  <span v-for="(port, index) in container.NetworkSettings.Ports" :key="index">{{port.length > 0 ? `${port[0].HostIp}:${port[0].HostPort}` : '无'}} -> {{index}}</span>
                </a-descriptions-item>
                <a-descriptions-item label="卷组" :span="3">
                  <li v-for="(volume, index) in container.Mounts" :key="index">
                    {{`${volume.Source.slice(0, 20)} -> ${volume.Destination}`}}
                    </li>
                </a-descriptions-item>
              </a-descriptions>
            </div>
          </a-spin>
          </a-drawer>
      </a-layout-content>
  </a-layout>
</template>

<script>
const columns = [
  {
    title: 'ID',
    dataIndex: 'Id',
    key: 'Id',
    scopedSlots: { customRender: 'id' }
  },
  {
    title: '镜像',
    dataIndex: 'Image',
    key: 'Image'
  },
  {
    title: '命令',
    dataIndex: 'Command',
    key: 'Command'
  },
  {
    title: '创建时间',
    key: 'Created',
    dataIndex: 'Created',
    scopedSlots: { customRender: 'Created' }
  },
  {
    title: '状态',
    key: 'Status',
    dataIndex: 'Status'
  },
  {
    title: '端口',
    key: 'Ports',
    dataIndex: 'Ports',
    scopedSlots: { customRender: 'Ports' }
  },
  {
    title: '名称',
    key: 'Names',
    dataIndex: 'Names',
    scopedSlots: { customRender: 'Name' }
  },
  {
    title: '操作',
    key: 'action',
    scopedSlots: { customRender: 'action' }
  }
]

export default {
  data () {
    return {
      columns,
      tableloading: false,
      drawerLoading: false,
      visible: false,
      data: [],
      container: {
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
      }
    }
  },
  created: function () {
    this.container_refush()
  },
  methods: {
    create_container () {
      this.$router.push({path: '/container/create'})
    },
    container_refush () {
      let that = this
      this.$data.tableloading = true
      this.$api.docker_container_list({all: true}).then(function (response) {
        that.$data.data = response.detail
        that.$data.tableloading = false
      }).catch(function(response) {
        that.$message.error("获取容器信息失败："+response.message)
      })
    },
    container_start (row) {
      let that = this
      this.$api.docker_container_start({id: row.Id}).then(function (response) {
        that.$message.success('启动成功')
        that.container_refush()
      }).catch(function (response) {
        that.$message.error(response.message)
      })
    },
    container_stop (row) {
      let that = this
      this.$api.docker_container_stop({id: row.Id}).then(function (response) {
        that.$message.success('停止成功')
        that.container_refush()
      }).catch(function (response) {
        that.$message.error(response.message)
      })
    },
    container_stats (row) {
      this.$router.push({path: '/docker/container/' + row.Id})
    },
    container_detail (id) {
      this.visible = true
      let that = this
      this.$data.drawerLoading = true
      this.$api.docker_container_info({id: id}).then(function (response) {
        that.$data.container = response.data
        that.$data.drawerLoading = false
      })
    },
    afterVisibleChange (val) {
      if (!val) {
        Object.assign(this.$data.container, this.$options.data().container)
      }
    },
    onClose () {
      this.visible = false
    }
  }
}
</script>
