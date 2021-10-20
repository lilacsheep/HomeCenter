<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href="">
        文件下载
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :gutter="16" style="border: 1px solid #f4f4f4;height: 100%">
      <a-col>
        <a-tabs default-active-key="1" @change="tabClick" tab-position="left">
          <a-tab-pane key="1" tab="下载列表">
            <a-button-group>
              <a-button type="primary" icon="edit" @click="download.create.visit = true">创建下载</a-button>
              <a-upload name="file" :loading="download.uploadloading" :multiple="false" @change="uploadChange" accept=".torrent" action="/api/download/torrent">
                <a-button  type="primary" icon="upload">上传种子</a-button>
              </a-upload>
            </a-button-group>
            <span style="float: right;">
              <a-radio-group v-model="task.query.status" size="small" button-style="solid">
                <a-radio-button value="全部">全部</a-radio-button>
                <a-radio-button value="下载中">下载中</a-radio-button>
                <a-radio-button value="已完成">已完成</a-radio-button>
                <a-radio-button value="其他">其他</a-radio-button>
              </a-radio-group>
              <a-tag color="green" size="small"><a-icon type="arrow-up"></a-icon>{{global.upload | diskSize}}/秒</a-tag>
              <a-tag color="orange" size="small"><a-icon type="arrow-down"></a-icon>{{global.download | diskSize}}/秒</a-tag>
            </span>
            <a-table :data-source="download.tasks" size="small" :columns="download.columns" style="margin-top: 10px;background: #FFFFFF">
              <span slot="gid" slot-scope="text, record">
                <a-button type="link" @click="taskInfoOpen(record)">{{getTaskName(record)}}</a-button>
              </span>
              <span slot="totalLength" slot-scope="text, record">
                {{record.totalLength | diskSize}}
              </span>
              <span slot="status" slot-scope="text, record">
                <span v-if="record.status == 'active'">{{(record.completedLength / record.totalLength * 100).toFixed(2)}}%</span>
                  <span v-else-if="record.status == 'stopped'">已停止</span>
                  <span v-else-if="record.status == 'paused'">已暂停</span>
                  <span v-else-if="record.status == 'complete'">已完成</span>
                  <span v-else>{{record.status}}</span>
              </span>
              <span slot="downloadSpeed" slot-scope="text, record">
                {{record.downloadSpeed | diskSize}}/秒
              </span>
              <span slot="uploadSpeed" slot-scope="text, record">
                {{record.uploadSpeed | diskSize}}/秒
              </span>
              <span slot="action" slot-scope="text, record">
                <a-popconfirm v-if="record.status == 'paused'" title="是否继续该任务？" @confirm="start_task(record)">
                  <a-button  style="color: green" type="link" icon="play-circle"></a-button>
                </a-popconfirm>
                <a-popconfirm v-if="record.status == 'error'" title="是否继续该任务？" @confirm="start_task(record)">
                  <a-button style="color: green" type="link" icon="play-circle"></a-button>
                </a-popconfirm>
                <a-popconfirm v-if="record.status == 'active'" title="是否暂停该任务？" @confirm="cancel_task(record)">
                  <a-button style="color: red" type="link" icon="pause-circle"></a-button>
                </a-popconfirm>
                <a-popconfirm title="是否删除该任务？" @confirm="remove_task(record)">
                  <a-button style="color: red" type="link" icon="delete"></a-button>
                </a-popconfirm>
              </span>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="2" tab="配置管理">
            <a-card style="padding: 0" >
              <a-form-model :model="settings.form" :label-col="labelCol" :wrapper-col="wrapperCol">
                <a-form-model-item label="启动">
                  <a-switch v-model="settings.form.auto_start">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">检查并启动Aria2</span>
                </a-form-model-item>
                <a-form-model-item label="RPC端口">
                  <a-input-number :min="80" :max="65535" v-model="settings.form.port"/>
                  <span style="color: #909399;font-size: 9px;margin-left: 5px ;">RPC连接端口</span>
                </a-form-model-item>
                <a-form-model-item label="下载端口">
                  <a-input addon-before="TCP" type="number" :min="80" :max="65535" style="width: 150px" v-model="settings.form.tcp_port"/>
                  <a-input addon-before="UDP" type="number" :min="80" :max="65535" style="width: 150px;margin-left: 5px;" v-model="settings.form.udp_port"/>
                  <span style="color: #909399;font-size: 9px;margin-left: 5px ;">Aria2下载端口</span>
                </a-form-model-item>
                <a-form-model-item label="BT端口">
                  <a-input-number :min="80" :max="65535" v-model="settings.form.bt_port"/>
                  <span style="color: #909399;font-size: 9px;margin-left: 5px ;">DHT和BT监听端口</span>
                </a-form-model-item>
                <a-form-model-item label="RPC Token">
                  <a-input v-model="settings.form.token" style="width: 300px"></a-input>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">Aria2 token</span>
                </a-form-model-item>
                <a-form-model-item label="WebUI">
                  <a-switch v-model="settings.form.webui">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">启用Aria2NG界面</span>
                </a-form-model-item>
                <a-form-model-item label="WebUI端口">
                  <a-input-number :min="80" :max="65535" v-model="settings.form.webui_port"/>
                  <span style="color: #909399;font-size: 9px;margin-left: 5px ;">Aria2NG端口</span>
                </a-form-model-item>
                <a-form-model-item label="磁盘模式">
                  <a-radio-group v-model="settings.form.fa">
                    <a-radio value="none">
                      None
                    </a-radio>
                    <a-radio value="falloc">
                      Falloc
                    </a-radio>
                    <a-radio value="trunc">
                      Trunc
                    </a-radio>
                    <a-radio value="prealloc">
                      Prealloc
                    </a-radio>
                  </a-radio-group>
                </a-form-model-item>
                <a-form-model-item label="自动同步">
                  <a-switch v-model="settings.form.rut">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">每天凌晨3点更新trackers</span>
                </a-form-model-item>
                <a-form-model-item label="种子保存">
                  <a-switch v-model="settings.form.smd">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">保存磁力链接为种子文件</span>
                </a-form-model-item>
                <a-form-model-item label="对外访问">
                  <a-switch v-model="settings.form.public_visit">
                    <a-icon slot="checkedChildren" type="check" />
                    <a-icon slot="unCheckedChildren" type="close" />
                  </a-switch>
                  <span style="color: #909399;font-size: 11px;margin-left: 5px ;">允许从外部访问Aria2和Aria2NG</span>
                </a-form-model-item>
                <a-form-model-item label="自动清理">
                  <a-input type="number" v-model="settings.form.auto_clean" style="width: 300px" :min="0" addon-after="MB"></a-input>
                  <span style="color: #909399;font-size: 9px;margin-left: 5px ;">清理BT不满足大小的文件, 0为关闭</span>
                </a-form-model-item>
                <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
                  <a-button type="primary" :loading="form_loading" @click="submit_update_settings">立即更新</a-button>
                </a-form-model-item>
              </a-form-model>
            </a-card>
          </a-tab-pane>
          <a-tab-pane key="3" tab="Aria2配置">
            <a-table :columns="aria2.columns"  :data-source="aria2.options" size="small" style="background: #FFFFFF">
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-col>
    </a-row>
    <a-modal title="新建下载" v-model="download.create.visit">
      <template slot="footer">
        <a-button @click="download.create.visit = false">取 消</a-button>
        <a-button type="primary" @click="submit_create_task">确 定</a-button>
      </template>
      <a-form-model :model="download.create.form" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item label="地址">
          <a-input type="textarea" v-model="download.create.form.url" :rows="3"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
    <a-drawer :withHeader="false" :visible="task.visible" @close="taskInfoClose" :closable="false" width="60%">
      <a-descriptions bordered size="small">
        <a-descriptions-item :span="3" label="文件名">
          {{task.info.filename}}
        </a-descriptions-item>
        <a-descriptions-item :span="3" label="GID">
          {{task.info.status.gid}}
        </a-descriptions-item>
        <a-descriptions-item v-if="task.info.status.errorMessage != ''" :span="3" label="错误">
          {{task.info.status.errorMessage}}
        </a-descriptions-item>
        <a-descriptions-item label="大小">
          {{task.info.status.totalLength | diskSize}}
        </a-descriptions-item>
        <a-descriptions-item :span="2" label="状态">
          {{task.info.status.status}}
        </a-descriptions-item>
        <a-descriptions-item label="已上传">
          {{task.info.status.uploadLength | diskSize}}
        </a-descriptions-item>
        <a-descriptions-item label="已完成">
          {{task.info.status.completedLength | diskSize}}
        </a-descriptions-item>
      </a-descriptions>
      <a-table :columns="task.info.columns" :data-source="task.info.status.files" size="small" style="margin-top: 10px;background: #FFFFFF" max-height="500">
        <span slot="path" slot-scope="text">
          {{text.split("/").slice(-1)[0]}}
        </span>
        <span slot="completedLength" slot-scope="text, record">
          {{(record.completedLength / record.length * 100).toFixed(2)}}%
        </span>
        <span slot="length" slot-scope="text, record">
          {{record.length | diskSize}}
        </span>
        <span slot="action" slot-scope="text, record">
          <a-button style="color: red" type="link" icon="download" @click="download_file(record)"></a-button>
        </span>
      </a-table>
    </a-drawer>
  </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      torrent: {
        sync: false
      },
      endpoint: "",
      connection: null,
      form_loading: false,
      labelCol: { span: 4},
      wrapperCol: { span: 18 },
      aria2: {
        options: [],
        columns: [
          {title: '键', dataIndex: 'key', key: 'key'},
          {title: '值', dataIndex: 'value', key: 'value'},
          ]
      },
      task: {
        visible: false,
        query: {
          status: "全部"
        },
        info: {
          filename: "",
          status: {
            files: [],
          },
          columns: [
            {title: '文件名', dataIndex: 'path', key: 'path', scopedSlots: { customRender: 'path' }},
            {title: '进度', dataIndex: 'completedLength', key: 'completedLength', scopedSlots: { customRender: 'completedLength' }},
            {title: '大小', dataIndex: 'length', key: 'length', scopedSlots: { customRender: 'length'}, sorter: (a, b) => a.length - b.length,},
            {title: '操作', key: 'action',scopedSlots: { customRender: 'action' }},
          ]
        }
      },
      global: {
        upload: 0,
        download: 0
      },
      download: {
        tasks: [],
        uploadloading: false,
        columns: [
          {title: '文件名', dataIndex: 'gid', key: 'gid', scopedSlots: { customRender: 'gid' }},
          {title: '大小', dataIndex: 'totalLength', key: 'totalLength', scopedSlots: { customRender: 'totalLength' }},
          {title: '状态', dataIndex: 'status', key: 'status', scopedSlots: { customRender: 'status' }},
          {title: '速度', dataIndex: 'downloadSpeed', key: 'downloadSpeed', scopedSlots: { customRender: 'downloadSpeed' }},
          {title: '上传', dataIndex: 'uploadSpeed', key: 'uploadSpeed', scopedSlots: { customRender: 'uploadSpeed' }},
          {title: '操作', key: 'action',scopedSlots: { customRender: 'action' }},
        ],
        create: {
          visit: false,
          form: {
            url: undefined,
            path: "",
            thread_num: 0,
          },
        }
      },
      settings: {
        form: {
          auto_clean: 0,
          bt_port: 0,
          config_path: "",
          container_id: "",
          ctu: "",
          download_path: "",
          fa: true,
          port: 6800,
          public_visit: true,
          rut: true,
          smd: true,
          tcp_port: 6881,
          token: "",
          udp_port: 6881,
          ut: true,
          webui: false,
          webui_port: 8080,
        },
        nodes: []
      }
    }
  },
  methods: {
    taskInfoClose() {
      this.task.visible = false
    },
    taskInfoOpen(info) {
      let that = this
      this.task.info.filename = this.getTaskName(info)
      this.$api.aria2.task.status(info.gid).then(function (response) {
        that.task.info.status = response.detail
        that.task.visible = true
      }).catch(function(response) {
        that.$message.error('获取任务信息失败：'+response.message)
      })
    },
    submit_create_task () {
      let that = this
      this.$api.aria2.task.add_uri(this.download.create.form.url).then(function (response) {
        that.download.create.visit = false
        that.refresh_tasks()
      })
    },
    remove_task (item) {
      this.$api.aria2.task.remove(item.gid)
    },
    cancel_task (item) {
      this.$api.aria2.task.pause(item.gid)
    },
    start_task (item) {
      console.log(item)
      this.$api.aria2.task.unpause(item.gid).then(function(response){}).catch(function(response) {
        console.log(response)
      })
    },
    refresh_settings () {
      let that = this
      this.form_loading = true
      this.$api.aria2.settings.info().then(function (response) {
        that.settings.form = response.detail
        that.download.create.form.thread_num = response.detail.thread_num
        that.download.create.form.path = response.detail.path
        that.form_loading = false
      }).catch(function(response) {
        that.form_loading = false
      })
    },
    submit_update_settings () {
      let that = this
      this.form_loading = true
      this.$api.aria2.settings.update(this.settings.form).then(function (response) {
        that.$message.success('更新成功')
        that.form_loading = false
      }).catch(function (resp) {
        that.$message.error('更新失败')
        that.form_loading = false
      })
    },
    refresh_global_options() {
      let that = this, options = []
      this.$api.aria2.global.options().then(function (response) {
        that.aria2.options = response.detail
      }).catch(function (response) {
        that.$message.error('获取全局配置失败:' + response.message)
      })
    },
    tabClick: function (key) {
      if (key === '1') {
        if (!this.timer) {
          this.timer = setInterval(() => {
            this.connection.send("tasks")
            this.connection.send("stats")
          }, 3000)
        }
      } else if (key === '2') {
        if (this.timer) {
          clearInterval(this.timer)
          this.timer = undefined
        } 
      } else if (key === '3') {
        if (this.timer) {
          clearInterval(this.timer)
          this.timer = undefined
        }
        this.refresh_global_options()
      }
    },
    getTaskName: function(info) {
      let taskName = "Unknown"
      if (info.bittorrent && info.bittorrent.info) {
        taskName = info.bittorrent.info.name;
      }
      if (!taskName && info.files && info.files.length > 0) {
        taskName = this.getFileName(info.files[0]);
      }
      return taskName
    },
    uploadChange(info) {
      if (info.file.status !== 'uploading') {
        this.download.uploadloading = true
      }
      if (info.file.status === 'done') {
        this.download.uploadloading = false
      }
    },
    getFileName: function (file) {
      if (!file) {
          return '';
      }

      var path = file.path;
      var needUrlDecode = false;

      if (!path && file.uris && file.uris.length > 0) {
          path = file.uris[0].uri;
          needUrlDecode = true;
      }

      var index = path.lastIndexOf('/');

      if (index <= 0 || index === path.length) {
          return path;
      }

      var fileNameAndQueryString = path.substring(index + 1);
      var queryStringStartPos = fileNameAndQueryString.indexOf('?');
      var fileName = fileNameAndQueryString;

      if (queryStringStartPos > 0) {
          fileName = fileNameAndQueryString.substring(0, queryStringStartPos);
      }

      if (needUrlDecode) {
        fileName = decodeURI(fileName);
      }

      return fileName;
    },
    download_file(row) {
      window.open(`/api/download/file?file_index=${row.index}&gid=${this.task.info.status.gid}`, '_blank')
    },
    onOpen() {

    },
    onError(error) {
      console.log(error)
    },
    onmessage(evt) {
       let data = JSON.parse(evt.data)
      switch (data.type) {
        case "tasks":
          this.download.tasks = data.data
        case "stats":
          if (data.data.uploadSpeed == "") {
            this.global.upload = 0
          } else {
            this.global.upload = parseInt(data.data.uploadSpeed)
          }
          if (data.data.downloadSpeed == "") {
            this.global.download = 0
          } else {
            this.global.download = parseInt(data.data.downloadSpeed)
          }
      }
    },
  },
  created: function () {
    let that = this
    if (window.location.protocol === 'https:') {
        this.protocol = 'wss://'
    } else {
        this.protocol = 'ws://'
    }
    let host = this.$apihost == ""? window.location.host: this.$apihost
    this.endpoint = `${this.protocol}${host}/api/download/query`
    let ws = new WebSocket(this.endpoint)//后端接口位置
    this.connection = ws
    this.connection.onmessage = this.onmessage
    this.connection.onerror = this.onError


    this.$api.aria2.settings.info().then(function (response) {
      that.settings.form = response.detail
      that.download.create.form.thread_num = response.detail.thread_num
      that.download.create.form.path = response.detail.path
    }).catch(function(response) {
      that.$message.error("加载配置失败: "+response.message)
    })
    this.timer = setInterval(() => {
      this.connection.send("tasks")
      this.connection.send("stats")
    }, 3000)
    
  },
  beforeDestroy () {
    clearInterval(this.timer)
  },
  mounted: function () {}
};
</script>

<style>

</style>