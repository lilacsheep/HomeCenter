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
            <a-table :columns="columns" :data-source="data" :loading="tableloading" size="small" style="background: #FFFFFF" :scroll="{ x: 1200 }" >
              <a-button slot="id" slot-scope="text" type="link" size="small" @click="container_detail(text)">
                {{ text.slice(0,12) }}
              </a-button>

              <span slot="Name" slot-scope="text">{{ text[0].slice(1) }}</span>
              <span slot="Created" slot-scope="text">{{ text * 1000 | dateformat }}</span>
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
            width="60%"
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
                <a-descriptions-item label="环境变量" :span="3">
                  <li v-for="(env, index) in container.Config.Env" :key="index">
                    {{env}}
                    </li>
                </a-descriptions-item>
                <a-descriptions-item label="端口" :span="3">
                  <li v-for="(port, index) in container.NetworkSettings.Ports" :key="index">
                    {{port ? `${port[0].HostIp}:${port[0].HostPort}` : '无'}} -> {{index}}</li>
                </a-descriptions-item>
                <a-descriptions-item label="卷组" :span="3">
                  <li v-for="(volume, index) in container.Mounts" :key="index">
                    <template v-if="volume.Source.length >= 20">
                      <a-tooltip>
                        <template slot="title">
                          {{volume.Source}}
                        </template>
                        {{`${volume.Source.slice(0, 20)}...  ->  ${volume.Destination}`}}
                      </a-tooltip>
                    </template>
                    <template v-else>
                      {{`${volume.Source} -> ${volume.Destination}`}}
                    </template>                    
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
              <a-form-model-item v-bind="create.formItemLayout" label="名称">
                <a-input v-model="create.form.name" />
              </a-form-model-item>
              <a-form-model-item label="镜像" v-bind="create.formItemLayout">
                <a-select v-model="create.form.config.Image">
                  <a-select-option v-for="i in create.images" :key="i.Id" :value="i.RepoTags[0]">{{i.RepoTags[0]}}</a-select-option>
                </a-select>
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="命令">
                <a-input v-model="create.form.config.Cmd" />
              </a-form-model-item>
              <a-form-model-item v-bind="create.formItemLayout" label="重启策略">
                <a-radio-group v-model="create.form.host_config.RestartPolicy.Name" size="small" button-style="solid">
                  <a-radio-button value="unless-stopped">
                    除非停止
                  </a-radio-button>
                  <a-radio-button value="always">
                    总是
                  </a-radio-button>
                  <a-radio-button value="on-failure">
                    失败
                  </a-radio-button>
                  <a-radio-button value="no">
                    从不
                  </a-radio-button>
                </a-radio-group>
              </a-form-model-item>
              <a-tabs type="card">
                <a-tab-pane key="1" tab="环境变量">
                  <a-form-model-item v-for="(env, index) in create.form.config.Env" :key="env.index" v-bind="index === 0 ? create.formItemLayout : {}"
                    :label="index === 0 ? '环境变量' : ''"
                    :prop="'Env.' + index + '.value'">
                    <a-input-group compact>
                      <a-input v-model="env.key" style=" width: 200px;" placeholder="KEY" />
                      <a-input style=" width: 30px; pointer-events: none; backgroundColor: #fff" placeholder="=" disabled/>
                      <a-input v-model="env.value" style="width: 200px;" placeholder="VALUE" />
                      <a-icon v-if="create.form.config.Env.length > 1" class="dynamic-delete-button" type="minus-circle-o"  :disabled="create.form.config.Env.length === 1" @click="removeDomain(env)" style="margin-left: 5px;margin-top: 7px"/>
                    </a-input-group>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayoutWithOutLabel">
                    <a-button type="dashed" style="width: 60%" @click="addDomain">
                      <a-icon type="plus" /> 新增变量
                    </a-button>
                  </a-form-model-item>
                </a-tab-pane>
                <a-tab-pane key="2" tab="磁盘卷组">
                  <a-form-model-item v-for="(b, index) in create.form.host_config.Mounts" :key="b.index" v-bind="index === 0 ? create.formItemLayout : {}"
                    :label="index === 0 ? '绑定卷组' : ''"
                    :prop="'Mounts.' + index + '.value'">
                    <a-input-group compact>
                      <a-input v-model="b.Source" style=" width: 200px;" placeholder="源路径"/>
                      <a-input style="width: 30px; pointer-events: none; backgroundColor: #fff" placeholder=":" disabled/>
                      <a-input v-model="b.Target" style="width: 200px;" placeholder="目标路径"/>
                      <a-checkbox v-model="b.ReadOnly" style="margin: 5px 0 0 5px;">只读</a-checkbox>
                      <a-icon v-if="create.form.host_config.Mounts.length > 1" class="dynamic-delete-button" type="minus-circle-o"  :disabled="create.form.host_config.Mounts.length === 1" @click="removeBind(b)" style="margin-left: 5px;margin-top: 7px"/>
                    </a-input-group>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayoutWithOutLabel">
                    <a-button type="dashed" style="width: 60%" @click="addBind">
                      <a-icon type="plus" /> 新增卷组
                    </a-button>
                  </a-form-model-item>
                </a-tab-pane>
                <a-tab-pane key="3" tab="网络配置">
                  <a-form-model-item v-bind="create.formItemLayout" label="端口随机">
                    <a-switch v-model="create.form.host_config.PublishAllPorts" @change="publish_all_ports_change"/>
                  </a-form-model-item>
                  <a-form-model-item :style="create.network_mode_style" v-bind="create.formItemLayout" label="网络模式">
                    <a-radio-group v-model="create.form.network.mode" size="small" button-style="solid" @change="network_change">
                      <a-radio-button value="host">
                        本地网络
                      </a-radio-button>
                      <a-radio-button value="bridge">
                        桥接模式
                      </a-radio-button>
                      <a-radio-button value="none">
                        None
                      </a-radio-button>
                    </a-radio-group>
                  </a-form-model-item>
                  <a-form-model-item :style="create.network_style" v-for="(port, index) in create.form.network.ports" :key="port.index" v-bind="index === 0 ? create.formItemLayout : {}"
                    :label="index === 0 ? '端口映射' : ''"
                    :prop="'ports.' + index + '.value'">
                    <a-input-group compact>
                      <a-select v-model="port.host.proto">
                        <a-select-option value="tcp">
                          TCP
                        </a-select-option>
                        <a-select-option value="udp">
                          UDP
                        </a-select-option>
                      </a-select>
                      <a-input type="number" style="width: 100px" v-model="port.host.HostPort" />
                      <a-icon type="arrow-right" style="margin-top: 8px;margin-left: 5px;margin-right: 5px" />
                      <a-select v-model="port.proto" >
                        <a-select-option value="tcp">
                          TCP
                        </a-select-option>
                        <a-select-option value="udp">
                          UDP
                        </a-select-option>
                      </a-select>
                      <a-input type="number" style="width: 100px" v-model="port.port" />
                       <a-icon v-if="create.form.network.ports.length > 1" class="dynamic-delete-button" type="minus-circle-o"  :disabled="create.form.network.ports.length === 1" @click="removePort(env)" style="margin-left: 5px;margin-top: 7px"/>
                    </a-input-group>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayoutWithOutLabel" :style="create.network_style">
                    <a-button type="dashed" style="width: 60%" @click="addPort">
                      <a-icon type="plus" /> 添加端口
                    </a-button>
                  </a-form-model-item>
                </a-tab-pane>
                <a-tab-pane key="4" tab="高级配置">
                  <a-form-model-item v-bind="create.formItemLayout" label="内存限制">
                    <a-input type="number" addon-after="MB" default-value="0" style="width: 200px;"/>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayout" label="swap限制">
                    <a-input type="number" addon-after="MB" default-value="-1" style="width: 200px;"/>
                  </a-form-model-item>
                  <a-form-model-item v-bind="create.formItemLayout" label="CPU限制">
                    <a-input-number :max="20" :min="0" default-value="0" style="width: 200px;" :formatte="value => `${value}个`"/>
                  </a-form-model-item>
                </a-tab-pane>
              </a-tabs>
              
              <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
                <a-button type="primary" @click="onSubmit">
                  创 建
                </a-button>
                <a-button style="margin-left: 10px;">
                  取 消
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
    scopedSlots: { customRender: 'id' },
    fixed: 'left',
    width: 120
  },
  {
    title: '镜像',
    dataIndex: 'Image',
    key: 'Image',
    width: 300
  },
  // {
  //   title: '命令',
  //   dataIndex: 'Command',
  //   key: 'Command'
  // },
  {
    title: '创建时间',
    key: 'Created',
    dataIndex: 'Created',
    scopedSlots: { customRender: 'Created' },
    width: 150
  },
  {
    title: '状态',
    key: 'Status',
    dataIndex: 'Status',
    width: 120
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
    scopedSlots: { customRender: 'action' },
    fixed: 'right'
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
          },
        }
      },
      create: {
        images: [],
        form: {
          config: {
            Image: "",
            Env: [],
            Cmd: ""
          },
          host_config: {
            RestartPolicy: {Name: "no", MaximumRetryCount: 3},
            PublishAllPorts: false,
            PortBindings: {},
            Mounts: [
              {Type: 'bind', Source: '', Target: '', ReadOnly: false, index: 0}
            ],
          },
          network_config: {},
          name: '',
          bind: [],
          type: [],
          resource: '',
          desc: '',
          restart_policy: '',
          network: {
            mode: 'bridge',
            ports: [
              {port: 0, proto: 'tcp', host: {HostPort: 0, proto: 'tcp', HostIP: '0.0.0.0'} ,index: 0},
            ]
          },
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
        visible: false,
        network_style: "",
        network_mode_style: ""
      }
    }
  },
  created: function () {
    this.container_refush()
  },
  methods: {
    create_container () {
      this.get_images()
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
      let params = {
        config: {
          Image: this.create.form.config.Image,
          Env: [],
          Cmd: [],
        },
        host_config: {
          RestartPolicy: this.create.form.host_config.RestartPolicy,
          PublishAllPorts: false,
          PortBindings: {},
          Mounts: [],
        },
      }
      let this_ = this
      this.create.visible = false
      this.create.form.config.Env.forEach(element => {
        params.config.Env.push(`${element.key}=${element.value}`)
      });
      this.create.form.host_config.Mounts.forEach(element => {
        if ((element.Source != "") && (element.Target != "")){
          params.host_config.Mounts.push(element)
        }
      });
      this.create.form.network.ports.forEach(element => {
        if (element.port != 0) {
          let port = `${element.port}/${element.proto}`
          params.host_config.PortBindings[port] = {HostIP: element.host.HostIP, HostPort: `${element.host.HostPort}/${element.host.proto}`}
        }
      });
      if (this.create.form.host_config.RestartPolicy.Name != "on-failure") {
        params.host_config.RestartPolicy.MaximumRetryCount = 0
      }
      if (this.create.form.config.Cmd != "") {
        params.config.Cmd = this.create.form.config.Cmd.split(" ")
      }
      if (this.create.form.name != "") {
        params.Name = this.create.form.name 
      }
      this.$api.docker.container.create(params).then(function(response) {
        this_.create.visible = false
        this_.$message.info("创建成功")
      }).catch(function(response) {
        this_.$message.error("创建失败: "+response.message)
      })
    },
    removeDomain (item) {
      let index = this.create.form.config.Env.indexOf(item)
      if (index !== -1) {
        this.create.form.config.Env.splice(index, 1)
      }
    },
    removeBind (item) {
      let index = this.create.form.host_config.Mounts.indexOf(item)
      if (index !== -1) {
        this.create.form.host_config.Mounts.splice(index, 1)
      }
    },
    removePort (item) {
      let index = this.create.form.network.ports.indexOf(item)
      if (index !== -1) {
        this.create.form.network.ports.splice(index, 1)
      }
    },
    addDomain () {
      let index = this.create.form.config.Env.length
      this.create.form.config.Env.push({
        value: '',
        index: index++,
        key: ''
      })
    },
    addBind() {
      let index = this.create.form.host_config.Mounts.length
      this.create.form.host_config.Mounts.push(
        {Type: 'bind', Source: '', Target: '', ReadOnly: false, index: index++}
      )
    },
    addPort() {
      let index = this.create.form.network.ports.length
      this.create.form.network.ports.push(
        {port: 0, proto: 'tcp', host: {HostPort: 0, proto: 'tcp', HostIP: ''} ,index: index++}
      )
    },
    onClose () {
      this.visible = false
    },
    createClose() {
      this.create.visible = false
      Object.assign(this.create.form, this.$options.data().create.form)
    },
    network_change(event) {
      if (this.create.form.network.mode == "bridge") {
        this.create.network_style = ""
      } else {
        this.create.network_style = "display: none"
      }
    },
    publish_all_ports_change(checked, event) {
      if (checked) {
        this.create.network_mode_style = "display: none"
        this.create.network_style = "display: none"
      } else {
        this.create.network_mode_style = ""
        this.network_change()
      }
    },
    get_images() {
      let this_ = this
      this.$api.docker.images.list().then(function (response) {
        this_.create.images = response.detail
      }).catch(function (response) {
        this_.$message.error('获取镜像失败: '+response.message)
      })
    }
  }
}
</script>
