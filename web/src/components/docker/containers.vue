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
                  <!-- <span v-for="(port, index) in container.NetworkSettings.Ports" :key="index">{{port.length > 0 ? `${port[0].HostIp}:${port[0].HostPort}` : '无'}} -> {{index}}</span> -->
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

          <a-drawer :visible="create.visible" width="60%" @close="createClose">
            <a-form-model :model="create.form" v-bind="create.formItemLayoutWithOutLabel" style="height: 100%;">
              <a-form-model-item label="模板" v-bind="create.formItemLayout">
                <a-select v-model="create.template">
                  <a-select-option key="aaa" value="vvv">1111</a-select-option>
                </a-select>
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="镜像">
                <a-input v-model="create.form.name" />
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="端口随机">
                <a-switch v-model="create.form.delivery" />
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="重启策略">
                <a-radio-group v-model="create.form.restart_policy" size="small" button-style="solid">
                  <a-radio-button value="Unless stopped">
                    除非停止
                  </a-radio-button>
                  <a-radio-button value="always">
                    总是
                  </a-radio-button>
                  <a-radio-button value="onfailure">
                    失败
                  </a-radio-button>
                  <a-radio-button value="never">
                    从不
                  </a-radio-button>
                </a-radio-group>
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="高级选项">
                <a-checkbox-group v-model="create.form.type">
                  <a-checkbox value="1" name="type">
                    本地网络
                  </a-checkbox>
                </a-checkbox-group>
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="Resources">
                <a-radio-group v-model="create.form.resource">
                  <a-radio value="1">
                    Sponsor
                  </a-radio>
                  <a-radio value="2">
                    Venue
                  </a-radio>
                </a-radio-group>
              </a-form-model-item>
              <a-tabs type="card">
                <a-tab-pane key="1" tab="环境变量">
                  <a-form-model-item v-for="(env, index) in create.form.environment" :key="env.index" v-bind="index === 0 ? create.formItemLayout : {}"
                    :label="index === 0 ? '环境变量' : ''"
                    :prop="'environment.' + index + '.value'">
                    <a-input-group compact>
                      <a-input style=" width: 200px;" placeholder="KEY" />
                      <a-input v-model="env.key" style=" width: 30px; pointer-events: none; backgroundColor: #fff" placeholder="=" disabled/>
                      <a-input v-model="env.value" style="width: 200px;" placeholder="VALUE" />
                      <a-icon v-if="create.form.environment.length > 1" class="dynamic-delete-button" type="minus-circle-o"  :disabled="create.form.environment.length === 1" @click="removeDomain(env)" style="margin-left: 5px;margin-top: 7px"/>
                    </a-input-group>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayoutWithOutLabel">
                    <a-button type="dashed" style="width: 60%" @click="addDomain">
                      <a-icon type="plus" /> Add field
                    </a-button>
                  </a-form-model-item>
                </a-tab-pane>
                <a-tab-pane key="2" tab="磁盘卷组">
                  <p>Content of Tab Pane 2</p>
                  <p>Content of Tab Pane 2</p>
                  <p>Content of Tab Pane 2</p>
                </a-tab-pane>
                <a-tab-pane key="3" tab="网络配置">
                  <p>Content of Tab Pane 3</p>
                  <p>Content of Tab Pane 3</p>
                  <p>Content of Tab Pane 3</p>
                </a-tab-pane>
                <a-tab-pane key="4" tab="资源限制">
                  <p>Content of Tab Pane 3</p>
                  <p>Content of Tab Pane 3</p>
                  <p>Content of Tab Pane 3</p>
                </a-tab-pane>
              </a-tabs>
              
              <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
                <a-button type="primary" @click="onSubmit">
                  Create
                </a-button>
                <a-button style="margin-left: 10px;">
                  Cancel
                </a-button>
              </a-form-model-item>
            </a-form-model>
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
      },
      create: {
        form: {
          name: '',
          region: undefined,
          environment: [],
          delivery: false,
          type: [],
          resource: '',
          desc: '',
          restart_policy: '',
        },
        template: '',
        formItemLayoutWithOutLabel: {
          wrapperCol: {
            xs: { span: 24, offset: 0 },
            sm: { span: 16, offset: 3 }
          }
        },
        formItemLayout: {
          labelCol: {
            xs: { span: 24 },
            sm: { span: 3 },
          },
          wrapperCol: {
            xs: { span: 24 },
            sm: { span: 16 },
          },
        },
        visible: false
      }
    }
  },
  created: function () {
    this.container_refush()
  },
  methods: {
    create_container () {
      this.create.visible = true
    },
    container_refush () {
      let that = this
      this.$data.tableloading = true
      this.$api.docker.container.list({all: true}).then(function (response) {
        that.$data.data = response.detail
        that.$data.tableloading = false
      }).catch(function(response) {
        that.$message.error("获取容器信息失败："+response.message)
      })
    },
    container_start (row) {
      let that = this
      this.$api.docker.container.start(row.Id).then(function (response) {
        that.$message.success('启动成功')
        that.container_refush()
      }).catch(function (response) {
        that.$message.error(response.message)
      })
    },
    container_stop (row) {
      let that = this
      this.$api.docker.container.stop(row.Id).then(function (response) {
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
      this.drawerLoading = true
      this.$api.docker.container.info(id).then(function (response) {
        that.container = response.detail
        that.drawerLoading = false
      })
    },
    afterVisibleChange (val) {
      if (!val) {
        Object.assign(this.container, this.$options.data().container)
      }
    },
    onSubmit() {
      this.create.visible = false
    },
    removeDomain (item) {
      let index = this.create.form.environment.indexOf(item)
      if (index !== -1) {
        this.create.form.environment.splice(index, 1)
      }
    },
    addDomain () {
      let index = this.create.form.environment.length
      this.create.form.environment.push({
        value: '',
        index: index++,
        key: ''
      })
    },
    onClose () {
      this.visible = false
    },
    createClose() {
      this.create.visible = false
    }
  }
}
</script>
